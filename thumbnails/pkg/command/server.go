package command

import (
	"context"
	"fmt"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/owncloud/ocis/thumbnails/pkg/flagset"
	"github.com/owncloud/ocis/thumbnails/pkg/metrics"
	"github.com/owncloud/ocis/thumbnails/pkg/server/debug"
	"github.com/owncloud/ocis/thumbnails/pkg/server/grpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Thumbnail.Resolutions = c.StringSlice("thumbnail-resolution")

			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
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

			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				metrics = metrics.New()
			)

			defer cancel()

			metrics.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			service := grpc.NewService(
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Server.Name),
				grpc.Namespace(cfg.Server.Namespace),
				grpc.Address(cfg.Server.Address),
				grpc.Metrics(metrics),
			)

			gr.Add(func() error {
				return service.Run()
			}, func(_ error) {
				fmt.Println("shutting down grpc server")
				cancel()
			})

			server, err := debug.Server(
				debug.Logger(logger),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().
					Err(err).
					Str("transport", "debug").
					Msg("Failed to initialize server")

				return err
			}

			gr.Add(func() error {
				return server.ListenAndServe()
			}, func(_ error) {
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}
