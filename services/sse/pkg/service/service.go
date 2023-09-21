package service

import (
	"encoding/json"
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/r3labs/sse/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/sse/pkg/config"
)

// ServerSentEvent is the data structure sent by the sse service
type ServerSentEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// SSE defines implements the business logic for Service.
type SSE struct {
	c         *config.Config
	l         log.Logger
	m         *chi.Mux
	sse       *sse.Server
	evChannel <-chan events.Event
}

// NewSSE returns a service implementation for Service.
func NewSSE(c *config.Config, l log.Logger, ch <-chan events.Event, mux *chi.Mux) (SSE, error) {
	s := SSE{
		c:         c,
		l:         l,
		m:         mux,
		sse:       sse.New(),
		evChannel: ch,
	}
	mux.Route("/ocs/v2.php/apps/notifications/api/v1/notifications", func(r chi.Router) {
		r.Get("/sse", s.HandleSSE)
	})

	go s.ListenForEvents()

	return s, nil
}

// ServeHTTP fulfills Handler interface
func (s SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

// ListenForEvents listens for events
func (s SSE) ListenForEvents() {
	for e := range s.evChannel {
		switch ev := e.Event.(type) {
		default:
			s.l.Error().Interface("event", ev).Msg("unhandled event")
		case events.SendSSE:
			b, err := json.Marshal(ServerSentEvent{
				Type: ev.Type,
				Data: ev.Message,
			})
			if err != nil {
				s.l.Error().Interface("event", ev).Msg("cannot marshal event")
				continue
			}
			s.sse.Publish(ev.UserID, &sse.Event{
				Data: b,
			})
		}
	}
}

// HandleSSE is the GET handler for events
func (s SSE) HandleSSE(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		s.l.Error().Msg("sse: no user in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uid := u.GetId().GetOpaqueId()
	if uid == "" {
		s.l.Error().Msg("sse: user in context is broken")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream := s.sse.CreateStream(uid)
	stream.AutoReplay = false

	// add stream to URL
	q := r.URL.Query()
	q.Set("stream", uid)
	r.URL.RawQuery = q.Encode()

	s.sse.ServeHTTP(w, r)
}
