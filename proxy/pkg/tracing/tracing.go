package tracing

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Configure(cfg *config.Config, logger log.Logger) error {
	if cfg.Tracing.Enabled {
		switch t := cfg.Tracing.Type; t {
		case "jaeger":
			exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Tracing.Collector)))
			if err != nil {
				return err
			}

			tp := sdktrace.NewTracerProvider(
				// Always be sure to batch in production.
				sdktrace.WithBatcher(exp),
				sdktrace.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("proxy"),
					// attribute.String("environment", "development"), TODO(refs) flip this bit
				)),
			)

			otel.SetTracerProvider(tp)
			propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
			otel.SetTextMapPropagator(propagator)
		case "zipkin":
			fallthrough
		case "agent":
			fallthrough
		default:
			logger.Warn().
				Str("type", t).
				Msg("unsupported tracing backend")
		}
		trace.ApplyConfig(
			trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			},
		)
	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
	return nil
}
