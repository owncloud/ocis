package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/auth-machine/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AuthMachineCommand is the entrypoint for the AuthMachine command.
func AuthMachineCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthMachine.Service.Name,
		Usage:    subcommandDescription(cfg.AuthMachine.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
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
