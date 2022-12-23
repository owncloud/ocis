package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/hub/pkg/command"
	"github.com/urfave/cli/v2"
)

// HubCommand is the entrypoint for the web command.
func HubCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Hub.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Hub.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.WebDAV.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Hub),
	}
}

func init() {
	register.AddCommand(HubCommand)
}
