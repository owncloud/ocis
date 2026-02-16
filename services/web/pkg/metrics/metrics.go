package metrics

// Metrics defines the available metrics of this service.
type Metrics struct {
	// Counter  *prometheus.CounterVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{}

	return m
}
