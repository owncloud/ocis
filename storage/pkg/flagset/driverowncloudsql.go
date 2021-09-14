package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverOwnCloudSQLWithConfig applies cfg to the root flagset
func DriverOwnCloudSQLWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-owncloudsql-datadir",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.Root, "/var/tmp/ocis/storage/owncloud"),
			Usage:       "the path to the owncloudsql data directory",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DATADIR"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.Root,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-uploadinfo-dir",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.UploadInfoDir, "/var/tmp/ocis/storage/uploadinfo"),
			Usage:       "the path to the tus upload info directory",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.UploadInfoDir,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.ShareFolder, "/Shares"),
			Usage:       "name of the shares folder",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_SHARE_FOLDER"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.ShareFolder,
		},
		&cli.BoolFlag{
			Name:        "storage-owncloudsql-enable-home",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Storages.OwnCloudSQL.EnableHome, false),
			Usage:       "enable the creation of home storages",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_ENABLE_HOME"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.EnableHome,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.UserLayout, "{{.Username}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_LAYOUT"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-dbusername",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.DBUsername, "owncloud"),
			Usage:       `"username for accessing the database" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DBUSERNAME"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.DBUsername,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-dbpassword",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.DBPassword, "owncloud"),
			Usage:       `"password for accessing the database" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DBPASSWORD"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.DBPassword,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-dbhost",
			Value:       cfg.Reva.Storages.OwnCloudSQL.DBHost,
			Usage:       `"the database hostname or IP address" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DBHOST"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.DBHost,
		},
		&cli.IntFlag{
			Name:        "storage-owncloudsql-dbport",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Storages.OwnCloudSQL.DBPort, 3306),
			Usage:       `"port the database listens on" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DBPORT"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.DBPort,
		},
		&cli.StringFlag{
			Name:        "storage-owncloudsql-dbname",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloudSQL.DBName, "owncloud"),
			Usage:       `"the database name" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUDSQL_DBNAME"},
			Destination: &cfg.Reva.Storages.OwnCloudSQL.DBName,
		},
	}
}
