package command

import (
	"os"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/version"

	// init store manager
	_ "github.com/owncloud/ocis-accounts/pkg/store"
)

// Execute is the entry point for the ocis-accounts command.
func Execute() error {
	app := &cli.App{
		Name:    "ocis-accounts",
		Version: version.String,
		Usage:   "Example service for Reva/oCIS",

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Commands: []*cli.Command{
			Server(config.New()),
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
