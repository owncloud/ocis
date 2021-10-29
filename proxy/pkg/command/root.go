package command

import (
	"context"
	"os"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis-proxy command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-proxy",
		Version:  version.String,
		Usage:    "proxy for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		//Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			cfg.Service.Version = version.String
			return nil
		},

		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
			PrintVersion(cfg),
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
		log.Name("proxy"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads proxy configuration from Viper known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	// create a new config and load files and env values onto it since this needs to be thread-safe.
	cnf := gofig.NewWithOptions("proxy", gofig.ParseEnv)

	// TODO(refs) add ENV + toml + json
	cnf.AddDriver(gooyaml.Driver)

	// TODO(refs) load from expected locations with the expected name
	err := cnf.LoadFiles("/Users/aunger/code/owncloud/ocis/proxy/pkg/command/proxy_example_config.yaml")
	if err != nil {
		panic(err)
	}

	// bind all keys to cfg, as we expect an entire proxy.[yaml, toml...] to define all keys and not only sub values.
	err = cnf.BindStruct("", cfg)

	// step 2: overwrite the config values with those from ENV variables. Sadly the library only parses config files and does
	// not support tags for env variables.

	return nil
}

// SutureService allows for the proxy command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new proxy.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Proxy.Supervised = true
	}
	cfg.Proxy.Log.File = cfg.Log.File
	return SutureService{
		cfg: cfg.Proxy,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
