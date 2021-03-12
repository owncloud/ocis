package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageMetadata applies cfg to the root flagset
func StorageMetadata(cfg *config.Config) []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.DebugAddr, "0.0.0.0:9217"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_METADATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "grpc-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_METADATA_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.GRPCAddr, "0.0.0.0:9215"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_METADATA_GRPC_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.DataServerURL, "0.0.0.0:9216"),
			Usage:       "URL of the data-provider the storage-provider uses",
			EnvVars:     []string{"STORAGE_METADATA_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageMetadata.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "http-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.HTTPNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_METADATA_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.HTTPAddr, "0.0.0.0:9216"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_METADATA_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.HTTPAddr,
		},
		&cli.StringFlag{
			Name:        "tmp-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.TempFolder, "/var/tmp/ocis/tmp/metadata"),
			Usage:       "path to tmp folder",
			EnvVars:     []string{"STORAGE_METADATA_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageMetadata.TempFolder,
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.Driver, "ocis"),
			Usage:       "storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER"},
			Destination: &cfg.Reva.StorageMetadata.Driver,
		},

		// some drivers need to look up users at the gateway

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "endpoint to use for the gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// User provider

		&cli.StringFlag{
			Name:        "userprovider-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.Endpoint, "localhost:9144"),
			Usage:       "endpoint to use for the userprovider service",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
		},
	}

	f = append(f, TracingWithConfig(cfg)...)
	f = append(f, DebugWithConfig(cfg)...)
	f = append(f, SecretWithConfig(cfg)...)
	f = append(f, DriverEOSWithConfig(cfg)...)
	f = append(f, DriverLocalWithConfig(cfg)...)
	f = append(f, DriverOwnCloudWithConfig(cfg)...)
	f = append(f, DriverOCISWithConfig(cfg)...)
	f = append(f,
		&cli.StringFlag{
			Name:        "storage-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.Common.Root, "/var/tmp/ocis/storage/metadata"),
			Usage:       "the path to the metadata storage root",
			EnvVars:     []string{"STORAGE_METADATA_ROOT"},
			Destination: &cfg.Reva.Storages.Common.Root,
		},
	)
	return f

}
