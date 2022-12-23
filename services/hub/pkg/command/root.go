package command

import (
	"context"
	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		Server(cfg),
	}
}

// Execute is the entry point for the web command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "hub",
		Usage:    "Serve ownCloud hub for oCIS",
		Commands: GetCommands(cfg),
	})

	return app.Run(os.Args)
}

// SutureService allows for the web command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new web.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	cfg.Hub.Commons = cfg.Commons
	return SutureService{
		cfg: cfg.Hub,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
