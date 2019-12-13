package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/micro/go-micro"
	ahandler "github.com/micro/go-micro/api/handler"
	aapi "github.com/micro/go-micro/api/handler/api"
	"github.com/micro/go-micro/api/resolver"
	rrmicro "github.com/micro/go-micro/api/resolver/micro"
	"github.com/micro/go-micro/api/router"
	regRouter "github.com/micro/go-micro/api/router/registry"
	httpapi "github.com/micro/go-micro/api/server/http"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/rs/zerolog/log"
)

// MicroGateway implements the handler interface
func MicroGateway(ctx context.Context, cancel context.CancelFunc, gr *run.Group, cfg *config.Config) error {
	// create the gateway service
	// var opts []server.Option
	var h http.Handler
	r := mux.NewRouter()
	h = r

	// return version and list of services
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "OPTIONS" {
			return
		}

		response := fmt.Sprintf(`{"version": "%s"}`, "[void]")
		w.Write([]byte(response))
	})

	// strip favicon.ico
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	var srvOpts []micro.Option
	srvOpts = append(srvOpts, micro.Name("com.micro.api"))

	// initialise service
	service := micro.NewService(srvOpts...)

	// resolver options
	ropts := []resolver.Option{
		resolver.WithNamespace("com.micro"),
		resolver.WithHandler("meta"),
	}

	// default resolver
	rr := rrmicro.NewResolver(ropts...)

	rt := regRouter.NewRouter(
		router.WithNamespace("com.micro"),
		router.WithHandler(aapi.Handler),
		router.WithResolver(rr),
		router.WithRegistry(service.Options().Registry),
	)
	ap := aapi.NewHandler(
		ahandler.WithNamespace("com.micro"),
		ahandler.WithRouter(rt),
		ahandler.WithService(service),
	)
	r.PathPrefix("/").Handler(ap)

	api := httpapi.NewServer(":8111")
	api.Init()
	api.Handle("/", h)

	gr.Add(func() error {
		return service.Run()
	}, func(err error) {
		log.Err(err)

		cancel()
	})

	// add it to the run group
	return nil
}

// func init() {
// 	fmt.Println("doing things")
// 	register.AddHandler(MicroGateway)
// }
