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

// RevaStoragePublicLinkDataCommand is the entrypoint for the reva-storage-public-link-data command.
func RevaStoragePublicLinkDataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-public-link-data",
		Usage:    "Start reva public link storage dataprovider",
		Category: "Extensions",
		Flags:    flagset.StoragePublicLinkDataWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStoragePublicLinkData(cfg)

			return cli.HandleAction(
				command.StoragePublicLinkData(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStoragePublicLinkData(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStoragePublicLinkDataCommand)
}
