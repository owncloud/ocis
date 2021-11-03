package config

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			goType:      "bool",
			env:         []string{"PROXY_ENABLE_BASIC_AUTH"},
			destination: &cfg.EnableBasicAuth,
		},
		{
			goType:      "string",
			env:         []string{"PROXY_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			destination: &cfg.Log.Level,
		},
		{
			goType:      "bool",
			env:         []string{"PROXY_LOG_COLOR", "OCIS_LOG_COLOR"},
			destination: &cfg.Log.Color,
		},
		{
			goType:      "bool",
			env:         []string{"PROXY_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			destination: &cfg.Log.Pretty,
		},
	}
}
