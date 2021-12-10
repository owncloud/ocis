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
			EnvVars:     []string{"OCIS_LOG_LEVEL", "WEB_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "WEB_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "WEB_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE", "WEB_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"WEB_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "WEB_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "WEB_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "WEB_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "WEB_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"WEB_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"WEB_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"WEB_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"WEB_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"WEB_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"WEB_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"WEB_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"WEB_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"WEB_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		{
			EnvVars:     []string{"WEB_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		{
			EnvVars:     []string{"WEB_UI_CONFIG"},
			Destination: &cfg.Web.Path,
		},
		{
			EnvVars:     []string{"OCIS_URL", "WEB_UI_CONFIG_SERVER"}, // WEB_UI_CONFIG_SERVER takes precedence over OCIS_URL
			Destination: &cfg.Web.Config.Server,
		},
		{
			EnvVars:     []string{"OCIS_URL", "WEB_UI_THEME_SERVER"}, // WEB_UI_THEME_SERVER takes precedence over OCIS_URL
			Destination: &cfg.Web.ThemeServer,
		},
		{
			EnvVars:     []string{"WEB_UI_THEME_PATH"},
			Destination: &cfg.Web.ThemePath,
		},
		{
			EnvVars:     []string{"WEB_UI_CONFIG_VERSION"},
			Destination: &cfg.Web.Config.Version,
		},
		{
			EnvVars:     []string{"WEB_OIDC_METADATA_URL"},
			Destination: &cfg.Web.Config.OpenIDConnect.MetadataURL,
		},
		{
			EnvVars:     []string{"OCIS_URL", "WEB_OIDC_AUTHORITY"}, // WEB_OIDC_AUTHORITY takes precedence over OCIS_URL
			Destination: &cfg.Web.Config.OpenIDConnect.Authority,
		},
		{
			EnvVars:     []string{"WEB_OIDC_CLIENT_ID"},
			Destination: &cfg.Web.Config.OpenIDConnect.ClientID,
		},
		{
			EnvVars:     []string{"WEB_OIDC_RESPONSE_TYPE"},
			Destination: &cfg.Web.Config.OpenIDConnect.ResponseType,
		},
		{
			EnvVars:     []string{"WEB_OIDC_SCOPE"},
			Destination: &cfg.Web.Config.OpenIDConnect.Scope,
		},
	}
}
