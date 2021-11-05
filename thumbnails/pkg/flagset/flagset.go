package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"

	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9189"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"THUMBNAILS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ListThumbnailsWithConfig applies the config to the flagset for listing thumbnails services.
func ListThumbnailsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "thumbnails"),
			Usage:       "Name of the service",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.Server.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
	}
}
