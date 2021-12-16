package command

import (
	"context"
	"os"

	"github.com/imdario/mergo"
	"github.com/thejerf/suture/v4"
	"github.com/wkloucek/envdecode"

	"github.com/owncloud/ocis/graph/pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis-graph command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-graph",
		Version:  version.String,
		Usage:    "Serve Graph API for oCIS",
		Compiled: version.Compiled(),
		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Before: func(c *cli.Context) error {
			cfg.Service.Version = version.String
			return ParseConfig(c, cfg)
		},
		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
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
		log.Name("graph"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads graph configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	_, err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	//if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
	//	cfg.Log = &shared.Log{
	//		Level:  cfg.Commons.Log.Level,
	//		Pretty: cfg.Commons.Log.Pretty,
	//		Color:  cfg.Commons.Log.Color,
	//		File:   cfg.Commons.Log.File,
	//	}
	//} else if cfg.Log == nil && cfg.Commons == nil {
	//	cfg.Log = &shared.Log{}
	//}

	// load all env variables relevant to the config in the current context.
	envCfg := config.Config{}
	if err := envdecode.Decode(&envCfg); err != nil && err.Error() != "none of the target fields were set from environment variables" {
		return err
	}

	// merge environment variable config on top of the current config
	if err := mergo.Merge(cfg, envCfg, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

// SutureService allows for the graph command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new graph.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	cfg.Graph.Commons = cfg.Commons
	return SutureService{
		cfg: cfg.Graph,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
