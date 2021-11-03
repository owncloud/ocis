package config

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			goType:      "bool",
			env:         "PROXY_ENABLE_BASIC_AUTH",
			destination: &cfg.EnableBasicAuth,
		},
		{
			goType:      "string",
			env:         "PROXY_LOG_LEVEL",
			destination: &cfg.Log.Level,
		},
		{
			goType:      "bool",
			env:         "PROXY_LOG_COLOR",
			destination: &cfg.Log.Color,
		},
		{
			goType:      "bool",
			env:         "PROXY_LOG_PRETTY",
			destination: &cfg.Log.Pretty,
		},
	}
}
