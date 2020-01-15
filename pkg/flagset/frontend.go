package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// FrontendWithConfig applies cfg to the root flagset
func FrontendWithConfig(cfg *config.Config) []cli.Flag {
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

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9141",
			Usage:       "Address to bind debug server",
			EnvVar:      "REVA_FRONTEND_DEBUG_ADDR",
			Destination: &cfg.Reva.Frontend.DebugAddr,
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

		// REVA

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVar:      "REVA_JWT_SECRET",
			Destination: &cfg.Reva.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVar:      "REVA_TRANSFER_SECRET",
			Destination: &cfg.Reva.TransferSecret,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       "http://localhost:9140",
			Usage:       "OIDC issuer",
			EnvVar:      "REVA_OIDC_ISSUER",
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Usage:       "OIDC allow insecure communication",
			EnvVar:      "REVA_OIDC_INSECURE",
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		&cli.StringFlag{
			Name:        "oidc-id-claim",
			Value:       "sub", // sub is stable and defined as unique. the user manager needs to take care of the sub to user metadata lookup
			Usage:       "OIDC id claim",
			EnvVar:      "REVA_OIDC_ID_CLAIM",
			Destination: &cfg.Reva.OIDC.IDClaim,
		},

		// TODO allow configuring clients

		// Services

		// Frontend

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_FRONTEND_NETWORK",
			Destination: &cfg.Reva.Frontend.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "http",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_FRONTEND_PROTOCOL",
			Destination: &cfg.Reva.Frontend.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9140",
			Usage:       "Address to bind reva service",
			EnvVar:      "REVA_FRONTEND_ADDR",
			Destination: &cfg.Reva.Frontend.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9140",
			Usage:       "URL to use for the reva service",
			EnvVar:      "REVA_FRONTEND_URL",
			Destination: &cfg.Reva.Frontend.URL,
		},
		&cli.StringFlag{
			Name:        "services",
			Value:       "datagateway,wellknown,oidcprovider,ocdav,ocs",
			Usage:       "comma separated list of services to include",
			EnvVar:      "REVA_FRONTEND_SERVICES",
			Destination: &cfg.Reva.Frontend.Services,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva gateway service",
			EnvVar:      "REVA_GATEWAY_URL",
			Destination: &cfg.Reva.Gateway.URL,
		},
	}
}
