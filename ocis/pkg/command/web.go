package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/web/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "web",
		Usage:    "Start web server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			if cfg.Commons != nil {
				cfg.Web.Commons = cfg.Commons
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(cfg.Web)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(WebCommand)
}
