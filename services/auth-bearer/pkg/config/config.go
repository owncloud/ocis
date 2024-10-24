package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"AUTH_BEARER_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the encoding of the user's group memberships in the reva access token. This reduces the token size, especially when users are members of a large number of groups." introductionVersion:"pre5.0"`

	OIDC OIDC `yaml:"oidc"`

	Context context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_BEARER_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_BEARER_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_BEARER_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_BEARER_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_BEARER_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"AUTH_BEARER_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"AUTH_BEARER_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"AUTH_BEARER_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"AUTH_BEARER_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;AUTH_BEARER_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"pre5.0"`
}

type OIDC struct {
	Issuer   string `yaml:"issuer" env:"OCIS_URL;OCIS_OIDC_ISSUER;AUTH_BEARER_OIDC_ISSUER" desc:"URL of the OIDC issuer. It defaults to URL of the builtin IDP." introductionVersion:"pre5.0"`
	Insecure bool   `yaml:"insecure" env:"OCIS_INSECURE;AUTH_BEARER_OIDC_INSECURE" desc:"Allow insecure connections to the OIDC issuer." introductionVersion:"pre5.0"`
	IDClaim  string `yaml:"id_claim" env:"AUTH_BEARER_OIDC_ID_CLAIM" desc:"Name of the claim, which holds the user identifier." introductionVersion:"pre5.0"`
	UIDClaim string `yaml:"uid_claim" env:"AUTH_BEARER_OIDC_UID_CLAIM" desc:"Name of the claim, which holds the UID." introductionVersion:"pre5.0"`
	GIDClaim string `yaml:"gid_claim" env:"AUTH_BEARER_OIDC_GID_CLAIM" desc:"Name of the claim, which holds the GID." introductionVersion:"pre5.0"`
}
