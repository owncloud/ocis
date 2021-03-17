package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"IDP_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"IDP_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"IDP_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9134"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"IDP_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       flags.OverrideDefaultString(cfg.File, ""),
			Usage:       "Path to config file",
			EnvVars:     []string{"IDP_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"IDP_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"IDP_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"IDP_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"IDP_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "idp"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"IDP_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9134"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"IDP_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"IDP_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"IDP_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9130"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"IDP_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"IDP_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for service discovery",
			EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "idp"),
			Usage:       "Service name",
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "identity-manager",
			Value:       flags.OverrideDefaultString(cfg.IDP.IdentityManager, "ldap"),
			Usage:       "Identity manager (one of ldap,kc,cookie,dummy)",
			EnvVars:     []string{"IDP_IDENTITY_MANAGER"},
			Destination: &cfg.IDP.IdentityManager,
		},
		&cli.StringFlag{
			Name:        "ldap-uri",
			Value:       flags.OverrideDefaultString(cfg.Ldap.URI, "ldap://localhost:9125"),
			Usage:       "URI of the LDAP server (glauth)",
			EnvVars:     []string{"IDP_LDAP_URI"},
			Destination: &cfg.Ldap.URI,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       flags.OverrideDefaultString(cfg.Ldap.BindDN, "cn=idp,ou=sysusers,dc=example,dc=org"),
			Usage:       "Bind DN for the LDAP server (glauth)",
			EnvVars:     []string{"IDP_LDAP_BIND_DN"},
			Destination: &cfg.Ldap.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       flags.OverrideDefaultString(cfg.Ldap.BindPassword, "idp"),
			Usage:       "Password for the Bind DN of the LDAP server (glauth)",
			EnvVars:     []string{"IDP_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Ldap.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       flags.OverrideDefaultString(cfg.Ldap.BaseDN, "ou=users,dc=example,dc=org"),
			Usage:       "LDAP base DN of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_BASE_DN"},
			Destination: &cfg.Ldap.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-scope",
			Value:       flags.OverrideDefaultString(cfg.Ldap.Scope, "sub"),
			Usage:       "LDAP scope of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_SCOPE"},
			Destination: &cfg.Ldap.Scope,
		},
		&cli.StringFlag{
			Name:        "ldap-login-attribute",
			Value:       flags.OverrideDefaultString(cfg.Ldap.LoginAttribute, "cn"),
			Usage:       "LDAP login attribute of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_LOGIN_ATTRIBUTE"},
			Destination: &cfg.Ldap.LoginAttribute,
		},
		&cli.StringFlag{
			Name:        "ldap-email-attribute",
			Value:       flags.OverrideDefaultString(cfg.Ldap.EmailAttribute, "mail"),
			Usage:       "LDAP email attribute of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_EMAIL_ATTRIBUTE"},
			Destination: &cfg.Ldap.EmailAttribute,
		},
		&cli.StringFlag{
			Name:        "ldap-name-attribute",
			Value:       flags.OverrideDefaultString(cfg.Ldap.NameAttribute, "sn"),
			Usage:       "LDAP name attribute of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_NAME_ATTRIBUTE"},
			Destination: &cfg.Ldap.NameAttribute,
		},
		&cli.StringFlag{
			Name:        "ldap-uuid-attribute",
			Value:       flags.OverrideDefaultString(cfg.Ldap.UUIDAttribute, "uid"),
			Usage:       "LDAP UUID attribute of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_UUID_ATTRIBUTE"},
			Destination: &cfg.Ldap.UUIDAttribute,
		},
		&cli.StringFlag{
			Name:        "ldap-uuid-attribute-type",
			Value:       flags.OverrideDefaultString(cfg.Ldap.UUIDAttributeType, "text"),
			Usage:       "LDAP UUID attribute type of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_UUID_ATTRIBUTE_TYPE"},
			Destination: &cfg.Ldap.UUIDAttributeType,
		},
		&cli.StringFlag{
			Name:        "ldap-filter",
			Value:       flags.OverrideDefaultString(cfg.Ldap.Filter, "(objectClass=posixaccount)"),
			Usage:       "LDAP filter of the oCIS users",
			EnvVars:     []string{"IDP_LDAP_FILTER"},
			Destination: &cfg.Ldap.Filter,
		},
		&cli.StringFlag{
			Name:        "transport-tls-cert",
			Value:       flags.OverrideDefaultString(cfg.HTTP.TLSCert, ""),
			Usage:       "Certificate file for transport encryption",
			EnvVars:     []string{"IDP_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		&cli.StringFlag{
			Name:        "transport-tls-key",
			Value:       flags.OverrideDefaultString(cfg.HTTP.TLSKey, ""),
			Usage:       "Secret file for transport encryption",
			EnvVars:     []string{"IDP_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		&cli.StringFlag{
			Name:        "iss",
			Value:       flags.OverrideDefaultString(cfg.IDP.Iss, "https://localhost:9200"),
			Usage:       "OIDC issuer URL",
			EnvVars:     []string{"IDP_ISS", "OCIS_URL"}, // IDP_ISS takes precedence over OCIS_URL
			Destination: &cfg.IDP.Iss,
		},
		&cli.StringSliceFlag{
			Name:    "signing-private-key",
			Usage:   "Full path to PEM encoded private key file (must match the --signing-method algorithm)",
			EnvVars: []string{"IDP_SIGNING_PRIVATE_KEY"},
			Value:   nil,
		},
		&cli.StringFlag{
			Name:        "signing-kid",
			Usage:       "Value of kid field to use in created tokens (uniquely identifying the signing-private-key)",
			EnvVars:     []string{"IDP_SIGNING_KID"},
			Value:       flags.OverrideDefaultString(cfg.IDP.SigningKid, ""),
			Destination: &cfg.IDP.SigningKid,
		},
		&cli.StringFlag{
			Name:        "validation-keys-path",
			Usage:       "Full path to a folder containg PEM encoded private or public key files used for token validaton (file name without extension is used as kid)",
			EnvVars:     []string{"IDP_VALIDATION_KEYS_PATH"},
			Value:       flags.OverrideDefaultString(cfg.IDP.ValidationKeysPath, ""),
			Destination: &cfg.IDP.ValidationKeysPath,
		},
		&cli.StringFlag{
			Name:        "encryption-secret",
			Usage:       "Full path to a file containing a %d bytes secret key",
			EnvVars:     []string{"IDP_ENCRYPTION_SECRET"},
			Value:       flags.OverrideDefaultString(cfg.IDP.EncryptionSecretFile, ""),
			Destination: &cfg.IDP.EncryptionSecretFile,
		},
		&cli.StringFlag{
			Name:        "signing-method",
			Usage:       "JWT default signing method",
			EnvVars:     []string{"IDP_SIGNING_METHOD"},
			Value:       flags.OverrideDefaultString(cfg.IDP.SigningMethod, "PS256"),
			Destination: &cfg.IDP.SigningMethod,
		},
		&cli.StringFlag{
			Name:        "uri-base-path",
			Usage:       "Custom base path for URI endpoints",
			EnvVars:     []string{"IDP_URI_BASE_PATH"},
			Value:       flags.OverrideDefaultString(cfg.IDP.URIBasePath, ""),
			Destination: &cfg.IDP.URIBasePath,
		},
		&cli.StringFlag{
			Name:        "sign-in-uri",
			Usage:       "Custom redirection URI to sign-in form",
			EnvVars:     []string{"IDP_SIGN_IN_URI"},
			Value:       flags.OverrideDefaultString(cfg.IDP.SignInURI, ""),
			Destination: &cfg.IDP.SignInURI,
		},
		&cli.StringFlag{
			Name:        "signed-out-uri",
			Usage:       "Custom redirection URI to signed-out goodbye page",
			EnvVars:     []string{"IDP_SIGN_OUT_URI"},
			Value:       flags.OverrideDefaultString(cfg.IDP.SignedOutURI, ""),
			Destination: &cfg.IDP.SignedOutURI,
		},
		&cli.StringFlag{
			Name:        "authorization-endpoint-uri",
			Usage:       "Custom authorization endpoint URI",
			EnvVars:     []string{"IDP_ENDPOINT_URI"},
			Value:       flags.OverrideDefaultString(cfg.IDP.AuthorizationEndpointURI, ""),
			Destination: &cfg.IDP.AuthorizationEndpointURI,
		},
		&cli.StringFlag{
			Name:        "endsession-endpoint-uri",
			Usage:       "Custom endsession endpoint URI",
			EnvVars:     []string{"IDP_ENDSESSION_ENDPOINT_URI"},
			Value:       flags.OverrideDefaultString(cfg.IDP.EndsessionEndpointURI, ""),
			Destination: &cfg.IDP.EndsessionEndpointURI,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       flags.OverrideDefaultString(cfg.Asset.Path, ""),
			Usage:       "Path to custom assets",
			EnvVars:     []string{"IDP_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "identifier-client-path",
			Usage:       "Path to the identifier web client base folder",
			EnvVars:     []string{"IDP_IDENTIFIER_CLIENT_PATH"},
			Value:       flags.OverrideDefaultString(cfg.IDP.IdentifierClientPath, "/var/tmp/ocis/idp"),
			Destination: &cfg.IDP.IdentifierClientPath,
		},
		&cli.StringFlag{
			Name:        "identifier-registration-conf",
			Usage:       "Path to a identifier-registration.yaml configuration file",
			EnvVars:     []string{"IDP_IDENTIFIER_REGISTRATION_CONF"},
			Value:       flags.OverrideDefaultString(cfg.IDP.IdentifierRegistrationConf, "./config/identifier-registration.yaml"),
			Destination: &cfg.IDP.IdentifierRegistrationConf,
		},
		&cli.StringFlag{
			Name:        "identifier-scopes-conf",
			Usage:       "Path to a scopes.yaml configuration file",
			EnvVars:     []string{"IDP_IDENTIFIER_SCOPES_CONF"},
			Value:       flags.OverrideDefaultString(cfg.IDP.IdentifierScopesConf, ""),
			Destination: &cfg.IDP.IdentifierScopesConf,
		},
		&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "Disable TLS certificate and hostname validation",
			EnvVars:     []string{"IDP_INSECURE"},
			Destination: &cfg.IDP.Insecure,
		},
		&cli.BoolFlag{
			Name:        "tls",
			Usage:       "Use TLS (disable only if idp is behind a TLS-terminating reverse-proxy).",
			EnvVars:     []string{"IDP_TLS"},
			Value:       flags.OverrideDefaultBool(cfg.HTTP.TLS, false),
			Destination: &cfg.HTTP.TLS,
		},
		&cli.StringSliceFlag{
			Name:    "trusted-proxy",
			Usage:   "Trusted proxy IP or IP network (can be used multiple times)",
			EnvVars: []string{"IDP_TRUSTED_PROXY"},
			Value:   nil,
		},
		&cli.StringSliceFlag{
			Name:    "allow-scope",
			Usage:   "Allow OAuth 2 scope (can be used multiple times, if not set default scopes are allowed)",
			EnvVars: []string{"IDP_ALLOW_SCOPE"},
			Value:   nil,
		},
		&cli.BoolFlag{
			Name:        "allow-client-guests",
			Usage:       "Allow sign in of client controlled guest users",
			EnvVars:     []string{"IDP_ALLOW_CLIENT_GUESTS"},
			Destination: &cfg.IDP.AllowClientGuests,
		},
		&cli.BoolFlag{
			Name:        "allow-dynamic-client-registration",
			Usage:       "Allow dynamic OAuth2 client registration",
			EnvVars:     []string{"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"},
			Value:       flags.OverrideDefaultBool(cfg.IDP.AllowDynamicClientRegistration, true),
			Destination: &cfg.IDP.AllowDynamicClientRegistration,
		},
		&cli.BoolFlag{
			Name:        "disable-identifier-webapp",
			Usage:       "Disable built-in identifier-webapp to use a frontend hosted elsewhere.",
			EnvVars:     []string{"IDP_DISABLE_IDENTIFIER_WEBAPP"},
			Value:       flags.OverrideDefaultBool(cfg.IDP.IdentifierClientDisabled, true),
			Destination: &cfg.IDP.IdentifierClientDisabled,
		},
		&cli.Uint64Flag{
			Name:        "access-token-expiration",
			Usage:       "Expiration time of access tokens in seconds since generated",
			EnvVars:     []string{"IDP_ACCESS_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.AccessTokenDurationSeconds,
			Value:       flags.OverrideDefaultUint64(cfg.IDP.AccessTokenDurationSeconds, 60*10), // 10 minutes
		},
		&cli.Uint64Flag{
			Name:        "id-token-expiration",
			Usage:       "Expiration time of id tokens in seconds since generated",
			EnvVars:     []string{"IDP_ID_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.IDTokenDurationSeconds,
			Value:       flags.OverrideDefaultUint64(cfg.IDP.IDTokenDurationSeconds, 60*60), // 1 hour
		},
		&cli.Uint64Flag{
			Name:        "refresh-token-expiration",
			Usage:       "Expiration time of refresh tokens in seconds since generated",
			EnvVars:     []string{"IDP_REFRESH_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.RefreshTokenDurationSeconds,
			Value:       flags.OverrideDefaultUint64(cfg.IDP.RefreshTokenDurationSeconds, 60*60*24*365*3), // 1 year

		},
	}
}

// ListIDPWithConfig applies the config to the list commands flags
func ListIDPWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:        "http-namespace",
		Value:       flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web"),
		Usage:       "Set the base namespace for service discovery",
		EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
		Destination: &cfg.Service.Namespace,
	},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "idp"),
			Usage:       "Service name",
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
