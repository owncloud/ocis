package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/webdav/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9119"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"WEBDAV_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ListWebdavWithConfig applies the config to the list commands flagset.
func ListWebdavWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for service discovery",
			EnvVars:     []string{"WEBDAV_HTTP_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "webdav"),
			Usage:       "Service name",
			EnvVars:     []string{"WEBDAV_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
