package command

import (
	"context"
	"os"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/idp/pkg/flagset"
	"github.com/owncloud/ocis/idp/pkg/version"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
)

// Execute is the entry point for the ocis-idp command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-idp",
		Version:  version.String,
		Usage:    "Serve IDP API for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

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
		log.Name("idp"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}

// ParseConfig load configuration for every extension
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	sync.ParsingViperConfig.Lock()
	defer sync.ParsingViperConfig.Unlock()
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("IDP")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName("idp")

		viper.AddConfigPath("/etc/ocis")
		viper.AddConfigPath("$HOME/.ocis")
		viper.AddConfigPath("./config")
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			logger.Debug().
				Msg("no config found on preconfigured location")
		case viper.UnsupportedConfigError:
			logger.Fatal().
				Err(err).
				Msg("unsupported config type")
		default:
			logger.Fatal().
				Err(err).
				Msg("failed to read config")
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to parse config")
	}

	return nil
}

// SutureService allows for the idp command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new idp.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.IDP.Supervised = true
	}
	return SutureService{
		cfg: cfg.IDP,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
