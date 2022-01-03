package clihelper

import (
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/urfave/cli/v2"
)

func DefaultApp(app *cli.App) *cli.App {
	// version info
	app.Version = version.String
	app.Compiled = version.Compiled()

	// author info
	app.Authors = []*cli.Author{
		{
			Name:  "ownCloud GmbH",
			Email: "support@owncloud.com",
		},
	}

	// disable global version flag
	// instead we provide the version command
	app.HideVersion = true

	return app
}
