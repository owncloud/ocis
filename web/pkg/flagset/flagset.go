package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/web/pkg/config"
	"github.com/urfave/cli/v2"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"WEB_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"WEB_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"WEB_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9104"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"WEB_DEBUG_ADDR"},
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
			EnvVars:     []string{"WEB_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"WEB_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"WEB_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"WEB_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"WEB_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"WEB_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "web"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"WEB_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9104"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"WEB_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"WEB_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"WEB_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"WEB_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9100"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"WEB_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"WEB_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"WEB_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.IntFlag{
			Name:        "http-cache-ttl",
			Value:       flags.OverrideDefaultInt(cfg.HTTP.CacheTTL, 604800), // 7 days
			Usage:       "Set the static assets caching duration in seconds",
			EnvVars:     []string{"WEB_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       flags.OverrideDefaultString(cfg.Asset.Path, ""),
			Usage:       "Path to custom assets",
			EnvVars:     []string{"WEB_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "web-config",
			Value:       flags.OverrideDefaultString(cfg.Web.Path, ""),
			Usage:       "Path to web config",
			EnvVars:     []string{"WEB_UI_CONFIG"},
			Destination: &cfg.Web.Path,
		},
		&cli.StringFlag{
			Name:        "web-config-server",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.Server, "https://localhost:9200"),
			Usage:       "Server URL",
			EnvVars:     []string{"WEB_UI_CONFIG_SERVER", "OCIS_URL"}, // WEB_UI_CONFIG_SERVER takes precedence over OCIS_URL
			Destination: &cfg.Web.Config.Server,
		},
		&cli.StringFlag{
			Name:        "web-config-theme",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.Theme, "https://localhost:9200/themes/owncloud/theme.json"),
			Usage:       "Theme",
			EnvVars:     []string{"WEB_UI_CONFIG_THEME"},
			Destination: &cfg.Web.Config.Theme,
		},
		&cli.StringFlag{
			Name:        "web-config-version",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.Version, "0.1.0"),
			Usage:       "Version",
			EnvVars:     []string{"WEB_UI_CONFIG_VERSION"},
			Destination: &cfg.Web.Config.Version,
		},
		&cli.StringSliceFlag{
			Name:    "web-config-app",
			Value:   cli.NewStringSlice("files", "search", "media-viewer", "external"),
			Usage:   `--web-config-app files [--web-config-app draw-io]`,
			EnvVars: []string{"WEB_UI_CONFIG_APPS"},
		},
		&cli.StringFlag{
			Name:        "oidc-metadata-url",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.OpenIDConnect.MetadataURL, ""),
			Usage:       "OpenID Connect metadata URL, defaults to <WEB_OIDC_AUTHORITY>/.well-known/openid-configuration",
			EnvVars:     []string{"WEB_OIDC_METADATA_URL"},
			Destination: &cfg.Web.Config.OpenIDConnect.MetadataURL,
		},
		&cli.StringFlag{
			Name:        "oidc-authority",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.OpenIDConnect.Authority, "https://localhost:9200"),
			Usage:       "OpenID Connect authority",                 // TODO rename to Issuer
			EnvVars:     []string{"WEB_OIDC_AUTHORITY", "OCIS_URL"}, // WEB_OIDC_AUTHORITY takes precedence over OCIS_URL
			Destination: &cfg.Web.Config.OpenIDConnect.Authority,
		},
		&cli.StringFlag{
			Name:        "oidc-client-id",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.OpenIDConnect.ClientID, "web"),
			Usage:       "OpenID Connect client ID",
			EnvVars:     []string{"WEB_OIDC_CLIENT_ID"},
			Destination: &cfg.Web.Config.OpenIDConnect.ClientID,
		},
		&cli.StringFlag{
			Name:        "oidc-response-type",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.OpenIDConnect.ResponseType, "code"),
			Usage:       "OpenID Connect response type",
			EnvVars:     []string{"WEB_OIDC_RESPONSE_TYPE"},
			Destination: &cfg.Web.Config.OpenIDConnect.ResponseType,
		},
		&cli.StringFlag{
			Name:        "oidc-scope",
			Value:       flags.OverrideDefaultString(cfg.Web.Config.OpenIDConnect.Scope, "openid profile email"),
			Usage:       "OpenID Connect scope",
			EnvVars:     []string{"WEB_OIDC_SCOPE"},
			Destination: &cfg.Web.Config.OpenIDConnect.Scope,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
		},
	}
}
