package command

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/oklog/run"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/audit/pkg/config"
	"github.com/owncloud/ocis/v2/services/audit/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/audit/pkg/logging"
	"github.com/owncloud/ocis/v2/services/audit/pkg/server/debug"
	svc "github.com/owncloud/ocis/v2/services/audit/pkg/service"
	"github.com/owncloud/ocis/v2/services/audit/pkg/types"
)

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
			var (
				gr     = run.Group{}
				logger = logging.Configure(cfg.Service.Name, cfg.Log)

				ctx, cancel = context.WithCancel(c.Context)
			)
			defer cancel()

			client, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, "audit", types.RegisteredEvents()...)
			if err != nil {
				return err
			}

			gr.Add(func() error {
				svc.AuditLoggerFromConfig(ctx, cfg.Auditlog, evts, logger)
				return nil
			}, func(err error) {
				if err == nil {
					logger.Info().
						Str("transport", "stream").
						Str("server", cfg.Service.Name).
						Msg("Shutting down server")
				} else {
					logger.Error().Err(err).
						Str("transport", "stream").
						Str("server", cfg.Service.Name).
						Msg("Shutting down server")
				}

				cancel()
			})

			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(debugServer.ListenAndServe, func(_ error) {
					_ = debugServer.Shutdown(ctx)
					cancel()
				})
			}
			return gr.Run()
		},
	}
}
