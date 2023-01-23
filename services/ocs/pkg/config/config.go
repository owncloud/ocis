package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing    *Tracing    `yaml:"tracing"`
	Log        *Log        `yaml:"log"`
	CacheStore *CacheStore `yaml:"cache_store"`
	Debug      Debug       `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	TokenManager *TokenManager `yaml:"token_manager"`

	Context context.Context `yaml:"-"`
}
