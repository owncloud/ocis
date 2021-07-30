package tracing

import (
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocs/pkg/config"
	"go.opentelemetry.io/otel/exporters/jaeger"
)

var TP *sdktrace.TracerProvider

func Configure(cfg *config.Config, logger log.Logger) error {
	if cfg.Tracing.Enabled {
		switch t := cfg.Tracing.Type; t {
		case "jaeger":
			exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Tracing.Collector)))
			if err != nil {
				return err
			}

			TP = sdktrace.NewTracerProvider(
				// Always be sure to batch in production.
				sdktrace.WithBatcher(exp),
				sdktrace.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("ocs"),
					// attribute.String("environment", "development"), TODO(refs) flip this bit
				)),
			)

			//propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
			//otel.SetTextMapPropagator(propagator)
		case "zipkin":
			fallthrough
		case "agent":
			fallthrough
		default:
			logger.Warn().
				Str("type", t).
				Msg("unsupported tracing backend")
		}
	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
	return nil
}
