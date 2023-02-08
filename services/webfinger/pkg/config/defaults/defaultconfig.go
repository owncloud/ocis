package defaults

import (
	"strings"

	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/relations"
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
			//Addr: "127.0.0.1:19119", // FIXME
			Addr:   "127.0.0.1:0", // :0 to pick any free local port
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			//Addr: "127.0.0.1:19115", // FIXME
			Addr:      "127.0.0.1:0", // :0 to pick any free local port
			Root:      "/",
			Namespace: "com.owncloud.web",
			CORS: config.CORS{
				AllowedOrigins: []string{"*"},
			},
		},
		Service: config.Service{
			Name: "webfinger",
		},

		Relations: []string{relations.OpenIDConnectRel, relations.OwnCloudInstanceRel},
		Instances: []config.Instance{
			{
				Claim: "email",
				Regex: "einstein@example\\.org", // only einstein
				Href:  "{{.OCIS_URL}}",
				Titles: map[string]string{
					"en": "oCIS Instance for Einstein",
					"de": "oCIS Instanz für Einstein",
				},
				Break: true,
			},
			{
				Claim: "email",
				Regex: "marie@example\\.org", // only marie
				Href:  "https://{{.preferred_username}}.cloud.ocis.test",
				Titles: map[string]string{
					"en": "oCIS Instance for Marie",
					"de": "oCIS Instanz für Marie",
				},
				// also continue with next instance
			},
			{
				Claim: "email",
				Regex: ".+@example\\.org", // example.org, including marie but not for einstein
				Href:  "{{.OCIS_URL}}",    // zb https://{{schoolid}}.cloud.ocis.de bei dem der schoolid claim dann genommen wird. templates?
				Titles: map[string]string{
					"en": "oCIS Instance for example.org",
					"de": "oCIS Instanz für example.org",
				},
				Break: true,
			},
			{
				Claim: "email",
				Regex: ".+@example\\.com", // example.com
				Href:  "{{.OCIS_URL}}",
				Titles: map[string]string{
					"en": "oCIS Instance for example.com",
					"de": "oCIS Instanz für example.com",
				},
				Break: true,
			},
			{
				Claim: "email",
				Regex: ".+",
				Href:  "{{.OCIS_URL}}",
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
