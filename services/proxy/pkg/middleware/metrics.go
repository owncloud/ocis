package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Instrumenter provides a middleware to create metrics
func Instrumenter(m metrics.Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			m.Requests.With(prometheus.Labels{"method": r.Method}).Inc()

			next.ServeHTTP(ww, r)

			m.Duration.With(prometheus.Labels{"method": r.Method}).Observe(float64(time.Since(start).Seconds()))
			if ww.Status() >= 500 {
				m.Errors.With(prometheus.Labels{"method": r.Method}).Inc()
			}
		})
	}
}
