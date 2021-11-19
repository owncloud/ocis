package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr    string `ocisConfig:"addr"`
	Root    string `ocisConfig:"root"`
	TLSCert string `ocisConfig:"tls_cert"`
	TLSKey  string `ocisConfig:"tls_key"`
	TLS     bool   `ocisConfig:"tls"`
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI               string `ocisConfig:"uri"`
	BindDN            string `ocisConfig:"bind_dn"`
	BindPassword      string `ocisConfig:"bind_password"`
	BaseDN            string `ocisConfig:"base_dn"`
	Scope             string `ocisConfig:"scope"`
	LoginAttribute    string `ocisConfig:"login_attribute"`
	EmailAttribute    string `ocisConfig:"email_attribute"`
	NameAttribute     string `ocisConfig:"name_attribute"`
	UUIDAttribute     string `ocisConfig:"uuid_attribute"`
	UUIDAttributeType string `ocisConfig:"uuid_attribute_type"`
	Filter            string `ocisConfig:"filter"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `ocisConfig:"name"`
	Namespace string `ocisConfig:"namespace"`
	Version   string `ocisConfig:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"asset"`
}

type Settings struct {
	Iss                               string   `ocisConfig:"iss"`
	IdentityManager                   string   `ocisConfig:"identity_manager"`
	URIBasePath                       string   `ocisConfig:"uri_base_path"`
	SignInURI                         string   `ocisConfig:"sign_in_uri"`
	SignedOutURI                      string   `ocisConfig:"signed_out_uri"`
	AuthorizationEndpointURI          string   `ocisConfig:"authorization_endpoint_uri"`
	EndsessionEndpointURI             string   `ocisConfig:"end_session_endpoint_uri"`
	Insecure                          bool     `ocisConfig:"insecure"`
	TrustedProxy                      []string `ocisConfig:"trusted_proxy"`
	AllowScope                        []string `ocisConfig:"allow_scope"`
	AllowClientGuests                 bool     `ocisConfig:"allow_client_guests"`
	AllowDynamicClientRegistration    bool     `ocisConfig:"allow_dynamic_client_registration"`
	EncryptionSecretFile              string   `ocisConfig:"encrypt_secret_file"`
	Listen                            string   `ocisConfig:"listen"`
	IdentifierClientDisabled          bool     `ocisConfig:"identifier_client_disabled"`
	IdentifierClientPath              string   `ocisConfig:"identifier_client_path"`
	IdentifierRegistrationConf        string   `ocisConfig:"identifier_registration_conf"`
	IdentifierScopesConf              string   `ocisConfig:"identifier_scopes_conf"`
	IdentifierDefaultBannerLogo       string   `ocisConfig:"identifier_default_banner_logo"`
	IdentifierDefaultSignInPageText   string   `ocisConfig:"identifier_default_sign_in_page_text"`
	IdentifierDefaultUsernameHintText string   `ocisConfig:"identifier_default_username_hint_text"`
	SigningKid                        string   `ocisConfig:"sign_in_kid"`
	SigningMethod                     string   `ocisConfig:"sign_in_method"`
	SigningPrivateKeyFiles            []string `ocisConfig:"sign_in_private_key_files"`
	ValidationKeysPath                string   `ocisConfig:"validation_keys_path"`
	CookieBackendURI                  string   `ocisConfig:"cookie_backend_uri"`
	CookieNames                       []string `ocisConfig:"cookie_names"`
	AccessTokenDurationSeconds        uint64   `ocisConfig:"access_token_duration_seconds"`
	IDTokenDurationSeconds            uint64   `ocisConfig:"id_token_duration_seconds"`
	RefreshTokenDurationSeconds       uint64   `ocisConfig:"refresh_token_duration_seconds"`
	DyamicClientSecretDurationSeconds uint64   `ocisConfig:"dynamic_client_secret_duration_seconds"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File    string      `ocisConfig:"file"`
	Log     *shared.Log `ocisConfig:"log"`
	Debug   Debug       `ocisConfig:"debug"`
	HTTP    HTTP        `ocisConfig:"http"`
	Tracing Tracing     `ocisConfig:"tracing"`
	Asset   Asset       `ocisConfig:"asset"`
	IDP     Settings    `ocisConfig:"idp"`
	Ldap    Ldap        `ocisConfig:"ldap"`
	Service Service     `ocisConfig:"service"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
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
