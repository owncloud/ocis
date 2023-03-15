package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/command"
	"github.com/urfave/cli/v2"
)

// AntivirusCommand is the entrypoint for the antivirus command.
func AntivirusCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Antivirus.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Antivirus.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			//cfg.Antivirus.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Antivirus),
	}
}

func init() {
	register.AddCommand(AntivirusCommand)
}
