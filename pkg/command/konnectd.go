package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-konnectd/pkg/command"
	svcconfig "github.com/owncloud/ocis-konnectd/pkg/config"
	"github.com/owncloud/ocis-konnectd/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// KonnectdCommand is the entrypoint for the konnectd command.
func KonnectdCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "konnectd",
		Usage:    "Start konnectd server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Konnectd),
		Action: func(c *cli.Context) error {
			konnectdCommand := command.Server(configureKonnectd(cfg))

			if err := konnectdCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(konnectdCommand.Action, c)
		},
	}
}

func configureKonnectd(cfg *config.Config) *svcconfig.Config {
	cfg.Konnectd.Log.Level = cfg.Log.Level
	cfg.Konnectd.Log.Pretty = cfg.Log.Pretty
	cfg.Konnectd.Log.Color = cfg.Log.Color
	cfg.Konnectd.HTTP.TLS = false

	return cfg.Konnectd
}

func init() {
	register.AddCommand(KonnectdCommand)
}
