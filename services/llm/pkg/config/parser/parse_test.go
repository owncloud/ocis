package parser

import (
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
)

func TestValidate_RequiresEndpoint(t *testing.T) {
	cfg := &config.Config{}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected an error when LLM.Endpoint is empty")
	}
}

func TestValidate_PassesWhenEndpointSet(t *testing.T) {
	cfg := validConfig()

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected no error when LLM.Endpoint is set, got: %v", err)
	}
}

func validConfig() *config.Config {
	return &config.Config{
		LLM: config.LLM{
			Endpoint:       "https://api.example.com/v1",
			MaxTokensLimit: 4096,
		},
	}
}

func TestValidate_RejectsInvalidValues(t *testing.T) {
	tests := map[string]func(*config.Config){
		"endpoint without scheme": func(c *config.Config) { c.LLM.Endpoint = "api.example.com/v1" },
		"endpoint wrong scheme":   func(c *config.Config) { c.LLM.Endpoint = "ftp://api.example.com" },
		"zero max tokens":         func(c *config.Config) { c.LLM.MaxTokensLimit = 0 },
		"negative timeout":        func(c *config.Config) { c.LLM.Timeout = -time.Second },
		"negative window":         func(c *config.Config) { c.RateLimit.Window = -time.Second },
		"store ttl below window": func(c *config.Config) {
			c.RateLimit.Window = time.Hour
			c.RateLimit.MaxRequests = 10
			c.Store.TTL = 10 * time.Minute
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := validConfig()
			mutate(cfg)
			if err := Validate(cfg); err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}
}

func TestValidate_AcceptsStoreTTLCoveringWindow(t *testing.T) {
	cfg := validConfig()
	cfg.RateLimit.Window = time.Minute
	cfg.RateLimit.MaxRequests = 20
	cfg.Store.TTL = 10 * time.Minute

	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
