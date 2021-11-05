package config

type mapping struct {
	EnvVars     []string    // name of the EnvVars var.
	Destination interface{} // memory address of the original config value to modify.
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			EnvVars:     []string{"STORE_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"STORE_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"STORE_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"STORE_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"STORE_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"STORE_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"STORE_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"STORE_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
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
