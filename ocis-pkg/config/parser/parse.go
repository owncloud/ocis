package parser

import (
	"errors"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// ParseConfig loads ocis configuration.
func ParseConfig(cfg *config.Config) error {
	// cfg.ConfigFile can be set via env variable, therefore we need to do a environment variable run first
	if err := loadEnv(cfg); err != nil {
		return err
	}

	_, err := config.BindSourcesToStructs("ocis", cfg.ConfigFile, cfg.ConfigFile != config.DefaultConfig().ConfigFile, cfg)
	if err != nil {
		return err
	}

	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &shared.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &shared.Log{}
	}

	if err := loadEnv(cfg); err != nil {
		return err
	}

	return nil
}

func loadEnv(cfg *config.Config) error {
	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}
	return nil
}
