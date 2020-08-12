// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-graph/pkg/command"
	svcconfig "github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// GraphCommand is the entrypoint for the graph command.
func GraphCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "graph",
		Usage:    "Start graph server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Graph),
		Action: func(ctx *cli.Context) error {
			graphCommand := command.Server(configureGraph(cfg))

			if err := graphCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(graphCommand.Action, ctx)
		},
	}
}

func configureGraph(cfg *config.Config) *svcconfig.Config {
	cfg.Graph.Log.Level = cfg.Log.Level
	cfg.Graph.Log.Pretty = cfg.Log.Pretty
	cfg.Graph.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Graph.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Graph.Tracing.Type = cfg.Tracing.Type
		cfg.Graph.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Graph.Tracing.Collector = cfg.Tracing.Collector
		cfg.Graph.Tracing.Service = cfg.Tracing.Service
	}

	return cfg.Graph
}

func init() {
	register.AddCommand(GraphCommand)
}
