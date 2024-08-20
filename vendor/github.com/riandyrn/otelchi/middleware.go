package otelchi

import (
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/riandyrn/otelchi"
)

// Middleware sets up a handler to start tracing the incoming
// requests. The serverName parameter should describe the name of the
// (virtual) server handling the request.
func Middleware(serverName string, opts ...Option) func(next http.Handler) http.Handler {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(Version()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}

	return func(handler http.Handler) http.Handler {
		return traceware{
			serverName:             serverName,
			tracer:                 tracer,
			propagators:            cfg.Propagators,
			handler:                handler,
			chiRoutes:              cfg.ChiRoutes,
			reqMethodInSpanName:    cfg.RequestMethodInSpanName,
			filters:                cfg.Filters,
			traceResponseHeaderKey: cfg.TraceResponseHeaderKey,
			publicEndpointFn:       cfg.PublicEndpointFn,
		}
	}
}

type traceware struct {
	serverName             string
	tracer                 oteltrace.Tracer
	propagators            propagation.TextMapPropagator
	handler                http.Handler
	chiRoutes              chi.Routes
	reqMethodInSpanName    bool
	filters                []Filter
	traceResponseHeaderKey string
	publicEndpointFn       func(r *http.Request) bool
}

type recordingResponseWriter struct {
	writer  http.ResponseWriter
	written bool
	status  int
}

var rrwPool = &sync.Pool{
	New: func() interface{} {
		return &recordingResponseWriter{}
	},
}

func getRRW(writer http.ResponseWriter) *recordingResponseWriter {
	rrw := rrwPool.Get().(*recordingResponseWriter)
	rrw.written = false
	rrw.status = http.StatusOK
	rrw.writer = httpsnoop.Wrap(writer, httpsnoop.Hooks{
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				if !rrw.written {
					rrw.written = true
				}
				return next(b)
			}
		},
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if !rrw.written {
					rrw.written = true
					rrw.status = statusCode
				}
				next(statusCode)
			}
		},
	})
	return rrw
}

func putRRW(rrw *recordingResponseWriter) {
	rrw.writer = nil
	rrwPool.Put(rrw)
}

// ServeHTTP implements the http.Handler interface. It does the actual
// tracing of the request.
func (tw traceware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// go through all filters if any
	for _, filter := range tw.filters {
		// if there is a filter that returns false, we skip tracing
		// and execute next handler
		if !filter(r) {
			tw.handler.ServeHTTP(w, r)
			return
		}
	}

	// extract tracing header using propagator
	ctx := tw.propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	// create span, based on specification, we need to set already known attributes
	// when creating the span, the only thing missing here is HTTP route pattern since
	// in go-chi/chi route pattern could only be extracted once the request is executed
	// check here for details:
	//
	// https://github.com/go-chi/chi/issues/150#issuecomment-278850733
	//
	// if we have access to chi routes, we could extract the route pattern beforehand.
	spanName := ""
	routePattern := ""
	spanAttributes := httpconv.ServerRequest(tw.serverName, r)

	if tw.chiRoutes != nil {
		rctx := chi.NewRouteContext()
		if tw.chiRoutes.Match(rctx, r.Method, r.URL.Path) {
			routePattern = rctx.RoutePattern()
			spanName = addPrefixToSpanName(tw.reqMethodInSpanName, r.Method, routePattern)
			spanAttributes = append(spanAttributes, semconv.HTTPRoute(routePattern))
		}
	}

	// define span start options
	spanOpts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(spanAttributes...),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	}

	if tw.publicEndpointFn != nil && tw.publicEndpointFn(r) {
		// mark span as the root span
		spanOpts = append(spanOpts, oteltrace.WithNewRoot())

		// linking incoming span context to the root span, we need to
		// ensure if the incoming span context is valid (because it is
		// possible for us to receive invalid span context due to various
		// reason such as bug or context propagation error) and it is
		// coming from another service (remote) before linking it to the
		// root span
		spanCtx := oteltrace.SpanContextFromContext(ctx)
		if spanCtx.IsValid() && spanCtx.IsRemote() {
			spanOpts = append(
				spanOpts,
				oteltrace.WithLinks(oteltrace.Link{
					SpanContext: spanCtx,
				}),
			)
		}
	}

	// start span
	ctx, span := tw.tracer.Start(ctx, spanName, spanOpts...)
	defer span.End()

	// put trace_id to response header only when WithTraceResponseHeaderKey is used
	if len(tw.traceResponseHeaderKey) > 0 && span.SpanContext().HasTraceID() {
		w.Header().Add(tw.traceResponseHeaderKey, span.SpanContext().TraceID().String())
	}

	// get recording response writer
	rrw := getRRW(w)
	defer putRRW(rrw)

	// execute next http handler
	r = r.WithContext(ctx)
	tw.handler.ServeHTTP(rrw.writer, r)

	// set span name & http route attribute if route pattern cannot be determined
	// during span creation
	if len(routePattern) == 0 {
		routePattern = chi.RouteContext(r.Context()).RoutePattern()
		span.SetAttributes(semconv.HTTPRoute(routePattern))

		spanName = addPrefixToSpanName(tw.reqMethodInSpanName, r.Method, routePattern)
		span.SetName(spanName)
	}

	// set status code attribute
	span.SetAttributes(semconv.HTTPStatusCode(rrw.status))

	// set span status
	span.SetStatus(httpconv.ServerStatus(rrw.status))
}

func addPrefixToSpanName(shouldAdd bool, prefix, spanName string) string {
	// in chi v5.0.8, the root route will be returned has an empty string
	// (see https://github.com/go-chi/chi/blob/v5.0.8/context.go#L126)
	if spanName == "" {
		spanName = "/"
	}

	if shouldAdd && len(spanName) > 0 {
		spanName = prefix + " " + spanName
	}
	return spanName
}
