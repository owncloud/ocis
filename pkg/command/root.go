package command

import (
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/spf13/viper"

	// init store manager
	_ "github.com/owncloud/ocis-accounts/pkg/store"
)

var (
	defaultConfigPaths = []string{"/etc/ocis", "$HOME/.ocis", "./config"}
	defaultFilename    = "accounts"
)

// Execute is the entry point for the ocis-accounts command.
func Execute() error {
	app := &cli.App{
		Name:    "ocis-accounts",
		Version: version.String,
		Usage:   "Example service for Reva/oCIS",

		Before: func(c *cli.Context) error {
			log := NewLogger(config.New())
			for _, v := range defaultConfigPaths {
				// location is the user's home
				if v[0] == '$' || v[0] == '~' {
					usr, _ := user.Current()
					err := godotenv.Load(path.Join(usr.HomeDir, ".ocis", defaultFilename+".env"))
					if err != nil {
						log.Debug().Msgf("ignoring missing env file on dir: %v", v)
					}
				} else {
					err := godotenv.Load(path.Join(v, defaultFilename+".env"))
					if err != nil {
						log.Debug().Msgf("ignoring missing env file on dir: %v", v)
					}
				}
			}
			return nil
		},

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Commands: []*cli.Command{
			Server(config.New()),
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
		log.Level("info"),
		log.Pretty(true),
		log.Color(true),
	)
}

// ParseConfig loads accounts configuration from Viper known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ACCOUNTS")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName(defaultFilename)

		for _, v := range defaultConfigPaths {
			viper.AddConfigPath(v)
		}
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
}
