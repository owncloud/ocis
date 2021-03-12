package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AuthBearerWithConfig applies cfg to the root flagset
func AuthBearerWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBearer.DebugAddr, "0.0.0.0:9149"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBearer.DebugAddr,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       flags.OverrideDefaultString(cfg.Reva.OIDC.Issuer, "https://localhost:9200"),
			Usage:       "OIDC issuer",
			EnvVars:     []string{"STORAGE_OIDC_ISSUER", "OCIS_URL"}, // STORAGE_OIDC_ISSUER takes precedence over OCIS_URL
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Value:       flags.OverrideDefaultBool(cfg.Reva.OIDC.Insecure, true),
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"STORAGE_OIDC_INSECURE"},
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		&cli.StringFlag{
			Name: "oidc-id-claim",
			// preferred_username is a workaround
			// the user manager needs to take care of the sub to user metadata lookup, which ldap cannot do
			// TODO sub is stable and defined as unique.
			// AFAICT we want to use the account id from ocis-accounts
			// TODO add an ocis middleware to storage that changes the users opaqueid?
			// TODO add an ocis-accounts backed user manager
			Value:       flags.OverrideDefaultString(cfg.Reva.OIDC.IDClaim, "preferred_username"),
			Usage:       "OIDC id claim",
			EnvVars:     []string{"STORAGE_OIDC_ID_CLAIM"},
			Destination: &cfg.Reva.OIDC.IDClaim,
		},
		&cli.StringFlag{
			Name:        "oidc-uid-claim",
			Value:       flags.OverrideDefaultString(cfg.Reva.OIDC.UIDClaim, ""),
			Usage:       "OIDC uid claim",
			EnvVars:     []string{"STORAGE_OIDC_UID_CLAIM"},
			Destination: &cfg.Reva.OIDC.UIDClaim,
		},
		&cli.StringFlag{
			Name:        "oidc-gid-claim",
			Value:       flags.OverrideDefaultString(cfg.Reva.OIDC.GIDClaim, ""),
			Usage:       "OIDC gid claim",
			EnvVars:     []string{"STORAGE_OIDC_GID_CLAIM"},
			Destination: &cfg.Reva.OIDC.GIDClaim,
		},

		// Services

		// AuthBearer

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBearer.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthBearer.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBearer.GRPCAddr, "0.0.0.0:9148"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthBearer.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"), // TODO preferences
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_AUTH_BEARER_SERVICES"},
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "0.0.0.0:9142"),
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
