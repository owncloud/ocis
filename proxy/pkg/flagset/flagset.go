package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/proxy/pkg/config"
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
			Value:       "0.0.0.0:9109",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
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
			EnvVars:     []string{"PROXY_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"PROXY_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"PROXY_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"PROXY_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"PROXY_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "proxy",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"PROXY_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9205",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"PROXY_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"PROXY_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"PROXY_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9200",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"PROXY_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"PROXY_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVars:     []string{"PROXY_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "service-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the service namespace",
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       "proxy",
			Usage:       "Service name",
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "transport-tls-cert",
			Value:       "",
			Usage:       "Certificate file for transport encryption",
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		&cli.StringFlag{
			Name:        "transport-tls-key",
			Value:       "",
			Usage:       "Secret file for transport encryption",
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		&cli.BoolFlag{
			Name:        "tls",
			Usage:       "Use TLS (disable only if proxy is behind a TLS-terminating reverse-proxy).",
			EnvVars:     []string{"PROXY_TLS"},
			Value:       true,
			Destination: &cfg.HTTP.TLS,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"PROXY_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       "127.0.0.1:9142",
			Usage:       "REVA Gateway Endpoint",
			EnvVars:     []string{"PROXY_REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Address,
		},
		&cli.BoolFlag{
			Name:        "insecure",
			Value:       false,
			Usage:       "allow insecure communication to upstream servers",
			EnvVars:     []string{"PROXY_INSECURE_BACKENDS"},
			Destination: &cfg.InsecureBackends,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       "https://localhost:9200",
			Usage:       "OIDC issuer",
			EnvVars:     []string{"PROXY_OIDC_ISSUER", "OCIS_URL"}, // PROXY_OIDC_ISSUER takes precedence over OCIS_URL
			Destination: &cfg.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Value:       true,
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"PROXY_OIDC_INSECURE"},
			Destination: &cfg.OIDC.Insecure,
		},
		&cli.IntFlag{
			Name:        "oidc-userinfo-cache-tll",
			Value:       10,
			Usage:       "Fallback TTL in seconds for caching userinfo, when no token lifetime can be identified",
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_TTL"},
			Destination: &cfg.OIDC.UserinfoCache.TTL,
		},
		&cli.IntFlag{
			Name:        "oidc-userinfo-cache-size",
			Value:       1024,
			Usage:       "Max entries for caching userinfo",
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_SIZE"},
			Destination: &cfg.OIDC.UserinfoCache.Size,
		},

		&cli.BoolFlag{
			Name:        "autoprovision-accounts",
			Value:       false,
			Usage:       "create accounts from OIDC access tokens to learn new users",
			EnvVars:     []string{"PROXY_AUTOPROVISION_ACCOUNTS"},
			Destination: &cfg.AutoprovisionAccounts,
		},

		// Pre Signed URLs
		&cli.StringSliceFlag{
			Name:    "presignedurl-allow-method",
			Value:   cli.NewStringSlice("GET"),
			Usage:   "--presignedurl-allow-method GET [--presignedurl-allow-method POST]",
			EnvVars: []string{"PRESIGNEDURL_ALLOWED_METHODS"},
		},
		&cli.BoolFlag{
			Name:        "enable-presignedurls",
			Value:       true,
			Usage:       "Enable or disable handling the presigned urls in the proxy",
			EnvVars:     []string{"PROXY_ENABLE_PRESIGNEDURLS"},
			Destination: &cfg.PreSignedURL.Enabled,
		},

		// Basic auth
		&cli.BoolFlag{
			Name:        "enable-basic-auth",
			Value:       false,
			Usage:       "enable basic authentication",
			EnvVars:     []string{"PROXY_ENABLE_BASIC_AUTH"},
			Destination: &cfg.EnableBasicAuth,
		},

		&cli.StringFlag{
			Name:        "account-backend-type",
			Value:       "accounts",
			Usage:       "account-backend-type",
			EnvVars:     []string{"PROXY_ACCOUNT_BACKEND_TYPE"},
			Destination: &cfg.AccountBackend,
		},

		// Reva Middlewares Config
		&cli.StringSliceFlag{
			Name:    "proxy-user-agent-lock-in",
			Usage:   "--user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of [UserAgent:challenge] it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the proxy it configures the authentication middlewares.",
			EnvVars: []string{"PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT"},
		},
	}
}

// ListProxyWithConfig applies the config to the list commands flags.
func ListProxyWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "service-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the service namespace",
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       "proxy",
			Usage:       "Service name",
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
