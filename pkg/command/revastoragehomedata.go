package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaStorageHomeDataCommand is the entrypoint for the reva-storage-home-data command.
func RevaStorageHomeDataCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-storage-home-data",
		Usage:    "Start reva home storage dataprovider",
		Category: "Extensions",
		Flags:    flagset.StorageHomeDataWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageHomeData(cfg)

			return cli.HandleAction(
				command.StorageHomeData(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageHomeData(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageHomeDataCommand)
}
