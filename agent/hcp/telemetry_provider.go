package hcp

import (
	"context"
<<<<<<< HEAD
=======
	"errors"
	"fmt"
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/armon/go-metrics"
<<<<<<< HEAD
	"github.com/mitchellh/hashstructure/v2"

	"github.com/hashicorp/go-hclog"
=======
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/hashstructure/v2"
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh

	"github.com/hashicorp/consul/agent/hcp/client"
	"github.com/hashicorp/consul/agent/hcp/telemetry"
)

var (
	// internalMetricRefreshFailure is a metric to monitor refresh failures.
	internalMetricRefreshFailure []string = []string{"hcp", "telemetry_config_provider", "refresh", "failure"}
	// internalMetricRefreshFailure is a metric to monitor refresh successes.
	internalMetricRefreshSuccess []string = []string{"hcp", "telemetry_config_provider", "refresh", "success"}
<<<<<<< HEAD

	defaultTelemetryConfigRefreshInterval = 1 * time.Minute
	defaultTelemetryConfigFilters         = regexp.MustCompile(".+")
=======
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
)

// Ensure hcpProviderImpl implements telemetry provider interfaces.
var _ telemetry.ConfigProvider = &hcpProviderImpl{}
var _ telemetry.EndpointProvider = &hcpProviderImpl{}

// hcpProviderImpl holds telemetry configuration and settings for continuous fetch of new config from HCP.
// it updates configuration, if changes are detected.
type hcpProviderImpl struct {
<<<<<<< HEAD
	// hcpClient is an authenticated client used to make HTTP requests to HCP.
	hcpClient client.Client

	// cfg holds configuration that can be dynamically updated.
	cfg *dynamicConfig

	// updateTickerCh is a test channel for triggering updates manually during testing.
	updateTickerCh <-chan time.Time

=======
	// cfg holds configuration that can be dynamically updated.
	cfg *dynamicConfig

>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
	// A reader-writer mutex is used as the provider is read heavy.
	// OTEL components access telemetryConfig during metrics collection and export (read).
	// Meanwhile, config is only updated when there are changes (write).
	rw sync.RWMutex
<<<<<<< HEAD
=======
	// hcpClient is an authenticated client used to make HTTP requests to HCP.
	hcpClient client.Client
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
}

// dynamicConfig is a set of configurable settings for metrics collection, processing and export.
// fields MUST be exported to compute hash for equals method.
type dynamicConfig struct {
	Endpoint *url.URL
	Labels   map[string]string
	Filters  *regexp.Regexp
	// refreshInterval controls the interval at which configuration is fetched from HCP to refresh config.
	RefreshInterval time.Duration
}

func (d *dynamicConfig) equals(newCfg *dynamicConfig) (bool, error) {
	currHash, err := hashstructure.Hash(*d, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}

	newHash, err := hashstructure.Hash(*newCfg, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}

	return currHash == newHash, err
}

<<<<<<< HEAD
// NewHCPProviderImpl initializes and starts a HCP Telemetry provider with provided params.
func NewHCPProviderImpl(ctx context.Context, hcpClient client.Client) *hcpProviderImpl {
	return newHCPProviderImpl(ctx, hcpClient, make(<-chan time.Time))
}

func newHCPProviderImpl(ctx context.Context, hcpClient client.Client, updateTickerCh <-chan time.Time) *hcpProviderImpl {
	t := &hcpProviderImpl{
		hcpClient:      hcpClient,
		updateTickerCh: updateTickerCh,
	}
	go t.run(ctx)

	return t
}

