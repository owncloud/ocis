package oidc

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// HTTPClient to use for requests
	HTTPClient *http.Client
	// Logger to use for logging, must be set
	Logger log.Logger
	// The OpenID Connect Issuer URL
	OidcIssuer string
	// JWKSOptions to use when retrieving keys
	JWKSOptions config.JWKS
	// AccessTokenVerifyMethod to use when verifying access tokens
	// TODO pass a function or interface to verify? an AccessTokenVerifier?
	AccessTokenVerifyMethod string
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
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

// WithAccessTokenVerifyMethod provides a function to set the accessTokenVerifyMethod option.
func WithAccessTokenVerifyMethod(val string) Option {
	return func(o *Options) {
		o.AccessTokenVerifyMethod = val
	}
}
func WithHTTPClient(val *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = val
	}
}
func WithJWKSOptions(val config.JWKS) Option {
	return func(o *Options) {
		o.JWKSOptions = val
	}
}
