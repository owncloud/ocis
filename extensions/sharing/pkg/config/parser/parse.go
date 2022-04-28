package parser

import (
	"errors"
	"fmt"

	"github.com/owncloud/ocis/extensions/sharing/pkg/config"
	"github.com/owncloud/ocis/extensions/sharing/pkg/config/defaults"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"

	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
)

// ParseConfig loads accounts configuration from known paths.
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
	if cfg.PublicSharingDrivers.CS3.MachineAuthAPIKey == "" {
		return fmt.Errorf("machine auth api key for the cs3 public sharing driver is not set up properly, bailing out (%s)", cfg.Service.Name)
	}

	if cfg.UserSharingDrivers.CS3.MachineAuthAPIKey == "" {
		return fmt.Errorf("machine auth api key for the cs3 user sharing driver is not set up properly, bailing out (%s)", cfg.Service.Name)
	}

	return nil
}
