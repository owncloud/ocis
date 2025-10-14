package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/reva/v2/cmd/revad/runtime"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/logging"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/revaconfig"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/server/debug"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			gr := runner.NewGroup()

			{
				rCfg, err := revaconfig.FrontendConfigFromStruct(cfg, logger)
				if err != nil {
					return err
				}

				// run the appropriate reva servers based on the config
				if rServer := runtime.NewDrivenHTTPServerWithOptions(rCfg,
					runtime.WithLogger(&logger.Logger),
					runtime.WithRegistry(registry.GetRegistry()),
					runtime.WithTraceProvider(traceProvider),
				); rServer != nil {
					gr.Add(runner.NewRevaServiceRunner(cfg.Service.Name+".rhttp", rServer))
				}
				if rServer := runtime.NewDrivenGRPCServerWithOptions(rCfg,
					runtime.WithLogger(&logger.Logger),
					runtime.WithRegistry(registry.GetRegistry()),
					runtime.WithTraceProvider(traceProvider),
				); rServer != nil {
					gr.Add(runner.NewRevaServiceRunner(cfg.Service.Name+".rgrpc", rServer))
				}
			}

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

			httpSvc := registry.BuildHTTPService(cfg.HTTP.Namespace+"."+cfg.Service.Name, cfg.HTTP.Addr, version.GetString())
			if err := registry.RegisterService(ctx, logger, httpSvc, cfg.Debug.Addr); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the http service")
			}

			// add event handler
			gr.Add(runner.New(cfg.Service.Name+".event",
				func() error {
					return ListenForEvents(ctx, cfg, logger)
				}, func() {
					logger.Info().Msg("stopping event handler")
				},
			))

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
