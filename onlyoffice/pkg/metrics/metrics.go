package metrics

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "onlyoffice"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	// Counter  *prometheus.CounterVec
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
	}

	// prometheus.Register(
	// 	m.Counter,
	// )

	return m
}
