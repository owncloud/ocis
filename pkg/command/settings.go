// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-settings/pkg/command"
	svcconfig "github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "settings",
		Usage:    "Start settings server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Settings),
		Action: func(ctx *cli.Context) error {
			settingsCommand := command.Server(configureSettings(cfg))

			if err := settingsCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(settingsCommand.Action, ctx)
		},
	}
}

func configureSettings(cfg *config.Config) *svcconfig.Config {
	cfg.Settings.Log.Level = cfg.Log.Level
	cfg.Settings.Log.Pretty = cfg.Log.Pretty
	cfg.Settings.Log.Color = cfg.Log.Color

	return cfg.Settings
}

func init() {
	register.AddCommand(SettingsCommand)
}
