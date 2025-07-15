package tracing

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	rtrace "github.com/owncloud/reva/v2/pkg/trace"
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
