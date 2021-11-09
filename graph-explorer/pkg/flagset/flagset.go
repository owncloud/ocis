package flagset

import (
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9136"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}
