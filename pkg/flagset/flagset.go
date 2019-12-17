package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVar:      "REVA_CONFIG_FILE",
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVar:      "REVA_LOG_LEVEL",
			Destination: &cfg.Log.Level,
		},
		&cli.BoolTFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVar:      "REVA_LOG_PRETTY",
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolTFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVar:      "REVA_LOG_COLOR",
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
			EnvVar:      "REVA_DEBUG_ADDR",
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
			EnvVar:      "REVA_TRACING_ENABLED",
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVar:      "REVA_TRACING_TYPE",
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVar:      "REVA_TRACING_ENDPOINT",
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVar:      "REVA_TRACING_COLLECTOR",
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVar:      "REVA_TRACING_SERVICE",
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9139",
			Usage:       "Address to bind debug server",
			EnvVar:      "REVA_DEBUG_ADDR",
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVar:      "REVA_DEBUG_TOKEN",
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVar:      "REVA_DEBUG_PPROF",
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVar:      "REVA_DEBUG_ZPAGES",
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "reva-http-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva http server, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_HTTP_NETWORK",
			Destination: &cfg.Reva.HTTP.Network,
		},
		&cli.StringFlag{
			Name:        "reva-http-addr",
			Value:       "0.0.0.0:9135",
			Usage:       "Address to bind http port of reva server",
			EnvVar:      "REVA_HTTP_ADDR",
			Destination: &cfg.Reva.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "reva-http-root",
			Value:       "/",
			Usage:       "Root path of reva server",
			EnvVar:      "REVA__HTTP_ROOT",
			Destination: &cfg.Reva.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "reva-grpc-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva grpc server, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_GRPC_NETWORK",
			Destination: &cfg.Reva.GRPC.Network,
		},
		&cli.StringFlag{
			Name:        "reva-grpc-addr",
			Value:       "0.0.0.0:9136",
			Usage:       "Address to bind grpc port of reva server",
			EnvVar:      "REVA_GRPC_ADDR",
			Destination: &cfg.Reva.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "reva-max-cpus",
			Value:       "2",
			Usage:       "Max number of cpus for reva server",
			EnvVar:      "REVA_MAX_CPUS",
			Destination: &cfg.Reva.MaxCPUs,
		},
		&cli.StringFlag{
			Name:        "reva-log-level",
			Value:       "info",
			Usage:       "Log level for reva server",
			EnvVar:      "REVA_LOG_LEVEL",
			Destination: &cfg.Reva.LogLevel,
		},
		&cli.StringFlag{
			Name:        "reva-jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVar:      "REVA_JWT_SECRET",
			Destination: &cfg.Reva.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "reva-authprovider-provider",
			Value:       "",
			Usage:       "URL of the OpenID Connect Provider",
			EnvVar:      "REVA_AUTHPROVIDER_PROVIDER",
			Destination: &cfg.AuthProvider.Provider,
		},
		&cli.BoolFlag{
			Name:        "reva-authprovider-insecure",
			Usage:       "Allow insecure certificates",
			EnvVar:      "REVA_AUTHPROVIDER_INSECURE",
			Destination: &cfg.AuthProvider.Insecure,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVar:      "REVA_ASSET_PATH",
			Destination: &cfg.Asset.Path,
		},
	}
}
