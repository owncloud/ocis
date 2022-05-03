package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/group/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GroupCommand is the entrypoint for the Group command.
func GroupCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Group.Service.Name,
		Usage:    subcommandDescription(cfg.Group.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Group.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Group),
	}
}

func init() {
	register.AddCommand(GroupCommand)
}
