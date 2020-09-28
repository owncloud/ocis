package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
	"path"
)

// StorageMetadata applies cfg to the root flagset
func StorageMetadata(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9184",
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
			Value:       "0.0.0.0:9185",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_METADATA_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.Addr,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9187",
			Usage:       "URL of the data-server the storage-provider uses",
			EnvVars:     []string{"REVA_STORAGE_METADATA_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageMetadata.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "data-server-addr",
			Value:       "0.0.0.0:9187",
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
			Destination: &cfg.Reva.StorageMetadata.Driver,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)
	flags = append(flags, DriverLocalWithConfig(cfg)...)
	flags = append(flags, DriverOwnCloudWithConfig(cfg)...)
	flags = append(flags, DriverOCISWithConfig(cfg)...)

	// Metadata storage needs its own root
	cfg.Reva.Storages.Common.Root = path.Join(cfg.Reva.Storages.Common.Root, "metadata")
	cfg.Reva.Storages.OwnCloud.Root = path.Join(cfg.Reva.Storages.OwnCloud.Root, "metadata")
	cfg.Reva.Storages.EOS.Root = path.Join(cfg.Reva.Storages.EOS.Root, "metadata")
	cfg.Reva.Storages.Local.Root = path.Join(cfg.Reva.Storages.Local.Root, "metadata")

	return flags

}
