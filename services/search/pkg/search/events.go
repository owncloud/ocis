package search

import (
	"context"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"google.golang.org/grpc/metadata"
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
		events.TagsAdded{},
		events.TagsRemoved{},
		events.SpaceRenamed{},
	}

	if cfg.Events.AsyncUploads {
		evts = append(evts, events.UploadReady{})
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

	getSpaceID := func(ref *provider.Reference) *provider.StorageSpaceId {
		return &provider.StorageSpaceId{
			OpaqueId: storagespace.FormatResourceID(
				&provider.ResourceId{
					StorageId: ref.GetResourceId().GetStorageId(),
					SpaceId:   ref.GetResourceId().GetSpaceId(),
				},
			),
		}
	}

	indexSpaceDebouncer := NewSpaceDebouncer(time.Duration(cfg.Events.DebounceDuration)*time.Millisecond, func(ctx context.Context, id *provider.StorageSpaceId) {
		if err := s.IndexSpace(ctx, id); err != nil {
			logger.Error().Err(err).Interface("spaceID", id).Msg("error while indexing a space")
		}
	})

	// trace provider is available here, otherwise the search service should have crashed
	tp, _ := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)

	for i := 0; i < cfg.Events.NumConsumers; i++ {
		go func(s Searcher, ch <-chan events.Event) {
			for event := range ch {
				e := event
				go func() {
					evCtx := context.Background()
					tracedEvCtx, span := events.TraceEventConsumer(evCtx, tp, e)
					tracedEvCtx = metadata.NewOutgoingContext(tracedEvCtx, e.ExtraInfo)
					defer span.End()

					logger.Debug().Interface("event", e).Msg("updating index")

					switch ev := e.Event.(type) {
					case events.ItemTrashed:
						s.TrashItem(tracedEvCtx, ev.ID)
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.ItemMoved:
						s.MoveItem(tracedEvCtx, ev.Ref)
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.ItemRestored:
						s.RestoreItem(tracedEvCtx, ev.Ref)
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.ContainerCreated:
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.FileTouched:
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.FileVersionRestored:
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.TagsAdded:
						s.UpdateTags(tracedEvCtx, ev.Ref)
					case events.TagsRemoved:
						s.UpdateTags(tracedEvCtx, ev.Ref)
					case events.FileUploaded:
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.Ref))
					case events.UploadReady:
						if ev.Failed {
							return
						}
						indexSpaceDebouncer.Debounce(tracedEvCtx, getSpaceID(ev.FileRef))
					case events.SpaceRenamed:
						indexSpaceDebouncer.Debounce(tracedEvCtx, ev.ID)
					}
				}()
			}
		}(
			s,
			ch,
		)
	}

	return nil
}
