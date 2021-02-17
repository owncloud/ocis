package http

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/registry"

	"github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/plugins/transport/tcp/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	mt "github.com/asim/go-micro/v3/transport"
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

	wopts := []micro.Option{
		micro.Server(http.NewServer(server.TLSConfig(sopts.TLSConfig))),
		micro.Registry(*registry.GetRegistry()),
		micro.Address(sopts.Address),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
	}

	return Service{micro.NewService(wopts...)}
}

func getTransport(tlscfg *tls.Config) mt.Transport {
	if tlscfg == nil {
		// return a default http transport
		return mt.NewHTTPTransport()
	}
	return tcp.NewTransport(mt.Secure(true), mt.TLSConfig(tlscfg))
}

func transport(secure *tls.Config) string {
	if secure != nil {
		return "https"
	}

	return "http"
}
