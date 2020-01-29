package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaStorageOCDataCommand is the entrypoint for the reva-storage-oc-data command.
func RevaStorageOCDataCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-storage-oc-data",
		Usage:    "Start reva oc storage dataprovider",
		Category: "Extensions",
		Flags:    flagset.StorageOCDataWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageOCData(cfg)

			return cli.HandleAction(
				command.StorageOCData(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageOCData(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageOCDataCommand)
}
