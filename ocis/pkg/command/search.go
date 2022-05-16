package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/search/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// SearchCommand is the entry point for the search command.
func SearchCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Search.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Search.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Search.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Search),
	}
}

func init() {
	register.AddCommand(SearchCommand)
}
