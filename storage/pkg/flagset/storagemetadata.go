package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset/metadatadrivers"
	"github.com/urfave/cli/v2"
)

// StorageMetadata applies cfg to the root flagset
func StorageMetadata(cfg *config.Config) []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.DebugAddr, "127.0.0.1:9217"),
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
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.GRPCAddr, "127.0.0.1:9215"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_METADATA_GRPC_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.DataServerURL, "http://localhost:9216"),
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
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageMetadata.HTTPAddr, "127.0.0.1:9216"),
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
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
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
	f = append(f, metadatadrivers.DriverEOSWithConfig(cfg)...)
	f = append(f, metadatadrivers.DriverLocalWithConfig(cfg)...)
	f = append(f, metadatadrivers.DriverOCISWithConfig(cfg)...)
	f = append(f, metadatadrivers.DriverS3NGWithConfig(cfg)...)
	f = append(f, metadatadrivers.DriverS3WithConfig(cfg)...)

	return f

}
