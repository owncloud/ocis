package http

import (
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/registry"

	mhttps "github.com/asim/go-micro/plugins/server/http/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

// Service simply wraps the go-micro web service.
type Service struct {
	micro.Service
}

// NewService initializes a new http service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	wopts := []micro.Option{
		micro.Server(mhttps.NewServer(server.TLSConfig(sopts.TLSConfig))),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(registry.GetRegistry()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
	}

	return Service{micro.NewService(wopts...)}
}
