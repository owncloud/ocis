//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/settings/pkg/command"
	"github.com/urfave/cli/v2"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "settings",
		Usage:    "Start settings server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Settings),
		},
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			globalLog = cfg.Log

			return nil
		},
		Action: func(c *cli.Context) error {
			// if accounts logging is empty in ocis.yaml
			if (cfg.Settings.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.Settings.Log = globalLog
			}
			origCmd := command.Server(cfg.Settings)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(SettingsCommand)
}
