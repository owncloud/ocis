package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flag-set.
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "127.0.0.1:9010",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"OCIS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}
