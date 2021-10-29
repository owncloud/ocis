package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/urfave/cli/v2"
)

// ListProxyWithConfig applies the config to the list commands flags.
func ListProxyWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "service-namespace",
			Value:       flags.OverrideDefaultString(cfg.OIDC.Issuer, "com.owncloud.web"),
			Usage:       "Set the base namespace for the service namespace",
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "proxy"),
			Usage:       "Service name",
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
