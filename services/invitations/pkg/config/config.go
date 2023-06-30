package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Keycloak     Keycloak      `yaml:"keycloak"`
	TokenManager *TokenManager `yaml:"token_manager"`

	Context context.Context `yaml:"-"`
}

// Keycloak configuration
type Keycloak struct {
	BasePath           string `yaml:"base_path" env:"OCIS_KEYCLOAK_BASE_PATH;INVITATIONS_KEYCLOAK_BASE_PATH" desc:"The URL to access keycloak."`
	ClientID           string `yaml:"client_id" env:"OCIS_KEYCLOAK_CLIENT_ID;INVITATIONS_KEYCLOAK_CLIENT_ID" desc:"The client ID to authenticate with keycloak."`
	ClientSecret       string `yaml:"client_secret" env:"OCIS_KEYCLOAK_CLIENT_SECRET;INVITATIONS_KEYCLOAK_CLIENT_SECRET" desc:"The client secret to use in authentication."`
	ClientRealm        string `yaml:"client_realm" env:"OCIS_KEYCLOAK_CLIENT_REALM;INVITATIONS_KEYCLOAK_CLIENT_REALM" desc:"The realm the client is defined in."`
	UserRealm          string `yaml:"user_realm" env:"OCIS_KEYCLOAK_USER_REALM;INVITATIONS_KEYCLOAK_USER_REALM" desc:"The realm users are defined."`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" env:"OCIS_KEYCLOAK_INSECURE_SKIP_VERIFY;INVITATIONS_KEYCLOAK_INSECURE_SKIP_VERIFY" desc:"Disable TLS certificate validation for Keycloak connections. Do not set this in production environments."`
}
