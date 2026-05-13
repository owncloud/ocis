package parser

import (
	"errors"
	"fmt"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/defaults"
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

// Validate validates the configuration
func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	if cfg.MountID == "" {
		return fmt.Errorf("the storage-kiteworks mount ID has not been configured for %s. "+
			"Make sure to set the STORAGE_KITEWORKS_MOUNT_ID environment variable or configure it in the config file.",
			cfg.Service.Name)
	}

	if cfg.Driver.Endpoint == "" {
		return fmt.Errorf("the Kiteworks endpoint has not been configured for %s. "+
			"Make sure to set the STORAGE_KITEWORKS_ENDPOINT environment variable or configure it in the config file.",
			cfg.Service.Name)
	}

	return nil
}
