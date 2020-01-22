package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis/pkg/config"
)

// LoginWithConfig applies cfg to the root flagset
func LoginWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "client-id",
			Value:       "cli",
			Usage:       "oidc cli client id",
			EnvVar:      "OIDC_CLI_CLIENT_ID",
			Destination: &cfg.OIDC.ClientID,
		},
		&cli.StringFlag{
			Name:        "secret",
			Value:       "foobar",
			Usage:       "oidc cli client secret",
			EnvVar:      "OIDC_CLI_CLIENT_SECRET",
			Destination: &cfg.OIDC.ClientSecret,
		},
		&cli.StringFlag{
			Name:        "callback-addr",
			Value:       "localhost:18080",
			Usage:       "oidc cli client callback addr",
			EnvVar:      "OIDC_CLI_CALLBACK_URL",
			Destination: &cfg.OIDC.CallbackAddr,
		},
		&cli.StringFlag{
			Name:        "callback-url",
			Value:       "http://localhost:18080/callback",
			Usage:       "oidc cli client callback url",
			EnvVar:      "OIDC_CLI_CALLBACK_URL",
			Destination: &cfg.OIDC.CallbackURL,
		},
		&cli.StringFlag{
			Name:        "auth-endpoint",
			Value:       "http://localhost:9140/oauth2/auth",
			Usage:       "oidc auth endpoint",
			EnvVar:      "OIDC_AUTH_ENDPOINT",
			Destination: &cfg.OIDC.AuthEndpoint,
		},
		&cli.StringFlag{
			Name:        "token-endpoint",
			Value:       "http://localhost:9140/oauth2/token",
			Usage:       "oidc token endpoint",
			EnvVar:      "OIDC_TOKEN_ENDPOINT",
			Destination: &cfg.OIDC.TokenEndpoint,
		},
	}
}
