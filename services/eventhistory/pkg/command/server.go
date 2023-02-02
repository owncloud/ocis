package command

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/config"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/logging"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/server/grpc"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/store"
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
			err := ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				metrics = metrics.New()
			)

			defer cancel()

			metrics.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			consumer, err := stream.NatsFromConfig(stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			// TODO: configure store
			st := store.DefaultStore

			service := grpc.NewService(
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Service.Name),
				grpc.Namespace(cfg.GRPC.Namespace),
				grpc.Address(cfg.GRPC.Addr),
				grpc.Metrics(metrics),
				grpc.Consumer(consumer),
				grpc.Store(st),
			)

			gr.Add(service.Run, func(_ error) {
				logger.Error().
					Err(err).
					Str("server", "grpc").
					Msg("Shutting down server")

				cancel()
			})

			return gr.Run()

		},
	}
}
