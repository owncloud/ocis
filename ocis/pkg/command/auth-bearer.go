package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/auth-bearer/pkg/command"
	"github.com/urfave/cli/v2"
)

// AuthBearerCommand is the entrypoint for the AuthBearer command.
func AuthBearerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthBearer.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AuthBearer.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.AuthBearer.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AuthBearer),
	}
}

func init() {
	register.AddCommand(AuthBearerCommand)
}
