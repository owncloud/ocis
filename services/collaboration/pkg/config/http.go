package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"COLLABORATION_HTTP_ADDR" desc:"The external address of the collaboration service wihout a leading scheme. Either use an IP address or a hostname (127.0.0.1:9301 or wopi.private.prv). The configured 'Scheme' in another envvar will be used to finally build the public URL along with this address."`
	BindAddr  string                `yaml:"bindaddr" env:"COLLABORATION_HTTP_BINDADDR" desc:"The bind address of the HTTP service. Use '<ip-address>:<port>', for example, '127.0.0.1:9301' or '0.0.0.0:9301'."`
	Namespace string                `yaml:"-"`
	Scheme    string                `yaml:"scheme" env:"COLLABORATION_HTTP_SCHEME" desc:"The scheme to use for the HTTP address, which is either 'http' or 'https'."`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
}
