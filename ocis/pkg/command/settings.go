package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/settings/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Settings.Service.Name,
		Usage:    subcommandDescription(cfg.Settings.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
			}
			return err
		},
		Subcommands: command.GetCommands(cfg.Settings),
	}
}

func init() {
	register.AddCommand(SettingsCommand)
}
