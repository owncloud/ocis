package flagset

import (
	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9134"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ListIDPWithConfig applies the config to the list commands flags
func ListIDPWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:        "http-namespace",
		Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
		Usage:       "Set the base namespace for service discovery",
		EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
		Destination: &cfg.Service.Namespace,
	},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "idp"),
			Usage:       "Service name",
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
