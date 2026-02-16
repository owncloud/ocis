package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/sse/pkg/config"
	"github.com/owncloud/ocis/v2/services/sse/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/sse/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/sse/pkg/server/http"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	events.SendSSE{},
}

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			logger := log.NewLogger(
				log.Name(cfg.Service.Name),
				log.Level(cfg.Log.Level),
				log.Pretty(cfg.Log.Pretty),
				log.Color(cfg.Log.Color),
				log.File(cfg.Log.File),
			)

			tracerProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			gr := runner.NewGroup()
			{
				connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
				natsStream, err := stream.NatsFromConfig(connName, true, stream.NatsConfig(cfg.Events))
				if err != nil {
					return err
				}

				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Consumer(natsStream),
					http.RegisteredEvents(_registeredEvents),
					http.TracerProvider(tracerProvider),
				)
				if err != nil {
					return err
				}

				gr.Add(runner.NewGoMicroHttpServerRunner(cfg.Service.Name+".http", server))
			}

			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))
			}

			logger.Warn().Msgf("starting service %s", cfg.Service.Name)
			grResults := gr.Run(ctx)

			if err := runner.ProcessResults(grResults); err != nil {
				logger.Error().Err(err).Msgf("service %s stopped with error", cfg.Service.Name)
				return err
			}
			logger.Warn().Msgf("service %s stopped without error", cfg.Service.Name)
			return nil
		},
	}
}
