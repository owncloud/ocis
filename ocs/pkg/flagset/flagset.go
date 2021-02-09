package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocs/pkg/config"
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
			Value:       false,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       false,
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
			Value:       "0.0.0.0:9114",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"OCS_DEBUG_ADDR"},
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
			EnvVars:     []string{"OCS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Value:       false,
			Usage:       "Enable sending traces",
			EnvVars:     []string{"OCS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"OCS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"OCS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"OCS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "ocs",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"OCS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9114",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"OCS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"OCS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"OCS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"OCS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9110",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"OCS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "ocs",
			Usage:       "Service name",
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/ocs",
			Usage:       "Root path of http server",
			EnvVars:     []string{"OCS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Used to dismantle the access token, should equal reva's jwt-secret",
			EnvVars:     []string{"OCS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
	}
}

// ListOcsWithConfig applies the config to the list commands flagset.
func ListOcsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "ocs",
			Usage:       "Service name",
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
