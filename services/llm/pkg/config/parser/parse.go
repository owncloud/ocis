package parser

import (
	"errors"
	"fmt"
	"net/url"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config/defaults"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

// Validate validates the config.
func Validate(cfg *config.Config) error {
	if cfg.LLM.Endpoint == "" {
		return errors.New("llm: LLM_ENDPOINT is required")
	}
	u, err := url.Parse(cfg.LLM.Endpoint)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("llm: LLM_ENDPOINT must be an http(s) URL, got %q", cfg.LLM.Endpoint)
	}
	if cfg.LLM.MaxTokensLimit <= 0 {
		return errors.New("llm: LLM_MAX_TOKENS_LIMIT must be greater than 0")
	}
	if cfg.LLM.Timeout < 0 {
		return errors.New("llm: LLM_TIMEOUT must not be negative")
	}
	if cfg.RateLimit.Window < 0 || cfg.RateLimit.MaxRequests < 0 {
		return errors.New("llm: LLM_RATE_LIMIT_WINDOW and LLM_RATE_LIMIT_MAX_REQUESTS must not be negative")
	}
	// Rate-limit records are expired by the store; a store TTL shorter than
	// the window would silently shrink the effective window.
	if cfg.RateLimit.Window > 0 && cfg.RateLimit.MaxRequests > 0 &&
		cfg.Store.TTL > 0 && cfg.Store.TTL < cfg.RateLimit.Window {
		return fmt.Errorf("llm: LLM_STORE_TTL (%s) must not be shorter than LLM_RATE_LIMIT_WINDOW (%s)", cfg.Store.TTL, cfg.RateLimit.Window)
	}
	return nil
}
