package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/ocs/pkg/config"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9114"),
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
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"OCS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"OCS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Value:       flags.OverrideDefaultBool(cfg.Tracing.Enabled, false),
			Usage:       "Enable sending traces",
			EnvVars:     []string{"OCS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"OCS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"OCS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"OCS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "ocs"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"OCS_TRACING_SERVICE", "OCIS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9114"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"OCS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
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
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9110"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"OCS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "ocs"),
			Usage:       "Service name",
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/ocs"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"OCS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       flags.OverrideDefaultString(cfg.TokenManager.JWTSecret, "Pive-Fumkiu4"),
			Usage:       "Used to dismantle the access token, should equal reva's jwt-secret",
			EnvVars:     []string{"OCS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},

		&cli.StringFlag{
			Name:        "account-backend-type",
			Value:       flags.OverrideDefaultString(cfg.AccountBackend, "accounts"),
			Usage:       "account-backend-type",
			EnvVars:     []string{"OCS_ACCOUNT_BACKEND_TYPE"},
			Destination: &cfg.AccountBackend,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.RevaAddress, "127.0.0.1:9142"),
			Usage:       "REVA Gateway Endpoint",
			EnvVars:     []string{"OCS_REVA_GATEWAY_ADDR"},
			Destination: &cfg.RevaAddress,
		},
		&cli.StringFlag{
			Name:        "idm-address",
			Value:       flags.OverrideDefaultString(cfg.IdentityManagement.Address, "https://localhost:9200"),
			EnvVars:     []string{"OCS_IDM_ADDRESS", "OCIS_URL"},
			Usage:       "keeps track of the IDM Address. Needed because of Reva requisite of uniqueness for users",
			Destination: &cfg.IdentityManagement.Address,
		},
		&cli.StringFlag{
			Name:        "users-driver",
			Value:       flags.OverrideDefaultString(cfg.StorageUsersDriver, "ocis"),
			Usage:       "storage driver for users mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"OCS_STORAGE_USERS_DRIVER", "STORAGE_USERS_DRIVER"},
			Destination: &cfg.StorageUsersDriver,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}

// ListOcsWithConfig applies the config to the list commands flagset.
func ListOcsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "ocs"),
			Usage:       "Service name",
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
