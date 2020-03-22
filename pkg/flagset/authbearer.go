package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// AuthBearerWithConfig applies cfg to the root flagset
func AuthBearerWithConfig(cfg *config.Config) []cli.Flag {
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
			Value:       "0.0.0.0:9149",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_AUTH_BEARER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBearer.DebugAddr,
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

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       "https://localhost:9200",
			Usage:       "OIDC issuer",
			EnvVars:     []string{"REVA_OIDC_ISSUER"},
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Value:       true,
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"REVA_OIDC_INSECURE"},
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		&cli.StringFlag{
			Name: "oidc-id-claim",
			// preferred_username is a workaround
			// the user manager needs to take care of the sub to user metadata lookup, which ldap cannot do
			// TODO sub is stable and defined as unique.
			// AFAICT we want to use the account id from ocis-accounts
			// TODO add an ocis middleware to reva that changes the users opaqueid?
			// TODO add an ocis-accounts backed user manager
			Value:       "preferred_username",
			Usage:       "OIDC id claim",
			EnvVars:     []string{"REVA_OIDC_ID_CLAIM"},
			Destination: &cfg.Reva.OIDC.IDClaim,
		},

		// Services

		// AuthBearer

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_AUTH_BEARER_NETWORK"},
			Destination: &cfg.Reva.AuthBearer.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_AUTH_BEARER_PROTOCOL"},
			Destination: &cfg.Reva.AuthBearer.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9148",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_AUTH_BEARER_ADDR"},
			Destination: &cfg.Reva.AuthBearer.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9148",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_AUTH_BEARER_URL"},
			Destination: &cfg.Reva.AuthBearer.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"), // TODO preferences
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"REVA_AUTH_BEARER_SERVICES"},
		},
	}
}
