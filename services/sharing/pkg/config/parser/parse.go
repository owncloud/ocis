package parser

import (
	"errors"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/sharing/pkg/config"
	"github.com/owncloud/ocis/v2/services/sharing/pkg/config/defaults"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
)

const (
	_backendCS3 = "cs3"
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

func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	if cfg.PublicSharingDriver == _backendCS3 && cfg.PublicSharingDrivers.CS3.SystemUserAPIKey == "" {
		return shared.MissingSystemUserApiKeyError(cfg.Service.Name)
	}

	if cfg.PublicSharingDriver == _backendCS3 && cfg.PublicSharingDrivers.CS3.SystemUserID == "" {
		return shared.MissingSystemUserID(cfg.Service.Name)
	}

	if cfg.UserSharingDriver == _backendCS3 && cfg.UserSharingDrivers.CS3.SystemUserAPIKey == "" {
		return shared.MissingSystemUserApiKeyError(cfg.Service.Name)
	}

	if cfg.UserSharingDriver == _backendCS3 && cfg.UserSharingDrivers.CS3.SystemUserID == "" {
		return shared.MissingSystemUserID(cfg.Service.Name)
	}

	return nil
}
