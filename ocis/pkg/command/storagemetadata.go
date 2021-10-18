package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// StorageMetadataCommand is the entrypoint for the storage-metadata command.
func StorageMetadataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-metadata",
		Usage:    "Start storage and data service for metadata",
		Category: "Extensions",
		Flags:    flagset.StorageMetadata(cfg.Storage),
		Action: func(c *cli.Context) error {
			origCmd := command.StorageMetadata(configureStorageMetadata(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureStorageMetadata(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(StorageMetadataCommand)
}
