package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/frontend/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageFrontendCommand is the entrypoint for the reva-frontend command.
func StorageFrontendCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-frontend",
		Usage:    "start storage frontend",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Frontend.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Frontend(cfg.Frontend)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageFrontendCommand)
}
