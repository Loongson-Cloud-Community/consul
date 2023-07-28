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
	for name, tc := range map[string]struct {
		inputs          *testConfig
		defaultEmptyCfg *dynamicConfig
	}{
		"initEmptyCfg": {
			defaultEmptyCfg: &dynamicConfig{
				Labels:          map[string]string{},
				Filters:         defaultTelemetryConfigFilters,
				RefreshInterval: defaultTelemetryConfigRefreshInterval,
			},
		},
		"initFirstClientFetch": {
			inputs: &testConfig{
				endpoint:        "http://test.com",
				labels:          map[string]string{"test": "123"},
				filters:         "test",
				refreshInterval: 1 * time.Second,
				enabled:         true,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mockClient := client.NewMockClient(t)
			if tc.defaultEmptyCfg != nil {
				mockClient.EXPECT().FetchTelemetryConfig(mock.Anything).Return(nil, errors.New("failed to fetch"))
				cfgProvider := NewHCPProvider(ctx, mockClient)
				require.Equal(t, cfgProvider.cfg, tc.defaultEmptyCfg)
				return
			}

			telemetryCfg, err := testTelemetryCfg(tc.inputs)
			require.NoError(t, err)
			mockClient.EXPECT().FetchTelemetryConfig(mock.Anything).Return(telemetryCfg, nil)

			expectedDynamicCfg, err := testDynamicCfg(tc.inputs)
			require.NoError(t, err)

			cfgProvider := NewHCPProvider(ctx, mockClient)
			require.Equal(t, expectedDynamicCfg, cfgProvider.cfg)
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

func TestDynamicConfigEquals(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		a        *testConfig
		b        *testConfig
		expected bool
	}{
		"same": {
			a: &testConfig{
				endpoint:        "test.com",
				filters:         "state|raft",
				labels:          map[string]string{"test": "123"},
				refreshInterval: 1 * time.Second,
				enabled:         true,
			},
			b: &testConfig{
				endpoint:        "test.com",
				filters:         "state|raft",
				labels:          map[string]string{"test": "123"},
				refreshInterval: 1 * time.Second,
				enabled:         true,
			},
			expected: true,
		},
		"different": {
			a: &testConfig{
				endpoint:        "newendpoint.com",
				filters:         "state|raft|extra",
				labels:          map[string]string{"test": "12334"},
				refreshInterval: 2 * time.Second,
			},
			b: &testConfig{
				endpoint:        "test.com",
				filters:         "state|raft",
				labels:          map[string]string{"test": "123"},
				refreshInterval: 1 * time.Second,
			},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			aCfg, err := testDynamicCfg(tc.a)
			require.NoError(t, err)
			bCfg, err := testDynamicCfg(tc.b)
			require.NoError(t, err)

			equal, err := aCfg.equals(bCfg)

			require.NoError(t, err)
			require.Equal(t, tc.expected, equal)
		})
	}
}

// mockClientRace returns new configuration everytime checkUpdate is called
// by creating unique labels using its counter. This allows us to induce
// race conditions with the changing dynamic config in the provider.
type mockClientRace struct {
	counter         int
	defaultEndpoint *url.URL
	defaultFilters  *regexp.Regexp
}

func (mc *mockClientRace) FetchBootstrap(ctx context.Context) (*client.BootstrapConfig, error) {
	return nil, nil
}
func (mc *mockClientRace) PushServerStatus(ctx context.Context, status *client.ServerStatus) error {
	return nil
}
func (mc *mockClientRace) DiscoverServers(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (mc *mockClientRace) FetchTelemetryConfig(ctx context.Context) (*client.TelemetryConfig, error) {
	return &client.TelemetryConfig{
		MetricsConfig: &client.MetricsConfig{
			Endpoint: mc.defaultEndpoint,
			Filters:  mc.defaultFilters,
			// Generate unique labels.
			Labels: map[string]string{fmt.Sprintf("label_%d", mc.counter): fmt.Sprintf("value_%d", mc.counter)},
		},
		RefreshConfig: &client.RefreshConfig{
			RefreshInterval: testRefreshInterval,
		},
	}, nil
}

func TestTelemetryConfigProvider_Race(t *testing.T) {
	testCfg := &testConfig{
		endpoint: "http://test.com/v1/metrics",
		filters:  "test",
		labels:   map[string]string{"test": "123"},
	}

	telemetryCfg, err := testTelemetryCfg(testCfg)
	require.NoError(t, err)

	dynamicCfg, err := testDynamicCfg(testCfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(testRefreshInterval)
	provider := &hcpProviderImpl{
		hcpClient: &mockClientRace{
			defaultEndpoint: telemetryCfg.MetricsConfig.Endpoint,
			defaultFilters:  telemetryCfg.MetricsConfig.Filters,
		},
		ticker: ticker,
		cfg:    dynamicCfg,
	}

	go provider.run(ctx, ticker.C)

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

	expectedEndpoint := telemetryCfg.MetricsConfig.Endpoint.String()
	endpointErr := fmt.Errorf("expected endpoint to be %s", expectedEndpoint)
	endpointErrCh := make(chan error, testRaceSampleCount)
	// Start 5000 goroutines that try to access endpoint configuration.
	kickOff(wg, endpointErrCh, provider, func(provider *hcpProviderImpl) bool {
		return provider.GetEndpoint().String() == expectedEndpoint
	}, endpointErr)

	expectedFilters := telemetryCfg.MetricsConfig.Filters.String()
	filtersErr := fmt.Errorf("expected filters to be %s", expectedFilters)
	filtersErrCh := make(chan error, testRaceSampleCount)
	// Start 5000 goroutines that try to access filter configuration.
	kickOff(wg, filtersErrCh, provider, func(provider *hcpProviderImpl) bool {
		return provider.GetFilters().String() == expectedFilters
	}, filtersErr)

	wg.Wait()

	require.Empty(t, labelErrCh)
	require.Empty(t, endpointErrCh)
	require.Empty(t, filtersErrCh)
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
		Enabled:         testCfg.enabled,
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
