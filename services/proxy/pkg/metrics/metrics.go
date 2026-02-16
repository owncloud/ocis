package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "proxy"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	Requests  *prometheus.CounterVec
	Errors    *prometheus.CounterVec
	Duration  *prometheus.HistogramVec
	BuildInfo *prometheus.GaugeVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "requests_total",
			Help:      "How many requests processed in total",
		}, []string{"method"}),
		Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "errors_total",
			Help:      "How many requests run into errors",
		}, []string{"method"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "duration_seconds",
			Help:      "request duration in seconds",
		}, []string{"method"}),
		BuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "build_info",
			Help:      "Build Information",
		}, []string{"version"}),
	}

	// Initialize the metrics with 0
	m.Requests.WithLabelValues("GET").Add(0)
	m.Errors.WithLabelValues("GET").Add(0)

	_ = prometheus.Register(m.Requests)
	_ = prometheus.Register(m.Errors)
	_ = prometheus.Register(m.Duration)
	_ = prometheus.Register(m.BuildInfo)
	return m
}
