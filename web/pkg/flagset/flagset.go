package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/web/pkg/config"
	"github.com/urfave/cli/v2"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCIS_LOG_LEVEL", "WEB_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY", "WEB_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCIS_LOG_COLOR", "WEB_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9104"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"WEB_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}
