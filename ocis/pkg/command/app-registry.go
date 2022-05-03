package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/app-registry/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AppRegistryCommand is the entrypoint for the AppRegistry command.
func AppRegistryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AppRegistry.Service.Name,
		Usage:    subcommandDescription(cfg.AppRegistry.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
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
