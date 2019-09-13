package command

import (
	"os"
	"strings"

	"github.com/owncloud/ocis-phoenix/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root is the entry point for the ocis-phoenix command.
func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ocis-phoenix",
		Short:   "Reva service for phoenix",
		Long:    ``,
		Version: version.String,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogger()
			setupConfig()
		},
	}

	cmd.PersistentFlags().String("log-level", "", "Set logging level")
	viper.BindPFlag("log.level", cmd.PersistentFlags().Lookup("log-level"))
	viper.SetDefault("log.level", "info")
	viper.BindEnv("log.level", "PHOENIX_LOG_LEVEL")

	cmd.PersistentFlags().Bool("log-pretty", false, "Enable pretty logging")
	viper.BindPFlag("log.pretty", cmd.PersistentFlags().Lookup("log-pretty"))
	viper.SetDefault("log.pretty", true)
	viper.BindEnv("log.pretty", "PHOENIX_LOG_PRETTY")

	cmd.PersistentFlags().Bool("log-color", false, "Enable colored logging")
	viper.BindPFlag("log.color", cmd.PersistentFlags().Lookup("log-color"))
	viper.SetDefault("log.color", true)
	viper.BindEnv("log.color", "PHOENIX_LOG_COLOR")

	cmd.AddCommand(Server())
	cmd.AddCommand(Health())

	return cmd
}

// setupLogger prepares the logger.
func setupLogger() {
	switch strings.ToLower(viper.GetString("log.level")) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if viper.GetBool("log.pretty") {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !viper.GetBool("log.color"),
			},
		)
	}
}

// setupConfig prepares the config.
func setupConfig() {
	viper.SetConfigName("phoenix")

	viper.AddConfigPath("/etc/ocis")
	viper.AddConfigPath("$HOME/.ocis")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Debug().
				Msg("Initializing without config file")
		case viper.UnsupportedConfigError:
			log.Fatal().
				Msg("Unsupported config type")
		default:
			if e := log.Debug(); e.Enabled() {
				log.Fatal().
					Err(err).
					Msg("Failed to read config")
			} else {
				log.Fatal().
					Msg("Failed to read config")
			}
		}
	}
}
