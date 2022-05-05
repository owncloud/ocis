package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/groups/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GroupsCommand is the entrypoint for the groups command.
func GroupsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Groups.Service.Name,
		Usage:    subcommandDescription(cfg.Groups.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Groups.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Groups),
	}
}

func init() {
	register.AddCommand(GroupsCommand)
}
