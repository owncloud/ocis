package command

import (
	"github.com/owncloud/ocis/extensions/group/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageGroupProviderCommand is the entrypoint for the storage-groupprovider command.
func StorageGroupProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-groupprovider",
		Usage:    "start storage groupprovider service",
		Category: "extensions",
		Action: func(c *cli.Context) error {
			origCmd := command.Groups(cfg.Group)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageGroupProviderCommand)
}
