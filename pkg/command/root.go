package command

import (
	"os"

	"github.com/micro/cli"
	// init store manager
	_ "github.com/owncloud/ocis-accounts/pkg/store"
	"github.com/owncloud/ocis-hello/pkg/version"
)

// Execute is the entry point for the ocis-accounts command.
func Execute() error {
	app := &cli.App{
		Name:    "ocis-accounts",
		Version: version.String,
		Usage:   "Example service for Reva/oCIS",

		Authors: []cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Commands: []cli.Command{
			Server(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version,v",
		Usage: "Print the version",
	}

	return app.Run(os.Args)
}
