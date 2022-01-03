package config

import "github.com/owncloud/ocis/ocis-pkg/version"

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9104",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9100",
			Root:      "/",
			Namespace: "com.owncloud.web",
			CacheTTL:  604800, // 7 days
		},
		Service: Service{
			Name:    "web",
			Version: version.String, // TODO: ensure everywhere or remove
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
		},
		Asset: Asset{
			Path: "",
		},
		Web: Web{
			Path:        "",
			ThemeServer: "https://localhost:9200",
			ThemePath:   "/themes/owncloud/theme.json",
			Config: WebConfig{
				Server:  "https://localhost:9200",
				Theme:   "",
				Version: "0.1.0",
				OpenIDConnect: OIDC{
					MetadataURL:  "",
					Authority:    "https://localhost:9200",
					ClientID:     "web",
					ResponseType: "code",
					Scope:        "openid profile email",
				},
				Apps: []string{"files", "search", "media-viewer", "external"},
			},
		},
	}
}
