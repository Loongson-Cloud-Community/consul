package hcp

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/agent/hcp/client"
)

const (
	testRefreshInterval = 100 * time.Millisecond
	testSinkServiceName = "test.telemetry_config_provider"
	testRaceSampleCount = 5000
)

var (
	// Test constants to verify inmem sink metrics.
	testMetricKeyFailure = testSinkServiceName + "." + strings.Join(internalMetricRefreshFailure, ".")
	testMetricKeySuccess = testSinkServiceName + "." + strings.Join(internalMetricRefreshSuccess, ".")
)

type testConfig struct {
	filters         string
	endpoint        string
	labels          map[string]string
	refreshInterval time.Duration
}

func TestNewTelemetryConfigProvider(t *testing.T) {
	for name, tc := range map[string]struct {
		inputs   *testConfig
		mock     func(*client.MockClient)
		emptyCfg *dynamicConfig
	}{
		"initWithoutFirstClientFetch": {
			mock: func(m *client.MockClient) {
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(nil, errors.New("failed to fetch"))
			},
			emptyCfg: &dynamicConfig{
				Labels: map[string]string{},
			},
		},
		"initWithFirstClientFetch": {
			mock: func(m *client.MockClient) {
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(nil, errors.New("failed to fetch"))
			},
			emptyCfg: &dynamicConfig{
				Labels: map[string]string{},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mockClient := client.NewMockClient(t)
			tc.mock(mockClient)

			cfgProvider := NewHCPProvider(ctx, mockClient)
			require.Equal(t, cfgProvider.cfg, expectedCfg)
		})
	}
}

func TestTelemetryConfigProviderGetUpdate(t *testing.T) {
	for name, tc := range map[string]struct {
		mockExpect func(*client.MockClient)
		metricKey  string
		optsInputs *testConfig
		expected   *testConfig
	}{
		"noChanges": {
			optsInputs: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: testRefreshInterval,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					endpoint: "http://test.com/v1/metrics",
					filters:  "test",
					labels: map[string]string{
						"test_label": "123",
					},
					refreshInterval: testRefreshInterval,
				})
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(mockCfg, nil)
			},
			expected: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				labels: map[string]string{
					"test_label": "123",
				},
				filters:         "test",
				refreshInterval: testRefreshInterval,
			},
			metricKey: testMetricKeySuccess,
		},
		"newConfig": {
			optsInputs: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: 2 * time.Second,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					endpoint: "http://newendpoint/v1/metrics",
					filters:  "consul",
					labels: map[string]string{
						"new_label": "1234",
					},
					refreshInterval: 2 * time.Second,
				})
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(mockCfg, nil)
			},
			expected: &testConfig{
				endpoint: "http://newendpoint/v1/metrics",
				filters:  "consul",
				labels: map[string]string{
					"new_label": "1234",
				},
				refreshInterval: 2 * time.Second,
			},
			metricKey: testMetricKeySuccess,
		},
		"sameConfigInvalidRefreshInterval": {
			optsInputs: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: testRefreshInterval,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					refreshInterval: 0 * time.Second,
				})
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(mockCfg, nil)
			},
			expected: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				labels: map[string]string{
					"test_label": "123",
				},
				filters:         "test",
				refreshInterval: testRefreshInterval,
			},
			metricKey: testMetricKeyFailure,
		},
		"sameConfigHCPClientFailure": {
			optsInputs: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: testRefreshInterval,
			},
			mockExpect: func(m *client.MockClient) {
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(nil, fmt.Errorf("failure"))
			},
			expected: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: testRefreshInterval,
			},
			metricKey: testMetricKeyFailure,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sink := initGlobalSink()
			mockClient := client.NewMockClient(t)
			tc.mockExpect(mockClient)

			dynamicCfg, err := testDynamicCfg(tc.optsInputs)
			require.NoError(t, err)

			provider := &hcpProviderImpl{
				hcpClient: mockClient,
				cfg:       dynamicCfg,
			}

			provider.getUpdate(context.Background())

			// Verify endpoint provider returns correct config values.
			require.Equal(t, tc.expected.endpoint, provider.GetEndpoint().String())
			require.Equal(t, tc.expected.filters, provider.GetFilters().String())
			require.Equal(t, tc.expected.labels, provider.GetLabels())

			// Verify count for transform success metric.
			interval := sink.Data()[0]
			require.NotNil(t, interval, 1)
			sv := interval.Counters[tc.metricKey]
			assert.NotNil(t, sv.AggregateSample)
			require.Equal(t, sv.AggregateSample.Count, 1)
		})
	}
}

