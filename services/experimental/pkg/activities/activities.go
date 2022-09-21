package activities

import (
	"encoding/json"
	"fmt"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/experimental/pkg/config"
	"net/http"
	"time"
)

// Storage defines how to store activities
type Storage interface {
	Add(Activity)
	List(uID string) []Activity
}

// Activity is used to build the activities timeline
type Activity struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	UserID   string    `json:"-"`
	DateTime time.Time `json:"dateTime"`
	Data     any       `json:"data"`
}

type activitiesService struct {
	logger  log.Logger
	storage Storage
}

// NewActivitiesService bootstraps the activities service
func NewActivitiesService(r chi.Router, es events.Stream, logger log.Logger, cfg config.Activities) error {
	var storage Storage
	switch cfg.Storage.Type {
	case "mem_storage":
		storage = NewMemStore(cfg.Storage.MemStore.Capacity)
	default:
		return fmt.Errorf("unknown activity storage: %s", cfg.Storage.Type)
	}

	svc := activitiesService{
		logger:  logger,
		storage: storage,
	}

	if err := svc.handleEvents(es); err != nil {
		return err
	}

	r.Get("/activities", svc.GetActivities)

	return nil
}

// GetActivities lists all available activities as json response.
func (s *activitiesService) GetActivities(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jm, err := json.Marshal(s.storage.List(u.Id.OpaqueId))
	if err != nil {
		s.logger.Error().Err(err).Msg("Could not read activities")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(jm); err != nil {
		s.logger.Error().Err(err).Msg("Could not write activities")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *activitiesService) handleEvents(bus events.Stream) error {
	evts := []events.Unmarshaller{
		events.VirusscanFinished{},
		events.ShareCreated{},
	}

	ch, err := events.Consume(bus, "experimental_activities", evts...)
	if err != nil {
		return err
	}

	go func(s *activitiesService, ch <-chan interface{}) {
		for e := range ch {
			switch ev := e.(type) {
			case events.VirusscanFinished:
				s.storage.Add(
					Activity{
						Type:     "VirusscanFinished",
						ID:       uuid.Must(uuid.NewV4()).String(),
						UserID:   ev.ExecutingUser.Id.OpaqueId,
						DateTime: ev.Scandate,
						Data: struct {
							Infected    bool                         `json:"infected"`
							Outcome     events.PostprocessingOutcome `json:"outcome"`
							Description string                       `json:"description"`
							Filename    string                       `json:"filename"`
							Scandate    time.Time                    `json:"scandate"`
							ResourceID  string                       `json:"resourceID"`
							Error       string                       `json:"error"`
						}{
							Infected:    ev.Infected,
							Outcome:     ev.Outcome,
							Description: ev.Description,
							Filename:    ev.Filename,
							Scandate:    ev.Scandate,
							ResourceID:  storagespace.FormatResourceID(*ev.ResourceID),
							Error:       ev.ErrorMsg,
						},
					},
				)
			}
		}
	}(s, ch)

	return nil
}
