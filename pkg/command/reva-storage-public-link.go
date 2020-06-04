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

// RevaPublicLinkStorage is the entry point for the proxy command.
func RevaPublicLinkStorage(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-public-link",
		Usage:    "Start public link storage driver",
		Category: "Extensions",
		Flags:    flagset.StoragePublicLink(cfg.Reva),
		Action: func(ctx *cli.Context) error {
			publicStorageCmd := command.StoragePublicLink(configurePublicStorage(cfg))
			return cli.HandleAction(publicStorageCmd.Action, ctx)
		},
	}
}

func configurePublicStorage(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	return cfg.Reva
}

func init() {
	register.AddCommand(RevaPublicLinkStorage)
}
