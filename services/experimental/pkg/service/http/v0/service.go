package svc

import (
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/tags"
	"net/http"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Experimental {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Experimental{
		mux: m,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		tags.NewTagsService(
			r,
			searchsvc.NewSearchProviderService("com.owncloud.api.search", grpc.DefaultClient),
			options.Logger,
		)
	})

	return svc
}

// Experimental implements the business logic for Service.
type Experimental struct {
	mux *chi.Mux
}

// ServeHTTP implements the Service interface.
func (s Experimental) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
