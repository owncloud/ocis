package command

import (
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
	"github.com/urfave/cli/v2"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "health",
		Usage: "Check health status",
		Action: func(c *cli.Context) error {
			// Not implemented
			return nil
		},
	}
}
