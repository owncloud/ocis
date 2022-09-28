package svc

import (
	"net/http"

	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	searchSvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/activities"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/tags"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (Experimental, error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config
	es, err := server.NewNatsStream(
		natsjs.Address(cfg.Events.Endpoint),
		natsjs.ClusterID(cfg.Events.Cluster),
	)
	if err != nil {
		return Experimental{}, err
	}

	r := chi.NewRouter()
	r.Use(options.Middleware...)

	svc := Experimental{
		r: r,
	}

	tags.NewTagsService(
		r,
		searchSvc.NewSearchProviderService("com.owncloud.api.search", grpc.DefaultClient()),
		logger,
	)

	if err := activities.NewActivitiesService(r, es, logger, cfg.Activities); err != nil {
		return svc, err
	}

	r.Mount(cfg.HTTP.Root, r)

	return svc, nil
}

// Experimental implements the business logic for Service.
type Experimental struct {
	r *chi.Mux
}

// ServeHTTP implements the Service interface.
func (s Experimental) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
