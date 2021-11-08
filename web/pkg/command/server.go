package command

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"

	gofig "github.com/gookit/config/v2"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/web/pkg/config"
	"github.com/owncloud/ocis/web/pkg/metrics"
	"github.com/owncloud/ocis/web/pkg/server/debug"
	"github.com/owncloud/ocis/web/pkg/server/http"
	"github.com/owncloud/ocis/web/pkg/tracing"
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
				cfg.HTTP.Root = strings.TrimRight(cfg.HTTP.Root, "/")
			}

			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			// build well known openid-configuration endpoint if it is not set
			if cfg.Web.Config.OpenIDConnect.MetadataURL == "" {
				cfg.Web.Config.OpenIDConnect.MetadataURL = strings.TrimRight(cfg.Web.Config.OpenIDConnect.Authority, "/") + "/.well-known/openid-configuration"
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

			// actually read the contents of the config file and override defaults
			if cfg.File != "" {
				contents, err := ioutil.ReadFile(cfg.File)
				if err != nil {
					logger.Err(err).Msg("error opening config file")
					return err
				}
				if err := json.Unmarshal(contents, &cfg.Web.Config); err != nil {
					logger.Fatal().Err(err).Msg("error unmarshalling config file")
					return err
				}
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

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Namespace(cfg.HTTP.Namespace),
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

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}
