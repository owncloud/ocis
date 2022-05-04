package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/idm/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// IDMCommand is the entrypoint for the idm server command.
func IDMCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "idm",
		Usage:    "idm extension commands",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.IDM.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.IDM),
	}
}

func init() {
	register.AddCommand(IDMCommand)
}
