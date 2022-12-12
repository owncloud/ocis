package search

import (
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
)

// HandleEvents listens to the needed events,
// it handles the whole resource indexing livecycle.
func HandleEvents(s Searcher, bus events.Consumer, logger log.Logger, cfg *config.Config) error {
	evts := []events.Unmarshaller{
		events.ItemTrashed{},
		events.ItemRestored{},
		events.ItemMoved{},
		events.ContainerCreated{},
		events.FileTouched{},
		events.FileVersionRestored{},
		//events.TagsAdded{},
		//events.TagsRemoved{},
	}

	if cfg.Events.AsyncUploads {
		// evts = append(evts, events.UploadReady{})
	} else {
		evts = append(evts, events.FileUploaded{})
	}

	ch, err := events.Consume(bus, "search", evts...)
	if err != nil {
		return err
	}

	if cfg.Events.NumConsumers == 0 {
		cfg.Events.NumConsumers = 1
	}

	spaceID := func(ref *provider.Reference) *provider.StorageSpaceId {
		return &provider.StorageSpaceId{
			OpaqueId: storagespace.FormatResourceID(
				provider.ResourceId{
					StorageId: ref.GetResourceId().GetStorageId(),
					SpaceId:   ref.GetResourceId().GetSpaceId(),
				},
			),
		}
	}

	indexSpaceDebouncer := NewSpaceDebouncer(time.Duration(cfg.Events.DebounceDuration)*time.Millisecond, func(id *provider.StorageSpaceId, userID *user.UserId) {
		if err := s.IndexSpace(id, userID); err != nil {
			logger.Error().Err(err).Interface("spaceID", id).Interface("userID", userID).Msg("error while indexing a space")
		}
	})

	for i := 0; i < cfg.Events.NumConsumers; i++ {
		go func(s Searcher, ch <-chan interface{}) {
			for e := range ch {
				logger.Debug().Interface("event", e).Msg("updating index")

				var err error

				switch ev := e.(type) {
				case events.ItemTrashed:
					s.TrashItem(ev.ID)
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.ItemMoved:
					s.MoveItem(ev.Ref, ev.Executant)
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.ItemRestored:
					s.RestoreItem(ev.Ref, ev.Executant)
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.ContainerCreated:
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.FileTouched:
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.FileVersionRestored:
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				//case events.TagsAdded:
				//	indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				//case events.TagsRemoved:
				//indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
				case events.FileUploaded:
					indexSpaceDebouncer.Debounce(spaceID(ev.Ref), ev.Executant)
					//case events.UploadReady:
					//indexSpaceDebouncer.Debounce(spaceID(ev.FileRef), ev.ExecutingUser.Id)
				}

				if err != nil {
					logger.Error().Err(err).Interface("event", e)
				}
			}
		}(
			s,
			ch,
		)
	}

	return nil
}
