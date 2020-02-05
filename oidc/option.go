package oidc

import (
	"github.com/owncloud/ocis-pkg/v2/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// Endpoint is the OpenID Connect provider URL
	Endpoint string
	// Realm to use in the WWW-Authenticate header, defaults to Endpoint
	Realm string
	// SigningAlgs to use when verifying jwt signatures, defaults to "RS256" & "PS256"
	SigningAlgs []string
	// Insecure can be used to disable http certificate checks
	Insecure bool
}

// Logger provides a function to set the logger option.
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Endpoint provides a function to set the endpoint option.
func Endpoint(e string) Option {
	return func(o *Options) {
		o.Endpoint = e
	}
}

// Realm provides a function to set the realm option.
func Realm(r string) Option {
	return func(o *Options) {
		o.Realm = r
	}
}

// SigningAlgs provides a function to set the signing algorithms option.
func SigningAlgs(sa []string) Option {
	return func(o *Options) {
		o.SigningAlgs = sa
	}
}

// Insecure provides a function to set the insecure option.
func Insecure(i bool) Option {
	return func(o *Options) {
		o.Insecure = i
	}
}
