package middleware

import (
	"fmt"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	pkgtrace "github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	proxytracing "github.com/owncloud/ocis/v2/services/proxy/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracer provides a middleware to start traces
func Tracer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &tracer{
			next: next,
		}
	}
}

type tracer struct {
	next http.Handler
}

func (m tracer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		span trace.Span
	)

	tracer := proxytracing.TraceProvider.Tracer("proxy")
	spanOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span = tracer.Start(ctx, fmt.Sprintf("%s %v", r.Method, r.URL.Path), spanOpts...)
	defer span.End()

	span.SetAttributes(
		attribute.KeyValue{
			Key:   "x-request-id",
			Value: attribute.StringValue(chimiddleware.GetReqID(r.Context())),
		})

	pkgtrace.Propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	m.next.ServeHTTP(w, r.WithContext(ctx))
}
