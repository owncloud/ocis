// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// StorageStorageHomeDataCommand is the entrypoint for the reva-storage-home-data command.
func StorageStorageHomeDataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-storage-home-data",
		Usage:    "Start storage storage data provider for home mount",
		Category: "Extensions",
		Flags:    flagset.StorageHomeDataWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			scfg := configureStorageStorageHomeData(cfg)

			return cli.HandleAction(
				command.StorageHomeData(scfg).Action,
				c,
			)
		},
	}
}

func configureStorageStorageHomeData(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
		cfg.Storage.Tracing.Service = cfg.Tracing.Service
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(StorageStorageHomeDataCommand)
}
