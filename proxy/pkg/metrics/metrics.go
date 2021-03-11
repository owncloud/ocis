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
	Counter   *prometheus.CounterVec
	Latency   *prometheus.SummaryVec
	Duration  *prometheus.HistogramVec
	BuildInfo *prometheus.GaugeVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "proxy_total",
			Help:      "How many proxy requests processed",
		}, []string{}),
		Latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "proxy_latency_microseconds",
			Help:      "proxy request latencies in microseconds",
		}, []string{}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "proxy_duration_seconds",
			Help:      "proxy method request time in seconds",
		}, []string{}),
		BuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "build_info",
			Help:      "Build Information",
		}, []string{"versions"}),
	}

	_ = prometheus.Register(m.Counter)
	_ = prometheus.Register(m.Latency)
	_ = prometheus.Register(m.Duration)
	_ = prometheus.Register(m.BuildInfo)
	return m
}
