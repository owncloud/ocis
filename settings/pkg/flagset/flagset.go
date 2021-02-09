package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/settings/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       true,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       true,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9194",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"SETTINGS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"SETTINGS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"SETTINGS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"SETTINGS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"SETTINGS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"SETTINGS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "settings",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"SETTINGS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9194",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"SETTINGS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"SETTINGS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"SETTINGS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"SETTINGS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9190",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"SETTINGS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"SETTINGS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"SETTINGS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.IntFlag{
			Name:        "http-cache-ttl",
			Value:       604800, // 7 days
			Usage:       "Set the static assets caching duration in seconds",
			EnvVars:     []string{"SETTINGS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       "0.0.0.0:9191",
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"SETTINGS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVars:     []string{"SETTINGS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"SETTINGS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "settings",
			Usage:       "service name",
			EnvVars:     []string{"SETTINGS_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "data-path",
			Value:       "/var/tmp/ocis/settings",
			Usage:       "Mount path for the storage",
			EnvVars:     []string{"SETTINGS_DATA_PATH"},
			Destination: &cfg.Service.DataPath,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"SETTINGS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
	}
}

// ListSettingsWithConfig applies list command flags to cfg
func ListSettingsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"SETTINGS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "settings",
			Usage:       "service name",
			EnvVars:     []string{"SETTINGS_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
