package config

type mapping struct {
	EnvVars     []string    // name of the EnvVars var.
	Destination interface{} // memory address of the original config value to modify.
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			EnvVars:     []string{"ACCOUNTS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"ACCOUNTS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HASH_DIFFICULTY"},
			Destination: &cfg.Server.HashDifficulty,
		},
		{
			EnvVars:     []string{"ACCOUNTS_DEMO_USERS_AND_GROUPS"},
			Destination: &cfg.Server.DemoUsersAndGroups,
		},
		{
			EnvVars:     []string{"ACCOUNTS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		{
			EnvVars:     []string{"ACCOUNTS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_BACKEND"},
			Destination: &cfg.Repo.Backend,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_DISK_PATH"},
			Destination: &cfg.Repo.Disk.Path,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"},
			Destination: &cfg.Repo.CS3.ProviderAddr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_URL"},
			Destination: &cfg.Repo.CS3.DataURL,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_PREFIX"},
			Destination: &cfg.Repo.CS3.DataPrefix,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.Repo.CS3.JWTSecret,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UUID"},
			Destination: &cfg.ServiceUser.UUID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_USERNAME"},
			Destination: &cfg.ServiceUser.Username,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UID"},
			Destination: &cfg.ServiceUser.UID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_GID"},
			Destination: &cfg.ServiceUser.GID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.UID.Lower,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.GID.Lower,
		},
		{
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.UID.Upper,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.GID.Upper,
		},
	}
}
