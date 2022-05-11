package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/ocdav/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// OCDavCommand is the entrypoint for the OCDav command.
func OCDavCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.OCDav.Service.Name,
		Usage:    subcommandDescription(cfg.OCDav.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.OCDav.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.OCDav),
	}
}

func init() {
	register.AddCommand(OCDavCommand)
}
