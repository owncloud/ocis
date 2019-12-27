package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-webdav/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVar:      "WEBDAV_CONFIG_FILE",
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVar:      "WEBDAV_LOG_LEVEL",
			Destination: &cfg.Log.Level,
		},
		&cli.BoolTFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVar:      "WEBDAV_LOG_PRETTY",
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolTFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVar:      "WEBDAV_LOG_COLOR",
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9119",
			Usage:       "Address to debug endpoint",
			EnvVar:      "WEBDAV_DEBUG_ADDR",
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
			EnvVar:      "WEBDAV_TRACING_ENABLED",
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVar:      "WEBDAV_TRACING_TYPE",
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVar:      "WEBDAV_TRACING_ENDPOINT",
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVar:      "WEBDAV_TRACING_COLLECTOR",
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "webdav",
			Usage:       "Service name for tracing",
			EnvVar:      "WEBDAV_TRACING_SERVICE",
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9119",
			Usage:       "Address to bind debug server",
			EnvVar:      "WEBDAV_DEBUG_ADDR",
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVar:      "WEBDAV_DEBUG_TOKEN",
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVar:      "WEBDAV_DEBUG_PPROF",
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVar:      "WEBDAV_DEBUG_ZPAGES",
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9115",
			Usage:       "Address to bind http server",
			EnvVar:      "WEBDAV_HTTP_ADDR",
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for service discovery",
			EnvVar:      "WEBDAV_HTTP_NAMESPACE",
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVar:      "WEBDAV_HTTP_ROOT",
			Destination: &cfg.HTTP.Root,
		},
	}
}
