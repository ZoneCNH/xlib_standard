package templatex

const (
	MetricClientCreatedTotal    = "client_created_total"
	MetricClientClosedTotal     = "client_closed_total"
	MetricClientErrorsTotal     = "client_errors_total"
	MetricClientHealthStatus    = "client_health_status"
	MetricClientHealthLatencyMS = "client_health_latency_ms"
)

// Metrics defines the interface for recording client metrics.
type Metrics interface {
	IncrCounter(name string, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
}

// NoopMetrics is a Metrics implementation that does nothing.
type NoopMetrics struct{}

func (NoopMetrics) IncrCounter(name string, labels map[string]string)                     {}
func (NoopMetrics) SetGauge(name string, value float64, labels map[string]string)         {}
func (NoopMetrics) ObserveHistogram(name string, value float64, labels map[string]string) {}
