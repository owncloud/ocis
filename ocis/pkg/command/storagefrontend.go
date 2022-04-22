package command

import (
	"github.com/owncloud/ocis/extensions/frontend/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageFrontendCommand is the entrypoint for the reva-frontend command.
func StorageFrontendCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-frontend",
		Usage:    "start storage frontend",
		Category: "extensions",
		Action: func(c *cli.Context) error {
			origCmd := command.Frontend(cfg.Frontend)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageFrontendCommand)
}
