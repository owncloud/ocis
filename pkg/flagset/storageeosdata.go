package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StorageEOSDataWithConfig applies cfg to the root flagset
func StorageEOSDataWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9161",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_OC_DATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageEOSData.DebugAddr,
		},

		// Services

		// Storage eos data

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_NETWORK"},
			Destination: &cfg.Reva.StorageEOSData.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "http",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_PROTOCOL"},
			Destination: &cfg.Reva.StorageEOSData.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9160",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_ADDR"},
			Destination: &cfg.Reva.StorageEOSData.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9160",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_URL"},
			Destination: &cfg.Reva.StorageEOSData.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("dataprovider"),
			Usage:   "--service dataprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_EOS_DATA_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       "eos",
			Usage:       "storage driver for eos data mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_DRIVER"},
			Destination: &cfg.Reva.StorageEOSData.Driver,
		},
		&cli.StringFlag{
			Name:        "prefix",
			Value:       "data",
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_PREFIX"},
			Destination: &cfg.Reva.StorageEOSData.Prefix,
		},
		&cli.StringFlag{
			Name:        "temp-folder",
			Value:       "/var/tmp/",
			Usage:       "temp folder",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_TEMP_FOLDER"},
			Destination: &cfg.Reva.StorageEOSData.TempFolder,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva gateway service",
			EnvVars:     []string{"REVA_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},

		// User provider

		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_USERS_URL"},
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
