package command

import (
	"context"
	"strings"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/idp/pkg/metrics"
	"github.com/owncloud/ocis/idp/pkg/server/debug"
	"github.com/owncloud/ocis/idp/pkg/server/http"
	"github.com/owncloud/ocis/idp/pkg/tracing"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Before: func(ctx *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			//if len(ctx.StringSlice("trusted-proxy")) > 0 {
			//	cfg.IDP.TrustedProxy = ctx.StringSlice("trusted-proxy")
			//}
			//
			//if len(ctx.StringSlice("allow-scope")) > 0 {
			//	cfg.IDP.AllowScope = ctx.StringSlice("allow-scope")
			//}
			//
			//if len(ctx.StringSlice("signing-private-key")) > 0 {
			//	cfg.IDP.SigningPrivateKeyFiles = ctx.StringSlice("signing-private-key")
			//}

			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			tracing.Configure(cfg)

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

			metrics.BuildInfo.WithLabelValues(cfg.Service.Version).Set(1)

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(metrics),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
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

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}
