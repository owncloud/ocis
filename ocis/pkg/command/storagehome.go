//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageHomeCommand is the entrypoint for the storage-home command.
func StorageHomeCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-home",
		Usage:    "Start storage and data provider for /home mount",
		Category: "Extensions",
		//Flags:    flagset.StorageHomeWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.StorageHome(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageHomeCommand)
}
