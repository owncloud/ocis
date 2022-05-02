package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/notifications/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// NatsServerCommand is the entrypoint for the nats server command.
func NotificationsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "notifications",
		Usage:    "start notifications service",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
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
