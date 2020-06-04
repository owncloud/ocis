// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaSharingCommand is the entrypoint for the reva-sharing command.
func RevaSharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-sharing",
		Usage:    "Start reva sharing service",
		Category: "Extensions",
		Flags:    flagset.SharingWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaSharing(cfg)

			return cli.HandleAction(
				command.Sharing(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaSharing(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaSharingCommand)
}