func TestTelemetryConfigProvider_Race(t *testing.T) {
	cfg, err := testTelemetryCfg(&testConfig{
		endpoint: "http://test.com/v1/metrics",
		filters:  "test",
		labels: map[string]string{
			"test_label": "123",
		},
		refreshInterval: testRefreshInterval,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := client.NewMockClient(t)
	mockClient.EXPECT().FetchTelemetryConfig(mock.Anything).Return(cfg, nil)

	// Start the provider goroutine
	// Every refresh interval, config will be modified.
	provider, err := NewHCPProvider(ctx, mockClient, cfg)
	require.NoError(t, err)

	// Every refresh interval, try to query config using Get* methods inducing a race condition.
	timer := time.NewTimer(testRefreshInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			wg := &sync.WaitGroup{}
			// Start goroutines that try to access label configuration.
			kickOff(wg, testRaceSampleCount, provider, func(provider *hcpProviderImpl) {
				require.Equal(t, provider.GetLabels(), cfg.MetricsConfig.Labels)
			})

			// Start goroutines that try to access endpoint configuration.
			kickOff(wg, testRaceSampleCount, provider, func(provider *hcpProviderImpl) {
				require.Equal(t, provider.GetFilters(), cfg.MetricsConfig.Filters)
			})

			// Start goroutines that try to access filter configuration.
			kickOff(wg, testRaceSampleCount, provider, func(provider *hcpProviderImpl) {
				require.Equal(t, provider.GetEndpoint(), cfg.MetricsConfig.Endpoint)
			})

			wg.Wait()
		// Stop after 10 refresh intervals.
		case <-time.After(10 * testRefreshInterval):
			return
		case <-ctx.Done():
			require.Fail(t, "Context cancelled before test finishes")
			return
		}
	}
}

func kickOff(wg *sync.WaitGroup, count int, provider *hcpProviderImpl, check func(cfgProvider *hcpProviderImpl)) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			check(provider)
		}()
	}
}

// initGlobalSink is a helper function to initialize a Go metrics inmemsink.
func initGlobalSink() *metrics.InmemSink {
	cfg := metrics.DefaultConfig(testSinkServiceName)
	cfg.EnableHostname = false

	sink := metrics.NewInmemSink(10*time.Second, 10*time.Second)
	metrics.NewGlobal(cfg, sink)

	return sink
}

// testDynamicCfg converts testConfig inputs to a dynamicConfig to be used in tests.
func testDynamicCfg(testCfg *testConfig) (*dynamicConfig, error) {
	filters, err := regexp.Compile(testCfg.filters)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(testCfg.endpoint)
	if err != nil {
		return nil, err
	}
	return &dynamicConfig{
		Endpoint:        endpoint,
		Filters:         filters,
		Labels:          testCfg.labels,
		RefreshInterval: testCfg.refreshInterval,
	}, nil
}

// testTelemetryCfg converts testConfig inputs to a TelemetryConfig to be used in tests.
func testTelemetryCfg(testCfg *testConfig) (*client.TelemetryConfig, error) {
	filters, err := regexp.Compile(testCfg.filters)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(testCfg.endpoint)
	if err != nil {
		return nil, err
	}
	return &client.TelemetryConfig{
		MetricsConfig: &client.MetricsConfig{
			Endpoint: endpoint,
			Filters:  filters,
			Labels:   testCfg.labels,
		},
		RefreshConfig: &client.RefreshConfig{
			RefreshInterval: testCfg.refreshInterval,
		},
	}, nil
}
