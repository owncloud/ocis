package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// UsersWithConfig applies cfg to the root flagset
func UsersWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.DebugAddr, "0.0.0.0:9145"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Users.DebugAddr,
		},

		// Services

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
			Usage:       "user driver: 'demo', 'json', 'ldap', or 'rest'",
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
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)
	flags = append(flags, RestWithConfig(cfg)...)

	return flags
}
