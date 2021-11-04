package config

type mapping struct {
	EnvVars     []string    // name of the EnvVars var.
	Destination interface{} // memory address of the original config value to modify.
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			EnvVars:     []string{"GRAPH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"GRAPH_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"GRAPH_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"GRAPH_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"GRAPH_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"GRAPH_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"GRAPH_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"GRAPH_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"GRAPH_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"GRAPH_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"GRAPH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"GRAPH_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"GRAPH_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"GRAPH_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"GRAPH_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"GRAPH_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"GRAPH_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"GRAPH_SPACES_WEBDAV_BASE", "OCIS_URL"},
			Destination: &cfg.Spaces.WebDavBase,
		},
		{
			EnvVars:     []string{"GRAPH_SPACES_WEBDAV_PATH"},
			Destination: &cfg.Spaces.WebDavPath,
		},
		{
			EnvVars:     []string{"GRAPH_SPACES_DEFAULT_QUOTA"},
			Destination: &cfg.Spaces.DefaultQuota,
		},
		{
			EnvVars:     []string{"GRAPH_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Address,
		},
	}
}
