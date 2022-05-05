package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/v2/extensions/idp/pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
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
			IdentifierRegistrationConf:        path.Join(defaults.BaseDataPath(), "idp", "tmp", "identifier-registration.yaml"),
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
		Clients: []config.Client{
			{
				ID:      "web",
				Name:    "ownCloud Web app",
				Trusted: true,
				RedirectURIs: []string{
					"{{OCIS_URL}}/",
					"{{OCIS_URL}}/oidc-callback.html",
					"{{OCIS_URL}}/oidc-silent-redirect.html",
				},
				Origins: []string{
					"{{OCIS_URL}}",
				},
			},
			{
				ID:      "ocis-explorer.js",
				Name:    "oCIS Graph Explorer",
				Trusted: true,
				RedirectURIs: []string{
					"{{OCIS_URL}}/graph-explorer/",
				},
				Origins: []string{
					"{{OCIS_URL}}",
				},
			},
			{
				ID:              "xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69",
				Secret:          "UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh",
				Name:            "ownCloud desktop app",
				ApplicationType: "native",
				RedirectURIs: []string{
					"http://127.0.0.1",
					"http://localhost",
				},
			},
			{
				ID:              "e4rAsNUSIUs0lF4nbv9FmCeUkTlV9GdgTLDH1b5uie7syb90SzEVrbN7HIpmWJeD",
				Secret:          "dInFYGV33xKzhbRmpqQltYNdfLdJIfJ9L5ISoKhNoT9qZftpdWSP71VrpGR9pmoD",
				Name:            "ownCloud Android app",
				ApplicationType: "native",
				RedirectURIs: []string{
					"oc://android.owncloud.com",
				},
			},
			{
				ID:              "mxd5OQDk6es5LzOzRvidJNfXLUZS2oN3oUFeXPP8LpPrhx3UroJFduGEYIBOxkY1",
				Secret:          "KFeFWWEZO9TkisIQzR3fo7hfiMXlOpaqP8CFuTbSHzV1TUuGECglPxpiVKJfOXIx",
				Name:            "ownCloud iOS app",
				ApplicationType: "native",
				RedirectURIs: []string{
					"oc://ios.owncloud.com",
					"oc.ios://ios.owncloud.com",
				},
			},
		},
		Ldap: config.Ldap{
			URI:               "ldaps://localhost:9235",
			TLSCACert:         path.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
			BindDN:            "uid=idp,ou=sysusers,o=libregraph-idm",
			BaseDN:            "ou=users,o=libregraph-idm",
			Scope:             "sub",
			LoginAttribute:    "uid",
			EmailAttribute:    "mail",
			NameAttribute:     "displayName",
			UUIDAttribute:     "uid",
			UUIDAttributeType: "text",
			Filter:            "",
			ObjectClass:       "inetOrgPerson",
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
