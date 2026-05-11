package command

import (
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		// start this service
		Server(cfg),

		// infos about this service
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the ocis-storage-kiteworks command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "storage-kiteworks",
		Usage:    "Provide storage backed by Kiteworks in oCIS",
		Commands: GetCommands(cfg),
	})

	return app.RunContext(cfg.Context, os.Args)
}
