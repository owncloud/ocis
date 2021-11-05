package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/ocs/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9114"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"OCS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ListOcsWithConfig applies the config to the list commands flagset.
func ListOcsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "ocs"),
			Usage:       "Service name",
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
