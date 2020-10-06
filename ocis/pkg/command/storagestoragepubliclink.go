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

// StorageStoragePublicLinkCommand is the entrypoint for the reva-storage-oc command.
func StorageStoragePublicLinkCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-storage-public-link",
		Usage:    "Start storage public link storage",
		Category: "Extensions",
		Flags:    flagset.StoragePublicLink(cfg.Storage),
		Action: func(c *cli.Context) error {
			scfg := configureStorageStoragePublicLink(cfg)

			return cli.HandleAction(
				command.StoragePublicLink(scfg).Action,
				c,
			)
		},
	}
}

func configureStorageStoragePublicLink(cfg *config.Config) *svcconfig.Config {
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
	register.AddCommand(StorageStoragePublicLinkCommand)
}
