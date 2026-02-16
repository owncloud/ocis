package parser

import (
	"errors"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config/defaults"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
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

// Validate validates the config
func Validate(cfg *config.Config) error {
	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}

	if cfg.ServiceAccount.ID == "" {
		return shared.MissingServiceAccountID(cfg.Service.Name)
	}
	if cfg.ServiceAccount.Secret == "" {
		return shared.MissingServiceAccountSecret(cfg.Service.Name)
	}

	return nil
}
