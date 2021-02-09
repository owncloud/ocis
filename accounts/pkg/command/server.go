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
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	"github.com/owncloud/ocis/accounts/pkg/metrics"
	"github.com/owncloud/ocis/accounts/pkg/server/grpc"
	"github.com/owncloud/ocis/accounts/pkg/server/http"
	svc "github.com/owncloud/ocis/accounts/pkg/service/v0"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "uses an LDAP server as the storage backend",
		Flags:       append(flagset.ServerWithConfig(cfg), flagset.RootWithConfig(cfg)...),
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
				ctx, cancel = context.WithCancel(context.Background())
				mtrcs       = metrics.New()
			)

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			handler, err := svc.New(svc.Logger(logger), svc.Config(cfg))
			if err != nil {
				logger.Fatal().Err(err).Msg("could not initialize service handler")
			}

			{
				server := http.Server(
					http.Logger(logger),
					http.Name(cfg.Server.Name),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					//http.Flags(flagset.RootWithConfig(config.New())),
					//http.Flags(flagset.ServerWithConfig(config.New())),
					http.Handler(handler),
				)

				gr.Add(server.Run, func(_ error) {
					logger.Info().
						Str("server", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server := grpc.Server(
					grpc.Logger(logger),
					grpc.Name(cfg.Server.Name),
					grpc.Context(ctx),
					grpc.Config(cfg),
					grpc.Metrics(mtrcs),
					grpc.Handler(handler),
				)

				gr.Add(func() error {
					logger.Info().Str("service", server.Name()).Msg("Reporting settings bundles to settings service")
					svc.RegisterSettingsBundles(&logger)
					svc.RegisterPermissions(&logger)
					return server.Run()
				}, func(_ error) {
					logger.Info().
						Str("server", "grpc").
						Msg("Shutting down server")

					cancel()
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
