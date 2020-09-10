package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// UsersWithConfig applies cfg to the root flagset
func UsersWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9145",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Users.DebugAddr,
		},

		// Services

		// Users

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_USERS_NETWORK"},
			Destination: &cfg.Reva.Users.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_USERS_PROTOCOL"},
			Destination: &cfg.Reva.Users.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9144",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_USERS_ADDR"},
			Destination: &cfg.Reva.Users.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("userprovider"), // TODO preferences
			Usage:   "--service userprovider [--service otherservice]",
			EnvVars: []string{"REVA_USERS_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "ldap",
			Usage:       "user driver: 'demo', 'json', 'ldap', or 'rest'",
			EnvVars:     []string{"REVA_USERS_DRIVER"},
			Destination: &cfg.Reva.Users.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVars:     []string{"REVA_USERS_JSON"},
			Destination: &cfg.Reva.Users.JSON,
		},

		// rest driver

		&cli.StringFlag{
			Name:        "rest-client-id",
			Value:       "",
			Usage:       "User rest driver Client ID",
			EnvVars:     []string{"REVA_REST_CLIENT_ID"},
			Destination: &cfg.Reva.UserRest.ClientID,
		},
		&cli.StringFlag{
			Name:        "rest-client-secret",
			Value:       "",
			Usage:       "User rest driver Client Secret",
			EnvVars:     []string{"REVA_REST_CLIENT_SECRET"},
			Destination: &cfg.Reva.UserRest.ClientSecret,
		},
		&cli.StringFlag{
			Name:        "rest-redis-address",
			Value:       "localhost:6379",
			Usage:       "Address for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_ADDRESS"},
			Destination: &cfg.Reva.UserRest.RedisAddress,
		},
		&cli.StringFlag{
			Name:        "rest-redis-username",
			Value:       "",
			Usage:       "Username for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_USERNAME"},
			Destination: &cfg.Reva.UserRest.RedisUsername,
		},
		&cli.StringFlag{
			Name:        "rest-redis-password",
			Value:       "",
			Usage:       "Password for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_PASSWORD"},
			Destination: &cfg.Reva.UserRest.RedisPassword,
		},
		&cli.IntFlag{
			Name:        "rest-user-groups-cache-expiration",
			Value:       5,
			Usage:       "Time in minutes for redis cache expiration.",
			EnvVars:     []string{"REVA_REST_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.UserRest.UserGroupsCacheExpiration,
		},
		&cli.StringFlag{
			Name:        "rest-id-provider",
			Value:       "",
			Usage:       "The OIDC Provider",
			EnvVars:     []string{"REVA_REST_ID_PROVIDER"},
			Destination: &cfg.Reva.UserRest.IDProvider,
		},
		&cli.StringFlag{
			Name:        "rest-api-base-url",
			Value:       "",
			Usage:       "Base API Endpoint",
			EnvVars:     []string{"REVA_REST_API_BASE_URL"},
			Destination: &cfg.Reva.UserRest.APIBaseURL,
		},
		&cli.StringFlag{
			Name:        "rest-oidc-token-endpoint",
			Value:       "",
			Usage:       "Endpoint to generate token to access the API",
			EnvVars:     []string{"REVA_REST_OIDC_TOKEN_ENDPOINT"},
			Destination: &cfg.Reva.UserRest.OIDCTokenEndpoint,
		},
		&cli.StringFlag{
			Name:        "rest-target-api",
			Value:       "",
			Usage:       "The target application",
			EnvVars:     []string{"REVA_REST_TARGET_API"},
			Destination: &cfg.Reva.UserRest.TargetAPI,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)

	return flags
}
