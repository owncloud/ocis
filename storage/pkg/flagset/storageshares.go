package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// StorageShares applies cfg to the root flagset
func StorageShares(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.DebugAddr, "0.0.0.0:9179"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_SHARES_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageShares.DebugAddr,
		},

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_SHARES_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageShares.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.GRPCAddr, "0.0.0.0:9182"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_SHARES_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageShares.GRPCAddr,
		},

		&cli.StringFlag{
			Name:        "mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.MountPath, "/home/Shares"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_SHARES_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageShares.MountPath,
		},

		&cli.StringFlag{
			Name:        "gateway-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "endpoint to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		&cli.StringFlag{
			Name:        "storage-mount-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserStorageMountId, "1284d238-aa92-42ce-bdc4-0b0000009157"),
			Usage:       "mount id of the storage that is used for accessing the shares",
			EnvVars:     []string{"STORAGE_SHARING_USER_STORAGE_MOUNT_ID"},
			Destination: &cfg.Reva.Sharing.UserStorageMountId,
		},

		&cli.StringFlag{
			Name:        "sharing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.Endpoint, "localhost:9150"),
			Usage:       "endpoint to use for the storage service",
			EnvVars:     []string{"STORAGE_SHARING_ENDPOINT"},
			Destination: &cfg.Reva.Sharing.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
