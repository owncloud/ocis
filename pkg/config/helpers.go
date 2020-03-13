package config

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/spf13/viper"
)

// Embedded allows for proxy configuration bypassing root initialization.
// It's usage is intended when embedding a subcommand, like we do on owncloud/ocis.
// TODO this should be a generic function and not hardcode prefixes.
func Embedded(c *cli.Context, cfg *Config) {
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("PROXY")
	viper.AutomaticEnv()

	if c.IsSet("embedded-config") {
		viper.SetConfigFile(c.String("embedded-config"))
	} else {
		viper.SetConfigName("proxy")

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

	logger.Info().
		Msg("Running subcommand on embedded mode")
}

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *Config) log.Logger {
	return log.NewLogger(
		log.Name("proxy"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}
