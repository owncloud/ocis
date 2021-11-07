package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

// StructMappings binds a set of environment variables to a destination on cfg.
func StructMappings(cfg *Config) []shared.EnvBinding {
	return structMappings(cfg)
}

func structMappings(cfg *Config) []shared.EnvBinding {
	return []shared.EnvBinding{
		// Logging
		{
			EnvVars:     []string{"PROXY_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"PROXY_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"PROXY_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"PROXY_LOG_FILE", "OCIS_LOG_FILE"},
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

		{
			EnvVars:     []string{"PROXY_CONFIG_FILE"},
			Destination: &cfg.File,
		},

		// Tracing
		{
			EnvVars:     []string{"PROXY_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"PROXY_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"PROXY_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"PROXY_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
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
			EnvVars:     []string{"PROXY_JWT_SECRET", "OCIS_JWT_SECRET"},
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
			EnvVars:     []string{"PROXY_OIDC_ISSUER", "OCIS_URL"},
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
			EnvVars:     []string{"PROXY_MACHINE_AUTH_API_KEY", "OCIS_MACHINE_AUTH_API_KEY"},
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
