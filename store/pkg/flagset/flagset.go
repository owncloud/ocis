package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/store/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9464"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"STORE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ListStoreWithConfig applies the config to the list commands flags.
func ListStoreWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:        "grpc-namespace",
		Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.api"),
		Usage:       "Set the base namespace for the grpc namespace",
		EnvVars:     []string{"STORE_GRPC_NAMESPACE"},
		Destination: &cfg.Service.Namespace,
	},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "store"),
			Usage:       "Service name",
			EnvVars:     []string{"STORE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
