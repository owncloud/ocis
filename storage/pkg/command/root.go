package command

import (
	"os"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the storage command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "storage",
		Version:  version.String,
		Usage:    "Storage service for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetEnvPrefix("STORAGE")
			viper.AutomaticEnv()

			if c.IsSet("config-file") {
				viper.SetConfigFile(c.String("config-file"))
			} else {
				viper.SetConfigName("storage")

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
		},

		Commands: []*cli.Command{
			Frontend(cfg),
			Gateway(cfg),
			Users(cfg),
			Groups(cfg),
			AppProvider(cfg),
			AuthBasic(cfg),
			AuthBearer(cfg),
			Sharing(cfg),
			StorageHome(cfg),
			StorageUsers(cfg),
			StoragePublicLink(cfg),
			StorageMetadata(cfg),
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
		log.Name("storage"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}
