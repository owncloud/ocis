package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-phoenix/pkg/command"
	svcconfig "github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis-phoenix/pkg/flagset"
	"github.com/owncloud/ocis-phoenix/pkg/metrics"
	"github.com/owncloud/ocis-phoenix/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
)

// PhoenixCommand is the entrypoint for the phoenix command.
func PhoenixCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "phoenix",
		Usage: "Start phoenix server",
		Flags: flagset.ServerWithConfig(cfg.Phoenix),
		Action: func(c *cli.Context) error {
			scfg := configurePhoenix(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// PhoenixHandler defines the direct server handler.
func PhoenixHandler(ctx context.Context, cancel context.CancelFunc, gr run.Group, cfg *config.Config) error {
	scfg := configurePhoenix(cfg)
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

func configurePhoenix(cfg *config.Config) *svcconfig.Config {
	cfg.Phoenix.Log.Level = cfg.Log.Level
	cfg.Phoenix.Log.Pretty = cfg.Log.Pretty
	cfg.Phoenix.Log.Color = cfg.Log.Color
	cfg.Phoenix.Tracing.Enabled = false
	cfg.Phoenix.HTTP.Root = "/"

	return cfg.Phoenix
}

// func init() {
// 	register.AddCommand(PhoenixCommand)
// 	register.AddHandler(PhoenixHandler)
// }
