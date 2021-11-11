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

func structMappings(cfg *Config) []shared.EnvBinding {
	return []shared.EnvBinding{
		// Logging
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL", "PROXY_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "PROXY_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "PROXY_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "PROXY_LOG_FILE"},
			Destination: &cfg.Log.File,
		},

		// Basic auth
		{
			EnvVars:     []string{"PROXY_ENABLE_BASIC_AUTH"},
			Destination: &cfg.EnableBasicAuth,
		},

		// Debug (health)
		{
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},

		// Tracing
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "PROXY_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "PROXY_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "PROXY_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "PROXY_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"PROXY_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},

		// Debug
		{
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"PROXY_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"PROXY_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"PROXY_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},

		// HTTP
		{
			EnvVars:     []string{"PROXY_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"PROXY_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},

		// Service
		{
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		{
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
		{
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		{
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		{
			EnvVars:     []string{"PROXY_TLS"},
			Destination: &cfg.HTTP.TLS,
		},
		{
			EnvVars:     []string{"OCIS_JWT_SECRET", "PROXY_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},

		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Address,
		},
		{
			EnvVars:     []string{"PROXY_INSECURE_BACKENDS"},
			Destination: &cfg.InsecureBackends,
		},
		{
			EnvVars:     []string{"OCIS_URL", "PROXY_OIDC_ISSUER"},
			Destination: &cfg.OIDC.Issuer,
		},
		{
			EnvVars:     []string{"PROXY_OIDC_INSECURE"},
			Destination: &cfg.OIDC.Insecure,
		},
		{
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_TTL"},
			Destination: &cfg.OIDC.UserinfoCache.TTL,
		},
		{
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_SIZE"},
			Destination: &cfg.OIDC.UserinfoCache.Size,
		},
		{
			EnvVars:     []string{"PROXY_AUTOPROVISION_ACCOUNTS"},
			Destination: &cfg.AutoprovisionAccounts,
		},
		{
			EnvVars:     []string{"PROXY_USER_OIDC_CLAIM"},
			Destination: &cfg.UserOIDCClaim,
		},
		{
			EnvVars:     []string{"PROXY_USER_CS3_CLAIM"},
			Destination: &cfg.UserCS3Claim,
		},
		{
			EnvVars:     []string{"PROXY_ENABLE_PRESIGNEDURLS"},
			Destination: &cfg.PreSignedURL.Enabled,
		},
		{
			EnvVars:     []string{"PROXY_ACCOUNT_BACKEND_TYPE"},
			Destination: &cfg.AccountBackend,
		},
		{
			EnvVars:     []string{"OCIS_MACHINE_AUTH_API_KEY", "PROXY_MACHINE_AUTH_API_KEY"},
			Destination: &cfg.MachineAuthAPIKey,
		},
		// there are 2 missing bindings:
		// EnvVars: []string{"PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT"},
		// EnvVars: []string{"PRESIGNEDURL_ALLOWED_METHODS"},
		// since they both have no destination
		// see https://github.com/owncloud/ocis/blob/52e5effa4fa05a1626d46f7d4cb574dde3a54593/proxy/pkg/flagset/flagset.go#L256-L261
		// and https://github.com/owncloud/ocis/blob/52e5effa4fa05a1626d46f7d4cb574dde3a54593/proxy/pkg/flagset/flagset.go#L295-L300
	}
}
