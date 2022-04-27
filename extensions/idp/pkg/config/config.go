package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	HTTP HTTP `yaml:"http,omitempty"`

	Asset Asset    `yaml:"asset,omitempty"`
	IDP   Settings `yaml:"idp,omitempty"`
	Ldap  Ldap     `yaml:"ldap,omitempty"`

	Context context.Context `yaml:"-"`
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI       string `yaml:"uri,omitempty" env:"LDAP_URI;IDP_LDAP_URI"`
	TLSCACert string `yaml:"cacert,omitempty" env:"LDAP_CACERT;IDP_LDAP_TLS_CACERT"`

	BindDN       string `yaml:"bind_dn,omitempty" env:"LDAP_BIND_DN;IDP_LDAP_BIND_DN"`
	BindPassword string `yaml:"bind_password,omitempty" env:"LDAP_BIND_PASSWORD;IDP_LDAP_BIND_PASSWORD"`

	BaseDN string `yaml:"base_dn,omitempty" env:"LDAP_USER_BASE_DN,IDP_LDAP_BASE_DN"`
	Scope  string `yaml:"scope,omitempty" env:"LDAP_USER_SCOPE;IDP_LDAP_SCOPE"`

	LoginAttribute    string `yaml:"login_attribute,omitempty" env:"IDP_LDAP_LOGIN_ATTRIBUTE"`
	EmailAttribute    string `yaml:"email_attribute,omitempty" env:"LDAP_USER_SCHEMA_MAIL;IDP_LDAP_EMAIL_ATTRIBUTE"`
	NameAttribute     string `yaml:"name_attribute,omitempty" env:"LDAP_USER_SCHEMA_USERNAME;IDP_LDAP_NAME_ATTRIBUTE"`
	UUIDAttribute     string `yaml:"uuid_attribute,omitempty" env:"LDAP_USER_SCHEMA_ID;IDP_LDAP_UUID_ATTRIBUTE"`
	UUIDAttributeType string `yaml:"uuid_attribute_type,omitempty" env:"IDP_LDAP_UUID_ATTRIBUTE_TYPE"`

	Filter      string `yaml:"filter,omitempty" env:"LDAP_USER_FILTER;IDP_LDAP_FILTER"`
	ObjectClass string `yaml:"objectclass,omitempty" env:"LDAP_USER_OBJECTCLASS;IDP_LDAP_OBJECTCLASS"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"asset,omitempty" env:"IDP_ASSET_PATH"`
}

type Settings struct {
	// don't change the order of elements in this struct
	// it needs to match github.com/libregraph/lico/bootstrap.Settings

	Iss string `yaml:"iss" env:"OCIS_URL;IDP_ISS"`

	IdentityManager string `yaml:"identity_manager" env:"IDP_IDENTITY_MANAGER"`

	URIBasePath string `yaml:"uri_base_path" env:"IDP_URI_BASE_PATH"`

	SignInURI    string `yaml:"sign_in_uri" env:"IDP_SIGN_IN_URI"`
	SignedOutURI string `yaml:"signed_out_uri" env:"IDP_SIGN_OUT_URI"`

	AuthorizationEndpointURI string `yaml:"authorization_endpoint_uri" env:"IDP_ENDPOINT_URI"`
	EndsessionEndpointURI    string `yaml:"end_session_endpoint_uri" env:"IDP_ENDSESSION_ENDPOINT_URI"`

	Insecure bool `yaml:"insecure" env:"IDP_INSECURE"`

	TrustedProxy []string `yaml:"trusted_proxy"` //TODO: how to configure this via env?

	AllowScope                     []string `yaml:"allow_scope"` // TODO: is this even needed?
	AllowClientGuests              bool     `yaml:"allow_client_guests" env:"IDP_ALLOW_CLIENT_GUESTS"`
	AllowDynamicClientRegistration bool     `yaml:"allow_dynamic_client_registration" env:"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"`

	EncryptionSecretFile string `yaml:"encrypt_secret_file" env:"IDP_ENCRYPTION_SECRET"`

	Listen string

	IdentifierClientDisabled          bool   `yaml:"identifier_client_disabled" env:"IDP_DISABLE_IDENTIFIER_WEBAPP"`
	IdentifierClientPath              string `yaml:"identifier_client_path" env:"IDP_IDENTIFIER_CLIENT_PATH"`
	IdentifierRegistrationConf        string `yaml:"identifier_registration_conf" env:"IDP_IDENTIFIER_REGISTRATION_CONF"`
	IdentifierScopesConf              string `yaml:"identifier_scopes_conf" env:"IDP_IDENTIFIER_SCOPES_CONF"`
	IdentifierDefaultBannerLogo       string
	IdentifierDefaultSignInPageText   string
	IdentifierDefaultUsernameHintText string
	IdentifierUILocales               []string

	SigningKid             string   `yaml:"signing_kid" env:"IDP_SIGNING_KID"`
	SigningMethod          string   `yaml:"signing_method" env:"IDP_SIGNING_METHOD"`
	SigningPrivateKeyFiles []string `yaml:"signing_private_key_files"` // TODO: is this even needed?
	ValidationKeysPath     string   `yaml:"validation_keys_path" env:"IDP_VALIDATION_KEYS_PATH"`

	CookieBackendURI string
	CookieNames      []string

	AccessTokenDurationSeconds        uint64 `yaml:"access_token_duration_seconds" env:"IDP_ACCESS_TOKEN_EXPIRATION"`
	IDTokenDurationSeconds            uint64 `yaml:"id_token_duration_seconds" env:"IDP_ID_TOKEN_EXPIRATION"`
	RefreshTokenDurationSeconds       uint64 `yaml:"refresh_token_duration_seconds" env:"IDP_REFRESH_TOKEN_EXPIRATION"`
	DyamicClientSecretDurationSeconds uint64 `yaml:"dynamic_client_secret_duration_seconds" env:""`
}
