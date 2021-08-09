package middleware

import (
	"net/http"

	ocstracing "github.com/owncloud/ocis/ocs/pkg/tracing"
	"go.opentelemetry.io/otel/propagation"
)

var propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// LogTrace Sets the initial trace in the ocs service.
func LogTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := ocstracing.TraceProvider.Tracer("ocs").Start(r.Context(), r.URL.Path)
		defer span.End()

		propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
