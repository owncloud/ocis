package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv(cfg *Config) []string {
	var r = make([]string, len(structMappings(cfg)))
	for i := range structMappings(cfg) {
		r = append(r, structMappings(cfg)[i].EnvVars...)
	}

	return r
}

// StructMappings binds a set of environment variables to a destination on cfg. Iterating over this set and editing the
// Destination value of a binding will alter the original value, as it is a pointer to its memory address. This lets
// us propagate changes easier.
func StructMappings(cfg *Config) []shared.EnvBinding {
	return structMappings(cfg)
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []shared.EnvBinding {
	return []shared.EnvBinding{
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL", "IDP_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "IDP_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "IDP_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "IDP_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"IDP_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "IDP_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "IDP_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "IDP_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "IDP_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"IDP_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"IDP_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"IDP_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"IDP_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"IDP_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"IDP_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		{
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
		{
			EnvVars:     []string{"IDP_IDENTITY_MANAGER"},
			Destination: &cfg.IDP.IdentityManager,
		},
		{
			EnvVars:     []string{"IDP_LDAP_URI"},
			Destination: &cfg.Ldap.URI,
		},
		{
			EnvVars:     []string{"IDP_LDAP_BIND_DN"},
			Destination: &cfg.Ldap.BindDN,
		},
		{
			EnvVars:     []string{"IDP_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Ldap.BindPassword,
		},
		{
			EnvVars:     []string{"IDP_LDAP_BASE_DN"},
			Destination: &cfg.Ldap.BaseDN,
		},
		{
			EnvVars:     []string{"IDP_LDAP_SCOPE"},
			Destination: &cfg.Ldap.Scope,
		},
		{
			EnvVars:     []string{"IDP_LDAP_LOGIN_ATTRIBUTE"},
			Destination: &cfg.Ldap.LoginAttribute,
		},
		{
			EnvVars:     []string{"IDP_LDAP_EMAIL_ATTRIBUTE"},
			Destination: &cfg.Ldap.EmailAttribute,
		},
		{
			EnvVars:     []string{"IDP_LDAP_NAME_ATTRIBUTE"},
			Destination: &cfg.Ldap.NameAttribute,
		},
		{
			EnvVars:     []string{"IDP_LDAP_UUID_ATTRIBUTE"},
			Destination: &cfg.Ldap.UUIDAttribute,
		},
		{
			EnvVars:     []string{"IDP_LDAP_UUID_ATTRIBUTE_TYPE"},
			Destination: &cfg.Ldap.UUIDAttributeType,
		},
		{
			EnvVars:     []string{"IDP_LDAP_FILTER"},
			Destination: &cfg.Ldap.Filter,
		},
		{
			EnvVars:     []string{"IDP_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		{
			EnvVars:     []string{"IDP_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		{
			EnvVars:     []string{"OCIS_URL", "IDP_ISS"}, // IDP_ISS takes precedence over OCIS_URL
			Destination: &cfg.IDP.Iss,
		},
		{
			EnvVars:     []string{"IDP_SIGNING_KID"},
			Destination: &cfg.IDP.SigningKid,
		},
		{
			EnvVars:     []string{"IDP_VALIDATION_KEYS_PATH"},
			Destination: &cfg.IDP.ValidationKeysPath,
		},
		{
			EnvVars:     []string{"IDP_ENCRYPTION_SECRET"},
			Destination: &cfg.IDP.EncryptionSecretFile,
		},
		{
			EnvVars:     []string{"IDP_SIGNING_METHOD"},
			Destination: &cfg.IDP.SigningMethod,
		},
		{
			EnvVars:     []string{"IDP_URI_BASE_PATH"},
			Destination: &cfg.IDP.URIBasePath,
		},
		{
			EnvVars:     []string{"IDP_SIGN_IN_URI"},
			Destination: &cfg.IDP.SignInURI,
		},
		{
			EnvVars:     []string{"IDP_SIGN_OUT_URI"},
			Destination: &cfg.IDP.SignedOutURI,
		},
		{
			EnvVars:     []string{"IDP_ENDPOINT_URI"},
			Destination: &cfg.IDP.AuthorizationEndpointURI,
		},
		{
			EnvVars:     []string{"IDP_ENDSESSION_ENDPOINT_URI"},
			Destination: &cfg.IDP.EndsessionEndpointURI,
		},
		{
			EnvVars:     []string{"IDP_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		{
			EnvVars:     []string{"IDP_IDENTIFIER_CLIENT_PATH"},
			Destination: &cfg.IDP.IdentifierClientPath,
		},
		{
			EnvVars:     []string{"IDP_IDENTIFIER_REGISTRATION_CONF"},
			Destination: &cfg.IDP.IdentifierRegistrationConf,
		},
		{
			EnvVars:     []string{"IDP_IDENTIFIER_SCOPES_CONF"},
			Destination: &cfg.IDP.IdentifierScopesConf,
		},
		{
			EnvVars:     []string{"IDP_INSECURE"},
			Destination: &cfg.IDP.Insecure,
		},
		{
			EnvVars:     []string{"IDP_TLS"},
			Destination: &cfg.HTTP.TLS,
		},
		{
			EnvVars:     []string{"IDP_ALLOW_CLIENT_GUESTS"},
			Destination: &cfg.IDP.AllowClientGuests,
		},
		{
			EnvVars:     []string{"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"},
			Destination: &cfg.IDP.AllowDynamicClientRegistration,
		},
		{
			EnvVars:     []string{"IDP_DISABLE_IDENTIFIER_WEBAPP"},
			Destination: &cfg.IDP.IdentifierClientDisabled,
		},
		{
			EnvVars:     []string{"IDP_ACCESS_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.AccessTokenDurationSeconds,
		},
		{
			EnvVars:     []string{"IDP_ID_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.IDTokenDurationSeconds,
		},
		{
			EnvVars:     []string{"IDP_REFRESH_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.RefreshTokenDurationSeconds,
		},
	}
}
