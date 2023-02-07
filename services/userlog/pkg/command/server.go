package command

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/server/http"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/store"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	events.UploadReady{},
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

			defer cancel()

			consumer, err := stream.NatsFromConfig(stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			var st store.Store
			switch cfg.Store.Type {
			case "inmemory":
				st = store.NewMemoryStore()
			default:
				return fmt.Errorf("unknown store '%s' configured", cfg.Store.Type)
			}

			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Store(st),
					http.Consumer(consumer),
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
