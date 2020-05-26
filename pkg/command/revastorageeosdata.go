package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaStorageEOSDataCommand is the entrypoint for the reva-storage-eos-data command.
func RevaStorageEOSDataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-eos-data",
		Usage:    "Start reva eos storage dataprovider",
		Category: "Extensions",
		Flags:    flagset.StorageEOSDataWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageEOSData(cfg)

			return cli.HandleAction(
				command.StorageEOSData(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageEOSData(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageEOSDataCommand)
}
