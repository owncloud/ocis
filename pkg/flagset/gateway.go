package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// GatewayWithConfig applies cfg to the root flagset
func GatewayWithConfig(cfg *config.Config) []cli.Flag {
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
			Value:       "0.0.0.0:9143",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_GATEWAY_DEBUG_ADDR"},
			Destination: &cfg.Reva.Gateway.DebugAddr,
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
		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"REVA_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},
		&cli.IntFlag{
			Name:        "transfer-expires",
			Value:       24 * 60 * 60, // one day
			Usage:       "Transfer token ttl in seconds",
			EnvVars:     []string{"REVA_TRANSFER_EXPIRES"},
			Destination: &cfg.Reva.TransferExpires,
		},

		// TODO allow configuring clients

		// Services

		// Gateway

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_GATEWAY_NETWORK"},
			Destination: &cfg.Reva.Gateway.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_GATEWAY_PROTOCOL"},
			Destination: &cfg.Reva.Gateway.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9142",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Gateway.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("gateway", "authregistry", "storageregistry"), // TODO appregistry
			Usage:   "--service gateway [--service authregistry]",
			EnvVars: []string{"REVA_GATEWAY_SERVICES"},
		},
		&cli.BoolFlag{
			Name:  "commit-share-to-storage-grant",
			Value: true,
			// TODO clarify
			Usage:       "Commit shares to the share manager",
			EnvVars:     []string{"REVA_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageGrant,
		},
		&cli.BoolFlag{
			Name:  "commit-share-to-storage-ref",
			Value: true,
			// TODO clarify
			Usage:       "Commit shares to the storage",
			EnvVars:     []string{"REVA_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageRef,
		},
		&cli.StringFlag{
			Name:        "share-folder",
			Value:       "Shares",
			Usage:       "mount shares in this folder of the home storage provider",
			EnvVars:     []string{"REVA_GATEWAY_SHARE_FOLDER"},
			Destination: &cfg.Reva.Gateway.ShareFolder,
		},
		// TODO(refs) temporary workaround needed for storing link grants.
		&cli.StringFlag{
			Name:        "link_grants_file",
			Value:       "/var/tmp/reva/link_grants.json",
			Usage:       "when using a json manager, file to use as json serialized database",
			EnvVars:     []string{"REVA_GATEWAY_LINK_GRANTS_FILE"},
			Destination: &cfg.Reva.Gateway.LinkGrants,
		},
		&cli.BoolFlag{
			Name:        "disable-home-creation-on-login",
			Usage:       "Disable creation of home folder on login",
			EnvVars:     []string{"REVA_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN"},
			Destination: &cfg.Reva.Gateway.DisableHomeCreationOnLogin,
		},

		// other services

		// storage registry

		&cli.StringFlag{
			Name:        "storage-home-provider",
			Value:       "/home",
			Usage:       "mount point of the storage provider for user homes in the global namespace",
			EnvVars:     []string{"REVA_STORAGE_HOME_PROVIDER"},
			Destination: &cfg.Reva.Gateway.HomeProvider,
		},

		&cli.StringFlag{
			Name:        "frontend-url",
			Value:       "https://localhost:9200",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_FRONTEND_URL"},
			Destination: &cfg.Reva.Frontend.URL,
		},
		&cli.StringFlag{
			Name:        "datagateway-url",
			Value:       "https://localhost:9200/data",
			Usage:       "URL to use for the reva datagateway",
			EnvVars:     []string{"REVA_DATAGATEWAY_URL"},
			Destination: &cfg.Reva.DataGateway.URL,
		},
		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
		},
		&cli.StringFlag{
			Name:        "auth-basic-url",
			Value:       "localhost:9146",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_AUTH_BASIC_URL"},
			Destination: &cfg.Reva.AuthBasic.URL,
		},
		&cli.StringFlag{
			Name:        "auth-bearer-url",
			Value:       "localhost:9148",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_AUTH_BEARER_URL"},
			Destination: &cfg.Reva.AuthBearer.URL,
		},
		&cli.StringFlag{
			Name:        "sharing-url",
			Value:       "localhost:9150",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_SHARING_URL"},
			Destination: &cfg.Reva.Sharing.URL,
		},

		&cli.StringFlag{
			Name:        "storage-root-url",
			Value:       "localhost:9152",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_ROOT_URL"},
			Destination: &cfg.Reva.StorageRoot.URL,
		},
		&cli.StringFlag{
			Name:        "storage-root-mount-path",
			Value:       "/",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_ROOT_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageRoot.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-root-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009152",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_ROOT_MOUNT_ID"},
			Destination: &cfg.Reva.StorageRoot.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-home-url",
			Value:       "localhost:9154",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_HOME_URL"},
			Destination: &cfg.Reva.StorageHome.URL,
		},
		&cli.StringFlag{
			Name:        "storage-home-mount-path",
			Value:       "/home",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name: "storage-home-mount-id",
			// This is tho mount id of the /oc storage
			// set it to 1284d238-aa92-42ce-bdc4-0b0000009158 for /eos
			// Value:       "1284d238-aa92-42ce-bdc4-0b0000009162", /os
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009154", // /home
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-eos-url",
			Value:       "localhost:9158",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_EOS_URL"},
			Destination: &cfg.Reva.StorageEOS.URL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-path",
			Value:       "/eos",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_EOS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageEOS.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009158",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_EOS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageEOS.MountID,
		},

		&cli.StringFlag{
			Name:        "storage-oc-url",
			Value:       "localhost:9162",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_OC_URL"},
			Destination: &cfg.Reva.StorageOC.URL,
		},
		&cli.StringFlag{
			Name:        "storage-oc-mount-path",
			Value:       "/oc",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_OC_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageOC.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-oc-mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009162",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_OC_MOUNT_ID"},
			Destination: &cfg.Reva.StorageOC.MountID,
		},
		&cli.StringFlag{
			Name:        "public-links-url",
			Value:       "localhost:10054",
			Usage:       "URL to use for the public links service",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_URL"},
			Destination: &cfg.Reva.StoragePublicLink.URL,
		},
	}
}
