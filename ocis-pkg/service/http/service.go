package http

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/web"
)

// Service simply wraps the go-micro web service.
type Service struct {
	web.Service
}

// NewService initializes a new http service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	sopts.Logger.Info().
		Str("transport", transport(sopts.TLSConfig)).
		Str("addr", sopts.Address).
		Msg("starting server")

	wopts := []web.Option{
		web.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		web.Version(sopts.Version),
		web.Address(sopts.Address),
		web.RegisterTTL(time.Second * 30),
		web.RegisterInterval(time.Second * 10),
		web.Context(sopts.Context),
		web.TLSConfig(sopts.TLSConfig),
		web.Handler(sopts.Handler),
		web.Flags(sopts.Flags...),
	}

	return Service{
		web.NewService(
			wopts...,
		),
	}
}

func transport(secure *tls.Config) string {
	if secure != nil {
		return "https"
	}

	return "http"
}
