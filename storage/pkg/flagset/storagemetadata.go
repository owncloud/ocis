package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageMetadata applies cfg to the root flagset
func StorageMetadata(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9217",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_METADATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "grpc-network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_METADATA_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       "0.0.0.0:9215",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_METADATA_GRPC_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9216",
			Usage:       "URL of the data-provider the storage-provider uses",
			EnvVars:     []string{"STORAGE_METADATA_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageMetadata.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "http-network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_METADATA_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9216",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_METADATA_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.HTTPAddr,
		},
		&cli.StringFlag{
			Name:        "tmp-folder",
			Value:       "/var/tmp/ocis/tmp/metadata",
			Usage:       "path to tmp folder",
			EnvVars:     []string{"STORAGE_METADATA_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageMetadata.TempFolder,
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       "ocis",
			Usage:       "storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER"},
			Destination: &cfg.Reva.StorageMetadata.Driver,
		},

		// some drivers need to look up users at the gateway

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-endpoint",
			Value:       "localhost:9142",
			Usage:       "endpoint to use for the gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// User provider

		&cli.StringFlag{
			Name:        "userprovider-endpoint",
			Value:       "localhost:9144",
			Usage:       "endpoint to use for the userprovider service",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)
	flags = append(flags, DriverLocalWithConfig(cfg)...)
	flags = append(flags, DriverOwnCloudWithConfig(cfg)...)
	flags = append(flags, DriverOCISWithConfig(cfg)...)
	flags = append(flags,
		&cli.StringFlag{
			Name:        "storage-root",
			Value:       "./data/storage/metadata",
			Usage:       "the path to the metadata storage root",
			EnvVars:     []string{"STORAGE_METADATA_ROOT"},
			Destination: &cfg.Reva.Storages.Common.Root,
		},
	)
	return flags

}
