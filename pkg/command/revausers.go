package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaUsersCommand is the entrypoint for the reva-users command.
func RevaUsersCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-users",
		Usage:    "Start reva users service",
		Category: "Extensions",
		Flags:    flagset.UsersWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaUsers(cfg)

			return cli.HandleAction(
				command.Users(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaUsers(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaUsersCommand)
}
