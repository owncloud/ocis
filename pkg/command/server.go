package command

import (
	"context"
	"os"
	"os/signal"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli"
	"github.com/micro/go-micro/util/log"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis-graph/pkg/server/debug"
	"github.com/owncloud/ocis-graph/pkg/server/http"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Action: func(c *cli.Context) error {
			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					exporter, err := ocagent.NewExporter(
						ocagent.WithReconnectionPeriod(5*time.Second),
						ocagent.WithAddress(cfg.Tracing.Endpoint),
						ocagent.WithServiceName(cfg.Tracing.Service),
					)

					if err != nil {
						log.Error(
							"Failed to create agent tracing on [%s] endpoint and [%s] collector: %w",
							cfg.Tracing.Endpoint,
							cfg.Tracing.Collector,
							err,
						)

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
						log.Error(
							"Failed to create jaeger tracing on [%s] endpoint and [%s] collector: %w",
							cfg.Tracing.Endpoint,
							cfg.Tracing.Collector,
							err,
						)

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						log.Error(
							"Failed to create zipkin tracing on [%s] endpoint and [%s] collector: %w",
							cfg.Tracing.Endpoint,
							cfg.Tracing.Collector,
							err,
						)

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
					log.Warnf("Unknown tracing backend [%s]", t)
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				log.Debug("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
			)

			{
				server, err := http.Server(
					http.Context(ctx),
					http.Config(cfg),
				)

				if err != nil {
					log.Errorf("Server [http] failed to initialize: %w", err)
					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(_ error) {
					log.Infof("Server [http] shutting down")
					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					log.Errorf("Server [debug] failed to initialize: %w", err)
					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Errorf("Server [debug] shutdown failed: %w", err)
					} else {
						log.Infof("Server [debug] shutting down")
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
