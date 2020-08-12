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

// RevaStorageHomeCommand is the entrypoint for the reva-storage-home command.
func RevaStorageHomeCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-home",
		Usage:    "Start reva home storage",
		Category: "Extensions",
		Flags:    flagset.StorageHomeWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageHome(cfg)

			return cli.HandleAction(
				command.StorageHome(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageHome(cfg *config.Config) *svcconfig.Config {
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
	register.AddCommand(RevaStorageHomeCommand)
}
