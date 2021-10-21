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

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"WEBDAV_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"WEBDAV_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"WEBDAV_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"WEBDAV_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"WEBDAV_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"WEBDAV_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"WEBDAV_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"WEBDAV_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"WEBDAV_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "webdav"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"WEBDAV_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9119"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"WEBDAV_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"WEBDAV_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"WEBDAV_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"WEBDAV_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "127.0.0.1:9115"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"WEBDAV_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for service discovery",
			EnvVars:     []string{"WEBDAV_HTTP_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-origins",
			Value:   cli.NewStringSlice("*"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"WEBDAV_CORS_ALLOW_ORIGINS", "OCIS_CORS_ALLOW_ORIGINS"},
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-methods",
			Value:   cli.NewStringSlice("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"WEBDAV_CORS_ALLOW_METHODS", "OCIS_CORS_ALLOW_METHODS"},
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-headers",
			Value:   cli.NewStringSlice("Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"WEBDAV_CORS_ALLOW_HEADERS", "OCIS_CORS_ALLOW_HEADERS"},
		},
		&cli.BoolFlag{
			Name:    "cors-allow-credentials",
			Value:   flags.OverrideDefaultBool(cfg.HTTP.CORS.AllowCredentials, true),
			Usage:   "Allow credentials for CORS",
			EnvVars: []string{"WEBDAV_CORS_ALLOW_CREDENTIALS", "OCIS_CORS_ALLOW_CREDENTIALS"},
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "webdav"),
			Usage:       "Service name",
			EnvVars:     []string{"WEBDAV_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"WEBDAV_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "ocis-public-url",
			Value:       flags.OverrideDefaultString(cfg.OcisPublicURL, "https://127.0.0.1:9200"),
			Usage:       "The domain under which oCIS is reachable",
			EnvVars:     []string{"OCIS_PUBLIC_URL", "OCIS_URL"},
			Destination: &cfg.OcisPublicURL,
		},
		&cli.StringFlag{
			Name:        "webdav-namespace",
			Value:       flags.OverrideDefaultString(cfg.WebdavNamespace, "/home"),
			Usage:       "Namespace prefix for the /webdav endpoint",
			EnvVars:     []string{"STORAGE_WEBDAV_NAMESPACE"},
			Destination: &cfg.WebdavNamespace,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
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
