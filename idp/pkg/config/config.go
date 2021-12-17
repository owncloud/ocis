package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"IDP_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"IDP_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"IDP_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"IDP_DEBUG_ZPAGES"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"IDP_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"IDP_HTTP_ROOT"`
	Namespace string
	TLSCert   string `ocisConfig:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT"`
	TLSKey    string `ocisConfig:"tls_key" env:"IDP_TRANSPORT_TLS_KEY"`
	TLS       bool   `ocisConfig:"tls" env:"IDP_TLS"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
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

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;IDP_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;IDP_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;IDP_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;IDP_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"IDP_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;IDP_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;IDP_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;IDP_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;IDP_LOG_FILE"`
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

	Listen string `ocisConfig:"listen"` //TODO: is this even needed?

	IdentifierClientDisabled          bool   `ocisConfig:"identifier_client_disabled" env:"IDP_DISABLE_IDENTIFIER_WEBAPP"`
	IdentifierClientPath              string `ocisConfig:"identifier_client_path" env:"IDP_IDENTIFIER_CLIENT_PATH"`
	IdentifierRegistrationConf        string `ocisConfig:"identifier_registration_conf" env:"IDP_IDENTIFIER_REGISTRATION_CONF"`
	IdentifierScopesConf              string `ocisConfig:"identifier_scopes_conf" env:"IDP_IDENTIFIER_SCOPES_CONF"`
	IdentifierDefaultBannerLogo       string `ocisConfig:"identifier_default_banner_logo"`        // TODO: is this even needed?
	IdentifierDefaultSignInPageText   string `ocisConfig:"identifier_default_sign_in_page_text"`  // TODO: is this even needed?
	IdentifierDefaultUsernameHintText string `ocisConfig:"identifier_default_username_hint_text"` // TODO: is this even needed?

	SigningKid             string   `ocisConfig:"sign_in_kid" env:"IDP_SIGNING_KID"`
	SigningMethod          string   `ocisConfig:"sign_in_method" env:"IDP_SIGNING_METHOD"`
	SigningPrivateKeyFiles []string `ocisConfig:"sign_in_private_key_files"` // TODO: is this even needed?
	ValidationKeysPath     string   `ocisConfig:"validation_keys_path" env:"IDP_VALIDATION_KEYS_PATH"`

	CookieBackendURI string   `ocisConfig:"cookie_backend_uri"` // TODO: is this even needed?
	CookieNames      []string `ocisConfig:"cookie_names"`       // TODO: is this even needed?

	AccessTokenDurationSeconds        uint64 `ocisConfig:"access_token_duration_seconds" env:"IDP_ACCESS_TOKEN_EXPIRATION"`
	IDTokenDurationSeconds            uint64 `ocisConfig:"id_token_duration_seconds" env:"IDP_ID_TOKEN_EXPIRATION"`
	RefreshTokenDurationSeconds       uint64 `ocisConfig:"refresh_token_duration_seconds" env:"IDP_REFRESH_TOKEN_EXPIRATION"`
	DyamicClientSecretDurationSeconds uint64 `ocisConfig:"dynamic_client_secret_duration_seconds" env:""`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Asset Asset    `ocisConfig:"asset"`
	IDP   Settings `ocisConfig:"idp"`
	Ldap  Ldap     `ocisConfig:"ldap"`

	Context    context.Context
	Supervised bool
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr: "127.0.0.1:9134",
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9130",
			Root:      "/",
			Namespace: "com.owncloud.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "idp", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "idp", "server.key"),
			TLS:       false,
		},
		Service: Service{
			Name: "idp",
		},
		Tracing: Tracing{
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "idp",
		},
		Asset: Asset{},
		IDP: Settings{
			Iss:                               "https://localhost:9200",
			IdentityManager:                   "ldap",
			URIBasePath:                       "",
			SignInURI:                         "",
			SignedOutURI:                      "",
			AuthorizationEndpointURI:          "",
			EndsessionEndpointURI:             "",
			Insecure:                          false,
			TrustedProxy:                      nil,
			AllowScope:                        nil,
			AllowClientGuests:                 false,
			AllowDynamicClientRegistration:    false,
			EncryptionSecretFile:              "",
			Listen:                            "",
			IdentifierClientDisabled:          true,
			IdentifierClientPath:              path.Join(defaults.BaseDataPath(), "idp"),
			IdentifierRegistrationConf:        path.Join(defaults.BaseDataPath(), "idp", "identifier-registration.yaml"),
			IdentifierScopesConf:              "",
			IdentifierDefaultBannerLogo:       "",
			IdentifierDefaultSignInPageText:   "",
			IdentifierDefaultUsernameHintText: "",
			SigningKid:                        "",
			SigningMethod:                     "PS256",
			SigningPrivateKeyFiles:            nil,
			ValidationKeysPath:                "",
			CookieBackendURI:                  "",
			CookieNames:                       nil,
			AccessTokenDurationSeconds:        60 * 10,                // 10 minutes
			IDTokenDurationSeconds:            60 * 60,                // 1 hour
			RefreshTokenDurationSeconds:       60 * 60 * 24 * 365 * 3, // 1 year
			DyamicClientSecretDurationSeconds: 0,
		},
		Ldap: Ldap{
			URI:               "ldap://localhost:9125",
			BindDN:            "cn=idp,ou=sysusers,dc=ocis,dc=test",
			BindPassword:      "idp",
			BaseDN:            "ou=users,dc=ocis,dc=test",
			Scope:             "sub",
			LoginAttribute:    "cn",
			EmailAttribute:    "mail",
			NameAttribute:     "sn",
			UUIDAttribute:     "uid",
			UUIDAttributeType: "text",
			Filter:            "(objectClass=posixaccount)",
		},
	}
}
