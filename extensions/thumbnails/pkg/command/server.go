package command

import (
	"context"
	"fmt"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/logging"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/metrics"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/server/grpc"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/server/http"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/tracing"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Value:       cfg.ConfigFile,
				Usage:       "config file to be loaded by the extension",
				Destination: &cfg.ConfigFile,
			},
		},
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			err := tracing.Configure(cfg)
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

			metrics.BuildInfo.WithLabelValues(version.String).Set(1)

			service := grpc.NewService(
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Service.Name),
				grpc.Namespace(cfg.GRPC.Namespace),
				grpc.Address(cfg.GRPC.Addr),
				grpc.Metrics(metrics),
			)

			gr.Add(service.Run, func(_ error) {
				fmt.Println("shutting down grpc server")
				cancel()
			})

			server, err := debug.Server(
				debug.Logger(logger),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(server.ListenAndServe, func(_ error) {
				_ = server.Shutdown(ctx)
				cancel()
			})

			httpServer, err := http.Server(
				http.Logger(logger),
				http.Context(ctx),
				http.Config(cfg),
				http.Metrics(metrics),
				http.Namespace(cfg.HTTP.Namespace),
			)

			if err != nil {
				logger.Info().
					Err(err).
					Str("transport", "http").
					Msg("Failed to initialize server")

				return err
			}

			gr.Add(httpServer.Run, func(_ error) {
				logger.Info().Str("server", "http").Msg("shutting down server")
				cancel()
			})

			return gr.Run()
		},
	}
}
