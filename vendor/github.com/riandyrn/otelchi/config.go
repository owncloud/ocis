package otelchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// config is used to configure the mux middleware.
type config struct {
	TracerProvider          oteltrace.TracerProvider
	Propagators             propagation.TextMapPropagator
	ChiRoutes               chi.Routes
	RequestMethodInSpanName bool
	Filter                  func(r *http.Request) bool
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		cfg.Propagators = propagators
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.TracerProvider = provider
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
		cfg.ChiRoutes = routes
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
		cfg.RequestMethodInSpanName = isActive
	})
}

// WithFilter is used for filtering request that should not be traced.
// This is useful for filtering health check request, etc.
// A Filter must return true if the request should be traced.
func WithFilter(filter func(r *http.Request) bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.Filter = filter
	})
}
