package proofkeys

import (
	"github.com/rs/zerolog"
)

// VerifyOption defines a single option function.
type VerifyOption func(o *VerifyOptions)

// VerifyOptions defines the available options for the Verify function.
type VerifyOptions struct {
	Logger *zerolog.Logger
}

// newOptions initializes the available default options.
func newOptions(opts ...VerifyOption) VerifyOptions {
	defaultLog := zerolog.Nop()
	opt := VerifyOptions{
		//Logger: log.NopLogger(), // use a NopLogger by default
		Logger: &defaultLog, // use a NopLogger by default
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// VerifyWithLogger provides a function to set the Logger option.
func VerifyWithLogger(val *zerolog.Logger) VerifyOption {
	return func(o *VerifyOptions) {
		o.Logger = val
	}
}
