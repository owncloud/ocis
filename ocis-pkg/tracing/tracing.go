package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	rtrace "github.com/owncloud/reva/v2/pkg/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ERR_INVALID_AGENT_ENDPOINT = "invalid agent endpoint %s. expected format: hostname:port"

// Propagator ensures the importer module uses the same trace propagation strategy.
var Propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// GetServiceTraceProvider returns a configured open-telemetry trace provider.
func GetServiceTraceProvider(c ConfigConverter, serviceName string) (trace.TracerProvider, error) {
	var cfg Config
	if c == nil || reflect.ValueOf(c).IsNil() {
		cfg = Config{Enabled: false}
	} else {
		cfg = c.Convert()
	}

	if cfg.Enabled {
		return GetTraceProvider(cfg.Endpoint, cfg.Collector, serviceName, cfg.Type)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.NeverSample()),
	)
	rtrace.SetDefaultTracerProvider(tp)

	return tp, nil
}

// GetPropagator gets a configured propagator.
func GetPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
}

// propagatorOnceFn is needed to ensure we set the global propagator only once
var propagatorOnceFn sync.Once

// setGlobalPropagatorOnce set the global propagator only once. This is needed
// because go-micro uses the global propagator to extract and inject data and
// we want to use the same.
// Note: in case of services running in different hosts, this needs to be run
// in all of them.
func setGlobalPropagatorOnce() {
	propagatorOnceFn.Do(func() {
		otel.SetTextMapPropagator(GetPropagator())
	})
}

// GetTraceProvider returns a configured open-telemetry trace provider.
func GetTraceProvider(endpoint, collector, serviceName, traceType string) (*sdktrace.TracerProvider, error) {
	switch t := traceType; t {
	case "", "jaeger":
		var (
			exp *jaeger.Exporter
			err error
		)

		if endpoint != "" {
			var agentHost string
			var agentPort string

			agentHost, agentPort, err = parseAgentConfig(endpoint)
			if err != nil {
				return nil, err
			}

			exp, err = jaeger.New(
				jaeger.WithAgentEndpoint(
					jaeger.WithAgentHost(agentHost),
					jaeger.WithAgentPort(agentPort),
				),
			)
		} else if collector != "" {
			exp, err = jaeger.New(
				jaeger.WithCollectorEndpoint(
					jaeger.WithEndpoint(collector),
				),
			)
		} else {
			return sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.NeverSample())), nil
		}
		if err != nil {
			return nil, err
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName)),
			),
		)
		rtrace.SetDefaultTracerProvider(tp)
		return tp, nil
	case "otlp":
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, endpoint,
			// Note the use of insecure transport here. TLS is recommended in production.
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}
		exporter, err := otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithGRPCConn(conn),
		)
		if err != nil {
			return nil, err
		}
		resources, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				attribute.String("service.name", serviceName),
				attribute.String("library.language", "go"),
			),
		)
		if err != nil {
			return nil, err
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		)
		rtrace.SetDefaultTracerProvider(tp)
		setGlobalPropagatorOnce()
		return tp, nil
	case "agent":
		fallthrough
	case "zipkin":
		fallthrough
	default:
		return nil, fmt.Errorf("unknown trace type %s", traceType)
	}
}

func parseAgentConfig(ae string) (string, string, error) {
	u, err := url.Parse(ae)
	// as per url.go:
	// [...] Trying to parse a hostname and path
	// without a scheme is invalid but may not necessarily return an
	// error, due to parsing ambiguities.
	if err == nil && u.Hostname() != "" && u.Port() != "" {
		return u.Hostname(), u.Port(), nil
	}

	p := strings.Split(ae, ":")
	if len(p) != 2 {
		return "", "", fmt.Errorf(ERR_INVALID_AGENT_ENDPOINT, ae)
	}

	switch {
	case p[0] == "" && p[1] == "": // case ae = ":"
		return "", "", fmt.Errorf(ERR_INVALID_AGENT_ENDPOINT, ae)
	case p[0] == "":
		return "", "", fmt.Errorf(ERR_INVALID_AGENT_ENDPOINT, ae)
	}
	return p[0], p[1], nil
}

// GetNewRequest gets a new HTTP request with tracing data coming from the
// provided context. Note that the provided context will NOT be used for the
// request, just to get the data.
// The request will have a new "context.Background()" context associated. This
// means that cancelling the provided context will NOT stop the request.
func GetNewRequest(injectCtx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return GetNewRequestWithDifferentContext(context.Background(), injectCtx, method, url, body)
}

// GetNewRequestWithContext gets a new HTTP request with tracing data coming
// from the provided context. The request will also have the same provided
// context associated (in case the context is cancelled)
func GetNewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return GetNewRequestWithDifferentContext(ctx, ctx, method, url, body)
}

// GetNewRequestWithDifferentContext gets a new HTTP request with tracing
// data coming from the "injectCtx" context. The "reqCtx" context will be
// associated with the request.
//
// This method is intended to be used if you want to associate a context
// with a request, and at the same time use a different context to get the
// tracing info.
func GetNewRequestWithDifferentContext(reqCtx, injectCtx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(reqCtx, method, url, body)
	if err != nil {
		return req, err
	}

	InjectTracingHeaders(injectCtx, req)
	return req, nil
}

// InjectTracingHeaders sets the tracing info from the context as HTTP headers
// in the provided request.
func InjectTracingHeaders(ctx context.Context, req *http.Request) {
	propagator := GetPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
}
