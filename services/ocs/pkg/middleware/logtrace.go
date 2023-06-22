package middleware

import (
	"net/http"

	ocstracing "github.com/owncloud/ocis/v2/services/ocs/pkg/tracing"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// LogTrace Sets the initial trace in the ocs service.
func LogTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanOpts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
		}
		ctx, span := ocstracing.TraceProvider.Tracer("ocs").Start(r.Context(), r.URL.Path, spanOpts...)
		defer span.End()

		propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
