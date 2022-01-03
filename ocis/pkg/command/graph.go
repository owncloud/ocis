package command

import (
	"github.com/owncloud/ocis/graph/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GraphCommand is the entrypoint for the graph command.
func GraphCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "graph",
		Usage:    "Start graph server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.Accounts.Commons = cfg.Commons
			}

			return nil
		},
		Subcommands: command.GetCommands(cfg.Graph),
	}
}

func init() {
	register.AddCommand(GraphCommand)
}
