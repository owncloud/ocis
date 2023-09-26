package parser

import (
	"errors"
	"fmt"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/logging"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
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
	logger := logging.Configure(cfg.Service.Name, cfg.Log)

	if cfg.Notifications.SMTP.Host != "" {
		switch cfg.Notifications.SMTP.Encryption {
		case "tls":
			logger.Warn().Msg("The smtp_encryption value 'tls' is deprecated. Please use the value 'starttls' instead.")
		case "ssl":
			logger.Warn().Msg("The smtp_encryption value 'ssl' is deprecated. Please use the value 'ssltls' instead.")
		case "starttls", "ssltls", "none":
			break
		default:
			return fmt.Errorf(
				"unknown value '%s' for 'smtp_encryption' in service %s. Allowed values are 'starttls', 'ssltls' or 'none'",
				cfg.Notifications.SMTP.Encryption, cfg.Service.Name,
			)
		}
	}
	return nil
}
