package parser

import (
	"errors"

	"github.com/owncloud/ocis/extensions/audit/pkg/config"
	"github.com/owncloud/ocis/extensions/audit/pkg/config/defaults"
	"github.com/owncloud/ocis/extensions/audit/pkg/logging"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"

	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
)

// ParseConfig loads accounts configuration from known paths.
func ParseConfig(cfg *config.Config) error {

	// cfg.ConfigFile can be set via env variable, therefore we need to do a environment variable run first
	if err := loadEnv(cfg); err != nil {
		return err
	}

	_, err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg.ConfigFile, cfg.ConfigFile != defaults.DefaultConfig().ConfigFile, cfg)
	if err != nil {
		logger := logging.Configure(cfg.Service.Name, &config.Log{})
		logger.Error().Err(err).Msg("couldn't find the specified config file")
		return err
	}

	defaults.EnsureDefaults(cfg)

	if err := loadEnv(cfg); err != nil {
		return err
	}

	defaults.Sanitize(cfg)

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
