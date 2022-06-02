package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/app-registry/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AppRegistryCommand is the entrypoint for the AppRegistry command.
func AppRegistryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AppRegistry.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AppRegistry.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.AppRegistry.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AppRegistry),
	}
}

func init() {
	register.AddCommand(AppRegistryCommand)
}
