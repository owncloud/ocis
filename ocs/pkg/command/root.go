package command

import (
	"os"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocs/pkg/config"
	"github.com/owncloud/ocis/ocs/pkg/flagset"
	"github.com/owncloud/ocis/ocs/pkg/version"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-ocs command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-ocs",
		Version:  version.String,
		Usage:    "Serve OCS API for oCIS",
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
		log.Name("ocs"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}

// ParseConfig reads ocs configuration from fs.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("OCS")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName("ocs")

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
