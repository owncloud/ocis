package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"go-micro.dev/v4/metadata"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// GetOtelhttpMiddleware gets a new tracing middleware based on otelhttp
// to trace the requests.
// This middleware will use the otelhttp middleware and then store the
// incoming data into the go-micro's metadata so it can be propagated through
// go-micro.
func GetOtelhttpMiddleware(service string, tp trace.TracerProvider) func(http.Handler) http.Handler {
	otelMid := otelhttp.NewMiddleware(
		service,
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithPropagators(tracing.GetPropagator()),
		otelhttp.WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
		otelhttp.WithSpanNameFormatter(func(name string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)

	httpToMicroMid := otelhttpToGoMicroGrpc()

	return func(next http.Handler) http.Handler {
		return otelMid(httpToMicroMid(next))
	}
}

func otelhttpToGoMicroGrpc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			propagator := tracing.GetPropagator()

			// based on go-micro plugin for opentelemetry
			// inject telemetry data into go-micro's metadata
			// in order to propagate the info to go-micro's calls
			md := make(metadata.Metadata)
			carrier := make(propagation.MapCarrier)
			propagator.Inject(ctx, carrier)
			for k, v := range carrier {
				md.Set(k, v)
			}
			mdCtx := metadata.NewContext(ctx, md)
			r = r.WithContext(mdCtx)

			next.ServeHTTP(w, r)
		})
	}
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
