package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the health flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9109"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"STORAGE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}
