package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Asset Asset    `ocisConfig:"asset"`
	IDP   Settings `ocisConfig:"idp"`
	Ldap  Ldap     `ocisConfig:"ldap"`

	Context context.Context
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI string `ocisConfig:"uri" env:"IDP_LDAP_URI"`

	BindDN       string `ocisConfig:"bind_dn" env:"IDP_LDAP_BIND_DN"`
	BindPassword string `ocisConfig:"bind_password" env:"IDP_LDAP_BIND_PASSWORD"`

	BaseDN string `ocisConfig:"base_dn" env:"IDP_LDAP_BASE_DN"`
	Scope  string `ocisConfig:"scope" env:"IDP_LDAP_SCOPE"`

	LoginAttribute    string `ocisConfig:"login_attribute" env:"IDP_LDAP_LOGIN_ATTRIBUTE"`
	EmailAttribute    string `ocisConfig:"email_attribute" env:"IDP_LDAP_EMAIL_ATTRIBUTE"`
	NameAttribute     string `ocisConfig:"name_attribute" env:"IDP_LDAP_NAME_ATTRIBUTE"`
	UUIDAttribute     string `ocisConfig:"uuid_attribute" env:"IDP_LDAP_UUID_ATTRIBUTE"`
	UUIDAttributeType string `ocisConfig:"uuid_attribute_type" env:"IDP_LDAP_UUID_ATTRIBUTE_TYPE"`

	Filter string `ocisConfig:"filter" env:"IDP_LDAP_FILTER"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"asset" env:"IDP_ASSET_PATH"`
}

type Settings struct {
	// don't change the order of elements in this struct
	// it needs to match github.com/libregraph/lico/bootstrap.Settings

	Iss string `ocisConfig:"iss" env:"OCIS_URL;IDP_ISS"`

	IdentityManager string `ocisConfig:"identity_manager" env:"IDP_IDENTITY_MANAGER"`

	URIBasePath string `ocisConfig:"uri_base_path" env:"IDP_URI_BASE_PATH"`

	SignInURI    string `ocisConfig:"sign_in_uri" env:"IDP_SIGN_IN_URI"`
	SignedOutURI string `ocisConfig:"signed_out_uri" env:"IDP_SIGN_OUT_URI"`

	AuthorizationEndpointURI string `ocisConfig:"authorization_endpoint_uri" env:"IDP_ENDPOINT_URI"`
	EndsessionEndpointURI    string `ocisConfig:"end_session_endpoint_uri" env:"IDP_ENDSESSION_ENDPOINT_URI"`

	Insecure bool `ocisConfig:"insecure" env:"IDP_INSECURE"`

	TrustedProxy []string `ocisConfig:"trusted_proxy"` //TODO: how to configure this via env?

	AllowScope                     []string `ocisConfig:"allow_scope"` // TODO: is this even needed?
	AllowClientGuests              bool     `ocisConfig:"allow_client_guests" env:"IDP_ALLOW_CLIENT_GUESTS"`
	AllowDynamicClientRegistration bool     `ocisConfig:"allow_dynamic_client_registration" env:"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"`

	EncryptionSecretFile string `ocisConfig:"encrypt_secret_file" env:"IDP_ENCRYPTION_SECRET"`

	Listen string

	IdentifierClientDisabled          bool   `ocisConfig:"identifier_client_disabled" env:"IDP_DISABLE_IDENTIFIER_WEBAPP"`
	IdentifierClientPath              string `ocisConfig:"identifier_client_path" env:"IDP_IDENTIFIER_CLIENT_PATH"`
	IdentifierRegistrationConf        string `ocisConfig:"identifier_registration_conf" env:"IDP_IDENTIFIER_REGISTRATION_CONF"`
	IdentifierScopesConf              string `ocisConfig:"identifier_scopes_conf" env:"IDP_IDENTIFIER_SCOPES_CONF"`
	IdentifierDefaultBannerLogo       string
	IdentifierDefaultSignInPageText   string
	IdentifierDefaultUsernameHintText string

	SigningKid             string   `ocisConfig:"signing_kid" env:"IDP_SIGNING_KID"`
	SigningMethod          string   `ocisConfig:"signing_method" env:"IDP_SIGNING_METHOD"`
	SigningPrivateKeyFiles []string `ocisConfig:"signing_private_key_files"` // TODO: is this even needed?
	ValidationKeysPath     string   `ocisConfig:"validation_keys_path" env:"IDP_VALIDATION_KEYS_PATH"`

	CookieBackendURI string
	CookieNames      []string

	AccessTokenDurationSeconds        uint64 `ocisConfig:"access_token_duration_seconds" env:"IDP_ACCESS_TOKEN_EXPIRATION"`
	IDTokenDurationSeconds            uint64 `ocisConfig:"id_token_duration_seconds" env:"IDP_ID_TOKEN_EXPIRATION"`
	RefreshTokenDurationSeconds       uint64 `ocisConfig:"refresh_token_duration_seconds" env:"IDP_REFRESH_TOKEN_EXPIRATION"`
	DyamicClientSecretDurationSeconds uint64 `ocisConfig:"dynamic_client_secret_duration_seconds" env:""`
}
