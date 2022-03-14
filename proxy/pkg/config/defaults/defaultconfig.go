package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9205",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "0.0.0.0:9200",
			Root:      "/",
			Namespace: "com.owncloud.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "proxy", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "proxy", "server.key"),
			TLS:       true,
		},
		Service: config.Service{
			Name: "proxy",
		},
		OIDC: config.OIDC{
			Issuer:   "https://localhost:9200",
			Insecure: true,
			//Insecure: true,
			UserinfoCache: config.UserinfoCache{
				Size: 1024,
				TTL:  10,
			},
		},
		TokenManager: config.TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		PolicySelector: nil,
		Reva: config.Reva{
			Address: "127.0.0.1:9142",
		},
		PreSignedURL: config.PreSignedURL{
			AllowedHTTPMethods: []string{"GET"},
			Enabled:            true,
		},
		AccountBackend:        "accounts",
		UserOIDCClaim:         "email",
		UserCS3Claim:          "mail",
		MachineAuthAPIKey:     "change-me-please",
		AutoprovisionAccounts: false,
		EnableBasicAuth:       false,
		InsecureBackends:      false,
		// TODO: enable
		//Policies: defaultPolicies(),
	}
}

func DefaultPolicies() []config.Policy {
	return []config.Policy{
		{
			Name: "ocis",
			Routes: []config.Route{
				{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/archiver",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     config.RegexRoute,
					Endpoint: "/ocs/v[12].php/cloud/(users?|groups)", // we have `user`, `users` and `groups` in ocis-ocs
					Backend:  "http://localhost:9110",
				},
				{
					Endpoint: "/ocs/",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/remote.php/?preview=1",
					Backend:  "http://localhost:9115",
				},
				{
					Endpoint: "/remote.php/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/dav/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/webdav/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/status.php",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/index.php/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/data",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/app/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/graph/",
					Backend:  "http://localhost:9120",
				},
				{
					Endpoint: "/graph-explorer",
					Backend:  "http://localhost:9135",
				},
				// if we were using the go micro api gateway we could look up the endpoint in the registry dynamically
				{
					Endpoint: "/api/v0/accounts",
					Backend:  "http://localhost:9181",
				},
				// TODO the lookup needs a better mechanism
				{
					Endpoint: "/accounts.js",
					Backend:  "http://localhost:9181",
				},
				{
					Endpoint: "/api/v0/settings",
					Backend:  "http://localhost:9190",
				},
				{
					Endpoint: "/settings.js",
					Backend:  "http://localhost:9190",
				},
			},
		},
		{
			Name: "oc10",
			Routes: []config.Route{
				{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/archiver",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint:    "/ocs/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/remote.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/dav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/webdav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/status.php",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/index.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/data",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
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
	if cfg.Policies == nil {
		cfg.Policies = DefaultPolicies()
	}

	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
