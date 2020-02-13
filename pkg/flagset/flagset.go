package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-phoenix/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"PHOENIX_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"PHOENIX_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       true,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"PHOENIX_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       true,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"PHOENIX_LOG_COLOR"},
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
			EnvVars:     []string{"PHOENIX_DEBUG_ADDR"},
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
			EnvVars:     []string{"PHOENIX_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"PHOENIX_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"PHOENIX_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"PHOENIX_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "phoenix",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"PHOENIX_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9104",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"PHOENIX_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"PHOENIX_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"PHOENIX_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"PHOENIX_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9100",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"PHOENIX_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"PHOENIX_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"PHOENIX_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVars:     []string{"PHOENIX_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "web-config",
			Value:       "",
			Usage:       "Path to phoenix config",
			EnvVars:     []string{"PHOENIX_WEB_CONFIG"},
			Destination: &cfg.Phoenix.Path,
		},
		&cli.StringFlag{
			Name:        "web-config-server",
			Value:       "http://localhost:9140",
			Usage:       "Server URL",
			EnvVars:     []string{"PHOENIX_WEB_CONFIG_SERVER"},
			Destination: &cfg.Phoenix.Config.Server,
		},
		&cli.StringFlag{
			Name:        "web-config-theme",
			Value:       "owncloud",
			Usage:       "Theme",
			EnvVars:     []string{"PHOENIX_WEB_CONFIG_THEME"},
			Destination: &cfg.Phoenix.Config.Theme,
		},
		&cli.StringFlag{
			Name:        "web-config-version",
			Value:       "0.1.0",
			Usage:       "Version",
			EnvVars:     []string{"PHOENIX_WEB_CONFIG_VERSION"},
			Destination: &cfg.Phoenix.Config.Version,
		},
		&cli.StringSliceFlag{
			Name:    "web-config-app",
			Value:   cli.NewStringSlice("files", "draw-io", "pdf-viewer", "markdown-editor", "media-viewer"),
			Usage:   `--web-config-app files [--web-config-app draw-io]`,
			EnvVars: []string{"PHOENIX_WEB_CONFIG_APPS"},
		},
		&cli.StringFlag{
			Name:        "oidc-metadata-url",
			Value:       "https://localhost:9130/.well-known/openid-configuration",
			Usage:       "OpenID Connect metadata URL",
			EnvVars:     []string{"PHOENIX_OIDC_METADATA_URL"},
			Destination: &cfg.Phoenix.Config.OpenIDConnect.MetadataURL,
		},
		&cli.StringFlag{
			Name:        "oidc-authority",
			Value:       "https://localhost:9130",
			Usage:       "OpenID Connect authority", // TODO rename to Issuer
			EnvVars:     []string{"PHOENIX_OIDC_AUTHORITY"},
			Destination: &cfg.Phoenix.Config.OpenIDConnect.Authority,
		},
		&cli.StringFlag{
			Name:        "oidc-client-id",
			Value:       "phoenix",
			Usage:       "OpenID Connect client ID",
			EnvVars:     []string{"PHOENIX_OIDC_CLIENT_ID"},
			Destination: &cfg.Phoenix.Config.OpenIDConnect.ClientID,
		},
		&cli.StringFlag{
			Name:        "oidc-response-type",
			Value:       "code",
			Usage:       "OpenID Connect response type",
			EnvVars:     []string{"PHOENIX_OIDC_RESPONSE_TYPE"},
			Destination: &cfg.Phoenix.Config.OpenIDConnect.ResponseType,
		},
		&cli.StringFlag{
			Name:        "oidc-scope",
			Value:       "openid profile email",
			Usage:       "OpenID Connect scope",
			EnvVars:     []string{"PHOENIX_OIDC_SCOPE"},
			Destination: &cfg.Phoenix.Config.OpenIDConnect.Scope,
		},
	}
}
