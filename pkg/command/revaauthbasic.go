package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaAuthBasicCommand is the entrypoint for the reva-auth-basic command.
func RevaAuthBasicCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-auth-basic",
		Usage:    "Start reva auth-basic service",
		Category: "Extensions",
		Flags:    flagset.AuthBasicWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaAuthBasic(cfg)

			return cli.HandleAction(
				command.AuthBasic(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaAuthBasic(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaAuthBasicCommand)
}
