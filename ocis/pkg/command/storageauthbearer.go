//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageAuthBearerCommand is the entrypoint for the reva-auth-bearer command.
func StorageAuthBearerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-bearer",
		Usage:    "Start storage auth-bearer service",
		Category: "Extensions",
		//Flags:    flagset.AuthBearerWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AuthBearer(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthBearerCommand)
}
