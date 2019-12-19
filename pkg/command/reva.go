package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/command"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// RevaCommand is the entrypoint for the konnectd command.
func RevaCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "reva",
		Usage:    "Start reva server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			return cli.HandleAction(
				command.Server(cfg.Reva).Action,
				c,
			)
		},
	}
}

func init() {
	register.AddCommand(RevaCommand)
}
