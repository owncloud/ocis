package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Propagator ensures the importer module uses the same trace propagation strategy.
var Propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// GetTraceProvider returns a configured open-telemetry trace provider.
func GetTraceProvider(collectorEndpoint, traceType, serviceName string) (*sdktrace.TracerProvider, error) {
	switch t := traceType; t {
	case "jaeger":
		if collectorEndpoint == "" {
			return sdktrace.NewTracerProvider(), nil
		}

		exp, err := jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(collectorEndpoint),
			),
		)
		if err != nil {
			return nil, err
		}

		return sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName)),
			),
		), nil

	case "agent":
		fallthrough
	case "zipkin":
		fallthrough
	default:
		return nil, fmt.Errorf("invalid trace configuration")
	}
}
