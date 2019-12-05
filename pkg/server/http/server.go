package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-graph/pkg/flagset"
	"github.com/owncloud/ocis-graph/pkg/version"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func handleMe(writer http.ResponseWriter, req *http.Request) {
	displayName := "Alice"
	id := "1234-5678-9000-000"
	me := &msgraph.User{
		DisplayName: &displayName,
		GivenName:   &displayName,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &id,
			},
		},
	}

	js, err := json.Marshal(me)
	if err != nil {
		//p.srv.Logger().Errorf("owncloud-plugin: error encoding response as json %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(js)
}

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
	service.HandleFunc("/v1.0/me", handleMe)
	return service, nil
}
