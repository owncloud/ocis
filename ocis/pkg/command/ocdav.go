package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/ocdav/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// OCDavCommand is the entrypoint for the ocdav command.
func OCDavCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "ocdav",
		Usage:    "start ocdav",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.OCDav.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.OCDav(cfg.OCDav)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(OCDavCommand)
}
