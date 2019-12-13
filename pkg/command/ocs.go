package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-ocs/pkg/command"
	svcconfig "github.com/owncloud/ocis-ocs/pkg/config"
	"github.com/owncloud/ocis-ocs/pkg/flagset"
	"github.com/owncloud/ocis-ocs/pkg/metrics"
	"github.com/owncloud/ocis-ocs/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "ocs",
		Usage: "Start ocs server",
		Flags: flagset.ServerWithConfig(cfg.OCS),
		Action: func(c *cli.Context) error {
			scfg := configureOCS(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// OCSHandler defines the direct server handler.
func OCSHandler(ctx context.Context, cancel context.CancelFunc, gr *run.Group, cfg *config.Config) error {
	scfg := configureOCS(cfg)
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

func configureOCS(cfg *config.Config) *svcconfig.Config {
	cfg.OCS.Log.Level = cfg.Log.Level
	cfg.OCS.Log.Pretty = cfg.Log.Pretty
	cfg.OCS.Log.Color = cfg.Log.Color
	cfg.OCS.Tracing.Enabled = false
	cfg.OCS.HTTP.Addr = "localhost:9109"
	cfg.OCS.HTTP.Root = "/"

	return cfg.OCS
}

func init() {
	register.AddCommand(OCSCommand)
	register.AddHandler(OCSHandler)
}
