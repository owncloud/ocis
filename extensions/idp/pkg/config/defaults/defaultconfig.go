package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/extensions/idp/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)
	Sanitize(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr: "127.0.0.1:9134",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9130",
			Root:      "/",
			Namespace: "com.owncloud.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "idp", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "idp", "server.key"),
			TLS:       false,
		},
		Service: config.Service{
			Name: "idp",
		},
		Asset: config.Asset{},
		IDP: config.Settings{
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
		Ldap: config.Ldap{
			URI:               "ldap://localhost:9125",
			BindDN:            "cn=idp,ou=sysusers,dc=ocis,dc=test",
			BindPassword:      "idp",
			BaseDN:            "ou=users,dc=ocis,dc=test",
			Scope:             "sub",
			LoginAttribute:    "cn",
			EmailAttribute:    "mail",
			NameAttribute:     "displayName",
			UUIDAttribute:     "uid",
			UUIDAttributeType: "text",
			Filter:            "",
			ObjectClass:       "posixAccount",
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for BindEnv.
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.Tracing{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
