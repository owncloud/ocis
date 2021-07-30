package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"
)

// Trace unpacks the request context looking for an existing trace id.
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		propagator := propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		)

		ctx = propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
