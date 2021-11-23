package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv(cfg *Config) []string {
	var r = make([]string, len(structMappings(cfg)))
	for i := range structMappings(cfg) {
		r = append(r, structMappings(cfg)[i].EnvVars...)
	}

	return r
}

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
			EnvVars:     []string{"OCIS_LOG_LEVEL", "GLAUTH_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "GLAUTH_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "GLAUTH_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "GLAUTH_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"GLAUTH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "GLAUTH_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "GLAUTH_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "GLAUTH_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "GLAUTH_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"GLAUTH_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"GLAUTH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"GLAUTH_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"GLAUTH_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"GLAUTH_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"GLAUTH_ROLE_BUNDLE_ID"},
			Destination: &cfg.RoleBundleUUID,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAP_ADDR"},
			Destination: &cfg.Ldap.Addr,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAP_ENABLED"},
			Destination: &cfg.Ldap.Enabled,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAPS_ADDR"},
			Destination: &cfg.Ldaps.Addr,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAPS_ENABLED"},
			Destination: &cfg.Ldaps.Enabled,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAPS_CERT"},
			Destination: &cfg.Ldaps.Cert,
		},
		{
			EnvVars:     []string{"GLAUTH_LDAPS_KEY"},
			Destination: &cfg.Ldaps.Key,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_BASEDN"},
			Destination: &cfg.Backend.BaseDN,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_NAME_FORMAT"},
			Destination: &cfg.Backend.NameFormat,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_GROUP_FORMAT"},
			Destination: &cfg.Backend.GroupFormat,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_SSH_KEY_ATTR"},
			Destination: &cfg.Backend.SSHKeyAttr,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_DATASTORE"},
			Destination: &cfg.Backend.Datastore,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_INSECURE"},
			Destination: &cfg.Backend.Insecure,
		},
		{
			EnvVars:     []string{"GLAUTH_BACKEND_USE_GRAPHAPI"},
			Destination: &cfg.Backend.UseGraphAPI,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_BASEDN"},
			Destination: &cfg.Fallback.BaseDN,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_NAME_FORMAT"},
			Destination: &cfg.Fallback.NameFormat,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_GROUP_FORMAT"},
			Destination: &cfg.Fallback.GroupFormat,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_SSH_KEY_ATTR"},
			Destination: &cfg.Fallback.SSHKeyAttr,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_DATASTORE"},
			Destination: &cfg.Fallback.Datastore,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_INSECURE"},
			Destination: &cfg.Fallback.Insecure,
		},
		{
			EnvVars:     []string{"GLAUTH_FALLBACK_USE_GRAPHAPI"},
			Destination: &cfg.Fallback.UseGraphAPI,
		},
	}
}
