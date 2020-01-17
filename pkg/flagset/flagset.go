package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-phoenix/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVar:      "PHOENIX_CONFIG_FILE",
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVar:      "PHOENIX_LOG_LEVEL",
			Destination: &cfg.Log.Level,
		},
		&cli.BoolTFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVar:      "PHOENIX_LOG_PRETTY",
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolTFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVar:      "PHOENIX_LOG_COLOR",
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9104",
			Usage:       "Address to debug endpoint",
			EnvVar:      "PHOENIX_DEBUG_ADDR",
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
			EnvVar:      "PHOENIX_TRACING_ENABLED",
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVar:      "PHOENIX_TRACING_TYPE",
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVar:      "PHOENIX_TRACING_ENDPOINT",
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVar:      "PHOENIX_TRACING_COLLECTOR",
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "phoenix",
			Usage:       "Service name for tracing",
			EnvVar:      "PHOENIX_TRACING_SERVICE",
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9104",
			Usage:       "Address to bind debug server",
			EnvVar:      "PHOENIX_DEBUG_ADDR",
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVar:      "PHOENIX_DEBUG_TOKEN",
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVar:      "PHOENIX_DEBUG_PPROF",
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVar:      "PHOENIX_DEBUG_ZPAGES",
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9100",
			Usage:       "Address to bind http server",
			EnvVar:      "PHOENIX_HTTP_ADDR",
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVar:      "PHOENIX_HTTP_ROOT",
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVar:      "PHOENIX_NAMESPACE",
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVar:      "PHOENIX_ASSET_PATH",
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "web-config",
			Value:       "",
			Usage:       "Path to phoenix config",
			EnvVar:      "PHOENIX_WEB_CONFIG",
			Destination: &cfg.Phoenix.Path,
		},
		&cli.StringFlag{
			Name:        "web-config-server",
			Value:       "http://localhost:9140",
			Usage:       "Server URL",
			EnvVar:      "PHOENIX_WEB_CONFIG_SERVER",
			Destination: &cfg.Phoenix.Config.Server,
		},
		&cli.StringFlag{
			Name:        "web-config-theme",
			Value:       "owncloud",
			Usage:       "Theme",
			EnvVar:      "PHOENIX_WEB_CONFIG_THEME",
			Destination: &cfg.Phoenix.Config.Theme,
		},
		&cli.StringFlag{
			Name:        "web-config-version",
			Value:       "0.1.0",
			Usage:       "Version",
			EnvVar:      "PHOENIX_WEB_CONFIG_VERSION",
			Destination: &cfg.Phoenix.Config.Version,
		},
		&cli.StringFlag{
			Name:   "web-config-apps",
			Value:  "files,pdf-viewer,markdown-editor,media-viewer",
			Usage:  `String with comma separated values. --web-config-apps "pdf-viewer, files, draw-io"`,
			EnvVar: "PHOENIX_WEB_CONFIG_APPS",
		},
		// TODO EXTERNAL APPS?
		&cli.StringFlag{
			Name:        "oidc-metadata-url",
			Value:       "http://localhost:9140/.well-known/openid-configuration",
			Usage:       "OpenID Connect metadata URL",
			EnvVar:      "PHOENIX_OIDC_METADATA_URL",
			Destination: &cfg.Phoenix.Config.OpenIDConnect.MetadataURL,
		},
		&cli.StringFlag{
			Name:        "oidc-authority",
			Value:       "http://localhost:9140",
			Usage:       "OpenID Connect authority", // TODO rename to Issuer
			EnvVar:      "PHOENIX_OIDC_AUTHORITY",
			Destination: &cfg.Phoenix.Config.OpenIDConnect.Authority,
		},
		&cli.StringFlag{
			Name:        "oidc-client-id",
			Value:       "phoenix",
			Usage:       "OpenID Connect client ID",
			EnvVar:      "PHOENIX_OIDC_CLIENT_ID",
			Destination: &cfg.Phoenix.Config.OpenIDConnect.ClientID,
		},
		&cli.StringFlag{
			Name:        "oidc-response-type",
			Value:       "code",
			Usage:       "OpenID Connect response type",
			EnvVar:      "PHOENIX_OIDC_RESPONSE_TYPE",
			Destination: &cfg.Phoenix.Config.OpenIDConnect.ResponseType,
		},
		&cli.StringFlag{
			Name:        "oidc-scope",
			Value:       "openid profile email",
			Usage:       "OpenID Connect scope",
			EnvVar:      "PHOENIX_OIDC_SCOPE",
			Destination: &cfg.Phoenix.Config.OpenIDConnect.Scope,
		},
	}
}
