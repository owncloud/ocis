package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebfingerCommand is the entrypoint for the webfinger command.
func WebfingerCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:     cfg.Webfinger.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Webfinger.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.Webfinger.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Webfinger),
	}
}

func init() {
	register.AddCommand(WebfingerCommand)
}
