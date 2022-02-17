package command

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/run"

	"github.com/owncloud/ocis/nats/pkg/config"
	"github.com/owncloud/ocis/nats/pkg/config/parser"
	"github.com/owncloud/ocis/nats/pkg/logging"
	"github.com/owncloud/ocis/nats/pkg/server/nats"
	"github.com/urfave/cli/v2"

	// TODO: .Logger Option on events/server would make this import redundant
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
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

			var natsServer *stanServer.StanServer

			gr.Add(func() error {
				var err error

				natsServer, err = nats.RunNatsServer(
					nats.Host(cfg.Nats.Host),
					nats.Port(cfg.Nats.Port),
					nats.StanOpts(
						func(o *stanServer.Options) {
							o.CustomLogger = logging.NewLogWrapper(logger)
						},
					),
				)

				if err != nil {
					return err
				}

				errChan := make(chan error)

				go func() {
					for {
						// check if NATs server has an encountered an error
						if err := natsServer.LastError(); err != nil {
							errChan <- err
							return
						}
						if ctx.Err() != nil {
							return // context closed
						}
						time.Sleep(1 * time.Second)
					}
				}()

				select {
				case <-ctx.Done():
					return nil
				case err = <-errChan:
					return err
				}

			}, func(_ error) {
				logger.Info().
					Msg("Shutting down server")

				natsServer.Shutdown()
				cancel()
			})

			return gr.Run()
		},
	}
}
