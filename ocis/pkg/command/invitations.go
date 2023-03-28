package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/command"
	"github.com/urfave/cli/v2"
)

// InvitationsCommand is the entrypoint for the invitations command.
func InvitationsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Invitations.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Invitations.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Invitations.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Invitations),
	}
}

func init() {
	register.AddCommand(InvitationsCommand)
}
