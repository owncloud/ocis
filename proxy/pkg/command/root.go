package command

import (
	"context"
	"os"

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

// ParseConfig loads proxy configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	conf, err := ociscfg.BindSourcesToStructs("proxy", cfg)
	if err != nil {
		return err
	}

	// load all env variables relevant to the config in the current context.
	conf.LoadOSEnv(config.GetEnv(), false)

	if err = config.UnmapEnv(conf, cfg); err != nil {
		return err
	}

	return nil
}

// SutureService allows for the proxy command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new proxy.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	inheritLogging(cfg)
	if cfg.Mode == 0 {
		cfg.Proxy.Supervised = true
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

// inheritLogging is a poor man's global logging state tip-toeing around circular dependencies. It sets the logging
// of the service to whatever is in the higher config (in this case coming from ocis.yaml) and sets them as defaults,
// being overwritten when the extension parses its config file / env variables.
func inheritLogging(cfg *ociscfg.Config) {
	cfg.Proxy.Log.File = cfg.Log.File
	cfg.Proxy.Log.Color = cfg.Log.Color
	cfg.Proxy.Log.Pretty = cfg.Log.Pretty
	cfg.Proxy.Log.Level = cfg.Log.Level
}
