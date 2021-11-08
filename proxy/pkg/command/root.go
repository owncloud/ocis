package command

import (
	"context"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/shared"

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

// ParseConfig loads proxy configuration. Loading will first attempt to parse config files in the expected locations
// and then parses environment variables. In the context of oCIS env variables will always overwrite values set
// in a config file.
// If this extension is run as a subcommand (i.e: ocis proxy) then there are 2 levels of config parsing:
// 1. ocis.yaml (if any)
// 2. proxy.yaml (if any)
// 3. environment variables.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	conf, err := ociscfg.BindSourcesToStructs("proxy", cfg)
	if err != nil {
		return err
	}

	conf.LoadOSEnv(config.GetEnv(), false)
	bindings := config.StructMappings(cfg)
	return ociscfg.BindEnv(conf, bindings)
}

// SutureService allows for the proxy command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new proxy.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if (cfg.Proxy.Log == shared.Log{}) {
		cfg.Proxy.Log = cfg.Log
	}
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
