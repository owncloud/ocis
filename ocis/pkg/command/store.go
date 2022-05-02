package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/store/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:     cfg.Store.Service.Name,
		Usage:    subcommandDescription(cfg.Store.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
			}
			return err
		},
		Subcommands: command.GetCommands(cfg.Store),
	}
}

func init() {
	register.AddCommand(StoreCommand)
}
