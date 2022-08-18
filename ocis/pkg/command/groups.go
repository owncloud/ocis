package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/groups/pkg/command"
	"github.com/urfave/cli/v2"
)

// GroupsCommand is the entrypoint for the groups command.
func GroupsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Groups.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Groups.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.Groups.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Groups),
	}
}

func init() {
	register.AddCommand(GroupsCommand)
}
