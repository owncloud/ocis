package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/auth-basic/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageAuthBasicCommand is the entrypoint for the reva-auth-basic command.
func StorageAuthBasicCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-basic",
		Usage:    "start storage auth-basic service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.AuthBasic.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.AuthBasic(cfg.AuthBasic)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageAuthBasicCommand)
}
