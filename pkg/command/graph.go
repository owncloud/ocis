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
func GraphCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "graph",
		Usage:    "Start graph server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Graph),
		Action: func(c *cli.Context) error {
			scfg := configureGraph(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureGraph(cfg *config.Config) *svcconfig.Config {
	cfg.Graph.Log.Level = cfg.Log.Level
	cfg.Graph.Log.Pretty = cfg.Log.Pretty
	cfg.Graph.Log.Color = cfg.Log.Color
	cfg.Graph.Tracing.Enabled = false
	cfg.Graph.HTTP.Addr = "localhost:9120"
	cfg.Graph.HTTP.Root = "/"

	return cfg.Graph
}

func init() {
	register.AddCommand(GraphCommand)
}
