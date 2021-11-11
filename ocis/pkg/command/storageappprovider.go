//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// StorageAppProviderCommand is the entrypoint for the reva-app-provider command.
func StorageAppProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-app-provider",
		Usage:    "Start storage app-provider service",
		Category: "Extensions",
		Flags:    flagset.AppProviderWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AppProvider(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAppProviderCommand)
}
