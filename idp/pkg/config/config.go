package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr    string `mapstructure:"addr"`
	Root    string `mapstructure:"root"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
	TLS     bool   `mapstructure:"tls"`
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI               string `mapstructure:"uri"`
	BindDN            string `mapstructure:"bind_dn"`
	BindPassword      string `mapstructure:"bind_password"`
	BaseDN            string `mapstructure:"base_dn"`
	Scope             string `mapstructure:"scope"`
	LoginAttribute    string `mapstructure:"login_attribute"`
	EmailAttribute    string `mapstructure:"email_attribute"`
	NameAttribute     string `mapstructure:"name_attribute"`
	UUIDAttribute     string `mapstructure:"uuid_attribute"`
	UUIDAttributeType string `mapstructure:"uuid_attribute_type"`
	Filter            string `mapstructure:"filter"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `mapstructure:"name"`
	Namespace string `mapstructure:"namespace"`
	Version   string `mapstructure:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `mapstructure:"asset"`
}

type Settings struct {
	Iss                               string   `mapstructure:"iss"`
	IdentityManager                   string   `mapstructure:"identity_manager"`
	URIBasePath                       string   `mapstructure:"uri_base_path"`
	SignInURI                         string   `mapstructure:"sign_in_uri"`
	SignedOutURI                      string   `mapstructure:"signed_out_uri"`
	AuthorizationEndpointURI          string   `mapstructure:"authorization_endpoint_uri"`
	EndsessionEndpointURI             string   `mapstructure:"end_session_endpoint_uri"`
	Insecure                          bool     `mapstructure:"insecure"`
	TrustedProxy                      []string `mapstructure:"trusted_proxy"`
	AllowScope                        []string `mapstructure:"allow_scope"`
	AllowClientGuests                 bool     `mapstructure:"allow_client_guests"`
	AllowDynamicClientRegistration    bool     `mapstructure:"allow_dynamic_client_registration"`
	EncryptionSecretFile              string   `mapstructure:"encrypt_secret_file"`
	Listen                            string   `mapstructure:"listen"`
	IdentifierClientDisabled          bool     `mapstructure:"identifier_client_disabled"`
	IdentifierClientPath              string   `mapstructure:"identifier_client_path"`
	IdentifierRegistrationConf        string   `mapstructure:"identifier_registration_conf"`
	IdentifierScopesConf              string   `mapstructure:"identifier_scopes_conf"`
	IdentifierDefaultBannerLogo       string   `mapstructure:"identifier_default_banner_logo"`
	IdentifierDefaultSignInPageText   string   `mapstructure:"identifier_default_sign_in_page_text"`
	IdentifierDefaultUsernameHintText string   `mapstructure:"identifier_default_username_hint_text"`
	SigningKid                        string   `mapstructure:"sign_in_kid"`
	SigningMethod                     string   `mapstructure:"sign_in_method"`
	SigningPrivateKeyFiles            []string `mapstructure:"sign_in_private_key_files"`
	ValidationKeysPath                string   `mapstructure:"validation_keys_path"`
	CookieBackendURI                  string   `mapstructure:"cookie_backend_uri"`
	CookieNames                       []string `mapstructure:"cookie_names"`
	AccessTokenDurationSeconds        uint64   `mapstructure:"access_token_duration_seconds"`
	IDTokenDurationSeconds            uint64   `mapstructure:"id_token_duration_seconds"`
	RefreshTokenDurationSeconds       uint64   `mapstructure:"refresh_token_duration_seconds"`
	DyamicClientSecretDurationSeconds uint64   `mapstructure:"dynamic_client_secret_duration_seconds"`
}

// Config combines all available configuration parts.
type Config struct {
	File    string     `mapstructure:"file"`
	Log     shared.Log `mapstructure:"log"`
	Debug   Debug      `mapstructure:"debug"`
	HTTP    HTTP       `mapstructure:"http"`
	Tracing Tracing    `mapstructure:"tracing"`
	Asset   Asset      `mapstructure:"asset"`
	IDP     Settings   `mapstructure:"idp"`
	Ldap    Ldap       `mapstructure:"ldap"`
	Service Service    `mapstructure:"service"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
		Debug: Debug{
			Addr: "127.0.0.1:9134",
		},
		HTTP: HTTP{
			Addr:    "127.0.0.1:9130",
			Root:    "/",
			TLSCert: path.Join(defaults.BaseDataPath(), "idp", "server.crt"),
			TLSKey:  path.Join(defaults.BaseDataPath(), "idp", "server.key"),
			TLS:     false,
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
		Service: Service{
			Name:      "idp",
			Namespace: "com.owncloud.web",
		},
	}
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(DefaultConfig())))
	for i := range structMappings(DefaultConfig()) {
		r = append(r, structMappings(DefaultConfig())[i].EnvVars...)
	}

	return r
}
