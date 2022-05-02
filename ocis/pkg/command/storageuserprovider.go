package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/user/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageUserProviderCommand is the entrypoint for the storage-userprovider command.
func StorageUserProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-userprovider",
		Usage:    "start storage userprovider service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.StorageUsers.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.User(cfg.User)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageUserProviderCommand)
}
