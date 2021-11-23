//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageGroupProviderCommand is the entrypoint for the storage-groupprovider command.
func StorageGroupProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-groupprovider",
		Usage:    "Start storage groupprovider service",
		Category: "Extensions",
		//Flags:    flagset.GroupsWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Groups(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageGroupProviderCommand)
}
