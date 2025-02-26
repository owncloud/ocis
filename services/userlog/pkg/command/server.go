package command

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/store"
	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/server/http"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	// file related
	events.PostprocessingStepFinished{},

	// space related
	events.SpaceDisabled{},
	events.SpaceDeleted{},
	events.SpaceShared{},
	events.SpaceUnshared{},
	events.SpaceMembershipExpired{},

	// share related
	events.ShareCreated{},
	events.ShareRemoved{},
	events.ShareExpired{},
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

			grpcClient, err := ogrpc.NewClient(
				append(ogrpc.GetClientOptions(cfg.GRPCClientTLS), ogrpc.WithTraceProvider(tracerProvider))...,
			)
			if err != nil {
				return err
			}

			gr := run.Group{}
			ctx, cancel := context.WithCancel(c.Context)

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			defer cancel()

			stream, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			st := store.Create(
				store.Store(cfg.Persistence.Store),
				store.TTL(cfg.Persistence.TTL),
				microstore.Nodes(cfg.Persistence.Nodes...),
				microstore.Database(cfg.Persistence.Database),
				microstore.Table(cfg.Persistence.Table),
				store.Authentication(cfg.Persistence.AuthUsername, cfg.Persistence.AuthPassword),
			)

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

			hClient := ehsvc.NewEventHistoryService("com.owncloud.api.eventhistory", grpcClient)
			vClient := settingssvc.NewValueService("com.owncloud.api.settings", grpcClient)
			rClient := settingssvc.NewRoleService("com.owncloud.api.settings", grpcClient)

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Store(st),
					http.Stream(stream),
					http.GatewaySelector(gatewaySelector),
					http.History(hClient),
					http.Value(vClient),
					http.Role(rClient),
					http.RegisteredEvents(_registeredEvents),
					http.TracerProvider(tracerProvider),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(err error) {
					if err == nil {
						logger.Info().
							Str("transport", "http").
							Str("server", cfg.Service.Name).
							Msg("Shutting down server")
					} else {
						logger.Error().Err(err).
							Str("transport", "http").
							Str("server", cfg.Service.Name).
							Msg("Shutting down server")
					}

					cancel()
				})
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

				gr.Add(debugServer.ListenAndServe, func(_ error) {
					_ = debugServer.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
