package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/command"
	"github.com/urfave/cli/v2"
)

// NotificationsCommand is the entrypoint for the notifications command.
func NotificationsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Notifications.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.Notifications.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.Notifications.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Notifications),
	}
}

func init() {
	register.AddCommand(NotificationsCommand)
}
