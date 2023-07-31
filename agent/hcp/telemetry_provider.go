package hcp

import (
	"context"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/hashicorp/consul/agent/hcp/client"
	"github.com/hashicorp/consul/agent/hcp/telemetry"
)

var (
	// internalMetricRefreshFailure is a metric to monitor refresh failures.
	internalMetricRefreshFailure []string = []string{"hcp", "telemetry_config_provider", "refresh", "failure"}
	// internalMetricRefreshSuccess is a metric to monitor refresh successes.
	internalMetricRefreshSuccess          []string = []string{"hcp", "telemetry_config_provider", "refresh", "success"}
	defaultTelemetryConfigRefreshInterval          = 1 * time.Minute
	defaultTelemetryConfigFilters                  = regexp.MustCompile(".+")
)

// Ensure hcpProviderImpl implements telemetry provider interfaces.
var _ telemetry.ConfigProvider = &hcpProviderImpl{}
var _ telemetry.EndpointProvider = &hcpProviderImpl{}

// hcpProviderImpl holds telemetry configuration and settings for continuous fetch of new config from HCP.
// it updates configuration, if changes are detected.
type hcpProviderImpl struct {
	// cfg holds configuration that can be dynamically updated.
	cfg *dynamicConfig

	// A reader-writer mutex is used as the provider is read heavy.
	// OTEL components access telemetryConfig during metrics collection and export (read).
	// Meanwhile, config is only updated when there are changes (write).
	rw sync.RWMutex
	// hcpClient is an authenticated client used to make HTTP requests to HCP.
	hcpClient client.Client
}

// dynamicConfig is a set of configurable settings for metrics collection, processing and export.
// fields MUST be exported to compute hash for equals method.
type dynamicConfig struct {
	Enabled  bool
	Endpoint *url.URL
	Labels   map[string]string
	Filters  *regexp.Regexp
	// refreshInterval controls the interval at which configuration is fetched from HCP to refresh config.
	RefreshInterval time.Duration
}

// equals returns true if two dynamicConfig objects are equal.
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

// NewHCPProvider initializes and starts a HCP Telemetry provider with provided params.
func NewHCPProvider(ctx context.Context, hcpClient client.Client) *hcpProviderImpl {
	ticker := time.NewTicker(defaultTelemetryConfigRefreshInterval)
	t := &hcpProviderImpl{
		cfg: &dynamicConfig{
			Labels: make(map[string]string),
		},
		hcpClient: hcpClient,
		ticker:    ticker,
	}

	// Try to initialize the config once before running periodic fetcher.
	newCfg, _ := t.checkUpdate(ctx)
	if newCfg != nil {
		t.cfg = newCfg
	}

	go t.run(ctx, ticker.C)

	return t
}

// run continously checks for updates to the telemetry configuration by making a request to HCP.
func (h *hcpProviderImpl) run(ctx context.Context, refreshInterval time.Duration) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if newCfg := h.getUpdate(ctx); newCfg != nil {
				ticker.Reset(newCfg.RefreshInterval)
			}
		case <-ctx.Done():
			return
		}
	}
}

// getUpdate makes a HTTP request to HCP to return a new metrics configuration
// and updates the hcpProviderImpl.
func (h *hcpProviderImpl) getUpdate(ctx context.Context) *dynamicConfig {
	logger := hclog.FromContext(ctx).Named("telemetry_config_provider")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	telemetryCfg, err := h.hcpClient.FetchTelemetryConfig(ctx)
	if err != nil {
		logger.Error("failed to fetch telemetry config from HCP", "error", err)
		metrics.IncrCounter(internalMetricRefreshFailure, 1)
		return nil
	}

	// newRefreshInterval of 0 or less can cause ticker Reset() panic.
	newRefreshInterval := telemetryCfg.RefreshConfig.RefreshInterval
	if newRefreshInterval <= 0 {
		logger.Error("invalid refresh interval duration", "refreshInterval", newRefreshInterval)
		metrics.IncrCounter(internalMetricRefreshFailure, 1)
		return nil
	}

	newDynamicConfig := &dynamicConfig{
		Filters:         telemetryCfg.MetricsConfig.Filters,
		Endpoint:        telemetryCfg.MetricsConfig.Endpoint,
		Labels:          telemetryCfg.MetricsConfig.Labels,
		RefreshInterval: telemetryCfg.RefreshConfig.RefreshInterval,
		Enabled:         telemetryCfg.MetricsEnabled(),
	}

	// Acquire write lock to update new configuration.
	h.rw.Lock()
	defer h.rw.Unlock()

	h.cfg = newDynamicConfig

	metrics.IncrCounter(internalMetricRefreshSuccess, 1)

	return newDynamicConfig
}

// GetEndpoint acquires a read lock to return endpoint configuration for consumers.
func (h *hcpProviderImpl) GetEndpoint() *url.URL {
	h.rw.RLock()
	defer h.rw.RUnlock()

	return h.cfg.Endpoint
}

// GetFilters acquires a read lock to return filters configuration for consumers.
func (h *hcpProviderImpl) GetFilters() *regexp.Regexp {
	h.rw.RLock()
	defer h.rw.RUnlock()

	return h.cfg.Filters
}

// GetLabels acquires a read lock to return labels configuration for consumers.
func (h *hcpProviderImpl) GetLabels() map[string]string {
	h.rw.RLock()
	defer h.rw.RUnlock()

	return h.cfg.Labels
}

func (t *hcpProviderImpl) Enabled() bool {
	return t.cfg.Enabled
}
