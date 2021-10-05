package command

import (
	"github.com/owncloud/ocis/glauth/pkg/command"
	svcconfig "github.com/owncloud/ocis/glauth/pkg/config"
	"github.com/owncloud/ocis/glauth/pkg/flagset"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "glauth",
		Usage:    "Start glauth server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.GLAuth),
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureGLAuth(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureGLAuth(cfg *config.Config) *svcconfig.Config {
	cfg.GLAuth.Log.Level = cfg.Log.Level
	cfg.GLAuth.Log.Pretty = cfg.Log.Pretty
	cfg.GLAuth.Log.Color = cfg.Log.Color
	cfg.GLAuth.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.GLAuth.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.GLAuth.Tracing.Type = cfg.Tracing.Type
		cfg.GLAuth.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.GLAuth.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.GLAuth
}

func init() {
	register.AddCommand(GLAuthCommand)
}
