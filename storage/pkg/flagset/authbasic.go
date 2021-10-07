package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// AuthBasicWithConfig applies cfg to the root flagset
func AuthBasicWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBasic.DebugAddr, "0.0.0.0:9147"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBasic.DebugAddr,
		},

		// Auth

		&cli.StringFlag{
			Name:        "auth-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthProvider.Driver, "ldap"),
			Usage:       "auth driver: 'demo', 'json' or 'ldap'",
			EnvVars:     []string{"STORAGE_AUTH_DRIVER"},
			Destination: &cfg.Reva.AuthProvider.Driver,
		},
		&cli.StringFlag{
			Name:        "auth-json",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthProvider.JSON, ""),
			Usage:       "Path to users.json file",
			EnvVars:     []string{"STORAGE_AUTH_JSON"},
			Destination: &cfg.Reva.AuthProvider.JSON,
		},

		// Services

		// AuthBasic

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBasic.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthBasic.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBasic.GRPCAddr, "0.0.0.0:9146"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthBasic.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"),
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_AUTH_BASIC_SERVICES"},
		},

		// Gateway

		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)

	return flags
}
