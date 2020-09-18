package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
)

// AuthBearerWithConfig applies cfg to the root flagset
func AuthBearerWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9149",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_AUTH_BEARER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBearer.DebugAddr,
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
		&cli.StringFlag{
			Name: "oidc-uid-claim",
			Value:       "",
			Usage:       "OIDC uid claim",
			EnvVars:     []string{"REVA_OIDC_UID_CLAIM"},
			Destination: &cfg.Reva.OIDC.UIDClaim,
		},
		&cli.StringFlag{
			Name: "oidc-gid-claim",
			Value:       "",
			Usage:       "OIDC gid claim",
			EnvVars:     []string{"REVA_OIDC_GID_CLAIM"},
			Destination: &cfg.Reva.OIDC.GIDClaim,
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

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
