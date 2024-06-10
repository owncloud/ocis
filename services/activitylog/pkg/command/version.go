package command

import (
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
	"github.com/urfave/cli/v2"
)

// Version prints the service versions of all running instances.
func Version(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "version",
		Usage:    "print the version of this binary and the running service instances",
		Category: "info",
		Action: func(c *cli.Context) error {
			// not implemented
			return nil
		},
	}
}
