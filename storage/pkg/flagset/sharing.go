package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
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
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, SharingSQLWithConfig(cfg)...)

	return flags
}
