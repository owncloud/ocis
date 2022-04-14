package command

import (
	"github.com/owncloud/ocis/extensions/search/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// SearchCommand is the entry point for the search command.
func SearchCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Search.Service.Name,
		Usage:    subcommandDescription(cfg.Search.Service.Name),
		Category: "extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.Search),
	}
}

func init() {
	register.AddCommand(SearchCommand)
}
