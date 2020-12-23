// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
)

// StorageSharingCommand is the entrypoint for the reva-sharing command.
func StorageSharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-sharing",
		Usage:    "Start storage sharing service",
		Category: "Extensions",
		Flags:    flagset.SharingWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			origCmd := command.Sharing(configureStorageSharing(cfg))

			if err := origCmd.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(origCmd.Action, c)
		},
	}
}

func configureStorageSharing(cfg *config.Config) *svcconfig.Config {
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
	register.AddCommand(StorageSharingCommand)
}
