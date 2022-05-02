package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/auth-bearer/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageAuthBearerCommand is the entrypoint for the reva-auth-bearer command.
func StorageAuthBearerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-bearer",
		Usage:    "Start storage auth-bearer service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.AuthBearer.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AuthBearer(cfg.AuthBearer)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthBearerCommand)
}
