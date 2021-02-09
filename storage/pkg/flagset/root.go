package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"OCIS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       false,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       false,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}
