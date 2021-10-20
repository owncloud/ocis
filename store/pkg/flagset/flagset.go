package flagset

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/store/pkg/config"
	"github.com/urfave/cli/v2"
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
			Usage:       "Set logging level",
			EnvVars:     []string{"STORE_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"STORE_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"STORE_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

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

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"STORE_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"STORE_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"STORE_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"STORE_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"STORE_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "store"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"STORE_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9464"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
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
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"STORE_GRPC_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.GRPC.Addr, "127.0.0.1:9460"),
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"STORE_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "store"),
			Usage:       "Service name",
			EnvVars:     []string{"STORE_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "data-path",
			Value:       flags.OverrideDefaultString(cfg.Datapath, path.Join(defaults.BaseDataPath(), "store")),
			Usage:       "location of the store data path",
			EnvVars:     []string{"STORE_DATA_PATH"},
			Destination: &cfg.Datapath,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
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
