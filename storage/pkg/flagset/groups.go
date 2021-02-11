package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// GroupsWithConfig applies cfg to the root flagset
func GroupsWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9161",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.Groups.DebugAddr,
		},

		// Services

		// Groupprovider

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_NETWORK"},
			Destination: &cfg.Reva.Groups.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9160",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ADDR"},
			Destination: &cfg.Reva.Groups.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "endpoint",
			Value:       "localhost:9160",
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
			Value:       "ldap",
			Usage:       "group driver: 'json', 'ldap', or 'rest'",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DRIVER"},
			Destination: &cfg.Reva.Groups.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       "",
			Usage:       "Path to groups.json file",
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_JSON"},
			Destination: &cfg.Reva.Groups.JSON,
		},
    &cli.IntFlag{
			Name:        "group-members-cache-expiration",
			Value:       5,
			Usage:       "Time in minutes for redis cache expiration.",
			EnvVars:     []string{"STORAGE_GROUP_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.Groups.GroupMembersCacheExpiration,
		},

    // rest driver

    &cli.StringFlag{
			Name:        "rest-client-id",
			Value:       "",
			Usage:       "User/group rest driver Client ID",
			EnvVars:     []string{"STORAGE_REST_CLIENT_ID"},
			Destination: &cfg.Reva.UserGroupRest.ClientID,
		},
		&cli.StringFlag{
			Name:        "rest-client-secret",
			Value:       "",
			Usage:       "User/group rest driver Client Secret",
			EnvVars:     []string{"STORAGE_REST_CLIENT_SECRET"},
			Destination: &cfg.Reva.UserGroupRest.ClientSecret,
		},
		&cli.StringFlag{
			Name:        "rest-redis-address",
			Value:       "localhost:6379",
			Usage:       "Address for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_ADDRESS"},
			Destination: &cfg.Reva.UserGroupRest.RedisAddress,
		},
		&cli.StringFlag{
			Name:        "rest-redis-username",
			Value:       "",
			Usage:       "Username for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_USERNAME"},
			Destination: &cfg.Reva.UserGroupRest.RedisUsername,
		},
		&cli.StringFlag{
			Name:        "rest-redis-password",
			Value:       "",
			Usage:       "Password for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_PASSWORD"},
			Destination: &cfg.Reva.UserGroupRest.RedisPassword,
		},
		&cli.StringFlag{
			Name:        "rest-id-provider",
			Value:       "",
			Usage:       "The OIDC Provider",
			EnvVars:     []string{"STORAGE_REST_ID_PROVIDER"},
			Destination: &cfg.Reva.UserGroupRest.IDProvider,
		},
		&cli.StringFlag{
			Name:        "rest-api-base-url",
			Value:       "",
			Usage:       "Base API Endpoint",
			EnvVars:     []string{"STORAGE_REST_API_BASE_URL"},
			Destination: &cfg.Reva.UserGroupRest.APIBaseURL,
		},
		&cli.StringFlag{
			Name:        "rest-oidc-token-endpoint",
			Value:       "",
			Usage:       "Endpoint to generate token to access the API",
			EnvVars:     []string{"STORAGE_REST_OIDC_TOKEN_ENDPOINT"},
			Destination: &cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
		},
		&cli.StringFlag{
			Name:        "rest-target-api",
			Value:       "",
			Usage:       "The target application",
			EnvVars:     []string{"STORAGE_REST_TARGET_API"},
			Destination: &cfg.Reva.UserGroupRest.TargetAPI,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)

	return flags
}
