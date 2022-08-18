package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/graph/pkg/command"
	"github.com/urfave/cli/v2"
)

// GraphCommand is the entrypoint for the graph command.
func GraphCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Graph.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Graph.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.LogError(parser.ParseConfig(cfg, true))
			cfg.Graph.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Graph),
	}
}

func init() {
	register.AddCommand(GraphCommand)
}
