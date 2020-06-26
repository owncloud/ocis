package flagset

import (
	"os"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StorageDriversWithConfig applies cfg to the root flagset
func storageDriversWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		// Eos

		&cli.StringFlag{
			Name:        "storage-eos-namespace",
			Value:       "/eos/dockertest/reva/users",
			Usage:       "Namespace for metadata operations",
			EnvVars:     []string{"REVA_STORAGE_EOS_NAMESPACE"},
			Destination: &cfg.Reva.Storages.EOS.Namespace,
		},
		&cli.StringFlag{
			Name: "storage-eos-shadow-namespace",
			// Defaults to path.Join(c.Namespace, ".shadow")
			Usage:       "Shadow namespace where share references are stored",
			EnvVars:     []string{"REVA_STORAGE_EOS_SHADOW_NAMESPACE"},
			Destination: &cfg.Reva.Storages.EOS.ShadowNamespace,
		},
		&cli.StringFlag{
			Name:        "storage-eos-share-folder",
			Value:       "/Shares",
			Usage:       "name of the share folder",
			EnvVars:     []string{"REVA_STORAGE_EOS_SHARE_FOLDER"},
			Destination: &cfg.Reva.Storages.EOS.ShareFolder,
		},
		&cli.StringFlag{
			Name:        "storage-eos-binary",
			Value:       "/usr/bin/eos",
			Usage:       "Location of the eos binary",
			EnvVars:     []string{"REVA_STORAGE_EOS_BINARY"},
			Destination: &cfg.Reva.Storages.EOS.EosBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-xrdcopy-binary",
			Value:       "/usr/bin/xrdcopy",
			Usage:       "Location of the xrdcopy binary",
			EnvVars:     []string{"REVA_STORAGE_EOS_XRDCOPY_BINARY"},
			Destination: &cfg.Reva.Storages.EOS.XrdcopyBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-master-url",
			Value:       "root://eos-mgm1.eoscluster.cern.ch:1094",
			Usage:       "URL of the Master EOS MGM",
			EnvVars:     []string{"REVA_STORAGE_EOS_MASTER_URL"},
			Destination: &cfg.Reva.Storages.EOS.MasterURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-slave-url",
			Value:       "root://eos-mgm1.eoscluster.cern.ch:1094",
			Usage:       "URL of the Slave EOS MGM",
			EnvVars:     []string{"REVA_STORAGE_EOS_SLAVE_URL"},
			Destination: &cfg.Reva.Storages.EOS.SlaveURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-cache-directory",
			Value:       os.TempDir(),
			Usage:       "Location on the local fs where to store reads",
			EnvVars:     []string{"REVA_STORAGE_EOS_CACHE_DIRECTORY"},
			Destination: &cfg.Reva.Storages.EOS.CacheDirectory,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-enable-logging",
			Value:       true,
			Usage:       "Enables logging of the commands executed",
			EnvVars:     []string{"REVA_STORAGE_EOS_ENABLE_LOGGING"},
			Destination: &cfg.Reva.Storages.EOS.EnableLogging,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-show-hidden-sysfiles",
			Usage:       "show internal EOS files like .sys.v# and .sys.a# files.",
			EnvVars:     []string{"REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES"},
			Destination: &cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-force-singleuser-mode",
			Usage:       "force connections to EOS to use SingleUsername",
			EnvVars:     []string{"REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE"},
			Destination: &cfg.Reva.Storages.EOS.ForceSingleUserMode,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-use-keytab",
			Usage:       "authenticate requests by using an EOS keytab",
			EnvVars:     []string{"REVA_STORAGE_EOS_USE_KEYTAB"},
			Destination: &cfg.Reva.Storages.EOS.UseKeytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-sec-protocol",
			Usage:       "the xrootd security protocol to use between the server and EOS",
			EnvVars:     []string{"REVA_STORAGE_EOS_SEC_PROTOCOL"},
			Destination: &cfg.Reva.Storages.EOS.SecProtocol,
		},
		&cli.StringFlag{
			Name:        "storage-eos-keytab",
			Usage:       "the location of the keytab to use to authenticate to EOS",
			EnvVars:     []string{"REVA_STORAGE_EOS_KEYTAB"},
			Destination: &cfg.Reva.Storages.EOS.Keytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-single-username",
			Usage:       "the username to use when SingleUserMode is enabled",
			EnvVars:     []string{"REVA_STORAGE_EOS_SINGLE_USERNAME"},
			Destination: &cfg.Reva.Storages.EOS.SingleUsername,
		},
		&cli.StringFlag{
			Name:        "storage-eos-user-layout",
			Value:       "{{substr 0 1 .Username}}/{{.Username}}",
			Usage:       `"layout of the users internal home dir path, in addition to {{.Username}}, and {{.Mail}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "E/Einstein" `,
			EnvVars:     []string{"REVA_STORAGE_EOS_USER_LAYOUT"},
			Destination: &cfg.Reva.Storages.EOS.UserLayout,
		},

		// local

		&cli.StringFlag{
			Name:        "storage-local-root",
			Value:       "/var/tmp/reva/root",
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"REVA_STORAGE_LOCAL_ROOT"},
			Destination: &cfg.Reva.Storages.Local.Root,
		},
		&cli.StringFlag{
			Name:        "storage-local-share-folder",
			Value:       "/MyShares",
			Usage:       "name of the share folder",
			EnvVars:     []string{"REVA_STORAGE_LOCAL_SHARE_FOLDER"},
			Destination: &cfg.Reva.Storages.Local.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-local-user-layout",
			Value:       "{{substr 0 1 .Username}}/{{.Username}}",
			Usage:       `"layout of the users internal home dir path, in addition to {{.Username}}, and {{.Mail}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "E/Einstein" `,
			EnvVars:     []string{"REVA_STORAGE_LOCAL_USER_LAYOUT"},
			Destination: &cfg.Reva.Storages.Local.UserLayout,
		},

		// owncloud

		&cli.StringFlag{
			Name:        "storage-owncloud-datadir",
			Value:       "/var/tmp/reva/data",
			Usage:       "the path to the owncloud data directory",
			EnvVars:     []string{"REVA_STORAGE_OWNCLOUD_DATADIR"},
			Destination: &cfg.Reva.Storages.OwnCloud.Datadirectory,
		},
		&cli.BoolFlag{
			Name:        "storage-owncloud-scan",
			Value:       true,
			Usage:       "scan files on startup to add fileids",
			EnvVars:     []string{"REVA_STORAGE_OWNCLOUD_SCAN"},
			Destination: &cfg.Reva.Storages.OwnCloud.Scan,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-redis",
			Value:       ":6379",
			Usage:       "the address of the redis server",
			EnvVars:     []string{"REVA_STORAGE_OWNCLOUD_REDIS_ADDR"},
			Destination: &cfg.Reva.Storages.OwnCloud.Redis,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-layout",
			Value:       "{{.Username}}",
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"REVA_STORAGE_OWNCLOUD_LAYOUT"},
			Destination: &cfg.Reva.Storages.OwnCloud.Layout,
		},
	}
}
