package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/version"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config"
	"github.com/urfave/cli/v2"
)

// Version prints the service versions of all running instances.
func Version(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "version",
		Usage:    "print the version of this binary and the running service instances",
		Category: "info",
		Action: func(c *cli.Context) error {
			fmt.Println("Version: " + version.GetString())
			fmt.Printf("Compiled: %s\n", version.Compiled())
			fmt.Println("")

			return nil
		},
	}
}
