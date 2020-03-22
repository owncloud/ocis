package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// SharingWithConfig applies cfg to the root flagset
func SharingWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"REVA_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"REVA_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"REVA_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"REVA_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"REVA_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9151",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Sharing.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"REVA_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"REVA_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"REVA_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},

		// REVA

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVars:     []string{"REVA_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},

		// Services

		// Sharing

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_SHARING_NETWORK"},
			Destination: &cfg.Reva.Sharing.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_SHARING_PROTOCOL"},
			Destination: &cfg.Reva.Sharing.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9150",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_SHARING_ADDR"},
			Destination: &cfg.Reva.Sharing.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9150",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_SHARING_URL"},
			Destination: &cfg.Reva.Sharing.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("usershareprovider", "publicshareprovider"), // TODO osmshareprovider
			Usage:   "--service usershareprovider [--service publicshareprovider]",
			EnvVars: []string{"REVA_SHARING_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "user-driver",
			Value:       "json",
			Usage:       "driver to use for the UserShareProvider",
			EnvVars:     []string{"REVA_SHARING_USER_DRIVER"},
			Destination: &cfg.Reva.Sharing.UserDriver,
		},
		&cli.StringFlag{
			Name:        "public-driver",
			Value:       "memory",
			Usage:       "driver to use for the PublicShareProvider",
			EnvVars:     []string{"REVA_SHARING_PUBLIC_DRIVER"},
			Destination: &cfg.Reva.Sharing.PublicDriver,
		},
	}
}
