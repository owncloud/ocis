package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/authz/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuthzCommand is the entrypoint for the web command.
func AuthzCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Authz.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Authz.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.WebDAV.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Authz),
	}
}

func init() {
	register.AddCommand(AuthzCommand)
}
