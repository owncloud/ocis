package command

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	gorun "github.com/micro/go-micro/runtime"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) cli.Command {
	app := cli.Command{
		Name:  "server",
		Usage: "Start fullstack server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			return nil
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

			// by the time we reach this point, all micro runtime services subcommands
			// should be available for being initialized
			muRuntime := cmd.DefaultCmd.Options().Runtime
			env := os.Environ()

			services := []string{
				"network",  // :8085
				"runtime",  // :8088
				"registry", // :8000
				"broker",   // :8001
				"store",    // :8002
				"tunnel",   // :8083
				"router",   // :8084
				"monitor",  // :????
				"debug",    // :????
				"proxy",    // :8081
				"api",      // :8080
				"web",      // :8082
				"bot",      // :????

				// ocis extensions
				"hello",
			}

			for _, service := range services {
				args := []gorun.CreateOption{
					gorun.WithCommand(os.Args[0], service),
					gorun.WithEnv(env),
					gorun.WithOutput(os.Stdout),
				}

				muService := &gorun.Service{Name: service}
				if err := (*muRuntime).Create(muService, args...); err != nil {
					logger.Error().Msgf("Failed to create runtime enviroment: %v", err)
				}
			}

			shutdown := make(chan os.Signal, 1)
			signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

			logger.Info().Msg("Starting service runtime")

			// start the runtime
			if err := (*muRuntime).Start(); err != nil {
				os.Exit(1)
			}

			logger.Info().Msgf("Service runtime started")

			select {
			case <-shutdown:
				logger.Info().Msg("shutdown signal received")
				logger.Info().Msg("stopping service runtime")
			}

			// stop all the things
			if err := (*muRuntime).Stop(); err != nil {
				logger.Err(err)
			}

			logger.Info().Msgf("Service runtime shutdown")

			// exit success
			os.Exit(0)

			return nil
		},
	}
	return app
}
