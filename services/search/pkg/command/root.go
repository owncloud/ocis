package command

import (
	"context"
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/thejerf/suture/v4"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		// start this service
		Server(cfg),

		// interaction with this service
		Index(cfg),

		// infos about this service
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the ocis-search command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "search",
		Usage:    "Serve search API for oCIS",
		Commands: GetCommands(cfg),
	})
	return app.Run(os.Args)
}

// SutureService allows for the search command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new search.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	cfg.Search.Commons = cfg.Commons
	return SutureService{
		cfg: cfg.Search,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
