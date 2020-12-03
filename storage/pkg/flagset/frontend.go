package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// FrontendWithConfig applies cfg to the root flagset
func FrontendWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9141",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_FRONTEND_DEBUG_ADDR"},
			Destination: &cfg.Reva.Frontend.DebugAddr,
		},

		// REVA

		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       "replace-me-with-a-transfer-secret",
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"STORAGE_TRANSFER_SECRET"},
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
			Value:       "/users/",
			Usage:       "Namespace prefix for the webdav /dav/files endpoint",
			EnvVars:     []string{"DAV_FILES_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.DavFilesNamespace,
		},

		// Services

		// Frontend

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_NETWORK"},
			Destination: &cfg.Reva.Frontend.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9140",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_ADDR"},
			Destination: &cfg.Reva.Frontend.HTTPAddr,
		},
		&cli.StringFlag{
			Name:        "public-url",
			Value:       "https://localhost:9200",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_PUBLIC_URL"},
			Destination: &cfg.Reva.Frontend.PublicURL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("datagateway", "ocdav", "ocs"),
			Usage:   "--service ocdav [--service ocs]",
			EnvVars: []string{"STORAGE_FRONTEND_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "datagateway-prefix",
			Value:       "data",
			Usage:       "datagateway prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_DATAGATEWAY_PREFIX"},
			Destination: &cfg.Reva.Frontend.DatagatewayPrefix,
		},
		&cli.StringFlag{
			Name:        "ocdav-prefix",
			Value:       "",
			Usage:       "owncloud webdav endpoint prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_OCDAV_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCDavPrefix,
		},
		&cli.StringFlag{
			Name:        "ocs-prefix",
			Value:       "ocs",
			Usage:       "open collaboration services endpoint prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSPrefix,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// Chunking
		&cli.BoolFlag{
			Name:        "upload-disable-tus",
			Value:       false,
			Usage:       "Disables TUS upload mechanism",
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_DISABLE_TUS"},
			Destination: &cfg.Reva.UploadDisableTus,
		},
		&cli.IntFlag{
			Name:        "upload-max-chunk-size",
			Value:       0,
			Usage:       "Max chunk size in bytes to advertise to clients through capabilities, or 0 for unlimited",
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_MAX_CHUNK_SIZE"},
			Destination: &cfg.Reva.UploadMaxChunkSize,
		},
		&cli.StringFlag{
			Name:        "upload-http-method-override",
			Value:       "",
			Usage:       "Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH",
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE"},
			Destination: &cfg.Reva.UploadHTTPMethodOverride,
		},

		// Reva Middlewares Config
		&cli.StringSliceFlag{
			Name:    "user-agent-whitelist-lock-in",
			Usage:   "--user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of comma separated [UserAgent:challenge] values, it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the storage-frontend it configures Reva.",
			EnvVars: []string{"STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT"},
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
