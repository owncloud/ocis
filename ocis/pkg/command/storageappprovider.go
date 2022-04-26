package command

import (
	"github.com/owncloud/ocis/extensions/appprovider/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageAppProviderCommand is the entrypoint for the reva-app-provider command.
func StorageAppProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-app-provider",
		Usage:    "start storage app-provider service",
		Category: "extensions",
		Action: func(c *cli.Context) error {
			origCmd := command.AppProvider(cfg.AppProvider)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAppProviderCommand)
}
