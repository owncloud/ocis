package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/idp/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDPCommand is the entrypoint for the idp command.
func IDPCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.IDP.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.IDP.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
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
