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

	Backend      string        `yaml:"backend" env:"INVITATIONS_BACKEND" desc:"The backend used to provision invited guests. Supported values are 'keycloak' (default; creates the guest in the configured Keycloak realm and sends a credential-setup email) and 'ldap' (creates the guest directly in the oCIS identity backend, so it is immediately usable as a share recipient; requires OCIS_LDAP_SERVER_WRITE_ENABLED)." introductionVersion:"NEXT"`
	Keycloak     Keycloak      `yaml:"keycloak"`
	LDAP         LDAP          `yaml:"ldap"`
	TokenManager *TokenManager `yaml:"token_manager"`

	Context context.Context `yaml:"-"`
}

// Keycloak configuration
type Keycloak struct {
	BasePath           string `yaml:"base_path" env:"OCIS_KEYCLOAK_BASE_PATH;INVITATIONS_KEYCLOAK_BASE_PATH" desc:"The URL to access keycloak." introductionVersion:"pre5.0"`
	ClientID           string `yaml:"client_id" env:"OCIS_KEYCLOAK_CLIENT_ID;INVITATIONS_KEYCLOAK_CLIENT_ID" desc:"The client ID to authenticate with keycloak." introductionVersion:"pre5.0"`
	ClientSecret       string `yaml:"client_secret" env:"OCIS_KEYCLOAK_CLIENT_SECRET;INVITATIONS_KEYCLOAK_CLIENT_SECRET" desc:"The client secret to use in authentication." introductionVersion:"pre5.0"`
	ClientRealm        string `yaml:"client_realm" env:"OCIS_KEYCLOAK_CLIENT_REALM;INVITATIONS_KEYCLOAK_CLIENT_REALM" desc:"The realm the client is defined in." introductionVersion:"pre5.0"`
	UserRealm          string `yaml:"user_realm" env:"OCIS_KEYCLOAK_USER_REALM;INVITATIONS_KEYCLOAK_USER_REALM" desc:"The realm users are defined." introductionVersion:"pre5.0"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" env:"OCIS_KEYCLOAK_INSECURE_SKIP_VERIFY;INVITATIONS_KEYCLOAK_INSECURE_SKIP_VERIFY" desc:"Disable TLS certificate validation for Keycloak connections. Do not set this in production environments." introductionVersion:"pre5.0"`
}

// LDAP configures the 'ldap' invitation backend, which provisions invited guests
// directly into the oCIS identity backend (the directory oCIS reads for
// share-recipient resolution). The connection settings reuse the shared
// OCIS_LDAP_* configuration of the graph service so both resolve the same entries.
type LDAP struct {
	URI          string `yaml:"uri" env:"OCIS_LDAP_URI;INVITATIONS_LDAP_URI" desc:"URI of the LDAP server to provision guests into. Only used when INVITATIONS_BACKEND is 'ldap'." introductionVersion:"NEXT"`
	CACert       string `yaml:"cacert" env:"OCIS_LDAP_CACERT;INVITATIONS_LDAP_CACERT" desc:"Path to the root CA certificate (PEM) used to validate the LDAP server's TLS certificate." introductionVersion:"NEXT"`
	Insecure     bool   `yaml:"insecure" env:"OCIS_LDAP_INSECURE;INVITATIONS_LDAP_INSECURE" desc:"Disable TLS certificate validation for the LDAP connection. Do not set this in production environments." introductionVersion:"NEXT"`
	BindDN       string `yaml:"bind_dn" env:"OCIS_LDAP_BIND_DN;INVITATIONS_LDAP_BIND_DN" desc:"LDAP DN used for simple bind authentication when provisioning guests." introductionVersion:"NEXT"`
	BindPassword string `yaml:"bind_password" env:"OCIS_LDAP_BIND_PASSWORD;INVITATIONS_LDAP_BIND_PASSWORD" desc:"Password for the 'bind_dn'." introductionVersion:"NEXT"`
	UserBaseDN   string `yaml:"user_base_dn" env:"OCIS_LDAP_USER_BASE_DN;INVITATIONS_LDAP_USER_BASE_DN" desc:"Search base DN under which invited guests are created." introductionVersion:"NEXT"`
	WriteEnabled bool   `yaml:"write_enabled" env:"OCIS_LDAP_SERVER_WRITE_ENABLED;INVITATIONS_LDAP_WRITE_ENABLED" desc:"Allow creating users in LDAP. Must be 'true' for the 'ldap' backend to provision guests." introductionVersion:"NEXT"`
}
