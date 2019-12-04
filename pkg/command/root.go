package command

import (
	"os"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/util/log"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis-graph/pkg/version"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-graph command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis-graph",
		Version:  version.String,
		Usage:    "Example service for Reva/oCIS",
		Compiled: version.Compiled(),

		Authors: []cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			NewLogger(cfg)

			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetEnvPrefix("graph")
			viper.AutomaticEnv()

			if c.IsSet("config-file") {
				viper.SetConfigFile(c.String("config-file"))
			} else {
				viper.SetConfigName("graph")

				viper.AddConfigPath("/etc/ocis")
				viper.AddConfigPath("$HOME/.ocis")
				viper.AddConfigPath("./config")
			}

			if err := viper.ReadInConfig(); err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					log.Info("Continue without config")
				case viper.UnsupportedConfigError:
					log.Fatalf("Unsupported config type: %w", err)
				default:
					log.Fatalf("Failed to read config: %w", err)
				}
			}

			if err := viper.Unmarshal(&cfg); err != nil {
				log.Fatalf("Failed to parse config: %w", err)
			}

			return nil
		},

		Commands: []cli.Command{
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

func NewLogger(cfg *config.Config) {
	switch strings.ToLower(cfg.Log.Level) {
	case "fatal":
		log.SetLevel(log.LevelFatal)
	case "error":
		log.SetLevel(log.LevelError)
	case "info":
		log.SetLevel(log.LevelInfo)
	case "warn":
		log.SetLevel(log.LevelWarn)
	case "debug":
		log.SetLevel(log.LevelDebug)
	case "trace":
		log.SetLevel(log.LevelTrace)
	default:
		log.SetLevel(log.LevelInfo)
	}

	log.Name("graph")
}
