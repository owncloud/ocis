package middleware

import (
	"net/http"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// The OpenID Connect Issuer URL
	OidcIssuer string
	// GatewayAPIClient is a reva gateway client
	GatewayAPIClient gatewayv1beta1.GatewayAPIClient
	// HttpClient is a http client
	HttpClient http.Client
}

// WithLogger provides a function to set the openid connect issuer option.
func WithOidcIssuer(val string) Option {
	return func(o *Options) {
		o.OidcIssuer = val
	}
}

// WithLogger provides a function to set the logger option.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// WithGatewayAPIClient provides a function to set the reva gateway client option.
func WithGatewayAPIClient(val gatewayv1beta1.GatewayAPIClient) Option {
	return func(o *Options) {
		o.GatewayAPIClient = val
	}
}

// HttpClient provides a function to set the http client option.
func WithHttpClient(val http.Client) Option {
	return func(o *Options) {
		o.HttpClient = val
	}
}
