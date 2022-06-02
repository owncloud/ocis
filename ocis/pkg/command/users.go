package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/users/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// UsersCommand is the entrypoint for the users command.
func UsersCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Users.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Users.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Users.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Users),
	}
}

func init() {
	register.AddCommand(UsersCommand)
}
