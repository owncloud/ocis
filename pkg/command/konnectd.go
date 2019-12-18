package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-konnectd/pkg/command"
	svcconfig "github.com/owncloud/ocis-konnectd/pkg/config"
	"github.com/owncloud/ocis-konnectd/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// KonnectdCommand is the entrypoint for the konnectd command.
func KonnectdCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "konnectd",
		Usage:    "Start konnectd server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Konnectd),
		Action: func(c *cli.Context) error {
			scfg := configureKonnectd(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureKonnectd(cfg *config.Config) *svcconfig.Config {
	cfg.Konnectd.Log.Level = cfg.Log.Level
	cfg.Konnectd.Log.Pretty = cfg.Log.Pretty
	cfg.Konnectd.Log.Color = cfg.Log.Color
	cfg.Konnectd.Tracing.Enabled = false
	cfg.Konnectd.HTTP.Addr = "localhost:9130"
	cfg.Konnectd.HTTP.Root = "/"

	return cfg.Konnectd
}

func init() {
	register.AddCommand(KonnectdCommand)
}
