package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

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
			EnvVars:     []string{"OCS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"OCS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"OCS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"OCS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"OCS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"OCS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"OCS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"OCS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"OCS_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		{
			EnvVars:     []string{"OCS_NAME"},
			Destination: &cfg.Service.Name,
		},
		{
			EnvVars:     []string{"OCS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"OCS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"OCS_ACCOUNT_BACKEND_TYPE"},
			Destination: &cfg.AccountBackend,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.RevaAddress,
		},
		{
			EnvVars:     []string{"OCS_MACHINE_AUTH_API_KEY", "OCIS_MACHINE_AUTH_API_KEY"},
			Destination: &cfg.MachineAuthAPIKey,
		},
		{
			EnvVars:     []string{"OCS_IDM_ADDRESS", "OCIS_URL"},
			Destination: &cfg.IdentityManagement.Address,
		},
		{
			EnvVars:     []string{"OCS_STORAGE_USERS_DRIVER", "STORAGE_USERS_DRIVER"},
			Destination: &cfg.StorageUsersDriver,
		},
	}
}
