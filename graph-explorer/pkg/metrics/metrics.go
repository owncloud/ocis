package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "graph-explorer"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	// Counter  *prometheus.CounterVec
	BuildInfo *prometheus.GaugeVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		BuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "build_info",
			Help:      "Build information",
		}, []string{"version"}),
	}

	_ = prometheus.Register(m.BuildInfo)
	// TODO: implement metrics
	return m
}
