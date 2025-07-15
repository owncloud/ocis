package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/service"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	events.UploadReady{},
	events.ItemTrashed{},
	events.ItemRestored{},
	events.ItemMoved{},
	events.ContainerCreated{},
	events.FileLocked{},
	events.FileUnlocked{},
	events.FileTouched{},
	events.SpaceShared{},
	events.SpaceShareUpdated{},
	events.SpaceUnshared{},
	events.ShareCreated{},
	events.ShareRemoved{},
	events.ShareUpdated{},
	events.LinkCreated{},
	events.LinkUpdated{},
	events.LinkRemoved{},
	events.BackchannelLogout{},
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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			tracerProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			s, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gatewaySelector, err := pool.GatewaySelector(
				cfg.RevaGateway,
				pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
				pool.WithTLSMode(tm),
				pool.WithRegistry(registry.GetRegistry()),
				pool.WithTracerProvider(tracerProvider),
			)
			if err != nil {
				return fmt.Errorf("could not get reva client selector: %s", err)
			}

			gr := runner.NewGroup()
			{
				svc, err := service.NewClientlogService(
					service.Logger(logger),
					service.Config(cfg),
					service.Stream(s),
					service.GatewaySelector(gatewaySelector),
					service.RegisteredEvents(_registeredEvents),
					service.TraceProvider(tracerProvider),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.New(cfg.Service.Name+".svc", func() error {
					return svc.Run()
				}, func() {
					svc.Close()
				}))
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
