package parser

import (
	"errors"
	"fmt"
	"net/url"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	ocisdefaults "github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config/defaults"
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
	if cfg.Wopi.Secret == "" {
		return shared.MissingWOPISecretError(cfg.Service.Name)
	}
	url, err := url.Parse(cfg.Wopi.WopiSrc)
	if err != nil {
		return fmt.Errorf("The WOPI Src has not been set properly in your config for %s. "+
			"Make sure your %s config contains the proper values "+
			"(e.g. by running ocis init or setting it manually in "+
			"the config/corresponding environment variable): %s",
			cfg.Service.Name, ocisdefaults.BaseConfigPath(), err.Error())
	}
	if url.Path != "" {
		return fmt.Errorf("The WOPI Src must not contain a path in your config for %s. "+
			"Make sure your %s config contains the proper values "+
			"(e.g. by running ocis init or setting it manually in "+
			"the config/corresponding environment variable)",
			cfg.Service.Name, ocisdefaults.BaseConfigPath())
	}

	return nil
}
