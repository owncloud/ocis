package prometheus

import (
	"context"

	"github.com/asim/go-micro/v3/server"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Namespace defines the namespace of the defined metrics.
	Namespace = "ocis"
)

// NewHandlerWrapper initializes the prometheus handler wrapper.
func NewHandlerWrapper(opts ...server.Option) server.HandlerWrapper {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "request_total",
			Help:      "How many service requests processed",
		},
		[]string{"method", "status"},
	)

	latency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: Namespace,
			Name:      "upstream_latency_microseconds",
			Help:      "Service method latencies in microseconds",
		},
		[]string{"method"},
	)

	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "request_duration_seconds",
			Help:      "Service method request time in seconds",
		},
		[]string{"method"},
	)

	prometheus.Register(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	prometheus.Register(
		prometheus.NewGoCollector(),
	)

	prometheus.Register(
		counter,
	)

	prometheus.Register(
		latency,
	)

	prometheus.Register(
		duration,
	)

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := req.Endpoint()

			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				us := v * 1000000

				latency.WithLabelValues(name).Observe(us)
				duration.WithLabelValues(name).Observe(v)
			}))

			defer timer.ObserveDuration()

			err := fn(ctx, req, rsp)

			if err == nil {
				counter.WithLabelValues(name, "success").Inc()
			} else {
				counter.WithLabelValues(name, "fail").Inc()
			}

			return err
		}
	}
}
