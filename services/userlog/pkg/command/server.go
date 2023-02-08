package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/store"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/server/http"
	"github.com/urfave/cli/v2"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	// file related
	events.UploadReady{},
	events.ContainerCreated{},
	events.FileTouched{},
	events.FileDownloaded{},
	events.FileVersionRestored{},
	events.ItemMoved{},
	events.ItemTrashed{},
	events.ItemPurged{},
	events.ItemRestored{},

	// space related
	events.SpaceCreated{},
	events.SpaceRenamed{},
	events.SpaceEnabled{},
	events.SpaceDisabled{},
	events.SpaceDeleted{},
	events.SpaceShared{},
	events.SpaceUnshared{},
	events.SpaceUpdated{},
	events.SpaceMembershipExpired{},

	// share related
	events.ShareCreated{},
	// events.ShareRemoved{}, // TODO: ShareRemoved doesn't hold sharee information
	events.ShareUpdated{},
	events.ShareExpired{},
	events.LinkCreated{},
	// events.LinkRemoved{}, // TODO: LinkRemoved doesn't hold sharee information
	events.LinkUpdated{},
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

			err := ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			defer cancel()

			consumer, err := stream.NatsFromConfig(stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			st := store.Create(
				store.Type(cfg.Store.Type),
				store.Addresses(strings.Split(cfg.Store.Addresses, ",")...),
				store.Database(cfg.Store.Database),
				store.Table(cfg.Store.Table),
			)

			tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gwclient, err := pool.GetGatewayServiceClient(
				cfg.RevaGateway,
				pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
				pool.WithTLSMode(tm),
			)
			if err != nil {
				return fmt.Errorf("could not get reva client: %s", err)
			}

			hClient := ehsvc.NewEventHistoryService("com.owncloud.api.eventhistory", grpc.DefaultClient())

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Store(st),
					http.Consumer(consumer),
					http.Gateway(gwclient),
					http.History(hClient),
					http.RegisteredEvents(_registeredEvents),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(err error) {
					logger.Error().
						Str("transport", "http").
						Err(err).
						Msg("Shutting down server")

					cancel()
				})
			}

			return gr.Run()
		},
	}
}
