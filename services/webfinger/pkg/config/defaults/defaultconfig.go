package defaults

import (
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config"
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
			Addr: "127.0.0.1:19119", // FIXME
			//Addr:      "127.0.0.1:0", // :0 to pick any free local port
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr: "127.0.0.1:19115", // FIXME
			//Addr:      "127.0.0.1:0", // :0 to pick any free local port
			Root:      "/",
			Namespace: "com.owncloud.web",
			CORS: config.CORS{
				AllowedOrigins: []string{"*"},
			},
		},
		Reva: shared.DefaultRevaConfig(),
		Service: config.Service{
			Name: "webfinger",
		},
		LookupChain: "openid-discovery,owncloud-status,owncloud-account,owncloud-instance",
		Instances: []config.Instance{
			{
				Claim: "mail",
				Regex: "einstein@example.com",
				Href:  "{{OCIS_URL}}",
				Titles: map[string]string{
					"en": "oCIS Instance for Einstein",
					"de": "oCIS Instanz für Einstein",
				},
			},
			{
				Claim: "mail",
				Regex: ".*@example.com",
				Href:  "{{OCIS_URL}}",
				Titles: map[string]string{
					"en": "oCIS Instance for example.org",
					"de": "oCIS Instanz für example.org",
				},
			},
			{
				Claim: "id",
				Regex: ".*",
				Href:  "{{OCIS_URL}}",
				Titles: map[string]string{
					"en": "oCIS Instance",
					"de": "oCIS Instanz",
				},
			},
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for "envdecode".
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
	// provide with defaults for shared tracing, since we need a valid destination address for "envdecode".
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

	if cfg.Commons != nil {
		cfg.HTTP.TLS = cfg.Commons.HTTPServiceTLS
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}

}
