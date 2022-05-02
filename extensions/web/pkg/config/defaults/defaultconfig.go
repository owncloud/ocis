package defaults

import (
	"strings"

	"github.com/owncloud/ocis/extensions/web/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9104",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9100",
			Root:      "/",
			Namespace: "com.owncloud.web",
			CacheTTL:  604800, // 7 days
		},
		Service: config.Service{
			Name: "web",
		},
		Asset: config.Asset{
			Path: "",
		},
		Web: config.Web{
			Path:        "",
			ThemeServer: "https://localhost:9200",
			ThemePath:   "/themes/owncloud/theme.json",
			Config: config.WebConfig{
				Server:  "https://localhost:9200",
				Theme:   "",
				Version: "0.1.0",
				OpenIDConnect: config.OIDC{
					MetadataURL:  "",
					Authority:    "https://localhost:9200",
					ClientID:     "web",
					ResponseType: "code",
					Scope:        "openid profile email",
				},
				Apps: []string{"files", "search", "preview", "text-editor", "pdf-viewer", "external", "user-management"},
			},
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for BindEnv.
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.Tracing{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimRight(cfg.HTTP.Root, "/")
	}
	// build well known openid-configuration endpoint if it is not set
	if cfg.Web.Config.OpenIDConnect.MetadataURL == "" {
		cfg.Web.Config.OpenIDConnect.MetadataURL = strings.TrimRight(cfg.Web.Config.OpenIDConnect.Authority, "/") + "/.well-known/openid-configuration"
	}
}
