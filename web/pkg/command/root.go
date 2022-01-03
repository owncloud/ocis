package command

import (
	"context"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/clihelper"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/web/pkg/config"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		// start this service
		Server(cfg),

		// interaction with this service

		// infos about this service
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the web command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "web",
		Usage:    "Serve ownCloud Web for oCIS",
		Commands: GetCommands(cfg),
	})

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	return app.Run(os.Args)
}

// SutureService allows for the web command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new web.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	cfg.Web.Commons = cfg.Commons
	return SutureService{
		cfg: cfg.Web,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
