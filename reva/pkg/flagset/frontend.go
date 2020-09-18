package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// FrontendWithConfig applies cfg to the root flagset
func FrontendWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9141",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_FRONTEND_DEBUG_ADDR"},
			Destination: &cfg.Reva.Frontend.DebugAddr,
		},

		// REVA

		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"REVA_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},

		// OCDav

		&cli.StringFlag{
			Name:        "webdav-namespace",
			Value:       "/home/",
			Usage:       "Namespace prefix for the /webdav endpoint",
			EnvVars:     []string{"WEBDAV_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.WebdavNamespace,
		},

		// the /dav/files endpoint expects a username as the first path segment
		// this can eg. be set to /eos/users
		&cli.StringFlag{
			Name:        "dav-files-namespace",
			Value:       "/oc/",
			Usage:       "Namespace prefix for the webdav /dav/files endpoint",
			EnvVars:     []string{"DAV_FILES_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.DavFilesNamespace,
		},

		// Services

		// Frontend

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_FRONTEND_NETWORK"},
			Destination: &cfg.Reva.Frontend.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "http",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_FRONTEND_PROTOCOL"},
			Destination: &cfg.Reva.Frontend.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9140",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_FRONTEND_ADDR"},
			Destination: &cfg.Reva.Frontend.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "https://localhost:9200",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_FRONTEND_URL"},
			Destination: &cfg.Reva.Frontend.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("datagateway", "ocdav", "ocs"),
			Usage:   "--service ocdav [--service ocs]",
			EnvVars: []string{"REVA_FRONTEND_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "datagateway-prefix",
			Value:       "data",
			Usage:       "datagateway prefix",
			EnvVars:     []string{"REVA_FRONTEND_DATAGATEWAY_PREFIX"},
			Destination: &cfg.Reva.Frontend.DatagatewayPrefix,
		},
		&cli.StringFlag{
			Name:        "ocdav-prefix",
			Value:       "",
			Usage:       "owncloud webdav endpoint prefix",
			EnvVars:     []string{"REVA_FRONTEND_OCDAV_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCDavPrefix,
		},
		&cli.StringFlag{
			Name:        "ocs-prefix",
			Value:       "ocs",
			Usage:       "open collaboration services endpoint prefix",
			EnvVars:     []string{"REVA_FRONTEND_OCS_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSPrefix,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva gateway service",
			EnvVars:     []string{"REVA_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},

		// Chunking
		&cli.BoolFlag{
			Name:        "upload-disable-tus",
			Value:       false,
			Usage:       "Disables TUS upload mechanism",
			EnvVars:     []string{"REVA_FRONTEND_UPLOAD_DISABLE_TUS"},
			Destination: &cfg.Reva.UploadDisableTus,
		},
		&cli.IntFlag{
			Name:        "upload-max-chunk-size",
			Value:       0,
			Usage:       "Max chunk size in bytes to advertise to clients through capabilities, or 0 for unlimited",
			EnvVars:     []string{"REVA_FRONTEND_UPLOAD_MAX_CHUNK_SIZE"},
			Destination: &cfg.Reva.UploadMaxChunkSize,
		},
		&cli.StringFlag{
			Name:        "upload-http-method-override",
			Value:       "",
			Usage:       "Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH",
			EnvVars:     []string{"REVA_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE"},
			Destination: &cfg.Reva.UploadHTTPMethodOverride,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
