package command

import (
	"context"
	"os"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"

	"github.com/owncloud/ocis/accounts/pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

var (
	defaultConfigPaths = []string{"/etc/ocis", "$HOME/.ocis", "./config"}
	defaultFilename    = "accounts"
)

// Execute is the entry point for the ocis-accounts command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-accounts",
		Version:  version.String,
		Usage:    "Provide accounts and groups for oCIS",
		Compiled: version.Compiled(),
		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Before: func(c *cli.Context) error {
			cfg.Server.Version = version.String
			return ParseConfig(c, cfg)
		},

		Commands: []*cli.Command{
			Server(cfg),
			AddAccount(cfg),
			UpdateAccount(cfg),
			ListAccounts(cfg),
			InspectAccount(cfg),
			RemoveAccount(cfg),
			PrintVersion(cfg),
			RebuildIndex(cfg),
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

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("accounts"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads accounts configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	// create a new config and load files and env values onto it since this needs to be thread-safe.
	cnf := gofig.NewWithOptions("accounts", gofig.ParseEnv)

	// TODO(refs) add ENV + toml + json
	cnf.AddDriver(gooyaml.Driver)

	// TODO(refs) load from expected locations with the expected name
	err := cnf.LoadFiles("/Users/aunger/code/owncloud/ocis/accounts/pkg/command/accounts_example_config.yaml")
	if err != nil {
		// we have to swallow the error, since it is not mission critical a
		// config file is missing and default values are loaded instead.
		//return err
	}

	err = cnf.BindStruct("", cfg)

	return nil
}

// SutureService allows for the accounts command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new accounts.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	inheritLogging(cfg)
	if cfg.Mode == 0 {
		cfg.Accounts.Supervised = true
	}
	cfg.Accounts.Log.File = cfg.Log.File
	return SutureService{
		cfg: cfg.Accounts,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}

// inheritLogging is a poor man's global logging state tip-toeing around circular dependencies. It sets the logging
// of the service to whatever is in the higher config (in this case coming from ocis.yaml) and sets them as defaults,
// being overwritten when the extension parses its config file / env variables.
func inheritLogging(cfg *ociscfg.Config) {
	cfg.Accounts.Log.File = cfg.Log.File
	cfg.Accounts.Log.Color = cfg.Log.Color
	cfg.Accounts.Log.Pretty = cfg.Log.Pretty
	cfg.Accounts.Log.Level = cfg.Log.Level
}
