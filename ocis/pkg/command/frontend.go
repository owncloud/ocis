package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/command"
	"github.com/urfave/cli/v2"
)

// FrontendCommand is the entrypoint for the Frontend command.
func FrontendCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Frontend.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Frontend.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Frontend.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Frontend),
	}
}

func init() {
	register.AddCommand(FrontendCommand)
}
