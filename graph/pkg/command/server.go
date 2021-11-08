package command

import (
	"context"
	"strings"

	gofig "github.com/gookit/config/v2"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/metrics"
	"github.com/owncloud/ocis/graph/pkg/server/debug"
	"github.com/owncloud/ocis/graph/pkg/server/http"
	"github.com/owncloud/ocis/graph/pkg/tracing"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Before: func(ctx *cli.Context) error {
			// remember shared logging info to prevent empty overwrites
			inLog := cfg.Log

			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			if (cfg.Log == shared.Log{}) && (inLog != shared.Log{}) {
				// set the default to the parent config
				cfg.Log = inLog

				// and parse the environment
				conf := &gofig.Config{}
				conf.LoadOSEnv(config.GetEnv(), false)
				bindings := config.StructMappings(cfg)
				if err := ociscfg.BindEnv(conf, bindings); err != nil {
					return err
				}
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if err := tracing.Configure(cfg); err != nil {
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

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
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
