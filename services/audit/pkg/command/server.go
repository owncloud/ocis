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
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
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
			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			gr := runner.NewGroup()

			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			client, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, "audit", types.RegisteredEvents()...)
			if err != nil {
				return err
			}

			// we need an additional context for the audit server in order to
			// cancel it anytime
			svcCtx, svcCancel := context.WithCancel(ctx)
			defer svcCancel()

			gr.Add(runner.New(cfg.Service.Name+".svc", func() error {
				svc.AuditLoggerFromConfig(svcCtx, cfg.Auditlog, evts, logger)
				return nil
			}, func() {
				svcCancel()
			}))

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
