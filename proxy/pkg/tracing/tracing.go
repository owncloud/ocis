package tracing

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	// Propagator ensures the entire module uses the same trace propagation strategy.
	Propagator propagation.TextMapPropagator

	// TraceProvider is the global trace provider for the proxy service.
	TraceProvider = sdktrace.NewTracerProvider()
)

func Configure(cfg *config.Config, logger log.Logger) error {
	if cfg.Tracing.Enabled {
		switch t := cfg.Tracing.Type; t {
		case "jaeger":
			{
				exp, err := jaeger.New(
					jaeger.WithCollectorEndpoint(
						jaeger.WithEndpoint(cfg.Tracing.Collector),
					),
				)
				if err != nil {
					panic(err)
				}

				// set package level trace provider and propagator.
				TraceProvider = sdktrace.NewTracerProvider(
					sdktrace.WithBatcher(exp),
					sdktrace.WithResource(resource.NewWithAttributes(
						semconv.SchemaURL,
						semconv.ServiceNameKey.String("proxy")),
					),
				)

				Propagator = propagation.NewCompositeTextMapPropagator(
					propagation.Baggage{},
					propagation.TraceContext{},
				)
			}
		case "agent":
			fallthrough
		case "zipkin":
			fallthrough
		default:
			logger.Warn().Str("type", t).Msg("Unknown tracing backend")
		}
	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
	return nil
}
