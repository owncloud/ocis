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

// StorageUserProviderCommand is the entrypoint for the storage-userprovider command.
func StorageUserProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-userprovider",
		Usage:    "Start storage userprovider service",
		Category: "Extensions",
		Flags:    flagset.UsersWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Users(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageUserProviderCommand)
}
