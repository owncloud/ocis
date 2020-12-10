package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/store/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"STORE_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"STORE_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       true,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"STORE_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       true,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"STORE_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9460",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"STORE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"STORE_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"STORE_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"STORE_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"STORE_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "store",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"STORE_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9460",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"STORE_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"STORE_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"STORE_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"STORE_GRPC_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "store",
			Usage:       "Service name",
			EnvVars:     []string{"STORE_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "data-path",
			Value:       "./data/store",
			Usage:       "location of the store data path",
			EnvVars:     []string{"STORE_DATA_PATH"},
			Destination: &cfg.Datapath,
		},
	}
}

// ListStoreWithConfig applies the config to the list commands flags.
func ListStoreWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:        "grpc-namespace",
		Value:       "com.owncloud.api",
		Usage:       "Set the base namespace for the grpc namespace",
		EnvVars:     []string{"STORE_GRPC_NAMESPACE"},
		Destination: &cfg.Service.Namespace,
	},
		&cli.StringFlag{
			Name:        "name",
			Value:       "store",
			Usage:       "Service name",
			EnvVars:     []string{"STORE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
