//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageAuthMachineCommand is the entrypoint for the reva-auth-machine command.
func StorageAuthMachineCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-machine",
		Usage:    "Start storage auth-machine service",
		Category: "Extensions",
		//Flags:    flagset.AuthBearerWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AuthMachine(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthMachineCommand)
}
