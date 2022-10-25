package svc

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/activities"
	"net/http"
	"os"

	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-micro/plugins/v4/events/natsjs"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	searchSvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/tags"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (Experimental, error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	evtsCfg := cfg.Events
	var rootCAPool *x509.CertPool
	if evtsCfg.TLSRootCACertificate != "" {
		rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
		if err != nil {
			return Experimental{}, err
		}

		rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
		if err != nil {
			return Experimental{}, err
		}
		evtsCfg.TLSInsecure = false
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: evtsCfg.TLSInsecure, //nolint:gosec
		RootCAs:            rootCAPool,
	}
	bus, err := server.NewNatsStream(
		natsjs.TLSConfig(tlsConf),
		natsjs.Address(evtsCfg.Endpoint),
		natsjs.ClusterID(evtsCfg.Cluster),
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

	if cfg.Activities.Enabled {
		err := activities.NewActivitiesService(r, bus, logger, cfg.Activities)
		if err != nil {
			return svc, err
		}
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
