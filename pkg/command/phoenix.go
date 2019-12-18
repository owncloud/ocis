package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-phoenix/pkg/command"
	svcconfig "github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis-phoenix/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// PhoenixCommand is the entrypoint for the phoenix command.
func PhoenixCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "phoenix",
		Usage:    "Start phoenix server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Phoenix),
		Action: func(c *cli.Context) error {
			scfg := configurePhoenix(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configurePhoenix(cfg *config.Config) *svcconfig.Config {
	cfg.Phoenix.Log.Level = cfg.Log.Level
	cfg.Phoenix.Log.Pretty = cfg.Log.Pretty
	cfg.Phoenix.Log.Color = cfg.Log.Color
	cfg.Phoenix.Tracing.Enabled = false
	cfg.Phoenix.HTTP.Addr = "localhost:9100"
	cfg.Phoenix.HTTP.Root = "/"

	return cfg.Phoenix
}

func init() {
	register.AddCommand(PhoenixCommand)
}
