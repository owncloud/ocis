package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaStorageRootCommand is the entrypoint for the reva-storage-root command.
func RevaStorageRootCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva-storage-root",
		Usage:    "Start reva root storage",
		Category: "Extensions",
		Flags:    flagset.StorageRootWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureRevaStorageRoot(cfg)

			return cli.HandleAction(
				command.StorageRoot(scfg).Action,
				c,
			)
		},
	}
}

func configureRevaStorageRoot(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaStorageRootCommand)
}
