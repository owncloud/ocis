package command

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis/onlyoffice/pkg/config"
	"github.com/owncloud/ocis/onlyoffice/pkg/flagset"
	"github.com/owncloud/ocis/onlyoffice/pkg/metrics"
	"github.com/owncloud/ocis/onlyoffice/pkg/server/debug"
	"github.com/owncloud/ocis/onlyoffice/pkg/server/http"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			return ParseConfig(ctx, cfg)
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
							ServiceName:       cfg.Tracing.Service,
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
				ctx, cancel = context.WithCancel(context.Background())
				metrics     = metrics.New()
			)

			defer cancel()

			{
				server, err := http.Server(
					http.Name("onlyoffice"),
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(metrics),
					http.Flags(flagset.RootWithConfig(config.New())),
					http.Flags(flagset.ServerWithConfig(config.New())),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(_ error) {
					logger.Info().
						Str("server", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Error().
							Err(err).
							Str("server", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", "debug").
							Msg("Shutting down server")
					}
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
