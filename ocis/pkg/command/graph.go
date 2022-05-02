package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/graph/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GraphCommand is the entrypoint for the graph command.
func GraphCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Graph.Service.Name,
		Usage:    subcommandDescription(cfg.Graph.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Graph.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Graph),
	}
}

func init() {
	register.AddCommand(GraphCommand)
}
