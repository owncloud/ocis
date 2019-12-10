package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-graph/pkg/command"
	svcconfig "github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis-graph/pkg/metrics"
	"github.com/owncloud/ocis-graph/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// GraphCommand is the entrypoint for the graph command.
func GraphCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "graph",
		Usage: "Start graph server",
		Flags: flagset.ServerWithConfig(cfg.Graph),
		Action: func(c *cli.Context) error {
			scfg := configureGraph(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// GraphHandler defines the direct server handler.
func GraphHandler(ctx context.Context, cancel context.CancelFunc, gr run.Group, cfg *config.Config) error {
	scfg := configureGraph(cfg)
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

func configureGraph(cfg *config.Config) *svcconfig.Config {
	cfg.Graph.Log.Level = cfg.Log.Level
	cfg.Graph.Log.Pretty = cfg.Log.Pretty
	cfg.Graph.Log.Color = cfg.Log.Color
	cfg.Graph.Tracing.Enabled = false
	cfg.Graph.HTTP.Root = "/"

	return cfg.Graph
}

func init() {
	register.AddCommand(GraphCommand)
	register.AddHandler(GraphHandler)
}
