package command

import (
	"context"
	"fmt"
	"os"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/logging"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/server/http"
	"github.com/owncloud/ocis/v2/services/graph-explorer/pkg/tracing"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
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
				mtrcs = metrics.New()
			)

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Namespace(cfg.HTTP.Namespace),
					http.Config(cfg),
					http.Metrics(mtrcs),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(func() error {
					err := server.Run()
					if err != nil {
						logger.Error().
							Err(err).
							Str("transport", "http").
							Msg("Failed to start server")
					}
					return err
				}, func(_ error) {
					logger.Info().
						Str("transport", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
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
			}

			return gr.Run()
		},
	}
}
