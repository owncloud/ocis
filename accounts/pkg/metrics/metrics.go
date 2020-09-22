package metrics

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "accounts"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	// Counter  *prometheus.CounterVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
	}
	// TODO: implement metrics
	return m
}
