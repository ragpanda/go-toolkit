package metrics

type MetricsHubConfig struct {
	ServiceName string
	BackendType MetricsBackendType

	ExpirationSec int64
}
