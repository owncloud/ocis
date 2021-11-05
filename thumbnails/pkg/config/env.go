package config

type mapping struct {
	EnvVars     []string    // name of the EnvVars var.
	Destination interface{} // memory address of the original config value to modify.
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			EnvVars:     []string{"THUMBNAILS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"THUMBNAILS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"THUMBNAILS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"THUMBNAILS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"THUMBNAILS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"THUMBNAILS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"THUMBNAILS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"THUMBNAILS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"THUMBNAILS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"THUMBNAILS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"THUMBNAILS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"THUMBNAILS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"THUMBNAILS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"THUMBNAILS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		{
			EnvVars:     []string{"THUMBNAILS_GRPC_ADDR"},
			Destination: &cfg.Server.Address,
		},
		{
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
		{
			EnvVars:     []string{"THUMBNAILS_FILESYSTEMSTORAGE_ROOT"},
			Destination: &cfg.Thumbnail.FileSystemStorage.RootDirectory,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Thumbnail.RevaGateway,
		},
		{
			EnvVars:     []string{"THUMBNAILS_WEBDAVSOURCE_INSECURE"},
			Destination: &cfg.Thumbnail.WebdavAllowInsecure,
		},
		{
			EnvVars:     []string{"STORAGE_WEBDAV_NAMESPACE"},
			Destination: &cfg.Thumbnail.WebdavNamespace,
		},
	}
}
