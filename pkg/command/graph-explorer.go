// +build !simple

package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-graph-explorer/pkg/command"
	svcconfig "github.com/owncloud/ocis-graph-explorer/pkg/config"
	"github.com/owncloud/ocis-graph-explorer/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// GraphExplorerCommand is the entry point for the graph command.
func GraphExplorerCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "graph-explorer",
		Usage:    "Start graph explorer",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.GraphExplorer),
		Action: func(c *cli.Context) error {
			scfg := configureGraphExplorer(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureGraphExplorer(cfg *config.Config) *svcconfig.Config {
	cfg.GraphExplorer.Log.Level = cfg.Log.Level
	cfg.GraphExplorer.Log.Pretty = cfg.Log.Pretty
	cfg.GraphExplorer.Log.Color = cfg.Log.Color
	cfg.GraphExplorer.Tracing.Enabled = false
	cfg.GraphExplorer.HTTP.Addr = "localhost:9135"
	cfg.GraphExplorer.HTTP.Root = "/"

	return cfg.GraphExplorer
}

func init() {
	register.AddCommand(GraphExplorerCommand)
}
