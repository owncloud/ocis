package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/web/pkg/command"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	var globalLog shared.Log

	return &cli.Command{
		Name:     "web",
		Usage:    "Start web server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}

			globalLog = cfg.Log

			return nil
		},
		Action: func(c *cli.Context) error {
			// if accounts logging is empty in ocis.yaml
			if (cfg.Web.Log == shared.Log{}) && (globalLog != shared.Log{}) {
				// we can safely inherit the global logging values.
				cfg.Web.Log = globalLog
			}
			origCmd := command.Server(cfg.Web)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(WebCommand)
}
