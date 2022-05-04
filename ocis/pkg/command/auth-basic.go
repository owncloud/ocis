package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/auth-basic/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AuthBasicCommand is the entrypoint for the AuthBasic command.
func AuthBasicCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AuthBasic.Service.Name,
		Usage:    subcommandDescription(cfg.AuthBasic.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
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
