package command

import (
	"context"
	"fmt"
	"os/signal"
	"reflect"

	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/reva/v2/pkg/store"
	microstore "go-micro.dev/v4/store"

	"github.com/urfave/cli/v2"

	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/logging"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/service"
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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			grpcClient, err := grpc.NewClient(
				append(
					grpc.GetClientOptions(&cfg.GRPCClientTLS),
					grpc.WithTraceProvider(traceProvider),
				)...,
			)
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

			// evs defines a list of events to subscribe to
			evs := []events.Unmarshaller{
				events.ShareCreated{},
				events.ShareExpired{},
				events.ShareRemoved{},
				events.SpaceShared{},
				events.SpaceUnshared{},
				events.SpaceMembershipExpired{},
				events.ScienceMeshInviteTokenGenerated{},
				events.SendEmailsEvent{},
			}
			registeredEvents := make(map[string]events.Unmarshaller)
			for _, e := range evs {
				typ := reflect.TypeOf(e)
				registeredEvents[typ.String()] = e
			}

			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			client, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Notifications.Events))
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, "notifications", evs...)
			if err != nil {
				return err
			}
			channel, err := channels.NewMailChannel(*cfg, logger)
			if err != nil {
				return err
			}
			tm, err := pool.StringToTLSMode(cfg.Notifications.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gatewaySelector, err := pool.GatewaySelector(
				cfg.Notifications.RevaGateway,
				pool.WithTLSCACert(cfg.Notifications.GRPCClientTLS.CACert),
				pool.WithTLSMode(tm),
				pool.WithRegistry(registry.GetRegistry()),
				pool.WithTracerProvider(traceProvider),
			)
			if err != nil {
				logger.Fatal().Err(err).Str("addr", cfg.Notifications.RevaGateway).Msg("could not get reva gateway selector")
			}
			valueService := settingssvc.NewValueService("com.owncloud.api.settings", grpcClient)
			historyClient := ehsvc.NewEventHistoryService("com.owncloud.api.eventhistory", grpcClient)

			notificationStore := store.Create(
				store.Store(cfg.Store.Store),
				store.TTL(cfg.Store.TTL),
				microstore.Nodes(cfg.Store.Nodes...),
				microstore.Database(cfg.Store.Database),
				microstore.Table(cfg.Store.Table),
				store.Authentication(cfg.Store.AuthUsername, cfg.Store.AuthPassword),
			)

			svc := service.NewEventsNotifier(evts, channel, logger, gatewaySelector, valueService,
				cfg.ServiceAccount.ServiceAccountID, cfg.ServiceAccount.ServiceAccountSecret,
				cfg.Notifications.EmailTemplatePath, cfg.Notifications.DefaultLanguage, cfg.WebUIURL,
				cfg.Notifications.TranslationPath, cfg.Notifications.SMTP.Sender, notificationStore, historyClient, registeredEvents)

			gr.Add(runner.New(cfg.Service.Name+".svc", func() error {
				return svc.Run()
			}, func() {
				svc.Close()
			}))

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
