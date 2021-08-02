package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// FrontendWithConfig applies cfg to the root flagset
func FrontendWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.DebugAddr, "0.0.0.0:9141"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_FRONTEND_DEBUG_ADDR"},
			Destination: &cfg.Reva.Frontend.DebugAddr,
		},

		// REVA

		&cli.StringFlag{
			Name:        "transfer-secret",
			Value:       flags.OverrideDefaultString(cfg.Reva.TransferSecret, "replace-me-with-a-transfer-secret"),
			Usage:       "Transfer secret for datagateway",
			EnvVars:     []string{"STORAGE_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},

		// OCDav

		//&cli.StringFlag{
		//	Name:        "chunk-folder",
		//	Value:       flags.OverrideDefaultString(cfg.Reva.OCDav.WebdavNamespace, "/var/tmp/ocis/tmp/chunks"),
		//	Usage:       "temp directory for chunked uploads",
		//	EnvVars:     []string{"STORAGE_CHUNK_FOLDER"},
		//	Destination: &cfg.Reva.OCDav.WebdavNamespace,
		//},

		&cli.StringFlag{
			Name:        "webdav-namespace",
			Value:       flags.OverrideDefaultString(cfg.Reva.OCDav.WebdavNamespace, "/home/"),
			Usage:       "Namespace prefix for the /webdav endpoint",
			EnvVars:     []string{"STORAGE_WEBDAV_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.WebdavNamespace,
		},

		// th/dav/files endpoint expects a username as the first path segment
		// this can eg. be set to /eos/users
		&cli.StringFlag{
			Name:        "dav-files-namespace",
			Value:       flags.OverrideDefaultString(cfg.Reva.OCDav.DavFilesNamespace, "/users/"),
			Usage:       "Namespace prefix for the webdav /dav/files endpoint",
			EnvVars:     []string{"STORAGE_DAV_FILES_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.DavFilesNamespace,
		},

		// Services

		// Frontend

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.HTTPNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_NETWORK"},
			Destination: &cfg.Reva.Frontend.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.HTTPAddr, "0.0.0.0:9140"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_ADDR"},
			Destination: &cfg.Reva.Frontend.HTTPAddr,
		},
		// please note that STORAGE_FRONTEND_PUBLIC_URL is also defined in
		// storage/pkg/flagset/gateway.go because this setting may be consumed
		// by both the gateway and frontend service
		&cli.StringFlag{
			Name:        "public-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.PublicURL, "https://localhost:9200"),
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_FRONTEND_PUBLIC_URL", "OCIS_URL"}, // STORAGE_FRONTEND_PUBLIC_URL takes precedence over OCIS_URL
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
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.DatagatewayPrefix, "data"),
			Usage:       "datagateway prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_DATAGATEWAY_PREFIX"},
			Destination: &cfg.Reva.Frontend.DatagatewayPrefix,
		},
		&cli.StringFlag{
			Name:        "ocdav-prefix",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.OCDavPrefix, ""),
			Usage:       "owncloud webdav endpoint prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_OCDAV_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCDavPrefix,
		},
		&cli.StringFlag{
			Name:        "ocs-prefix",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.OCSPrefix, "ocs"),
			Usage:       "open collaboration services endpoint prefix",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSPrefix,
		},
		&cli.StringFlag{
			Name:        "ocs-share-prefix",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.OCSSharePrefix, "/Shares"),
			Usage:       "the prefix prepended to the path of shared files",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_SHARE_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSSharePrefix,
		},
		&cli.StringFlag{
			Name:        "ocs-home-namespace",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.OCSHomeNamespace, "/home"),
			Usage:       "the prefix prepended to the incoming requests in OCS",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_HOME_NAMESPACE"},
			Destination: &cfg.Reva.Frontend.OCSHomeNamespace,
		},
		&cli.IntFlag{
			Name:        "ocs-resource-info-cache-ttl",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Frontend.OCSResourceInfoCacheTTL, 0),
			Usage:       "the TTL for statted resources in the share cache",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL"},
			Destination: &cfg.Reva.Frontend.OCSResourceInfoCacheTTL,
		},
		&cli.StringFlag{
			Name:        "ocs-cache-warmup-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Frontend.OCSCacheWarmupDriver, ""),
			Usage:       "the driver to be used for warming up the share cache",
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_CACHE_WARMUP_DRIVER"},
			Destination: &cfg.Reva.Frontend.OCSCacheWarmupDriver,
		},
		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// Chunking
		&cli.StringFlag{
			Name:        "default-upload-protocol",
			Value:       flags.OverrideDefaultString(cfg.Reva.DefaultUploadProtocol, "tus"),
			Usage:       "Default upload chunking protocol to be used out of tus/v1/ng",
			EnvVars:     []string{"STORAGE_FRONTEND_DEFAULT_UPLOAD_PROTOCOL"},
			Destination: &cfg.Reva.DefaultUploadProtocol,
		},
		&cli.IntFlag{
			Name:        "upload-max-chunk-size",
			Value:       flags.OverrideDefaultInt(cfg.Reva.UploadMaxChunkSize, 0),
			Usage:       "Max chunk size in bytes to advertise to clients through capabilities, or 0 for unlimited",
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_MAX_CHUNK_SIZE"},
			Destination: &cfg.Reva.UploadMaxChunkSize,
		},
		&cli.StringFlag{
			Name:        "upload-http-method-override",
			Value:       flags.OverrideDefaultString(cfg.Reva.UploadHTTPMethodOverride, ""),
			Usage:       "Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH",
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE"},
			Destination: &cfg.Reva.UploadHTTPMethodOverride,
		},
		&cli.StringSliceFlag{
			Name:    "checksum-suppored-type",
			Value:   cli.NewStringSlice("sha1", "md5", "adler32"),
			Usage:   "--checksum-suppored-type sha1 [--checksum-suppored-type adler32]",
			EnvVars: []string{"STORAGE_FRONTEND_CHECKSUM_SUPPORTED_TYPES"},
		},
		&cli.StringFlag{
			Name:        "checksum-preferred-upload-type",
			Value:       flags.OverrideDefaultString(cfg.Reva.ChecksumPreferredUploadType, ""),
			Usage:       "Specify the preferred checksum algorithm used for uploads",
			EnvVars:     []string{"STORAGE_FRONTEND_CHECKSUM_PREFERRED_UPLOAD_TYPE"},
			Destination: &cfg.Reva.ChecksumPreferredUploadType,
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
	flags = append(flags, SharingSQLWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)

	return flags
}
