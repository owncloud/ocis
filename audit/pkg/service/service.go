package svc

import (
	"github.com/go-chi/chi"
	"github.com/owncloud/ocis/audit/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Service defines the extension handlers.
type Service interface {
	ListenForEvents()
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Audit{
		logger: options.Logger,
		config: options.Config,
		mux:    m,
	}

	go svc.ListenForEvents()
	return svc
}

// Audit defines implements the business logic for Service.
type Audit struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
}

// ListenForEvents hooks into event queue and logs interesting events
func (g Audit) ListenForEvents() {
	log := g.logger
	ch, err := startConsumer(g.config.Eventstream, log)
	if err != nil {
		log.Fatal().Err(err).Msg("can't listen for events")
		return
	}

	startAuditLogger(g.config.Auditlog, ch, log)
}
