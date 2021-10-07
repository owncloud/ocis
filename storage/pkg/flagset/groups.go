package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// GroupsWithConfig applies cfg to the root flagset
func GroupsWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.DebugAddr, "0.0.0.0:9161"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.Groups.DebugAddr,
		},

		// Services

		// Gateway

		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// Groupprovider

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_NETWORK"},
			Destination: &cfg.Reva.Groups.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.GRPCAddr, "0.0.0.0:9160"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ADDR"},
			Destination: &cfg.Reva.Groups.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.Endpoint, "localhost:9160"),
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Groups.Endpoint,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("groupprovider"), // TODO preferences
			Usage:   "--service groupprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_GROUPPROVIDER_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.Driver, "ldap"),
			Usage:       "group driver: 'json', 'ldap', or 'rest'",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DRIVER"},
			Destination: &cfg.Reva.Groups.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       flags.OverrideDefaultString(cfg.Reva.Groups.JSON, ""),
			Usage:       "Path to groups.json file",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_JSON"},
			Destination: &cfg.Reva.Groups.JSON,
		},
		&cli.IntFlag{
			Name:        "group-members-cache-expiration",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Groups.GroupMembersCacheExpiration, 5),
			Usage:       "Time in minutes for redis cache expiration.",
			EnvVars:     []string{"STORAGE_GROUP_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.Groups.GroupMembersCacheExpiration,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)
	flags = append(flags, RestWithConfig(cfg)...)

	return flags
}
