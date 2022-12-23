package config

import (
	"context"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons      *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	HTTP         HTTP            `yaml:"http"`
	Service      Service         `yaml:"-"`
	TokenManager *TokenManager   `yaml:"token_manager"`
	Context      context.Context `yaml:"-"`
}

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"-"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"HUB_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"HUB_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;HUB_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}
