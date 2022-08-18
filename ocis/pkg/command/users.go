package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/users/pkg/command"
	"github.com/urfave/cli/v2"
)

// UsersCommand is the entrypoint for the users command.
func UsersCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Users.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Users.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.Users.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Users),
	}
}

func init() {
	register.AddCommand(UsersCommand)
}
