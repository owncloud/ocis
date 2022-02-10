package command

import (
	"fmt"

	"github.com/cs3org/reva/pkg/events/server"
	"github.com/owncloud/ocis/nats/pkg/config"
	"github.com/owncloud/ocis/nats/pkg/config/parser"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			err := server.RunNatsServer()
			if err != nil {
				return err
			}
			for {
			}
		},
	}
}
