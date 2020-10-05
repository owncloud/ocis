package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageOCDataWithConfig applies cfg to the root flagset
func StorageOCDataWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9165",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageOCData.DebugAddr,
		},

		// Services

		// Storage oc data

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_NETWORK"},
			Destination: &cfg.Reva.StorageOCData.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "http",
			Usage:       "protocol for storage service, can be 'http' or 'grpc'",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_PROTOCOL"},
			Destination: &cfg.Reva.StorageOCData.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9164",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_ADDR"},
			Destination: &cfg.Reva.StorageOCData.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9164",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_URL"},
			Destination: &cfg.Reva.StorageOCData.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("dataprovider"),
			Usage:   "--service dataprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_STORAGE_OC_DATA_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       "owncloud",
			Usage:       "storage driver for oc data mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_DRIVER"},
			Destination: &cfg.Reva.StorageOCData.Driver,
		},
		&cli.StringFlag{
			Name:        "prefix",
			Value:       "data",
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_PREFIX"},
			Destination: &cfg.Reva.StorageOCData.Prefix,
		},
		&cli.StringFlag{
			Name:        "temp-folder",
			Value:       "/var/tmp/",
			Usage:       "temp folder",
			EnvVars:     []string{"STORAGE_STORAGE_OC_DATA_TEMP_FOLDER"},
			Destination: &cfg.Reva.StorageOCData.TempFolder,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},

		// User provider

		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
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
