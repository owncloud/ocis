package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// AuthMachineWithConfig applies cfg to the root flagset
func AuthMachineWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthMachine.DebugAddr, "127.0.0.1:9167"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthMachine.DebugAddr,
		},

		// Machine Auth

		&cli.StringFlag{
			Name:        "machine-auth-api-key",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthMachineConfig.MachineAuthAPIKey, "change-me-please"),
			Usage:       "the API key to be used for the machine auth driver in reva",
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_AUTH_API_KEY", "OCIS_MACHINE_AUTH_API_KEY"},
			Destination: &cfg.Reva.AuthMachineConfig.MachineAuthAPIKey,
		},

		// Services

		// AuthMachine

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthMachine.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthMachine.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthMachine.GRPCAddr, "127.0.0.1:9166"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthMachine.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"), // TODO preferences
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_AUTH_MACHINE_SERVICES"},
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

	return flags
}
