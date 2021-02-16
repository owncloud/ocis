package http

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/owncloud/ocis/ocis-pkg/registry"
)

// Service simply wraps the go-micro web service.
type Service struct {
	micro.Service
}

// NewService initializes a new http service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	sopts.Logger.Info().
		Str("transport", transport(sopts.TLSConfig)).
		Str("addr", sopts.Address).
		Msg("starting server")

	r := registry.GetRegistry()
	srv := http.NewServer(
		// Broker
		server.Registry(*r),
		server.Address(sopts.Address),
		server.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		// Id
		server.Version(sopts.Version),
	)

	srv.Handle(srv.NewHandler(sopts.Handler))

	wopts := []micro.Option{
		micro.Server(srv),
		//micro.TLSConfig(sopts.TLSConfig),
		micro.Address(sopts.Address),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
	}

	return Service{
		micro.NewService(
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
