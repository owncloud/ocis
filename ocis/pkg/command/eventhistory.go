package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/command"
	"github.com/urfave/cli/v2"
)

// EventHistoryCommand is the entrypoint for the eventhistory command.
func EventHistoryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.EventHistory.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.EventHistory.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.EventHistory.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.EventHistory),
	}
}

func init() {
	register.AddCommand(EventHistoryCommand)
}
