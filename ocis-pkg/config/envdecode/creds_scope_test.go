package envdecode

import (
	"testing"
)

type credsScopeEvents struct {
	Endpoint  string `env:"TEST_EVENTS_ENDPOINT"`
	CredsFile string `env:"TEST_NATS_CREDS_FILE;TEST_SVC_EVENTS_CREDS_FILE"`
	Creds     string `env:"TEST_NATS_CREDS;TEST_SVC_EVENTS_CREDS"`
}

type credsScopeCache struct {
	ProviderCacheCredsFile string `env:"TEST_NATS_CREDS_FILE;TEST_SVC_CACHE_CREDS_FILE"`
	ProviderCacheCreds     string `env:"TEST_NATS_CREDS;TEST_SVC_CACHE_CREDS"`
}

type credsScopeConfig struct {
	Events credsScopeEvents
	Cache  *credsScopeCache
}

func TestResolveCredsScope(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		preset    credsScopeEvents // values from a config file
		wantFile  string
		wantCreds string
	}{
		{
			name:     "service file wins over global content",
			env:      map[string]string{"TEST_NATS_CREDS": "global", "TEST_SVC_EVENTS_CREDS_FILE": "/svc.creds"},
			wantFile: "/svc.creds",
		},
		{
			name:      "service content wins over global file",
			env:       map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_SVC_EVENTS_CREDS": "svc"},
			wantFile:  "/global.creds",
			wantCreds: "svc",
		},
		{
			name:      "global scope keeps both, the content wins later",
			env:       map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_NATS_CREDS": "global"},
			wantFile:  "/global.creds",
			wantCreds: "global",
		},
		{
			name:      "service scope keeps both, the content wins later",
			env:       map[string]string{"TEST_SVC_EVENTS_CREDS_FILE": "/svc.creds", "TEST_SVC_EVENTS_CREDS": "svc"},
			wantFile:  "/svc.creds",
			wantCreds: "svc",
		},
		{
			name:      "global content env wins over a file from the config file",
			env:       map[string]string{"TEST_NATS_CREDS": "global"},
			preset:    credsScopeEvents{CredsFile: "/yaml.creds"},
			wantFile:  "/yaml.creds",
			wantCreds: "global",
		},
		{
			name:     "service file env wins over content from the config file",
			env:      map[string]string{"TEST_SVC_EVENTS_CREDS_FILE": "/svc.creds"},
			preset:   credsScopeEvents{Creds: "yaml"},
			wantFile: "/svc.creds",
		},
		{
			name:      "only the file",
			env:       map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds"},
			wantFile:  "/global.creds",
			wantCreds: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TEST_EVENTS_ENDPOINT", "nats")
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg := credsScopeConfig{Events: tc.preset}
			if err := Decode(&cfg); err != nil {
				t.Fatal(err)
			}
			if cfg.Events.CredsFile != tc.wantFile || cfg.Events.Creds != tc.wantCreds {
				t.Errorf("got file %q creds %q, want file %q creds %q", cfg.Events.CredsFile, cfg.Events.Creds, tc.wantFile, tc.wantCreds)
			}
		})
	}
}

type credsScopeAuth struct {
	AuthUsername string `env:"TEST_EVENTS_AUTH_USERNAME;TEST_SVC_EVENTS_AUTH_USERNAME"`
	AuthPassword string `env:"TEST_EVENTS_AUTH_PASSWORD;TEST_SVC_EVENTS_AUTH_PASSWORD"`
	CredsFile    string `env:"TEST_NATS_CREDS_FILE;TEST_SVC_EVENTS_CREDS_FILE"`
	Creds        string `env:"TEST_NATS_CREDS;TEST_SVC_EVENTS_CREDS"`
}

func TestResolveCredsScopeUsernameAndPassword(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		preset    credsScopeAuth // values from a config file
		wantFile  string
		wantCreds string
	}{
		{
			name:     "service username and password win over a global file",
			env:      map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_SVC_EVENTS_AUTH_USERNAME": "user", "TEST_SVC_EVENTS_AUTH_PASSWORD": "pass"},
			wantFile: "",
		},
		{
			name: "service username and password win over global content and a global file",
			env:  map[string]string{"TEST_NATS_CREDS": "global", "TEST_NATS_CREDS_FILE": "/global.creds", "TEST_SVC_EVENTS_AUTH_USERNAME": "user", "TEST_SVC_EVENTS_AUTH_PASSWORD": "pass"},
		},
		{
			name:     "service username and global password win over a global file",
			env:      map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_SVC_EVENTS_AUTH_USERNAME": "user", "TEST_EVENTS_AUTH_PASSWORD": "pass"},
			wantFile: "",
		},
		{
			name:      "service password alone does not win over global content",
			env:       map[string]string{"TEST_NATS_CREDS": "global", "TEST_SVC_EVENTS_AUTH_PASSWORD": "pass"},
			wantCreds: "global",
		},
		{
			name:     "service username alone does not win over a global file",
			env:      map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_SVC_EVENTS_AUTH_USERNAME": "user"},
			wantFile: "/global.creds",
		},
		{
			name:     "global username does not win over a global file",
			env:      map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds", "TEST_EVENTS_AUTH_USERNAME": "user"},
			wantFile: "/global.creds",
		},
		{
			name:      "service username does not win over service content",
			env:       map[string]string{"TEST_SVC_EVENTS_CREDS": "svc", "TEST_SVC_EVENTS_AUTH_USERNAME": "user"},
			wantCreds: "svc",
		},
		{
			name:     "service file wins over a global username",
			env:      map[string]string{"TEST_SVC_EVENTS_CREDS_FILE": "/svc.creds", "TEST_EVENTS_AUTH_USERNAME": "user"},
			wantFile: "/svc.creds",
		},
		{
			name:     "username from the config file does not win over a global file",
			env:      map[string]string{"TEST_NATS_CREDS_FILE": "/global.creds"},
			preset:   credsScopeAuth{AuthUsername: "user", AuthPassword: "pass"},
			wantFile: "/global.creds",
		},
		{
			name:     "service username and password do not clear a file from the config file",
			env:      map[string]string{"TEST_SVC_EVENTS_AUTH_USERNAME": "user", "TEST_SVC_EVENTS_AUTH_PASSWORD": "pass"},
			preset:   credsScopeAuth{CredsFile: "/yaml.creds"},
			wantFile: "/yaml.creds",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg := tc.preset
			if err := Decode(&cfg); err != nil {
				t.Fatal(err)
			}
			if cfg.CredsFile != tc.wantFile || cfg.Creds != tc.wantCreds {
				t.Errorf("got file %q creds %q, want file %q creds %q", cfg.CredsFile, cfg.Creds, tc.wantFile, tc.wantCreds)
			}
		})
	}
}

func TestResolveCredsScopePrefixedPointerStruct(t *testing.T) {
	t.Setenv("TEST_NATS_CREDS", "global")
	t.Setenv("TEST_SVC_CACHE_CREDS_FILE", "/cache.creds")
	cfg := credsScopeConfig{Cache: &credsScopeCache{}}
	if err := Decode(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Cache.ProviderCacheCredsFile != "/cache.creds" || cfg.Cache.ProviderCacheCreds != "" {
		t.Errorf("got file %q creds %q", cfg.Cache.ProviderCacheCredsFile, cfg.Cache.ProviderCacheCreds)
	}
	// the events struct has only the global content, so it keeps it
	if cfg.Events.Creds != "global" {
		t.Errorf("events creds %q, want %q", cfg.Events.Creds, "global")
	}
}
