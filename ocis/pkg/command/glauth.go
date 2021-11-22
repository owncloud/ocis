package command

import (
	"github.com/owncloud/ocis/glauth/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "glauth",
		Usage:    "Start glauth server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if cfg.Commons != nil {
				cfg.GLAuth.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.GLAuth)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(GLAuthCommand)
}
