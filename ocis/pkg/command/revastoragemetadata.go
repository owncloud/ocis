// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis/ocis-reva/pkg/config"
	"github.com/owncloud/ocis/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// RevaStorageMetadataCommand is the entrypoint for the reva-storage-metadata command.
func RevaStorageMetadataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-metadata",
		Usage:    "Start reva storage service for metadata mount",
		Category: "Extensions",
		Flags:    flagset.StorageMetadata(cfg.Reva),
		Action: func(c *cli.Context) error {
			revaStorageMetadataCommand := command.StorageMetadata(configureRevaStorageMetadata(cfg))

			if err := revaStorageMetadataCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(revaStorageMetadataCommand.Action, c)
		},
	}
}

func configureRevaStorageMetadata(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Reva.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Reva.Tracing.Type = cfg.Tracing.Type
		cfg.Reva.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Reva.Tracing.Collector = cfg.Tracing.Collector
		cfg.Reva.Tracing.Service = cfg.Tracing.Service
	}

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageMetadataCommand)
}
