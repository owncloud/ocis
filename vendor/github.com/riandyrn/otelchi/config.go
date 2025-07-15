package otelchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// These defaults are used in `TraceHeaderConfig`.
const (
	DefaultTraceIDResponseHeaderKey      = "X-Trace-Id"
	DefaultTraceSampledResponseHeaderKey = "X-Trace-Sampled"
)

// config is used to configure the mux middleware.
type config struct {
	tracerProvider                oteltrace.TracerProvider
	propagators                   propagation.TextMapPropagator
	chiRoutes                     chi.Routes
	requestMethodInSpanName       bool
	filters                       []Filter
	traceIDResponseHeaderKey      string
	traceSampledResponseHeaderKey string
	publicEndpointFn              func(r *http.Request) bool
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// Filter is a predicate used to determine whether a given http.Request should
// be traced. A Filter must return true if the request should be traced.
type Filter func(*http.Request) bool

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		cfg.propagators = propagators
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.tracerProvider = provider
	})
}

// WithChiRoutes specified the routes that being used by application. Its main
// purpose is to provide route pattern as span name during span creation. If this
// option is not set, by default the span will be given name at the end of span
// execution. For some people, this behavior is not desirable since they want
// to override the span name on underlying handler. By setting this option, it
// is possible for them to override the span name.
func WithChiRoutes(routes chi.Routes) Option {
	return optionFunc(func(cfg *config) {
		cfg.chiRoutes = routes
	})
}

// WithRequestMethodInSpanName is used for adding http request method to span name.
// While this is not necessary for vendors that properly implemented the tracing
// specs (e.g Jaeger, AWS X-Ray, etc...), but for other vendors such as Elastic
// and New Relic this might be helpful.
//
// See following threads for details:
//
// - https://github.com/riandyrn/otelchi/pull/3#issuecomment-1005883910
// - https://github.com/riandyrn/otelchi/issues/6#issuecomment-1034461912
func WithRequestMethodInSpanName(isActive bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.requestMethodInSpanName = isActive
	})
}

// WithFilter adds a filter to the list of filters used by the handler.
// If any filter indicates to exclude a request then the request will not be
// traced. All filters must allow a request to be traced for a Span to be created.
// If no filters are provided then all requests are traced.
// Filters will be invoked for each processed request, it is advised to make them
// simple and fast.
func WithFilter(filter Filter) Option {
	return optionFunc(func(cfg *config) {
		cfg.filters = append(cfg.filters, filter)
	})
}

// WithTraceIDResponseHeader enables adding trace id into response header.
// It accepts a function that generates the header key name. If this parameter
// function set to `nil` the default header key which is `X-Trace-Id` will be used.
//
// Deprecated: use `WithTraceResponseHeaders` instead.
func WithTraceIDResponseHeader(headerKeyFunc func() string) Option {
	cfg := TraceHeaderConfig{
		TraceIDHeader:      "",
		TraceSampledHeader: "",
	}
	if headerKeyFunc != nil {
		cfg.TraceIDHeader = headerKeyFunc()
	}
	return WithTraceResponseHeaders(cfg)
}

// TraceHeaderConfig is configuration for trace headers in the response.
type TraceHeaderConfig struct {
	TraceIDHeader      string // if non-empty overrides the default of X-Trace-ID
	TraceSampledHeader string // if non-empty overrides the default of X-Trace-Sampled
}

// WithTraceResponseHeaders configures the response headers for trace information.
// It accepts a TraceHeaderConfig struct that contains the keys for the Trace ID
// and Trace Sampled headers. If the provided keys are empty, default values will
// be used for the respective headers.
func WithTraceResponseHeaders(cfg TraceHeaderConfig) Option {
	return optionFunc(func(c *config) {
		c.traceIDResponseHeaderKey = cfg.TraceIDHeader
		if c.traceIDResponseHeaderKey == "" {
			c.traceIDResponseHeaderKey = DefaultTraceIDResponseHeaderKey
		}

		c.traceSampledResponseHeaderKey = cfg.TraceSampledHeader
		if c.traceSampledResponseHeaderKey == "" {
			c.traceSampledResponseHeaderKey = DefaultTraceSampledResponseHeaderKey
		}
	})
}

// WithPublicEndpoint is used for marking every endpoint as public endpoint.
// This means if the incoming request has span context, it won't be used as
// parent span by the span generated by this middleware, instead the generated
// span will be the root span (new trace) and then linked to the span from the
// incoming request.
//
// Let say we have the following scenario:
//
//  1. We have 2 systems: `SysA` & `SysB`.
//  2. `SysA` has the following services: `SvcA.1` & `SvcA.2`.
//  3. `SysB` has the following services: `SvcB.1` & `SvcB.2`.
//  4. `SvcA.2` is used internally only by `SvcA.1`.
//  5. `SvcB.2` is used internally only by `SvcB.1`.
//  6. All of these services already instrumented otelchi & using the same collector (e.g Jaeger).
//  7. In `SvcA.1` we should set `WithPublicEndpoint()` since it is the entry point (a.k.a "public endpoint") for entering `SysA`.
//  8. In `SvcA.2` we should not set `WithPublicEndpoint()` since it is only used internally by `SvcA.1` inside `SysA`.
//  9. Point 7 & 8 also applies to both services in `SysB`.
//
// Now, whenever `SvcA.1` calls `SvcA.2` there will be only a single trace generated. This trace will contain 2 spans: root span from `SvcA.1` & child span from `SvcA.2`.
//
// But if let say `SvcA.2` calls `SvcB.1`, then there will be 2 traces generated: trace from `SysA` & trace from `SysB`. But in trace generated in `SysB` there will be like a marking that this trace is actually related to trace in `SysA` (a.k.a linked with the trace from `SysA`).
func WithPublicEndpoint() Option {
	return WithPublicEndpointFn(func(r *http.Request) bool { return true })
}

// WithPublicEndpointFn runs with every request, and allows conditionally
// configuring the Handler to link the generated span with an incoming span
// context.
//
// If the function return `true` the generated span will be linked with the
// incoming span context. Otherwise, the generated span will be set as the
// child span of the incoming span context.
//
// Essentially it has the same functionality as `WithPublicEndpoint` but with
// more flexibility.
func WithPublicEndpointFn(fn func(r *http.Request) bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.publicEndpointFn = fn
	})
}
