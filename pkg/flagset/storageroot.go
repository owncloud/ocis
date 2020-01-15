package flagset

import (
	"os"

	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StorageRootWithConfig applies cfg to the root flagset
func StorageRootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVar:      "REVA_TRACING_ENABLED",
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVar:      "REVA_TRACING_TYPE",
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVar:      "REVA_TRACING_ENDPOINT",
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVar:      "REVA_TRACING_COLLECTOR",
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVar:      "REVA_TRACING_SERVICE",
			Destination: &cfg.Tracing.Service,
		},

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9153",
			Usage:       "Address to bind debug server",
			EnvVar:      "REVA_STORAGE_ROOT_DEBUG_ADDR",
			Destination: &cfg.Reva.StorageRoot.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVar:      "REVA_DEBUG_TOKEN",
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVar:      "REVA_DEBUG_PPROF",
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVar:      "REVA_DEBUG_ZPAGES",
			Destination: &cfg.Debug.Zpages,
		},

		// REVA

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVar:      "REVA_JWT_SECRET",
			Destination: &cfg.Reva.JWTSecret,
		},

		// Services

		// Storage root

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_STORAGE_ROOT_NETWORK",
			Destination: &cfg.Reva.StorageRoot.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_STORAGE_ROOT_PROTOCOL",
			Destination: &cfg.Reva.StorageRoot.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9152",
			Usage:       "Address to bind reva service",
			EnvVar:      "REVA_STORAGE_ROOT_ADDR",
			Destination: &cfg.Reva.StorageRoot.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9152",
			Usage:       "URL to use for the reva service",
			EnvVar:      "REVA_STORAGE_ROOT_URL",
			Destination: &cfg.Reva.StorageRoot.URL,
		},
		&cli.StringFlag{
			Name:        "services",
			Value:       "storageprovider",
			Usage:       "comma separated list of services to include in the storage-root service",
			EnvVar:      "REVA_STORAGE_ROOT_SERVICES",
			Destination: &cfg.Reva.StorageRoot.Services,
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "local",
			Usage:       "storage driver, eg. local, eos, owncloud or s3",
			EnvVar:      "REVA_STORAGE_ROOT_DRIVER",
			Destination: &cfg.Reva.StorageRoot.Driver,
		},
		&cli.StringFlag{
			Name:        "path-wrapper",
			Value:       "",
			Usage:       "path wrapper",
			EnvVar:      "REVA_STORAGE_ROOT_PATH_WRAPPER",
			Destination: &cfg.Reva.StorageRoot.PathWrapper,
		},
		&cli.StringFlag{
			Name:        "path-wrapper-context-prefix",
			Value:       "",
			Usage:       "path wrapper context prefix",
			EnvVar:      "REVA_STORAGE_ROOT_PATH_WRAPPER_CONTEXT_PREFIX",
			Destination: &cfg.Reva.StorageRoot.PathWrapperContext.Prefix,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/",
			Usage:       "mount path",
			EnvVar:      "REVA_STORAGE_ROOT_MOUNT_PATH",
			Destination: &cfg.Reva.StorageRoot.MountPath,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       "123e4567-e89b-12d3-a456-426655440001",
			Usage:       "mount id",
			EnvVar:      "REVA_STORAGE_ROOT_MOUNT_ID",
			Destination: &cfg.Reva.StorageRoot.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Usage:       "exposes a dedicated data server",
			EnvVar:      "REVA_STORAGE_ROOT_EXPOSE_DATA_SERVER",
			Destination: &cfg.Reva.StorageRoot.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "",
			Usage:       "data server url",
			EnvVar:      "REVA_STORAGE_ROOT_DATA_SERVER_URL",
			Destination: &cfg.Reva.StorageRoot.DataServerURL,
		},

		// Storage drivers

		// Eos

		&cli.StringFlag{
			Name:        "storage-eos-namespace",
			Value:       "",
			Usage:       "Namespace for metadata operations",
			EnvVar:      "REVA_STORAGE_EOS_NAMESPACE",
			Destination: &cfg.Reva.Storages.EOS.Namespace,
		},
		&cli.StringFlag{
			Name:        "storage-eos-binary",
			Value:       "/usr/bin/eos",
			Usage:       "Location of the eos binary",
			EnvVar:      "REVA_STORAGE_EOS_BINARY",
			Destination: &cfg.Reva.Storages.EOS.EosBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-xrdcopy-binary",
			Value:       "/usr/bin/xrdcopy",
			Usage:       "Location of the xrdcopy binary",
			EnvVar:      "REVA_STORAGE_EOS_XRDCOPY_BINARY",
			Destination: &cfg.Reva.Storages.EOS.XrdcopyBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-master-url",
			Value:       "root://eos-example.org",
			Usage:       "URL of the Master EOS MGM",
			EnvVar:      "REVA_STORAGE_EOS_MASTER_URL",
			Destination: &cfg.Reva.Storages.EOS.MasterURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-slave-url",
			Value:       "root://eos-example.org",
			Usage:       "URL of the Slave EOS MGM",
			EnvVar:      "REVA_STORAGE_EOS_SLAVE_URL",
			Destination: &cfg.Reva.Storages.EOS.SlaveURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-cache-directory",
			Value:       os.TempDir(),
			Usage:       "Location on the local fs where to store reads",
			EnvVar:      "REVA_STORAGE_EOS_CACHE_DIRECTORY",
			Destination: &cfg.Reva.Storages.EOS.CacheDirectory,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-enable-logging",
			Usage:       "Enables logging of the commands executed",
			EnvVar:      "REVA_STORAGE_EOS_ENABLE_LOGGING",
			Destination: &cfg.Reva.Storages.EOS.EnableLogging,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-show-hidden-sysfiles",
			Usage:       "show internal EOS files like .sys.v# and .sys.a# files.",
			EnvVar:      "REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES",
			Destination: &cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-force-singleuser-mode",
			Usage:       "force connections to EOS to use SingleUsername",
			EnvVar:      "REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE",
			Destination: &cfg.Reva.Storages.EOS.ForceSingleUserMode,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-use-keytab",
			Usage:       "authenticate requests by using an EOS keytab",
			EnvVar:      "REVA_STORAGE_EOS_USE_KEYTAB",
			Destination: &cfg.Reva.Storages.EOS.UseKeytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-sec-protocol",
			Value:       "",
			Usage:       "the xrootd security protocol to use between the server and EOS",
			EnvVar:      "REVA_STORAGE_EOS_SEC_PROTOCOL",
			Destination: &cfg.Reva.Storages.EOS.SecProtocol,
		},
		&cli.StringFlag{
			Name:        "storage-eos-keytab",
			Value:       "",
			Usage:       "the location of the keytab to use to authenticate to EOS",
			EnvVar:      "REVA_STORAGE_EOS_KEYTAB",
			Destination: &cfg.Reva.Storages.EOS.Keytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-single-username",
			Value:       "",
			Usage:       "the username to use when SingleUserMode is enabled",
			EnvVar:      "REVA_STORAGE_EOS_SINGLE_USERNAME",
			Destination: &cfg.Reva.Storages.EOS.SingleUsername,
		},

		// local

		&cli.StringFlag{
			Name:        "storage-local-root",
			Value:       "/var/tmp/reva/root",
			Usage:       "the path to the local storage root",
			EnvVar:      "REVA_STORAGE_LOCAL_ROOT",
			Destination: &cfg.Reva.Storages.Local.Root,
		},

		// owncloud

		&cli.StringFlag{
			Name:        "storage-owncloud-datadir",
			Value:       "/var/tmp/reva/data",
			Usage:       "the path to the owncloud data directory",
			EnvVar:      "REVA_STORAGE_OWNCLOUD_DATADIR",
			Destination: &cfg.Reva.Storages.OwnCloud.Datadirectory,
		},
	}
}
