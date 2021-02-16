package http

import (
	"crypto/tls"
	"strings"
	"time"

	web "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3/server"
)

// Service simply wraps the go-micro web service.
type Service struct {
	server.Server
}

// NewService initializes a new http service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	sopts.Logger.Info().
		Str("transport", transport(sopts.TLSConfig)).
		Str("addr", sopts.Address).
		Msg("starting server")

	wopts := []server.Option{
		server.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		server.Version(sopts.Version),
		server.Address(sopts.Address),
		server.RegisterTTL(time.Second * 30),
		server.RegisterInterval(time.Second * 10),
		server.Context(sopts.Context),
		server.TLSConfig(sopts.TLSConfig),
		//server.Handler(sopts.Handler),
		//server.Flags(sopts.Flags...),
	}

	return Service{
		web.NewServer(
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
