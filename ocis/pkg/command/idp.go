package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/idp/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDPCommand is the entrypoint for the idp command.
func IDPCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.IDP.Service.Name,
		Usage:    subcommandDescription(cfg.IDP.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.IDP.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.IDP),
	}
}

func init() {
	register.AddCommand(IDPCommand)
}
