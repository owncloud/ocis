package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// GatewayWithConfig applies cfg to the root flagset
func GatewayWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9143",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_GATEWAY_DEBUG_ADDR"},
			Destination: &cfg.Reva.Gateway.DebugAddr,
		},

		// REVA

		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"STORAGE_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},
		&cli.IntFlag{
			Name:        "transfer-expires",
			Value:       24 * 60 * 60, // one day
			Usage:       "Transfer token ttl in seconds",
			EnvVars:     []string{"STORAGE_TRANSFER_EXPIRES"},
			Destination: &cfg.Reva.TransferExpires,
		},

		// TODO allow configuring clients

		// Services

		// Gateway

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_GATEWAY_NETWORK"},
			Destination: &cfg.Reva.Gateway.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for storage service, can be 'http' or 'grpc'",
			EnvVars:     []string{"STORAGE_GATEWAY_PROTOCOL"},
			Destination: &cfg.Reva.Gateway.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9142",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Gateway.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("gateway", "authregistry", "storageregistry"), // TODO appregistry
			Usage:   "--service gateway [--service authregistry]",
			EnvVars: []string{"STORAGE_GATEWAY_SERVICES"},
		},
		&cli.BoolFlag{
			Name:  "commit-share-to-storage-grant",
			Value: true,
			// TODO clarify
			Usage:       "Commit shares to the share manager",
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageGrant,
		},
		&cli.BoolFlag{
			Name:  "commit-share-to-storage-ref",
			Value: true,
			// TODO clarify
			Usage:       "Commit shares to the storage",
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageRef,
		},
		&cli.StringFlag{
			Name:        "share-folder",
			Value:       "Shares",
			Usage:       "mount shares in this folder of the home storage provider",
			EnvVars:     []string{"STORAGE_GATEWAY_SHARE_FOLDER"},
			Destination: &cfg.Reva.Gateway.ShareFolder,
		},
		&cli.BoolFlag{
			Name:        "disable-home-creation-on-login",
			Usage:       "Disable creation of home folder on login",
			EnvVars:     []string{"STORAGE_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN"},
			Destination: &cfg.Reva.Gateway.DisableHomeCreationOnLogin,
		},

		// other services

		// storage registry

		&cli.StringFlag{
			Name:        "storage-registry-driver",
			Value:       "static",
			Usage:       "driver of the storage registry",
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_DRIVER"},
			Destination: &cfg.Reva.StorageRegistry.Driver,
		},
		&cli.StringSliceFlag{
			Name:    "storage-registry-rule",
			Value:   cli.NewStringSlice(),
			Usage:   `Replaces the generated storage registry rules with this set: --storage-registry-rule "/eos=localhost:9158" [--storage-registry-rule "1284d238-aa92-42ce-bdc4-0b0000009162=localhost:9162"]`,
			EnvVars: []string{"STORAGE_STORAGE_REGISTRY_RULES"},
		},

		&cli.StringFlag{
			Name:        "storage-home-provider",
			Value:       "/home",
			Usage:       "mount point of the storage provider for user homes in the global namespace",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_PROVIDER"},
			Destination: &cfg.Reva.StorageRegistry.HomeProvider,
		},

		&cli.StringFlag{
			Name:        "frontend-url",
			Value:       "https://localhost:9200",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_URL"},
			Destination: &cfg.Reva.Frontend.URL,
		},
		&cli.StringFlag{
			Name:        "datagateway-url",
			Value:       "https://localhost:9200/data",
			Usage:       "URL to use for the storage datagateway",
			EnvVars:     []string{"STORAGE_DATAGATEWAY_URL"},
			Destination: &cfg.Reva.DataGateway.URL,
		},
		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
		},
		&cli.StringFlag{
			Name:        "auth-basic-url",
			Value:       "localhost:9146",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_URL"},
			Destination: &cfg.Reva.AuthBasic.URL,
		},
		&cli.StringFlag{
			Name:        "auth-bearer-url",
			Value:       "localhost:9148",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_URL"},
			Destination: &cfg.Reva.AuthBearer.URL,
		},
		&cli.StringFlag{
			Name:        "sharing-url",
			Value:       "localhost:9150",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_SHARING_URL"},
			Destination: &cfg.Reva.Sharing.URL,
		},

		&cli.StringFlag{
			Name:        "storage-root-url",
			Value:       "localhost:9152",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_ROOT_URL"},
			Destination: &cfg.Reva.StorageRoot.URL,
		},
		&cli.StringFlag{
			Name:        "storage-root-mount-path",
			Value:       "/",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_ROOT_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageRoot.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-root-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009152",
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_STORAGE_ROOT_MOUNT_ID"},
			Destination: &cfg.Reva.StorageRoot.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-home-url",
			Value:       "localhost:9154",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_URL"},
			Destination: &cfg.Reva.StorageHome.URL,
		},
		&cli.StringFlag{
			Name:        "storage-home-mount-path",
			Value:       "/home",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-home-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009154",
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-eos-url",
			Value:       "localhost:9158",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_EOS_URL"},
			Destination: &cfg.Reva.StorageEOS.URL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-path",
			Value:       "/eos",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_EOS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageEOS.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009158",
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_STORAGE_EOS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageEOS.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-oc-url",
			Value:       "localhost:9162",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_OC_URL"},
			Destination: &cfg.Reva.StorageOC.URL,
		},
		&cli.StringFlag{
			Name:        "storage-oc-mount-path",
			Value:       "/oc",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_OC_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageOC.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-oc-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009162",
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_STORAGE_OC_MOUNT_ID"},
			Destination: &cfg.Reva.StorageOC.MountID,
		},

		&cli.StringFlag{
			Name:        "public-link-url",
			Value:       "localhost:9178",
			Usage:       "URL to use for the public links service",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_URL"},
			Destination: &cfg.Reva.StoragePublicLink.URL,
		},
		&cli.StringFlag{
			Name:        "storage-public-link-mount-path",
			Value:       "/public/",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
		},
		// public-link has no mount id
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
