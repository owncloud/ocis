package userdrivers

import (
	"os"

	"github.com/owncloud/ocis/ocis-pkg/flags"

	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverEOSWithConfig applies cfg to the root flagset
func DriverEOSWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		&cli.StringFlag{
			Name:        "storage-eos-namespace",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.Root, "/eos/dockertest/reva"),
			Usage:       "Namespace for metadata operations",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.Root,
		},
		&cli.StringFlag{
			Name: "storage-eos-shadow-namespace",
			// Defaults to path.Join(c.Namespace, ".shadow")
			Usage:       "Shadow namespace where share references are stored",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHADOW_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.ShadowNamespace,
		},
		&cli.StringFlag{
			Name: "storage-eos-uploads-namespace",
			// Defaults to path.Join(c.Namespace, ".uploads")
			Usage:       "Uploads namespace",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_UPLOADS_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.UploadsNamespace,
		},
		&cli.StringFlag{
			Name:        "storage-eos-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.ShareFolder, "/Shares"),
			Usage:       "name of the share folder",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.EOS.ShareFolder,
		},
		&cli.StringFlag{
			Name:        "storage-eos-binary",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.EosBinary, "/usr/bin/eos"),
			Usage:       "Location of the eos binary",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_BINARY"},
			Destination: &cfg.Reva.UserStorage.EOS.EosBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-xrdcopy-binary",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.XrdcopyBinary, "/usr/bin/xrdcopy"),
			Usage:       "Location of the xrdcopy binary",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_XRDCOPY_BINARY"},
			Destination: &cfg.Reva.UserStorage.EOS.XrdcopyBinary,
		},
		&cli.StringFlag{
			Name:        "storage-eos-master-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.MasterURL, "root://eos-mgm1.eoscluster.cern.ch:1094"),
			Usage:       "URL of the Master EOS MGM",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_MASTER_URL"},
			Destination: &cfg.Reva.UserStorage.EOS.MasterURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-slave-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.SlaveURL, "root://eos-mgm1.eoscluster.cern.ch:1094"),
			Usage:       "URL of the Slave EOS MGM",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SLAVE_URL"},
			Destination: &cfg.Reva.UserStorage.EOS.SlaveURL,
		},
		&cli.StringFlag{
			Name:        "storage-eos-cache-directory",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.CacheDirectory, os.TempDir()),
			Usage:       "Location on the local fs where to store reads",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_CACHE_DIRECTORY"},
			Destination: &cfg.Reva.UserStorage.EOS.CacheDirectory,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-enable-logging",
			Usage:       "Enables logging of the commands executed",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_ENABLE_LOGGING"},
			Destination: &cfg.Reva.UserStorage.EOS.EnableLogging,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-show-hidden-sysfiles",
			Usage:       "show internal EOS files like .sys.v# and .sys.a# files.",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHOW_HIDDEN_SYSFILES"},
			Destination: &cfg.Reva.UserStorage.EOS.ShowHiddenSysFiles,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-force-singleuser-mode",
			Usage:       "force connections to EOS to use SingleUsername",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_FORCE_SINGLEUSER_MODE"},
			Destination: &cfg.Reva.UserStorage.EOS.ForceSingleUserMode,
		},
		&cli.BoolFlag{
			Name:        "storage-eos-use-keytab",
			Usage:       "authenticate requests by using an EOS keytab",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_USE_KEYTAB"},
			Destination: &cfg.Reva.UserStorage.EOS.UseKeytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-sec-protocol",
			Usage:       "the xrootd security protocol to use between the server and EOS",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SEC_PROTOCOL"},
			Destination: &cfg.Reva.UserStorage.EOS.SecProtocol,
		},
		&cli.StringFlag{
			Name:        "storage-eos-keytab",
			Usage:       "the location of the keytab to use to authenticate to EOS",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_KEYTAB"},
			Destination: &cfg.Reva.UserStorage.EOS.Keytab,
		},
		&cli.StringFlag{
			Name:        "storage-eos-single-username",
			Usage:       "the username to use when SingleUserMode is enabled",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SINGLE_USERNAME"},
			Destination: &cfg.Reva.UserStorage.EOS.SingleUsername,
		},
		&cli.StringFlag{
			Name:        "storage-eos-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.UserLayout, "{{substr 0 1 .Username}}/{{.Username}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.EOS.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-eos-reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.EOS.GatewaySVC, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.UserStorage.EOS.GatewaySVC,
		},
	}
}
