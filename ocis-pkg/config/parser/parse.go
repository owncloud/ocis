package parser

import (
	"errors"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// ParseConfig loads the ocis configuration and
// copies applicable parts into the commons part, from
// where the extensions can copy it into their own config
func ParseConfig(cfg *config.Config) error {
	_, err := config.BindSourcesToStructs("ocis", cfg)
	if err != nil {
		return err
	}

	EnsureDefaults(cfg)

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	EnsureCommons(cfg)

	return Validate(cfg)
}

// EnsureDefaults, ensures that all pointers in the
// oCIS config (not the extensions configs) are initialized
func EnsureDefaults(cfg *config.Config) {
	if cfg.Tracing == nil {
		cfg.Tracing = &shared.Tracing{}
	}
	if cfg.Log == nil {
		cfg.Log = &shared.Log{}
	}
	if cfg.TokenManager == nil {
		cfg.TokenManager = &shared.TokenManager{}
	}
}

// EnsureCommons copies applicable parts of the oCIS config into the commons part
func EnsureCommons(cfg *config.Config) {
	// ensure the commons part is initialized
	if cfg.Commons == nil {
		cfg.Commons = &shared.Commons{}
	}

	// copy config to the commons part if set
	if cfg.Log != nil {
		cfg.Commons.Log = &shared.Log{
			Level:  cfg.Log.Level,
			Pretty: cfg.Log.Pretty,
			Color:  cfg.Log.Color,
			File:   cfg.File,
		}
	} else {
		cfg.Commons.Log = &shared.Log{}
	}

	// copy tracing to the commons part if set
	if cfg.Tracing != nil {
		cfg.Commons.Tracing = &shared.Tracing{
			Enabled:   cfg.Tracing.Enabled,
			Type:      cfg.Tracing.Type,
			Endpoint:  cfg.Tracing.Endpoint,
			Collector: cfg.Tracing.Collector,
		}
	} else {
		cfg.Commons.Tracing = &shared.Tracing{}
	}

	// copy token manager to the commons part if set
	if cfg.TokenManager != nil {
		cfg.Commons.TokenManager = cfg.TokenManager
	} else {
		cfg.Commons.TokenManager = &shared.TokenManager{}
	}

	// copy machine auth api key to the commons part if set
	if cfg.MachineAuthAPIKey != "" {
		cfg.Commons.MachineAuthAPIKey = cfg.MachineAuthAPIKey
	}

	// copy transfer secret to the commons part if set
	if cfg.TransferSecret != "" {
		cfg.Commons.TransferSecret = cfg.TransferSecret
	}

	// copy metadata user id to the commons part if set
	if cfg.SystemUserID != "" {
		cfg.Commons.SystemUserID = cfg.SystemUserID
	}
}

func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError("ocis")
	}

	if cfg.TransferSecret == "" {
		return shared.MissingRevaTransferSecretError("ocis")
	}

	if cfg.MachineAuthAPIKey == "" {
		return shared.MissingMachineAuthApiKeyError("ocis")
	}

	if cfg.SystemUserID == "" {
		return shared.MissingSystemUserID("ocis")
	}

	return nil
}
