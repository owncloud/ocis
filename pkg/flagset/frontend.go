package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// FrontendWithConfig applies cfg to the root flagset
func FrontendWithConfig(cfg *config.Config) []cli.Flag {
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
			Value:       "0.0.0.0:9141",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_FRONTEND_DEBUG_ADDR"},
			Destination: &cfg.Reva.Frontend.DebugAddr,
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
		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"REVA_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},

		// OCDav

		&cli.StringFlag{
			Name:        "webdav-namespace-jail",
			Value:       "/home/",
			Usage:       "Namespace prefix for the webdav endpoints /dav/files and /webdav",
			EnvVars:     []string{"WEBDAV_NAMESPACE_JAIL"},
			Destination: &cfg.Reva.OCDav.NamespaceJail,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       "http://localhost:9140",
			Usage:       "OIDC issuer",
			EnvVars:     []string{"REVA_OIDC_ISSUER"},
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"REVA_OIDC_INSECURE"},
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		&cli.StringFlag{
			Name:        "oidc-id-claim",
			Value:       "sub", // sub is stable and defined as unique. the user manager needs to take care of the sub to user metadata lookup
			Usage:       "OIDC id claim",
			EnvVars:     []string{"REVA_OIDC_ID_CLAIM"},
			Destination: &cfg.Reva.OIDC.IDClaim,
		},

		// TODO allow configuring clients

		// Services

		// Frontend

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_FRONTEND_NETWORK"},
			Destination: &cfg.Reva.Frontend.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "http",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_FRONTEND_PROTOCOL"},
			Destination: &cfg.Reva.Frontend.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9140",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_FRONTEND_ADDR"},
			Destination: &cfg.Reva.Frontend.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9140",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_FRONTEND_URL"},
			Destination: &cfg.Reva.Frontend.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("datagateway", "wellknown", "oidcprovider", "ocdav", "ocs"),
			Usage:   "--service datagateway [--service wellknown]",
			EnvVars: []string{"REVA_FRONTEND_SERVICES"},
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva gateway service",
			EnvVars:     []string{"REVA_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},
	}
}
