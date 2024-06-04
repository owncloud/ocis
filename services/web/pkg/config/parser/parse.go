package parser

import (
	"errors"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/web/pkg/config"
	"github.com/owncloud/ocis/v2/services/web/pkg/config/defaults"
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

	// apps are a special case, as they are not part of the main config, but are loaded from a separate config file
	err = ociscfg.BindSourcesToStructs("apps", &cfg.Apps)
	if err != nil {
		return err
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

// Validate validates the configuration
func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	// deprecation: migration requested
	// check if the config still uses the deprecated asset path, if so,
	// log a warning and copy the value to the setting that is actually used
	// this is to ensure a smooth transition from the old to the new core asset path (pre 5.1 to 5.1)
	if cfg.Asset.DeprecatedPath != "" {
		if cfg.Asset.CorePath == "" {
			cfg.Asset.CorePath = cfg.Asset.DeprecatedPath
		}

		// message should be logged to the console,
		// do not use a logger here because the message MUST be visible independent of the log level
		log.Deprecation("WEB_ASSET_PATH is deprecated and will be removed in the future. Use WEB_ASSET_CORE_PATH instead.")
	}

	return nil
}
