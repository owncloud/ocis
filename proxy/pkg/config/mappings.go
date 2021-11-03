package config

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		// Logging
		{
			env:         []string{"PROXY_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			destination: &cfg.Log.Level,
		},
		{
			env:         []string{"PROXY_LOG_COLOR", "OCIS_LOG_COLOR"},
			destination: &cfg.Log.Color,
		},
		{
			env:         []string{"PROXY_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			destination: &cfg.Log.Pretty,
		},
		{
			env:         []string{"PROXY_LOG_FILE", "OCIS_LOG_FILE"},
			destination: &cfg.Log.File,
		},

		// Basic auth
		{
			env:         []string{"PROXY_ENABLE_BASIC_AUTH"},
			destination: &cfg.EnableBasicAuth,
		},

		// Debug (health)
		{
			env:         []string{"PROXY_DEBUG_ADDR"},
			destination: &cfg.Debug.Addr,
		},

		{
			env:         []string{"PROXY_CONFIG_FILE"},
			destination: &cfg.File,
		},

		// Tracing
		{
			env:         []string{"PROXY_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			destination: &cfg.Tracing.Enabled,
		},
		{
			env:         []string{"PROXY_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			destination: &cfg.Tracing.Type,
		},
		{
			env:         []string{"PROXY_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			destination: &cfg.Tracing.Endpoint,
		},
		{
			env:         []string{"PROXY_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			destination: &cfg.Tracing.Collector,
		},
		{
			env:         []string{"PROXY_TRACING_SERVICE"},
			destination: &cfg.Tracing.Service,
		},

		// Debug
		{
			env:         []string{"PROXY_DEBUG_ADDR"},
			destination: &cfg.Debug.Addr,
		},
		{
			env:         []string{"PROXY_DEBUG_TOKEN"},
			destination: &cfg.Debug.Token,
		},
		{
			env:         []string{"PROXY_DEBUG_PPROF"},
			destination: &cfg.Debug.Pprof,
		},
		{
			env:         []string{"PROXY_DEBUG_ZPAGES"},
			destination: &cfg.Debug.Zpages,
		},

		// HTTP
		{
			env:         []string{"PROXY_HTTP_ADDR"},
			destination: &cfg.HTTP.Addr,
		},
		{
			env:         []string{"PROXY_HTTP_ROOT"},
			destination: &cfg.HTTP.Root,
		},

		// Service
		{
			env:         []string{"PROXY_SERVICE_NAMESPACE"},
			destination: &cfg.Service.Name,
		},
		{
			env:         []string{"PROXY_SERVICE_NAME"},
			destination: &cfg.Service.Namespace,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
		{
			env:         nil,
			destination: nil,
		},
	}
}
