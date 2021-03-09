// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/settings/pkg/command"
	svcconfig "github.com/owncloud/ocis/settings/pkg/config"
	"github.com/owncloud/ocis/settings/pkg/flagset"
)

// SettingsCommand is the entry point for the settings command.
func SettingsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "settings",
		Usage:    "Start settings server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Settings),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.Settings),
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureSettings(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureSettings(cfg *config.Config) *svcconfig.Config {
	cfg.Settings.Log.Level = cfg.Log.Level
	cfg.Settings.Log.Pretty = cfg.Log.Pretty
	cfg.Settings.Log.Color = cfg.Log.Color
	cfg.Settings.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.Settings.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Settings.Tracing.Type = cfg.Tracing.Type
		cfg.Settings.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Settings.Tracing.Collector = cfg.Tracing.Collector
	}

	if cfg.TokenManager.JWTSecret != "" {
		cfg.Settings.TokenManager.JWTSecret = cfg.TokenManager.JWTSecret
	}

	return cfg.Settings
}

func init() {
	register.AddCommand(SettingsCommand)
}
