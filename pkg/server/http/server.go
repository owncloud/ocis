package http

import (
	"time"

	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis-graph/pkg/version"
)

func Server(opts ...Option) (web.Service, error) {
	options := newOptions(opts...)
	log.Infof("Server [http] listening on [%s]", options.Config.HTTP.Addr)

	// &cli.StringFlag{
	// 	Name:        "http-addr",
	// 	Value:       "0.0.0.0:8380",
	// 	Usage:       "Address to bind http server",
	// 	EnvVar:      "GRAPH_HTTP_ADDR",
	// 	Destination: &cfg.HTTP.Addr,
	// },

	service := web.NewService(
		web.Name("go.micro.web.graph"),
		web.Version(version.String),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*10),
		web.Context(options.Context),
		web.Flags(append(
			flagset.RootWithConfig(config.New()),
			flagset.ServerWithConfig(config.New())...,
		)...),
	)

	service.Init()
	return service, nil
}
