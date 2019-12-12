package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-webdav/pkg/command"
	svcconfig "github.com/owncloud/ocis-webdav/pkg/config"
	"github.com/owncloud/ocis-webdav/pkg/flagset"
	"github.com/owncloud/ocis-webdav/pkg/metrics"
	"github.com/owncloud/ocis-webdav/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "webdav",
		Usage: "Start webdav server",
		Flags: flagset.ServerWithConfig(cfg.WebDAV),
		Action: func(c *cli.Context) error {
			scfg := configureWebDAV(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// WebDAVHandler defines the direct server handler.
func WebDAVHandler(ctx context.Context, cancel context.CancelFunc, gr run.Group, cfg *config.Config) error {
	scfg := configureWebDAV(cfg)
	logger := command.NewLogger(scfg)
	m := metrics.New()

	{
		server, err := http.Server(
			http.Logger(logger),
			http.Context(ctx),
			http.Config(scfg),
			http.Metrics(m),
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

	return nil
}

func configureWebDAV(cfg *config.Config) *svcconfig.Config {
	cfg.WebDAV.Log.Level = cfg.Log.Level
	cfg.WebDAV.Log.Pretty = cfg.Log.Pretty
	cfg.WebDAV.Log.Color = cfg.Log.Color
	cfg.WebDAV.Tracing.Enabled = false
	cfg.WebDAV.HTTP.Root = "/"

	return cfg.WebDAV
}

// func init() {
// 	register.AddCommand(WebDAVCommand)
// 	register.AddHandler(WebDAVHandler)
// }
