package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// GetOtelhttpMiddleware gets a new tracing middleware based on otelhttp
// to trace the requests.
func GetOtelhttpMiddleware(service string, tp trace.TracerProvider) func(http.Handler) http.Handler {
	return otelhttp.NewMiddleware(
		service,
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithPropagators(tracing.GetPropagator()),
		otelhttp.WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
		otelhttp.WithSpanNameFormatter(func(name string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// GetOtelhttpClient will get a new HTTP client that will use telemetry and
// automatically set the telemetry headers. It will wrap the default transport
// in order to use telemetry.
func GetOtelhttpClient(tp trace.TracerProvider) *http.Client {
	return &http.Client{
		Transport: GetOtelhttpClientTransport(http.DefaultTransport, tp),
	}
}

// GetOtelhttpClientTransport will get a new wrapped transport that will
// include telemetry automatically. You can use the http.DefaultTransport
// as base transport
func GetOtelhttpClientTransport(baseTransport http.RoundTripper, tp trace.TracerProvider) http.RoundTripper {
	return otelhttp.NewTransport(
		baseTransport,
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithPropagators(tracing.GetPropagator()),
		otelhttp.WithSpanNameFormatter(func(name string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}
