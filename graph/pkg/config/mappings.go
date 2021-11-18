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
			EnvVars:     []string{"GRAPH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL", "GRAPH_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "GRAPH_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "GRAPH_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "GRAPH_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "GRAPH_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "GRAPH_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "GRAPH_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "GRAPH_TRACING_COLLECTOR"},
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
			EnvVars:     []string{"OCIS_URL", "GRAPH_SPACES_WEBDAV_BASE"},
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
			EnvVars:     []string{"OCIS_JWT_SECRET", "GRAPH_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Address,
		},
		{
			EnvVars:     []string{"GRAPH_IDENTITY_BACKEND"},
			Destination: &cfg.Identity.Backend,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_URI"},
			Destination: &cfg.Identity.LDAP.URI,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_BIND_DN"},
			Destination: &cfg.Identity.LDAP.BindDN,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Identity.LDAP.BindPassword,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_BASE_DN"},
			Destination: &cfg.Identity.LDAP.UserBaseDN,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_EMAIL_ATTRIBUTE"},
			Destination: &cfg.Identity.LDAP.UserEmailAttribute,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE"},
			Destination: &cfg.Identity.LDAP.UserDisplayNameAttribute,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_NAME_ATTRIBUTE"},
			Destination: &cfg.Identity.LDAP.UserNameAttribute,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_UID_ATTRIBUTE"},
			Destination: &cfg.Identity.LDAP.UserIDAttribute,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_FILTER"},
			Destination: &cfg.Identity.LDAP.UserFilter,
		},
		{
			EnvVars:     []string{"GRAPH_LDAP_USER_SCOPE"},
			Destination: &cfg.Identity.LDAP.UserSearchScope,
		},
	}
}
