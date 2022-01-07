package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/settings/pkg/command"
	"github.com/urfave/cli/v2"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Settings.Service.Name,
		Usage:    subcommandDescription(cfg.Settings.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Settings),
	}
}

func init() {
	register.AddCommand(SettingsCommand)
}
