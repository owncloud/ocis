package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "idp"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	// Counter  *prometheus.CounterVec
	BuildInfo *prometheus.GaugeVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		// Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
		// 	Namespace: Namespace,
		// 	Subsystem: Subsystem,
		// 	Name:      "greet_total",
		// 	Help:      "How many greeting requests processed",
		// }, []string{}),
		BuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "build_info",
			Help:      "Build Information",
		}, []string{"version"}),
	}

	// prometheus.Register(
	// 	m.Counter,
	// )

	_ = prometheus.Register(
		m.BuildInfo,
	)

	return m
}
