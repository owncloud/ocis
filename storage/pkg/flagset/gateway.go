package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// GatewayWithConfig applies cfg to the root flagset
func GatewayWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.DebugAddr, "0.0.0.0:9143"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_GATEWAY_DEBUG_ADDR"},
			Destination: &cfg.Reva.Gateway.DebugAddr,
		},

		// REVA

		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       flags.OverrideDefaultString(cfg.Reva.TransferSecret, "replace-me-with-a-transfer-secret"),
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"STORAGE_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},
		&cli.IntFlag{
			Name:        "transfer-expires",
			Value:       flags.OverrideDefaultInt(cfg.Reva.TransferExpires, 24*60*60), // one day
			Usage:       "Transfer token ttl in seconds",
			EnvVars:     []string{"STORAGE_TRANSFER_EXPIRES"},
			Destination: &cfg.Reva.TransferExpires,
		},

		// Services

		// Gateway
		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_GATEWAY_GRPC_NETWORK"},
			Destination: &cfg.Reva.Gateway.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.GRPCAddr, "0.0.0.0:9142"),
			Usage:       "Address to bind REVA service",
			EnvVars:     []string{"STORAGE_GATEWAY_GRPC_ADDR"},
			Destination: &cfg.Reva.Gateway.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("gateway", "authregistry", "storageregistry", "appregistry"),
			Usage:   "--service gateway [--service authregistry]",
			EnvVars: []string{"STORAGE_GATEWAY_SERVICES"},
		},
		&cli.BoolFlag{
			Name:        "commit-share-to-storage-grant",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Gateway.CommitShareToStorageGrant, true),
			Usage:       "Commit shares to the share manager",
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageGrant,
		},
		&cli.BoolFlag{
			Name:        "commit-share-to-storage-ref",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Gateway.CommitShareToStorageRef, true),
			Usage:       "Commit shares to the storage",
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageRef,
		},
		&cli.StringFlag{
			Name:        "share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.ShareFolder, "Shares"),
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
		&cli.StringFlag{
			Name:        "storage-home-mapping",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.HomeMapping, ""),
			Usage:       "mapping template for user home paths to user-specific mount points, e.g. /home/{{substr 0 1 .Username}}",
			EnvVars:     []string{"STORAGE_GATEWAY_HOME_MAPPING"},
			Destination: &cfg.Reva.Gateway.HomeMapping,
		},
		&cli.IntFlag{
			Name:        "etag-cache-ttl",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Gateway.EtagCacheTTL, 0),
			Usage:       "TTL for the home and shares directory etags cache",
			EnvVars:     []string{"STORAGE_GATEWAY_ETAG_CACHE_TTL"},
			Destination: &cfg.Reva.Gateway.EtagCacheTTL,
		},

		// other services

		&cli.StringFlag{
			Name:        "auth-basic-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBasic.Endpoint, "localhost:9146"),
			Usage:       "endpoint to use for the basic auth provider",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_ENDPOINT"},
			Destination: &cfg.Reva.AuthBasic.Endpoint,
		},
		&cli.StringFlag{
			Name:        "auth-bearer-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.AuthBearer.Endpoint, "localhost:9148"),
			Usage:       "endpoint to use for the bearer auth provider",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_ENDPOINT"},
			Destination: &cfg.Reva.AuthBearer.Endpoint,
		},

		// storage registry

		&cli.StringFlag{
			Name:        "storage-registry-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageRegistry.Driver, "static"),
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
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageRegistry.HomeProvider, "/home"),
			Usage:       "mount point of the storage provider for user homes in the global namespace",
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_HOME_PROVIDER"},
			Destination: &cfg.Reva.StorageRegistry.HomeProvider,
		},
		&cli.StringFlag{
			Name:        "storage-registry-json",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageRegistry.JSON, ""),
			Usage:       "JSON file containing the storage registry rules",
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_JSON"},
			Destination: &cfg.Reva.StorageRegistry.JSON,
		},

		// app registry

		&cli.StringFlag{
			Name:        "app-registry-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppRegistry.Driver, "static"),
			Usage:       "driver of the app registry",
			EnvVars:     []string{"STORAGE_APP_REGISTRY_DRIVER"},
			Destination: &cfg.Reva.AppRegistry.Driver,
		},
		&cli.StringFlag{
			Name:        "app-registry-mimetypes-json",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppRegistry.MimetypesJSON, ""),
			Usage:       "JSON file containing the storage registry rules",
			EnvVars:     []string{"STORAGE_APP_REGISTRY_MIMETYPES_JSON"},
			Destination: &cfg.Reva.AppRegistry.MimetypesJSON,
		},

		// please note that STORAGE_FRONTEND_PUBLIC_URL is also defined in
		// storage/pkg/flagset/frontend.go because this setting may be consumed
		// by both the gateway and frontend service
		&cli.StringFlag{
			Name:        "public-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.PublicURL, "https://localhost:9200"),
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_PUBLIC_URL", "OCIS_URL"}, // STORAGE_FRONTEND_PUBLIC_URL takes precedence over OCIS_URL
			Destination: &cfg.Reva.Frontend.PublicURL,
		},
		&cli.StringFlag{
			Name:        "datagateway-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.DataGateway.PublicURL, ""),
			Usage:       "URL to use for the storage datagateway, defaults to <STORAGE_FRONTEND_PUBLIC_URL>/data",
			EnvVars:     []string{"STORAGE_DATAGATEWAY_PUBLIC_URL"},
			Destination: &cfg.Reva.DataGateway.PublicURL,
		},
		&cli.StringFlag{
			Name:        "userprovider-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.Endpoint, "localhost:9144"),
			Usage:       "endpoint to use for the userprovider",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
		},
		&cli.StringFlag{
			Name:        "groupprovider-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.Endpoint, "localhost:9160"),
			Usage:       "endpoint to use for the groupprovider",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Groups.Endpoint,
		},
		&cli.StringFlag{
			Name:        "sharing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.Endpoint, "localhost:9150"),
			Usage:       "endpoint to use for the storage service",
			EnvVars:     []string{"STORAGE_SHARING_ENDPOINT"},
			Destination: &cfg.Reva.Sharing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "appprovider-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.Endpoint, "localhost:9164"),
			Usage:       "endpoint to use for the app provider",
			EnvVars:     []string{"STORAGE_APPPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.AppProvider.Endpoint,
		},

		// register home storage

		&cli.StringFlag{
			Name:        "storage-home-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.Endpoint, "localhost:9154"),
			Usage:       "endpoint to use for the home storage",
			EnvVars:     []string{"STORAGE_HOME_ENDPOINT"},
			Destination: &cfg.Reva.StorageHome.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-home-mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.MountPath, "/home"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-home-mount-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.MountID, "1284d238-aa92-42ce-bdc4-0b0000009154"),
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},

		// register users storage

		&cli.StringFlag{
			Name:        "storage-users-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.Endpoint, "localhost:9157"),
			Usage:       "endpoint to use for the users storage",
			EnvVars:     []string{"STORAGE_USERS_ENDPOINT"},
			Destination: &cfg.Reva.StorageUsers.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-users-mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.MountPath, "/users"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_USERS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageUsers.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-users-mount-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.MountID, "1284d238-aa92-42ce-bdc4-0b0000009157"),
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_USERS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageUsers.MountID,
		},

		// register public link storage

		&cli.StringFlag{
			Name:        "public-link-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.StoragePublicLink.Endpoint, "localhost:9178"),
			Usage:       "endpoint to use for the public links service",
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_ENDPOINT"},
			Destination: &cfg.Reva.StoragePublicLink.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-public-link-mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StoragePublicLink.MountPath, "/public"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
		},
		// public-link has no mount id
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
