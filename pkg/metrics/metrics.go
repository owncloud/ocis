package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "thumbnails"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	Counter  *prometheus.CounterVec
	Latency  *prometheus.SummaryVec
	Duration *prometheus.HistogramVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "getthumbnail_total",
			Help:      "How many GetThumbnail requests processed",
		}, []string{}),
		Latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "getthumbnail_latency_microseconds",
			Help:      "GetThumbnail request latencies in microseconds",
		}, []string{}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "getthumbnail_duration_seconds",
			Help:      "GetThumbnail method requests time in seconds",
		}, []string{}),
	}

	_ = prometheus.Register(
		m.Counter,
	)

	_ = prometheus.Register(
		m.Latency,
	)

	_ = prometheus.Register(
		m.Duration,
	)

	return m
}
