package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-phoenix/pkg/command"
	"github.com/owncloud/ocis-phoenix/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// PhoenixCommand is the entrypoint for the phoenix command.
func PhoenixCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "phoenix",
		Usage:    "Start phoenix server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Phoenix),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			cfg.Phoenix.Phoenix.Config.Apps = c.StringSlice("web-config-app")

			return nil
		},
		Action: func(c *cli.Context) error {

			scfg := configurePhoenix(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func init() {
	register.AddCommand(PhoenixCommand)
}
