package command

import (
	"context"
	"fmt"

	"github.com/oklog/run"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/nats/pkg/logging"
	"github.com/owncloud/ocis/v2/services/nats/pkg/server/nats"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.LogReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			defer cancel()

			natsServer, err := nats.NewNATSServer(
				ctx,
				logging.NewLogWrapper(logger),
				nats.Host(cfg.Nats.Host),
				nats.Port(cfg.Nats.Port),
				nats.ClusterID(cfg.Nats.ClusterID),
				nats.StoreDir(cfg.Nats.StoreDir),
			)
			if err != nil {
				return err
			}

			gr.Add(func() error {
				err := make(chan error)
				select {
				case <-ctx.Done():
					return nil
				case err <- natsServer.ListenAndServe():
					return <-err
				}

			}, func(err error) {
				logger.Error().
					Err(err).
					Msg("Shutting down server")

				natsServer.Shutdown()
				cancel()
			})

			return gr.Run()
		},
	}
}
