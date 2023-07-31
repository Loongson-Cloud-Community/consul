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
	enabled         bool
}

func TestNewTelemetryConfigProvider(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		mock            func(*client.MockClient)
		expectedTestCfg *testConfig
	}{
		"initWithDefaultConfig": {
			mock: func(mockClient *client.MockClient) {
				mockClient.EXPECT().FetchTelemetryConfig(mock.Anything).Return(nil, errors.New("failed to fetch config"))
			},
			expectedTestCfg: &testConfig{
				refreshInterval: defaultTelemetryConfigRefreshInterval,
				filters:         client.DefaultMetricFilters.String(),
				labels:          map[string]string{},
				enabled:         false,
				endpoint:        "",
			},
		},
		"initWithFirstUpdate": {
			mock: func(mockClient *client.MockClient) {
				cfg, err := testTelemetryCfg(&testConfig{
					refreshInterval: 2 * time.Second,
					filters:         "consul",
					labels:          map[string]string{"test_label": "123"},
					endpoint:        "https://test.com/v1/metrics",
				})
				require.NoError(t, err)
				mockClient.EXPECT().FetchTelemetryConfig(mock.Anything).Return(cfg, nil)
			},
			expectedTestCfg: &testConfig{
				refreshInterval: 2 * time.Second,
				filters:         "consul",
				labels:          map[string]string{"test_label": "123"},
				enabled:         true,
				endpoint:        "https://test.com/v1/metrics",
			},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mc := client.NewMockClient(t)
			tc.mock(mc)

			provider := NewHCPProvider(ctx, mc)

			expectedCfg, err := testDynamicCfg(tc.expectedTestCfg)
			require.NoError(t, err)

			require.Equal(t, expectedCfg, provider.cfg)
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
				enabled:         true,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					endpoint: "http://test.com/v1/metrics",
					filters:  "test",
					labels: map[string]string{
						"test_label": "123",
					},
					refreshInterval: testRefreshInterval,
					enabled:         true,
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
				enabled:         true,
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
				enabled:         true,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					endpoint: "http://newendpoint/v1/metrics",
					filters:  "consul",
					labels: map[string]string{
						"new_label": "1234",
					},
					refreshInterval: 2 * time.Second,
					enabled:         true,
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
				enabled:         true,
			},
			metricKey: testMetricKeySuccess,
		},
		"newConfigMetricsDisabled": {
			optsInputs: &testConfig{
				endpoint: "http://test.com/v1/metrics",
				filters:  "test",
				labels: map[string]string{
					"test_label": "123",
				},
				refreshInterval: 2 * time.Second,
				enabled:         true,
			},
			mockExpect: func(m *client.MockClient) {
				mockCfg, _ := testTelemetryCfg(&testConfig{
					endpoint: "",
					filters:  "consul",
					labels: map[string]string{
						"new_label": "1234",
					},
					refreshInterval: 2 * time.Second,
				})
				m.EXPECT().FetchTelemetryConfig(mock.Anything).Return(mockCfg, nil)
			},
			expected: &testConfig{
				endpoint: "",
				filters:  "consul",
				labels: map[string]string{
					"new_label": "1234",
				},
				refreshInterval: 2 * time.Second,
				enabled:         false,
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
				enabled:         true,
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
				enabled:         true,
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
				enabled:         true,
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
				enabled:         true,
			},
			metricKey: testMetricKeyFailure,
		},
	} {
		tc := tc
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
			expectedCfg, err := testDynamicCfg(tc.expected)
			require.NoError(t, err)

			require.Equal(t, expectedCfg.Endpoint, provider.GetEndpoint())
			require.Equal(t, expectedCfg.Filters, provider.GetFilters())
			require.Equal(t, expectedCfg.Labels, provider.GetLabels())
			require.Equal(t, expectedCfg.Enabled, provider.Enabled())

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
	provider := NewHCPProvider(ctx, mockClient)

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

	var endpoint *url.URL
	if testCfg.endpoint != "" {
		u, err := url.Parse(testCfg.endpoint)
		if err != nil {
			return nil, err
		}
		endpoint = u
	}
	return &dynamicConfig{
		Endpoint:        endpoint,
		Filters:         filters,
		Labels:          testCfg.labels,
		RefreshInterval: testCfg.refreshInterval,
		Enabled:         testCfg.enabled,
	}, nil
}

// testTelemetryCfg converts testConfig inputs to a TelemetryConfig to be used in tests.
func testTelemetryCfg(testCfg *testConfig) (*client.TelemetryConfig, error) {
	filters, err := regexp.Compile(testCfg.filters)
	if err != nil {
		return nil, err
	}

	var endpoint *url.URL
	if testCfg.endpoint != "" {
		u, err := url.Parse(testCfg.endpoint)
		if err != nil {
			return nil, err
		}
		endpoint = u
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
