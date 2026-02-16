package command

import (
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"

	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return append([]*cli.Command{
		// start this service
		Server(cfg),

		// interaction with this service

		// infos about this service
		Health(cfg),
		Version(cfg),
	}, UnifiedRoles(cfg)...)
}

// Execute is the entry point for the ocis-graph command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "graph",
		Usage:    "Serve Graph API for oCIS",
		Commands: GetCommands(cfg),
	})
	return app.RunContext(cfg.Context, os.Args)
}
