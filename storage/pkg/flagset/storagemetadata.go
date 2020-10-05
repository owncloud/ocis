package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
)

// StorageMetadata applies cfg to the root flagset
func StorageMetadata(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9217",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_METADATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_METADATA_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.Network,
		},
		&cli.StringFlag{
			Name:        "provider-addr",
			Value:       "0.0.0.0:9215",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_METADATA_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.Addr,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9216",
			Usage:       "URL of the data-server the storage-provider uses",
			EnvVars:     []string{"REVA_STORAGE_METADATA_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageMetadata.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "data-server-addr",
			Value:       "0.0.0.0:9216",
			Usage:       "Address to bind the metadata data-server to",
			EnvVars:     []string{"REVA_STORAGE_METADATA_DATA_SERVER_ADDR"},
			Destination: &cfg.Reva.StorageMetadataData.Addr,
		},
		&cli.StringFlag{
			Name:        "storage-provider-driver",
			Value:       "local",
			Usage:       "storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"REVA_STORAGE_METADATA_PROVIDER_DRIVER"},
			Destination: &cfg.Reva.StorageMetadata.Driver,
		},
		&cli.StringFlag{
			Name:        "data-provider-driver",
			Value:       "local",
			Usage:       "storage driver for data-provider mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"REVA_STORAGE_METADATA_DATA_PROVIDER_DRIVER"},
			Destination: &cfg.Reva.StorageMetadataData.Driver,
		},
		&cli.StringFlag{
			Name:        "storage-root",
			Value:       "/var/tmp/ocis/metadata",
			Usage:       "the path to the metadata storage root",
			EnvVars:     []string{"REVA_STORAGE_METADATA_ROOT"},
			Destination: &cfg.Reva.Storages.Common.Root,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)
	flags = append(flags, DriverLocalWithConfig(cfg)...)
	flags = append(flags, DriverOwnCloudWithConfig(cfg)...)
	flags = append(flags, DriverOCISWithConfig(cfg)...)

	return flags

}