// run continuously checks for updates to the telemetry configuration by making a request to HCP.
// Modification of config only occurs if changes are detected to decrease write locks that block read locks.
func (t *hcpProviderImpl) run(ctx context.Context) {
	ticker := time.NewTicker(defaultTelemetryConfigRefreshInterval)

	// Try to initialize the config.
	newCfg, _ := t.checkUpdate(ctx)
	if newCfg != nil {
		t.modifyTelemetryConfig(newCfg)
		ticker.Reset(newCfg.RefreshInterval)
	}

	for {
		select {
		case <-ticker.C:
		case <-t.updateTickerCh:
			// Check whether we need to update the config.
=======
// providerParams is used to initialize a hcpProviderImpl.
type providerParams struct {
	metricsConfig   *client.MetricsConfig
	refreshInterval time.Duration
	hcpClient       client.Client
}

// NewHCPProvider initializes and starts a HCP Telemetry provider with provided params.
func NewHCPProvider(ctx context.Context, params *providerParams) (*hcpProviderImpl, error) {
	if params.hcpClient == nil {
		return nil, errors.New("missing HCP client")
	}

	if params.metricsConfig == nil {
		return nil, errors.New("missing metrics config")
	}

	if params.refreshInterval <= 0 {
		return nil, fmt.Errorf("invalid refresh interval: %d", params.refreshInterval)
	}

	cfg := &dynamicConfig{
		Endpoint:        params.metricsConfig.Endpoint,
		Labels:          params.metricsConfig.Labels,
		Filters:         params.metricsConfig.Filters,
		RefreshInterval: params.refreshInterval,
	}

	t := &hcpProviderImpl{
		cfg:       cfg,
		hcpClient: params.hcpClient,
	}

	go t.run(ctx)

	return t, nil
}

// run continously checks for updates to the telemetry configuration by making a request to HCP.
// Modification of config only occurs if changes are detected to decrease write locks that block read locks.
func (t *hcpProviderImpl) run(ctx context.Context) {
	ticker := time.NewTicker(t.cfg.RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
			if newCfg, sameCfg := t.checkUpdate(ctx); !sameCfg {
				t.modifyTelemetryConfig(newCfg)
				ticker.Reset(newCfg.RefreshInterval)
			}
		case <-ctx.Done():
			return
		}
	}
}

// checkUpdate makes a HTTP request to HCP to return a new metrics configuration and true, if config changed.
// checkUpdate does not update the metricsConfig field to prevent acquiring the write lock unnecessarily.
func (t *hcpProviderImpl) checkUpdate(ctx context.Context) (*dynamicConfig, bool) {
	logger := hclog.FromContext(ctx).Named("telemetry_config_provider")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	telemetryCfg, err := t.hcpClient.FetchTelemetryConfig(ctx)
	if err != nil {
		logger.Error("failed to fetch telemetry config from HCP", "error", err)
		metrics.IncrCounter(internalMetricRefreshFailure, 1)
		return nil, false
	}

	newDynamicConfig := &dynamicConfig{
		Filters:         telemetryCfg.MetricsConfig.Filters,
		Endpoint:        telemetryCfg.MetricsConfig.Endpoint,
		Labels:          telemetryCfg.MetricsConfig.Labels,
		RefreshInterval: telemetryCfg.RefreshConfig.RefreshInterval,
	}

	t.rw.RLock()
	defer t.rw.RUnlock()

	equal, err := t.cfg.equals(newDynamicConfig)
	if err != nil {
		logger.Error("failed to calculate hash for new config", "error", err)
		metrics.IncrCounter(internalMetricRefreshFailure, 1)
		return nil, false
	}

	metrics.IncrCounter(internalMetricRefreshSuccess, 1)

	return newDynamicConfig, equal
}

// modifynewTelemetryConfig acquires a write lock to modify it with a given newTelemetryConfig object.
func (t *hcpProviderImpl) modifyTelemetryConfig(newCfg *dynamicConfig) {
	t.rw.Lock()
	defer t.rw.Unlock()

	t.cfg = newCfg
}

// GetEndpoint acquires a read lock to return endpoint configuration for consumers.
<<<<<<< HEAD
func (t *hcpProviderImpl) GetEndpoint() (*url.URL, bool) {
	t.rw.RLock()
	defer t.rw.RUnlock()

	if t.cfg == nil {
		return nil, false
	}

	return t.cfg.Endpoint, true
=======
func (t *hcpProviderImpl) GetEndpoint() *url.URL {
	t.rw.RLock()
	defer t.rw.RUnlock()

	return t.cfg.Endpoint
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
}

// GetFilters acquires a read lock to return filters configuration for consumers.
func (t *hcpProviderImpl) GetFilters() *regexp.Regexp {
	t.rw.RLock()
	defer t.rw.RUnlock()

<<<<<<< HEAD
	if t.cfg == nil {
		return defaultTelemetryConfigFilters
	}

=======
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
	return t.cfg.Filters
}

// GetLabels acquires a read lock to return labels configuration for consumers.
func (t *hcpProviderImpl) GetLabels() map[string]string {
	t.rw.RLock()
	defer t.rw.RUnlock()

<<<<<<< HEAD
	if t.cfg == nil {
		return map[string]string{}
	}

=======
>>>>>>> cc-4960/hcp-telemetry-periodic-refresh
	return t.cfg.Labels
}
