//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageSharesCommand is the entrypoint for the storage-shares command.
func StorageSharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-shares",
		Usage:    "Start storage and data provider for /home/Shares mount",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.StorageShares(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageSharesCommand)
}
