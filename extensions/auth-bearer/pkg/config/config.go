package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Log             *Log     `yaml:"log"`
	Debug           Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"AUTH_BEARER_SKIP_USER_GROUPS_IN_TOKEN"`

	OIDC OIDC `yaml:"oidc"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_BEARER_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_BEARER_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_BEARER_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_BEARER_TRACING_COLLECTOR"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_BEARER_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_BEARER_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_BEARER_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_BEARER_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_BEARER_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"AUTH_BEARER_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"AUTH_BEARER_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"AUTH_BEARER_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"AUTH_BEARER_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"AUTH_BEARER_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type OIDC struct {
	Issuer   string `yaml:"issuer" env:"OCIS_URL;OCIS_OIDC_ISSUER;AUTH_BEARER_OIDC_ISSUER"`
	Insecure bool   `yaml:"insecure" env:"OCIS_INSECURE;AUTH_BEARER_OIDC_INSECURE"`
	IDClaim  string `yaml:"id_claim" env:"AUTH_BEARER_OIDC_ID_CLAIM"`
	UIDClaim string `yaml:"uid_claim" env:"AUTH_BEARER_OIDC_UID_CLAIM"`
	GIDClaim string `yaml:"gid_claim" env:"AUTH_BEARER_OIDC_GID_CLAIM"`
}
