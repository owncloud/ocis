package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/auth-basic/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuthBasicCommand is the entrypoint for the AuthBasic command.
func AuthBasicCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthBasic.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AuthBasic.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.AuthBasic.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AuthBasic),
	}
}

func init() {
	register.AddCommand(AuthBasicCommand)
}
