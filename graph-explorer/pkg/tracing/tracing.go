package tracing

import (
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func Configure(cfg *config.Config, logger log.Logger) error {
	if cfg.Tracing.Enabled {
		switch t := cfg.Tracing.Type; t {
		case "agent":
			exporter, err := ocagent.NewExporter(
				ocagent.WithReconnectionPeriod(5*time.Second),
				ocagent.WithAddress(cfg.Tracing.Endpoint),
				ocagent.WithServiceName(cfg.Tracing.Service),
			)
			if err != nil {
				logger.Error().
					Err(err).
					Str("endpoint", cfg.Tracing.Endpoint).
					Str("collector", cfg.Tracing.Collector).
					Msg("Failed to create agent tracing")
				return err
			}
			trace.RegisterExporter(exporter)
			view.RegisterExporter(exporter)
		case "jaeger":
			exporter, err := jaeger.NewExporter(
				jaeger.Options{
					AgentEndpoint:     cfg.Tracing.Endpoint,
					CollectorEndpoint: cfg.Tracing.Collector,
					Process: jaeger.Process{
						ServiceName: cfg.Tracing.Service,
					},
				},
			)
			if err != nil {
				logger.Error().
					Err(err).
					Str("endpoint", cfg.Tracing.Endpoint).
					Str("collector", cfg.Tracing.Collector).
					Msg("Failed to create jaeger tracing")
				return err
			}
			trace.RegisterExporter(exporter)
		case "zipkin":
			endpoint, err := openzipkin.NewEndpoint(
				cfg.Tracing.Service,
				cfg.Tracing.Endpoint,
			)
			if err != nil {
				logger.Error().
					Err(err).
					Str("endpoint", cfg.Tracing.Endpoint).
					Str("collector", cfg.Tracing.Collector).
					Msg("Failed to create zipkin tracing")
				return err
			}
			exporter := zipkin.NewExporter(
				zipkinhttp.NewReporter(
					cfg.Tracing.Collector,
				),
				endpoint,
			)
			trace.RegisterExporter(exporter)
		default:
			logger.Warn().
				Str("type", t).
				Msg("Unknown tracing backend")
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
