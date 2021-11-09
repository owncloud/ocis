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
			EnvVars:     []string{"OCIS_LOG_LEVEL", "STORE_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "STORE_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "STORE_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "STORE_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "STORE_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "STORE_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "STORE_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "STORE_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"STORE_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"STORE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"STORE_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"STORE_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"STORE_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"STORE_GRPC_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		{
			EnvVars:     []string{"STORE_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		{
			EnvVars:     []string{"STORE_NAME"},
			Destination: &cfg.Service.Name,
		},
		{
			EnvVars:     []string{"STORE_DATA_PATH"},
			Destination: &cfg.Datapath,
		},
	}
}
