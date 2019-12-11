package oidc

import (
	"github.com/owncloud/ocis-pkg/log"
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
	// Audience to use when checking jwt based tokens
	Audience string
	// SigningAlgs to use when verifying jwt signatures, defaults to "RS256" & "PS256"
	SigningAlgs []string
	// ClientId to use as username for basic auth against the introspection endpoint
	ClientID string
	// ClientSecret to use as password for basic auth against the introspection endpoint
	ClientSecret string
	// Insecure can be used to disable http certificate checks
	Insecure bool
	// SkipCheck can be used to further reduce security. Fix that!
	SkipChecks bool
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

// Audience provides a function to set the audience option.
func Audience(a string) Option {
	return func(o *Options) {
		o.Audience = a
	}
}

// SigningAlgs provides a function to set the signing algorithms option.
func SigningAlgs(sa []string) Option {
	return func(o *Options) {
		o.SigningAlgs = sa
	}
}

// ClientID provides a function to set the client id option.
func ClientID(ci string) Option {
	return func(o *Options) {
		o.ClientID = ci
	}
}

// ClientSecret provides a function to set the client secret option.
func ClientSecret(cs string) Option {
	return func(o *Options) {
		o.ClientSecret = cs
	}
}

// Insecure provides a function to set the insecure option.
func Insecure(i bool) Option {
	return func(o *Options) {
		o.Insecure = i
	}
}

// SkipChecks provides a function to set the ready option.
func SkipChecks(sc bool) Option {
	return func(o *Options) {
		o.SkipChecks = sc
	}
}
