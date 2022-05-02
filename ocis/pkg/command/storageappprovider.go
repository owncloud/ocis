package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/appprovider/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageAppProviderCommand is the entrypoint for the reva-app-provider command.
func StorageAppProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-app-provider",
		Usage:    "start storage app-provider service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.AppProvider.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AppProvider(cfg.AppProvider)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAppProviderCommand)
}
