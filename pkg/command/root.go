package command

import (
	"os"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/api"
	"github.com/owncloud/ocis-pkg/log"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"github.com/owncloud/ocis/pkg/register"
	"github.com/owncloud/ocis/pkg/version"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-ocis command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis",
		Version:  version.String,
		Usage:    "ownCloud Infinite Scale Stack",
		Compiled: version.Compiled(),

		Authors: []cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetEnvPrefix("OCIS")
			viper.AutomaticEnv()

			if c.IsSet("config-file") {
				viper.SetConfigFile(c.String("config-file"))
			} else {
				viper.SetConfigName("ocis")

				viper.AddConfigPath("/etc/ocis")
				viper.AddConfigPath("$HOME/.ocis")
				viper.AddConfigPath("./config")
			}

			if err := viper.ReadInConfig(); err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					logger.Info().
						Msg("Continue without config")
				case viper.UnsupportedConfigError:
					logger.Fatal().
						Err(err).
						Msg("Unsupported config type")
				default:
					logger.Fatal().
						Err(err).
						Msg("Failed to read config")
				}
			}

			if err := viper.Unmarshal(&cfg); err != nil {
				logger.Fatal().
					Err(err).
					Msg("Failed to parse config")
			}

			return nil
		},

		Commands: []cli.Command{
			Server(cfg),
			Health(cfg),
		},
	}

	// register micro runtime services.
	// TODO(refs) use the registry? Some interfaces would need to be tweaked
	app.Commands = append(app.Commands, api.Commands()...)

	for _, fn := range register.Commands {
		app.Commands = append(
			app.Commands,
			fn(cfg),
		)
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
		log.Name("ocis"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}
