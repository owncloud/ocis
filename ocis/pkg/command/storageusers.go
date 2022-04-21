package command

import (
	"github.com/owncloud/ocis/extensions/storage-users/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageUsersCommand is the entrypoint for the storage-users command.
func StorageUsersCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-users",
		Usage:    "start storage and data provider for /users mount",
		Category: "extensions",
		//Flags:    flagset.StorageUsersWithConfig(cfg.Storage),
		// Before: func(ctx *cli.Context) error {
		// 	return ParseStorageCommon(ctx, cfg)
		// },
		Action: func(c *cli.Context) error {
			origCmd := command.StorageUsers(cfg.StorageUsers)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageUsersCommand)
}
