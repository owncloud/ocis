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
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_ISSUER", "OCIS_URL"},
			Destination: &cfg.GraphExplorer.Issuer,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_CLIENT_ID"},
			Destination: &cfg.GraphExplorer.ClientID,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_GRAPH_URL_BASE", "OCIS_URL"},
			Destination: &cfg.GraphExplorer.GraphURLBase,
		},
		{
			EnvVars:     []string{"GRAPH_EXPLORER_GRAPH_URL_PATH"},
			Destination: &cfg.GraphExplorer.GraphURLPath,
		},
	}
}
