package parser

import (
	"errors"

	"github.com/owncloud/ocis/extensions/idm/pkg/config"
	"github.com/owncloud/ocis/extensions/idm/pkg/config/defaults"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	_, err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
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
	if cfg.AdminUserID == "" {
		return shared.MissingAdminUserID(cfg.Service.Name)
	}

	if cfg.ServiceUserPasswords.Idm == "" {
		return shared.MissingServiceUserPassword(cfg.Service.Name, "IDM")
	}

	if cfg.ServiceUserPasswords.OcisAdmin == "" {
		return shared.MissingServiceUserPassword(cfg.Service.Name, "admin")
	}

	if cfg.ServiceUserPasswords.Idp == "" {
		return shared.MissingServiceUserPassword(cfg.Service.Name, "IDP")
	}

	if cfg.ServiceUserPasswords.Reva == "" {
		return shared.MissingServiceUserPassword(cfg.Service.Name, "REVA")
	}

	return nil
}
