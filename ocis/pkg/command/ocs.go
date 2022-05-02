package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/ocs/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// OCSCommand is the entrypoint for the ocs command.
func OCSCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.OCS.Service.Name,
		Usage:    subcommandDescription(cfg.OCS.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.OCS.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.OCS),
	}
}

func init() {
	register.AddCommand(OCSCommand)
}
