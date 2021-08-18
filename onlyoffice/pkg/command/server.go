package command

import (
	"context"
	"strings"

	"github.com/owncloud/ocis/onlyoffice/pkg/tracing"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/onlyoffice/pkg/config"
	"github.com/owncloud/ocis/onlyoffice/pkg/flagset"
	"github.com/owncloud/ocis/onlyoffice/pkg/metrics"
	"github.com/owncloud/ocis/onlyoffice/pkg/server/debug"
	"github.com/owncloud/ocis/onlyoffice/pkg/server/http"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			logger := NewLogger(cfg)
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			if !cfg.Supervised {
				return ParseConfig(ctx, cfg)
			}
			logger.Debug().Str("service", "onlyoffice").Msg("ignoring config file parsing when running supervised")
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if err := tracing.Configure(cfg); err != nil {
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

			if ctx == nil {
				ctx, cancel = context.WithCancel(context.Background())
			} else {
				ctx, cancel = context.WithCancel(cfg.Context)
			}

			defer cancel()

			{
				server, err := http.Server(
					http.Name("onlyoffice"),
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(metrics),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(_ error) {
					logger.Info().
						Str("server", "http").
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
					logger.Error().Err(err).Str("server", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}
