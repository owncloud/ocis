// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-graph-explorer/pkg/command"
	svcconfig "github.com/owncloud/ocis-graph-explorer/pkg/config"
	"github.com/owncloud/ocis-graph-explorer/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// GraphExplorerCommand is the entry point for the graph command.
func GraphExplorerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "graph-explorer",
		Usage:    "Start graph explorer",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.GraphExplorer),
		Action: func(ctx *cli.Context) error {
			graphExplorerCommand := command.Server(configureGraphExplorer(cfg))

			if err := graphExplorerCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(graphExplorerCommand.Action, ctx)
		},
	}
}

func configureGraphExplorer(cfg *config.Config) *svcconfig.Config {
	cfg.GraphExplorer.Log.Level = cfg.Log.Level
	cfg.GraphExplorer.Log.Pretty = cfg.Log.Pretty
	cfg.GraphExplorer.Log.Color = cfg.Log.Color

	return cfg.GraphExplorer
}

func init() {
	register.AddCommand(GraphExplorerCommand)
}
