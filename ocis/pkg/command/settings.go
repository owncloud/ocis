//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/settings/pkg/command"
	"github.com/urfave/cli/v2"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "settings",
		Usage:    "Start settings server",
		Category: "Extensions",
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Settings),
		},
		Before: func(ctx *cli.Context) error {
			if cfg.Commons != nil {
				cfg.Settings.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.Settings)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(SettingsCommand)
}
