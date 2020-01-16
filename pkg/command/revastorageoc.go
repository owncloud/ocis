package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaStorageOCCommand is the entrypoint for the reva-storage-oc command.
func RevaStorageOCCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-storage-oc",
		Usage:    "Start reva oc storage",
		Category: "Extensions",
		Flags:    flagset.StorageOCWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageOC(cfg)

			return cli.HandleAction(
				command.StorageOC(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageOC(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	cfg.Reva.Reva.StorageOC.ExposeDataServer = true

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageOCCommand)
}
