package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/settings/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Settings.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Settings.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Settings.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Settings),
	}
}

func init() {
	register.AddCommand(SettingsCommand)
}
