package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-graph-explorer/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9136",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
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
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "graph-explorer",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9136",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
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
			Value:       "0.0.0.0:9135",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"GRAPH_EXPLORER_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "issuer",
			Value:       "https://localhost:9130",
			Usage:       "Set the OpenID Connect Provider",
			EnvVars:     []string{"GRAPH_EXPLORER_ISSUER"},
			Destination: &cfg.GraphExplorer.Issuer,
		},
		&cli.StringFlag{
			Name:        "client-id",
			Value:       "ocis-explorer.js",
			Usage:       "Set the OpenID Client ID to send to the issuer",
			EnvVars:     []string{"GRAPH_EXPLORER_CLIENT_ID"},
			Destination: &cfg.GraphExplorer.ClientID,
		},
		&cli.StringFlag{
			Name:        "graph-url",
			Value:       "http://localhost:9120",
			Usage:       "Set the url to the graph api service",
			EnvVars:     []string{"GRAPH_EXPLORER_GRAPH_URL"},
			Destination: &cfg.GraphExplorer.GraphURL,
		},
	}
}
