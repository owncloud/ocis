package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-glauth/pkg/command"
	svcconfig "github.com/owncloud/ocis-glauth/pkg/config"
	"github.com/owncloud/ocis-glauth/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "glauth",
		Usage:    "Start glauth server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.GLAuth),
		Action: func(c *cli.Context) error {
			scfg := configureGLAuth(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureGLAuth(cfg *config.Config) *svcconfig.Config {
	cfg.GLAuth.Log.Level = cfg.Log.Level
	cfg.GLAuth.Log.Pretty = cfg.Log.Pretty
	cfg.GLAuth.Log.Color = cfg.Log.Color
	return cfg.GLAuth
}

func init() {
	register.AddCommand(GLAuthCommand)
}
