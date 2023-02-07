package middleware

import (
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// GatewayAPIClient is a reva gateway client
	GatewayAPIClient gatewayv1beta1.GatewayAPIClient
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
