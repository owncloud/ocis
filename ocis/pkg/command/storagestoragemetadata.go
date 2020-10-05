package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// StorageStorageMetadataCommand is the entrypoint for the reva-storage-metadata command.
func StorageStorageMetadataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-storage-metadata",
		Usage:    "Start storage storage service for metadata mount",
		Category: "Extensions",
		Flags:    flagset.StorageMetadata(cfg.Storage),
		Action: func(c *cli.Context) error {
			revaStorageMetadataCommand := command.StorageMetadata(configureStorageStorageMetadata(cfg))

			if err := revaStorageMetadataCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(revaStorageMetadataCommand.Action, c)
		},
	}
}

func configureStorageStorageMetadata(cfg *config.Config) *svcconfig.Config {
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
	register.AddCommand(StorageStorageMetadataCommand)
}
