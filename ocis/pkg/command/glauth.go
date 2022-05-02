package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/glauth/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GLAuthCommand is the entrypoint for the glauth command.
func GLAuthCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.GLAuth.Service.Name,
		Usage:    subcommandDescription(cfg.GLAuth.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.GLAuth.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.GLAuth),
	}
}

func init() {
	register.AddCommand(GLAuthCommand)
}
