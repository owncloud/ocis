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
			EnvVars:     []string{"OCIS_LOG_FILE", "THUMBNAILS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL", "THUMBNAILS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "THUMBNAILS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "THUMBNAILS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"THUMBNAILS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "THUMBNAILS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "THUMBNAILS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "THUMBNAILS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "THUMBNAILS_TRACING_COLLECTOR"},
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
