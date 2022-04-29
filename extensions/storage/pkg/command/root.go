package command

import (
	"os"

	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		Health(cfg),
	}
}

// Execute is the entry point for the storage command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:  "storage",
		Usage: "Storage service for oCIS",

		Commands: GetCommands(cfg),
	})

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	return app.Run(os.Args)
}

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("storage"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}
