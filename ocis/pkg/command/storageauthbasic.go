package command

import (
	"github.com/owncloud/ocis/extensions/auth-basic/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageAuthBasicCommand is the entrypoint for the reva-auth-basic command.
func StorageAuthBasicCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-basic",
		Usage:    "start storage auth-basic service",
		Category: "extensions",
		Action: func(c *cli.Context) error {
			origCmd := command.AuthBasic(cfg.AuthBasic)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthBasicCommand)
}
