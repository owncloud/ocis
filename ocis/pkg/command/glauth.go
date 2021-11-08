package command

import (
	"github.com/owncloud/ocis/glauth/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "glauth",
		Usage:    "Start glauth server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}
			globalLog = cfg.Log
			return nil
		},
		Action: func(c *cli.Context) error {
			// if Glauth logging is empty in ocis.yaml
			if (cfg.GLAuth.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.GLAuth.Log = globalLog
			}
			origCmd := command.Server(cfg.GLAuth)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(GLAuthCommand)
}
