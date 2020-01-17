package flagset

/* TODO move this into dedicated flagsets, along with storage commands

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "storage-eos-debug-addr",
			Value:       "0.0.0.0:9159",
			Usage:       "Address to bind storage eos debug server",
			EnvVar:      "REVA_STORAGE_EOS_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageEOS.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-debug-addr",
			Value:       "0.0.0.0:9161",
			Usage:       "Address to bind storage eos data debug server",
			EnvVar:      "REVA_STORAGE_HOME_DATA_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageEOSData.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "storage-s3-debug-addr",
			Value:       "0.0.0.0:9167",
			Usage:       "Address to bind storage s3 debug server",
			EnvVar:      "REVA_STORAGE_S3_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageS3.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-debug-addr",
			Value:       "0.0.0.0:9169",
			Usage:       "Address to bind storage s3 data debug server",
			EnvVar:      "REVA_STORAGE_S3_DATA_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageS3Data.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "storage-custom-debug-addr",
			Value:       "0.0.0.0:9171",
			Usage:       "Address to bind storage custom debug server",
			EnvVar:      "REVA_STORAGE_CUSTOM_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageCustom.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-debug-addr",
			Value:       "0.0.0.0:9173",
			Usage:       "Address to bind storage custom data debug server",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageCustomData.DebugAddr,
		},

		// Services

		// Storage eos

		&cli.StringFlag{
			Name:        "storage-eos-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-eos service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_EOS_NETWORK",
			Destination: &cfg.Reva.StorageEOS.Network,
		},
		&cli.StringFlag{
			Name:        "storage-eos-protocol",
			Value:       "grpc",
			Usage:       "protocol for reva storage-eos service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_EOS_PROTOCOL",
			Destination: &cfg.Reva.StorageEOS.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-eos-addr",
			Value:       "0.0.0.0:9158",
			Usage:       "Address to bind reva storage-eos service",
			EnvVar:      "REVA_STORAGE_EOS_ADDR",
			Destination: &cfg.Reva.StorageEOS.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-eos-url",
			Value:       "localhost:9158",
			Usage:       "URL to use for the reva storage-eos service",
			EnvVar:      "REVA_STORAGE_EOS_URL",
			Destination: &cfg.Reva.StorageEOS.URL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-services",
			Value:       "storageprovider",
			Usage:       "comma separated list of services to include in the storage-eos service",
			EnvVar:      "REVA_STORAGE_EOS_SERVICES",
			Destination: &cfg.Reva.StorageEOS.Services,
		},

		&cli.StringFlag{
			Name:        "storage-eos-driver",
			Value:       "local",
			Usage:       "eos storage driver",
			EnvVar:      "REVA_STORAGE_EOS_DRIVER",
			Destination: &cfg.Reva.StorageEOS.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-eos-path-wrapper",
			Value:       "",
			Usage:       "eos storage path wrapper",
			EnvVar:      "REVA_STORAGE_EOS_PATH_WRAPPER",
			Destination: &cfg.Reva.StorageEOS.PathWrapper,
		},
		&cli.StringFlag{
			Name:        "storage-eos-path-wrapper-context-prefix",
			Value:       "",
			Usage:       "eos storage path wrapper context prefix",
			EnvVar:      "REVA_STORAGE_EOS_PATH_WRAPPER_CONTEXT_PREFIX",
			Destination: &cfg.Reva.StorageEOS.PathWrapperContext.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-path",
			Value:       "/eos",
			Usage:       "eos storage mount path",
			EnvVar:      "REVA_STORAGE_EOS_MOUNT_PATH",
			Destination: &cfg.Reva.StorageEOS.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-eos-mount-id",
			Value:       "",
			Usage:       "eos storage mount id",
			EnvVar:      "REVA_STORAGE_EOS_MOUNT_ID",
			Destination: &cfg.Reva.StorageEOS.MountID,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-expose-data-server",
			Usage:       "eos storage exposes a dedicated data server",
			EnvVar:      "REVA_STORAGE_EOS_EXPOSE_DATA_SERVER",
			Destination: &cfg.Reva.StorageEOS.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-server-url",
			Value:       "",
			Usage:       "eos storage data server url",
			EnvVar:      "REVA_STORAGE_EOS_DATA_SERVER_URL",
			Destination: &cfg.Reva.StorageEOS.DataServerURL,
		},

		// Storage eos data

		&cli.StringFlag{
			Name:        "storage-eos-data-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-eos data service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_EOS_DATA_NETWORK",
			Destination: &cfg.Reva.StorageEOSData.Network,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-protocol",
			Value:       "http",
			Usage:       "protocol for reva storage-eos data service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_EOS_DATA_PROTOCOL",
			Destination: &cfg.Reva.StorageEOSData.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-addr",
			Value:       "0.0.0.0:9160",
			Usage:       "Address to bind reva storage-eos data service",
			EnvVar:      "REVA_STORAGE_EOS_DATA_ADDR",
			Destination: &cfg.Reva.StorageEOSData.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-url",
			Value:       "localhost:9160",
			Usage:       "URL to use for the reva storage-eos data service",
			EnvVar:      "REVA_STORAGE_EOS_DATA_URL",
			Destination: &cfg.Reva.StorageEOSData.URL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-services",
			Value:       "dataprovider",
			Usage:       "comma separated list of services to include in the storage-eos data service",
			EnvVar:      "REVA_STORAGE_EOS_DATA_SERVICES",
			Destination: &cfg.Reva.StorageEOSData.Services,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-driver",
			Value:       "eos",
			Usage:       "eos data storage driver",
			EnvVar:      "REVA_STORAGE_EOS_DATA_DRIVER",
			Destination: &cfg.Reva.StorageEOSData.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-prefix",
			Value:       "data",
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVar:      "REVA_STORAGE_EOS_DATA_PREFIX",
			Destination: &cfg.Reva.StorageEOSData.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-eos-data-temp-folder",
			Value:       "/var/tmp/",
			Usage:       "storage eos data temp folder",
			EnvVar:      "REVA_STORAGE_HOME_DATA_TEMP_FOLDER",
			Destination: &cfg.Reva.StorageEOSData.TempFolder,
		},

		// Storage s3

		&cli.StringFlag{
			Name:        "storage-s3-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-oc service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_S3_NETWORK",
			Destination: &cfg.Reva.StorageS3.Network,
		},
		&cli.StringFlag{
			Name:        "storage-s3-protocol",
			Value:       "grpc",
			Usage:       "protocol for reva storage-s3 service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_S3_PROTOCOL",
			Destination: &cfg.Reva.StorageS3.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-s3-addr",
			Value:       "0.0.0.0:9166",
			Usage:       "Address to bind reva storage-s3 service",
			EnvVar:      "REVA_STORAGE_S3_ADDR",
			Destination: &cfg.Reva.StorageS3.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-s3-url",
			Value:       "localhost:9166",
			Usage:       "URL to use for the reva storage-s3 service",
			EnvVar:      "REVA_STORAGE_S3_URL",
			Destination: &cfg.Reva.StorageS3.URL,
		},
		&cli.StringFlag{
			Name:        "storage-s3-services",
			Value:       "storageprovider",
			Usage:       "comma separated list of services to include in the storage-s3 service",
			EnvVar:      "REVA_STORAGE_S3_SERVICES",
			Destination: &cfg.Reva.StorageS3.Services,
		},

		&cli.StringFlag{
			Name:        "storage-s3-driver",
			Value:       "local",
			Usage:       "s3 storage driver",
			EnvVar:      "REVA_STORAGE_S3_DRIVER",
			Destination: &cfg.Reva.StorageS3.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-s3-path-wrapper",
			Value:       "",
			Usage:       "s3 storage path wrapper",
			EnvVar:      "REVA_STORAGE_S3_PATH_WRAPPER",
			Destination: &cfg.Reva.StorageS3.PathWrapper,
		},
		&cli.StringFlag{
			Name:        "storage-s3-path-wrapper-context-prefix",
			Value:       "",
			Usage:       "s3 storage path wrapper context prefix",
			EnvVar:      "REVA_STORAGE_S3_PATH_WRAPPER_CONTEXT_PREFIX",
			Destination: &cfg.Reva.StorageS3.PathWrapperContext.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-s3-mount-path",
			Value:       "",
			Usage:       "s3 storage mount path",
			EnvVar:      "REVA_STORAGE_S3_MOUNT_PATH",
			Destination: &cfg.Reva.StorageS3.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-s3-mount-id",
			Value:       "",
			Usage:       "s3 storage mount id",
			EnvVar:      "REVA_STORAGE_S3_MOUNT_ID",
			Destination: &cfg.Reva.StorageS3.MountID,
		},
		&cli.BoolFlag{
			Name:        "storage-s3-expose-data-server",
			Usage:       "s3 storage exposes a dedicated data server",
			EnvVar:      "REVA_STORAGE_S3_EXPOSE_DATA_SERVER",
			Destination: &cfg.Reva.StorageS3.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-server-url",
			Value:       "",
			Usage:       "s3 storage data server url",
			EnvVar:      "REVA_STORAGE_S3_DATA_SERVER_URL",
			Destination: &cfg.Reva.StorageS3.DataServerURL,
		},

		// Storage s3 data

		&cli.StringFlag{
			Name:        "storage-s3-data-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-s3 data service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_S3_DATA_NETWORK",
			Destination: &cfg.Reva.StorageS3Data.Network,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-protocol",
			Value:       "http",
			Usage:       "protocol for reva storage-s3 data service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_S3_DATA_PROTOCOL",
			Destination: &cfg.Reva.StorageS3Data.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-addr",
			Value:       "0.0.0.0:9168",
			Usage:       "Address to bind reva storage-s3 data service",
			EnvVar:      "REVA_STORAGE_S3_DATA_ADDR",
			Destination: &cfg.Reva.StorageS3Data.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-url",
			Value:       "localhost:9168",
			Usage:       "URL to use for the reva storage-s3 data service",
			EnvVar:      "REVA_STORAGE_S3_DATA_URL",
			Destination: &cfg.Reva.StorageS3Data.URL,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-services",
			Value:       "dataprovider",
			Usage:       "comma separated list of services to include in the storage-s3 data service",
			EnvVar:      "REVA_STORAGE_S3_DATA_SERVICES",
			Destination: &cfg.Reva.StorageS3Data.Services,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-driver",
			Value:       "s3",
			Usage:       "s3 data storage driver",
			EnvVar:      "REVA_STORAGE_S3_DATA_DRIVER",
			Destination: &cfg.Reva.StorageS3Data.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-prefix",
			Value:       "data",
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVar:      "REVA_STORAGE_S3_DATA_PREFIX",
			Destination: &cfg.Reva.StorageS3Data.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-s3-data-temp-folder",
			Value:       "/var/tmp/",
			Usage:       "storage s3 data temp folder",
			EnvVar:      "REVA_STORAGE_S3_DATA_TEMP_FOLDER",
			Destination: &cfg.Reva.StorageS3Data.TempFolder,
		},

		// Storage custom

		&cli.StringFlag{
			Name:        "storage-custom-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-custom service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_CUSTOM_NETWORK",
			Destination: &cfg.Reva.StorageCustom.Network,
		},
		&cli.StringFlag{
			Name:        "storage-custom-protocol",
			Value:       "grpc",
			Usage:       "protocol for reva storage-custom service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_CUSTOM_PROTOCOL",
			Destination: &cfg.Reva.StorageCustom.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-custom-addr",
			Value:       "0.0.0.0:9170",
			Usage:       "Address to bind reva storage-custom service",
			EnvVar:      "REVA_STORAGE_CUSTOM_ADDR",
			Destination: &cfg.Reva.StorageCustom.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-custom-url",
			Value:       "localhost:9170",
			Usage:       "URL to use for the reva storage-custom service",
			EnvVar:      "REVA_STORAGE_CUSTOM_URL",
			Destination: &cfg.Reva.StorageCustom.URL,
		},
		&cli.StringFlag{
			Name:        "storage-custom-services",
			Value:       "storageprovider",
			Usage:       "comma separated list of services to include in the storage-custom service",
			EnvVar:      "REVA_STORAGE_CUSTOM_SERVICES",
			Destination: &cfg.Reva.StorageCustom.Services,
		},

		&cli.StringFlag{
			Name:        "storage-custom-driver",
			Value:       "local",
			Usage:       "custom storage driver",
			EnvVar:      "REVA_STORAGE_CUSTOM_DRIVER",
			Destination: &cfg.Reva.StorageCustom.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-custom-path-wrapper",
			Value:       "",
			Usage:       "custom storage path wrapper",
			EnvVar:      "REVA_STORAGE_CUSTOM_PATH_WRAPPER",
			Destination: &cfg.Reva.StorageCustom.PathWrapper,
		},
		&cli.StringFlag{
			Name:        "storage-custom-path-wrapper-context-prefix",
			Value:       "",
			Usage:       "custom storage path wrapper context prefix",
			EnvVar:      "REVA_STORAGE_CUSTOM_PATH_WRAPPER_CONTEXT_PREFIX",
			Destination: &cfg.Reva.StorageCustom.PathWrapperContext.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-custom-mount-path",
			Value:       "",
			Usage:       "custom storage mount path",
			EnvVar:      "REVA_STORAGE_CUSTOM_MOUNT_PATH",
			Destination: &cfg.Reva.StorageCustom.MountPath,
		},
		&cli.StringFlag{
			Name:        "storage-custom-mount-id",
			Value:       "",
			Usage:       "custom storage mount id",
			EnvVar:      "REVA_STORAGE_CUSTOM_MOUNT_ID",
			Destination: &cfg.Reva.StorageCustom.MountID,
		},
		&cli.BoolFlag{
			Name:        "storage-custom-expose-data-server",
			Usage:       "custom storage exposes a dedicated data server",
			EnvVar:      "REVA_STORAGE_CUSTOM_EXPOSE_DATA_SERVER",
			Destination: &cfg.Reva.StorageCustom.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-server-url",
			Value:       "",
			Usage:       "custom storage data server url",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_SERVER_URL",
			Destination: &cfg.Reva.StorageCustom.DataServerURL,
		},

		// Storage custom data

		&cli.StringFlag{
			Name:        "storage-custom-data-network",
			Value:       "tcp",
			Usage:       "Network to use for the reva storage-custom data service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_NETWORK",
			Destination: &cfg.Reva.StorageCustomData.Network,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-protocol",
			Value:       "http",
			Usage:       "protocol for reva storage-custom data service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_PROTOCOL",
			Destination: &cfg.Reva.StorageCustomData.Protocol,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-addr",
			Value:       "0.0.0.0:9172",
			Usage:       "Address to bind reva storage-custom data service",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_ADDR",
			Destination: &cfg.Reva.StorageCustomData.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-url",
			Value:       "localhost:9172",
			Usage:       "URL to use for the reva storage-custom data service",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_URL",
			Destination: &cfg.Reva.StorageCustomData.URL,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-services",
			Value:       "dataprovider",
			Usage:       "comma separated list of services to include in the storage-custom data service",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_SERVICES",
			Destination: &cfg.Reva.StorageCustomData.Services,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-driver",
			Value:       "",
			Usage:       "custom data storage driver",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_DRIVER",
			Destination: &cfg.Reva.StorageCustomData.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-prefix",
			Value:       "data",
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVar:      "REVA_STORAGE_S3_DATA_PREFIX",
			Destination: &cfg.Reva.StorageCustomData.Prefix,
		},
		&cli.StringFlag{
			Name:        "storage-custom-data-temp-folder",
			Value:       "/var/tmp/",
			Usage:       "storage custom data temp folder",
			EnvVar:      "REVA_STORAGE_CUSTOM_DATA_TEMP_FOLDER",
			Destination: &cfg.Reva.StorageCustomData.TempFolder,
		},

		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVar:      "REVA_ASSET_PATH",
			Destination: &cfg.Asset.Path,
		},
	}
}
*/
