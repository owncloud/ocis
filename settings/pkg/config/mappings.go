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
			EnvVars:     []string{"SETTINGS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"SETTINGS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"SETTINGS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"SETTINGS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"SETTINGS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"SETTINGS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"SETTINGS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"SETTINGS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"SETTINGS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"SETTINGS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"SETTINGS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"SETTINGS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"SETTINGS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"SETTINGS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"SETTINGS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"SETTINGS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"SETTINGS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		{
			EnvVars:     []string{"SETTINGS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		{
			EnvVars:     []string{"SETTINGS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		{
			EnvVars:     []string{"SETTINGS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		{
			EnvVars:     []string{"SETTINGS_NAME"},
			Destination: &cfg.Service.Name,
		},
		{
			EnvVars:     []string{"SETTINGS_DATA_PATH"},
			Destination: &cfg.Service.DataPath,
		},
		{
			EnvVars:     []string{"SETTINGS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
	}
}
