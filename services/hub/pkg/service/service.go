package service

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/go-chi/chi/v5"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	opkgm "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
)

// Service defines the service handlers.
type Service struct {
	m *chi.Mux
}

// New returns a service implementation for Service.
func New(cfg *config.Config) Service {
	m := chi.NewMux()
	m.Use(
		opkgm.ExtractAccountUUID(
			account.JWTSecret(cfg.TokenManager.JWTSecret),
		),
	)

	s, err := NewSSE(cfg)
	if err != nil {
		log.Fatal("cant initiate sse", err)
	}

	ch, err := eventsConsumer(cfg.Events)
	if err != nil {
		log.Fatal("cant consume events", err)
	}

	go s.ListenForEvents(ch)

	m.Route("/hub", func(r chi.Router) {
		r.Route("/sse", func(r chi.Router) {
			r.Get("/", s.ServeHTTP)
		})
	})

	svc := Service{
		m: m,
	}

	return svc
}

func eventsConsumer(evtsCfg config.Events) (<-chan events.Event, error) {
	var tlsConf *tls.Config
	if evtsCfg.EnableTLS {
		var rootCAPool *x509.CertPool
		if evtsCfg.TLSRootCACertificate != "" {
			rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
			if err != nil {
				return nil, err
			}

			rootCAPool, err = crypto.NewCertPoolFromPEM(rootCrtFile)
			if err != nil {
				return nil, err
			}
			evtsCfg.TLSInsecure = false
		}

		tlsConf = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: evtsCfg.TLSInsecure, //nolint:gosec
			RootCAs:            rootCAPool,
		}
	}
	client, err := stream.Nats(
		natsjs.TLSConfig(tlsConf),
		natsjs.Address(evtsCfg.Endpoint),
		natsjs.ClusterID(evtsCfg.Cluster),
	)
	if err != nil {
		return nil, err
	}

	evts, err := events.Consume(client, "hub", events.UploadReady{})
	if err != nil {
		return nil, err
	}

	return evts, nil
}

// ServeHTTP implements the Service interface.
func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}
