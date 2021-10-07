package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// SharingWithConfig applies cfg to the root flagset
func SharingWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.DebugAddr, "0.0.0.0:9151"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Sharing.DebugAddr,
		},

		// Services

		// Gateway

		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// Sharing

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_SHARING_GRPC_NETWORK"},
			Destination: &cfg.Reva.Sharing.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.GRPCAddr, "0.0.0.0:9150"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_SHARING_GRPC_ADDR"},
			Destination: &cfg.Reva.Sharing.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("usershareprovider", "publicshareprovider"), // TODO osmshareprovider
			Usage:   "--service usershareprovider [--service publicshareprovider]",
			EnvVars: []string{"STORAGE_SHARING_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "user-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserDriver, "json"),
			Usage:       "driver to use for the UserShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_USER_DRIVER"},
			Destination: &cfg.Reva.Sharing.UserDriver,
		},
		&cli.StringFlag{
			Name:        "user-json-file",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserJSONFile, "/var/tmp/ocis/storage/shares.json"),
			Usage:       "file used to persist shares for the UserShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_USER_JSON_FILE"},
			Destination: &cfg.Reva.Sharing.UserJSONFile,
		},
		&cli.StringFlag{
			Name:        "public-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.PublicDriver, "json"),
			Usage:       "driver to use for the PublicShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_DRIVER"},
			Destination: &cfg.Reva.Sharing.PublicDriver,
		},
		&cli.StringFlag{
			Name:        "public-json-file",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.PublicJSONFile, "/var/tmp/ocis/storage/publicshares.json"),
			Usage:       "file used to persist shares for the PublicShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_JSON_FILE"},
			Destination: &cfg.Reva.Sharing.PublicJSONFile,
		},
		&cli.IntFlag{
			Name:        "public-password-hash-cost",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Sharing.PublicPasswordHashCost, 11),
			Usage:       "the cost of hashing the public shares passwords",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_PASSWORD_HASH_COST"},
			Destination: &cfg.Reva.Sharing.PublicPasswordHashCost,
		},
		&cli.BoolFlag{
			Name:        "public-enable-expired-shares-cleanup",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup, true),
			Usage:       "whether to periodically delete expired public shares",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_ENABLE_EXPIRED_SHARES_CLEANUP"},
			Destination: &cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
		},
		&cli.IntFlag{
			Name:        "public-janitor-run-interval",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Sharing.PublicJanitorRunInterval, 60),
			Usage:       "the time period in seconds after which to start a janitor run",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_JANITOR_RUN_INTERVAL"},
			Destination: &cfg.Reva.Sharing.PublicJanitorRunInterval,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, SharingSQLWithConfig(cfg)...)

	return flags
}
