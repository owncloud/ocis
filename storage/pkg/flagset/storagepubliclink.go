package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StoragePublicLink applies cfg to the root flagset
func StoragePublicLink(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9179",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_DEBUG_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.DebugAddr,
		},

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_NETWORK"},
			Destination: &cfg.Reva.StoragePublicLink.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for storage service, can be 'http' or 'grpc'",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_PROTOCOL"},
			Destination: &cfg.Reva.StoragePublicLink.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9178",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9178",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_URL"},
			Destination: &cfg.Reva.StoragePublicLink.URL,
		},

		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/public/",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
		},

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
