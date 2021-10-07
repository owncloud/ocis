package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// UsersWithConfig applies cfg to the root flagset
func UsersWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.DebugAddr, "0.0.0.0:9145"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_USERPROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.Users.DebugAddr,
		},

		// Services

		// Gateway

		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// Userprovider

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_USERPROVIDER_NETWORK"},
			Destination: &cfg.Reva.Users.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.GRPCAddr, "0.0.0.0:9144"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ADDR"},
			Destination: &cfg.Reva.Users.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.Endpoint, "localhost:9144"),
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("userprovider"), // TODO preferences
			Usage:   "--service userprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_USERPROVIDER_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.Driver, "ldap"),
			Usage:       "user driver: 'demo', 'json', 'ldap', 'owncloudsql' or 'rest'",
			EnvVars:     []string{"STORAGE_USERPROVIDER_DRIVER"},
			Destination: &cfg.Reva.Users.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.JSON, ""),
			Usage:       "Path to users.json file",
			EnvVars:     []string{"STORAGE_USERPROVIDER_JSON"},
			Destination: &cfg.Reva.Users.JSON,
		},
		&cli.IntFlag{
			Name:        "user-groups-cache-expiration",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Users.UserGroupsCacheExpiration, 5),
			Usage:       "Time in minutes for redis cache expiration.",
			EnvVars:     []string{"STORAGE_USER_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.Users.UserGroupsCacheExpiration,
		},

		// user owncloudsql

		&cli.StringFlag{
			Name:        "owncloudsql-dbhost",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserOwnCloudSQL.DBHost, "mysql"),
			Usage:       "hostname of the mysql db",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBHOST"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBHost,
		},
		&cli.IntFlag{
			Name:        "owncloudsql-dbport",
			Value:       flags.OverrideDefaultInt(cfg.Reva.UserOwnCloudSQL.DBPort, 3306),
			Usage:       "port of the mysql db",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPORT"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBPort,
		},
		&cli.StringFlag{
			Name:        "owncloudsql-dbname",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserOwnCloudSQL.DBName, "owncloud"),
			Usage:       "database name of the owncloud db",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBNAME"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBName,
		},
		&cli.StringFlag{
			Name:        "owncloudsql-dbuser",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserOwnCloudSQL.DBUsername, "owncloud"),
			Usage:       "user name to use when connecting to the mysql owncloud db",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBUSER"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBUsername,
		},
		&cli.StringFlag{
			Name:        "owncloudsql-dbpass",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserOwnCloudSQL.DBPassword, "secret"),
			Usage:       "password to use when connecting to the mysql owncloud db",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPASS"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBPassword,
		},
		&cli.StringFlag{
			Name:        "owncloudsql-idp",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserOwnCloudSQL.Idp, "https://localhost:9200"),
			Usage:       "Identity provider to use for users",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_IDP", "OCIS_URL"},
			Destination: &cfg.Reva.UserOwnCloudSQL.Idp,
		},
		&cli.Int64Flag{
			Name:        "owncloudsql-nobody",
			Value:       flags.OverrideDefaultInt64(cfg.Reva.UserOwnCloudSQL.Nobody, 99),
			Usage:       "fallback user id to use when user has no id",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_NOBODY"},
			Destination: &cfg.Reva.UserOwnCloudSQL.Nobody,
		},
		&cli.BoolFlag{
			Name:        "owncloudsql-join-username",
			Value:       flags.OverrideDefaultBool(cfg.Reva.UserOwnCloudSQL.JoinUsername, false),
			Usage:       "join the username from the oc_preferences table",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_USERNAME"},
			Destination: &cfg.Reva.UserOwnCloudSQL.JoinUsername,
		},
		&cli.BoolFlag{
			Name:        "owncloudsql-join-ownclouduuid",
			Value:       flags.OverrideDefaultBool(cfg.Reva.UserOwnCloudSQL.JoinOwnCloudUUID, false),
			Usage:       "join the ownclouduuid from the oc_preferences table",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_OWNCLOUDUUID"},
			Destination: &cfg.Reva.UserOwnCloudSQL.JoinOwnCloudUUID,
		},
		&cli.BoolFlag{
			Name:        "owncloudsql-enable-medial-search",
			Value:       flags.OverrideDefaultBool(cfg.Reva.UserOwnCloudSQL.EnableMedialSearch, false),
			Usage:       "enable medial search when finding users",
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH"},
			Destination: &cfg.Reva.UserOwnCloudSQL.EnableMedialSearch,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)
	flags = append(flags, RestWithConfig(cfg)...)

	return flags
}
