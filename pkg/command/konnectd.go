package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-konnectd/pkg/command"
	svcconfig "github.com/owncloud/ocis-konnectd/pkg/config"
	"github.com/owncloud/ocis-konnectd/pkg/flagset"
	"github.com/owncloud/ocis-konnectd/pkg/metrics"
	"github.com/owncloud/ocis-konnectd/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// KonnectdCommand is the entrypoint for the konnectd command.
func KonnectdCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "konnectd",
		Usage: "Start konnectd server",
		Flags: flagset.ServerWithConfig(cfg.Konnectd),
		Action: func(c *cli.Context) error {
			scfg := configureKonnectd(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// KonnectdHandler defines the direct server handler.
func KonnectdHandler(ctx context.Context, cancel context.CancelFunc, gr *run.Group, cfg *config.Config) error {
	scfg := configureKonnectd(cfg)
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

func configureKonnectd(cfg *config.Config) *svcconfig.Config {
	cfg.Konnectd.Log.Level = cfg.Log.Level
	cfg.Konnectd.Log.Pretty = cfg.Log.Pretty
	cfg.Konnectd.Log.Color = cfg.Log.Color
	cfg.Konnectd.Tracing.Enabled = false
	cfg.Konnectd.HTTP.Addr = "localhost:9011"
	cfg.Konnectd.HTTP.Root = "/"

	return cfg.Konnectd
}

func init() {
	register.AddCommand(KonnectdCommand)
	register.AddHandler(KonnectdHandler)
}
