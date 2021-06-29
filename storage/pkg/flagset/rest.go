package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// RestWithConfig applies REST user/group provider cfg to the flagset
func RestWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "rest-client-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.ClientID, ""),
			Usage:       "User/group rest driver Client ID",
			EnvVars:     []string{"STORAGE_REST_CLIENT_ID"},
			Destination: &cfg.Reva.UserGroupRest.ClientID,
		},
		&cli.StringFlag{
			Name:        "rest-client-secret",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.ClientSecret, ""),
			Usage:       "User/group rest driver Client Secret",
			EnvVars:     []string{"STORAGE_REST_CLIENT_SECRET"},
			Destination: &cfg.Reva.UserGroupRest.ClientSecret,
		},
		&cli.StringFlag{
			Name:        "rest-redis-address",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.RedisAddress, "localhost:6379"),
			Usage:       "Address for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_ADDRESS"},
			Destination: &cfg.Reva.UserGroupRest.RedisAddress,
		},
		&cli.StringFlag{
			Name:        "rest-redis-username",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.RedisUsername, ""),
			Usage:       "Username for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_USERNAME"},
			Destination: &cfg.Reva.UserGroupRest.RedisUsername,
		},
		&cli.StringFlag{
			Name:        "rest-redis-password",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.RedisPassword, ""),
			Usage:       "Password for redis server",
			EnvVars:     []string{"STORAGE_REST_REDIS_PASSWORD"},
			Destination: &cfg.Reva.UserGroupRest.RedisPassword,
		},
		&cli.StringFlag{
			Name:        "rest-id-provider",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.IDProvider, ""),
			Usage:       "The OIDC Provider",
			EnvVars:     []string{"STORAGE_REST_ID_PROVIDER"},
			Destination: &cfg.Reva.UserGroupRest.IDProvider,
		},
		&cli.StringFlag{
			Name:        "rest-api-base-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.APIBaseURL, ""),
			Usage:       "Base API Endpoint",
			EnvVars:     []string{"STORAGE_REST_API_BASE_URL"},
			Destination: &cfg.Reva.UserGroupRest.APIBaseURL,
		},
		&cli.StringFlag{
			Name:        "rest-oidc-token-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.OIDCTokenEndpoint, ""),
			Usage:       "Endpoint to generate token to access the API",
			EnvVars:     []string{"STORAGE_REST_OIDC_TOKEN_ENDPOINT"},
			Destination: &cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
		},
		&cli.StringFlag{
			Name:        "rest-target-api",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserGroupRest.TargetAPI, ""),
			Usage:       "The target application",
			EnvVars:     []string{"STORAGE_REST_TARGET_API"},
			Destination: &cfg.Reva.UserGroupRest.TargetAPI,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}
