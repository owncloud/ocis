package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9136"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
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
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "graph-explorer"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9136"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9135"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/graph-explorer"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"GRAPH_EXPLORER_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "issuer",
			Value:       flags.OverrideDefaultString(cfg.GraphExplorer.Issuer, "https://localhost:9200"),
			Usage:       "Set the OpenID Connect Provider",
			EnvVars:     []string{"GRAPH_EXPLORER_ISSUER", "OCIS_URL"},
			Destination: &cfg.GraphExplorer.Issuer,
		},
		&cli.StringFlag{
			Name:        "client-id",
			Value:       flags.OverrideDefaultString(cfg.GraphExplorer.ClientID, "ocis-explorer.js"),
			Usage:       "Set the OpenID Client ID to send to the issuer",
			EnvVars:     []string{"GRAPH_EXPLORER_CLIENT_ID"},
			Destination: &cfg.GraphExplorer.ClientID,
		},
		&cli.StringFlag{
			Name:        "graph-url",
			Value:       flags.OverrideDefaultString(cfg.GraphExplorer.GraphURL, "https://localhost:9200/graph"),
			Usage:       "Set the url to the graph api service",
			EnvVars:     []string{"GRAPH_EXPLORER_GRAPH_URL"},
			Destination: &cfg.GraphExplorer.GraphURL,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
		},
	}
}
