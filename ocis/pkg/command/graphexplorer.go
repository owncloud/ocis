package command

import (
	"github.com/owncloud/ocis/graph-explorer/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GraphExplorerCommand is the entrypoint for the graph-explorer command.
func GraphExplorerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "graph-explorer",
		Usage:    "Start graph-explorer server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.Graph.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.GraphExplorer)
			return handleOriginalAction(c, origCmd)
		},
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.GraphExplorer),
		},
	}
}

func init() {
	register.AddCommand(GraphExplorerCommand)
}
