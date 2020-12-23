// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/owncloud/ocis/ocs/pkg/command"
	svcconfig "github.com/owncloud/ocis/ocs/pkg/config"
	"github.com/owncloud/ocis/ocs/pkg/flagset"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocs",
		Usage:    "Start ocs server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.OCS),
		Subcommands: []*cli.Command{
			command.PrintVersion(cfg.OCS),
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureOCS(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureOCS(cfg *config.Config) *svcconfig.Config {
	cfg.OCS.Log.Level = cfg.Log.Level
	cfg.OCS.Log.Pretty = cfg.Log.Pretty
	cfg.OCS.Log.Color = cfg.Log.Color
	cfg.OCS.Service.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.OCS.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.OCS.Tracing.Type = cfg.Tracing.Type
		cfg.OCS.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.OCS.Tracing.Collector = cfg.Tracing.Collector
	}

	if cfg.TokenManager.JWTSecret != "" {
		cfg.OCS.TokenManager.JWTSecret = cfg.TokenManager.JWTSecret
	}

	return cfg.OCS
}

func init() {
	register.AddCommand(OCSCommand)
}
