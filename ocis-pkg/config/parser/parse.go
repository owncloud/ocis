package parser

import (
	"errors"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// ParseConfig loads ocis configuration.
func ParseConfig(cfg *config.Config) error {
	_, err := config.BindSourcesToStructs("ocis", cfg)
	if err != nil {
		return err
	}

	if cfg.Commons == nil {
		cfg.Commons = &shared.Commons{}
	}

	if cfg.Log != nil {
		cfg.Commons.Log = &shared.Log{
			Level:  cfg.Log.Level,
			Pretty: cfg.Log.Pretty,
			Color:  cfg.Log.Color,
			File:   cfg.File,
		}
	} else {
		cfg.Commons.Log = &shared.Log{}
		cfg.Log = &shared.Log{}
	}

	if cfg.Tracing != nil {
		cfg.Commons.Tracing = &shared.Tracing{
			Enabled:   cfg.Tracing.Enabled,
			Type:      cfg.Tracing.Type,
			Endpoint:  cfg.Tracing.Endpoint,
			Collector: cfg.Tracing.Collector,
		}
	} else {
		cfg.Commons.Tracing = &shared.Tracing{}
		cfg.Tracing = &shared.Tracing{}
	}

	if cfg.TokenManager != nil {
		cfg.Commons.TokenManager = cfg.TokenManager
	} else {
		cfg.Commons.TokenManager = &shared.TokenManager{}
		cfg.TokenManager = cfg.Commons.TokenManager
	}

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	return nil
}
