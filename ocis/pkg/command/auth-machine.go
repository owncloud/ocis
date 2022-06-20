package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/auth-machine/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuthMachineCommand is the entrypoint for the AuthMachine command.
func AuthMachineCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthMachine.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AuthMachine.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.AuthMachine.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AuthMachine),
	}
}

func init() {
	register.AddCommand(AuthMachineCommand)
}
