package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/web/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// WebCommand is the entrypoint for the web command.
func WebCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Web.Service.Name,
		Usage:    subcommandDescription(cfg.Web.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Web.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Web),
	}
}

func init() {
	register.AddCommand(WebCommand)
}
