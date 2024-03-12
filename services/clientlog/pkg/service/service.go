package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"go.opentelemetry.io/otel/trace"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/config"
)

// ClientlogService is the service responsible for user activities
type ClientlogService struct {
	log              log.Logger
	cfg              *config.Config
	gatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	registeredEvents map[string]events.Unmarshaller // ?
	tp               trace.TracerProvider
	tracer           trace.Tracer
	publisher        events.Publisher
	ch               <-chan events.Event
}

// NewClientlogService returns a clientlog service
func NewClientlogService(opts ...Option) (*ClientlogService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.Stream == nil {
		return nil, fmt.Errorf("need non nil stream (%v) to work properly", o.Stream)
	}

	ch, err := events.Consume(o.Stream, "clientlog", o.RegisteredEvents...)
	if err != nil {
		return nil, err
	}

	cl := &ClientlogService{
		log:              o.Logger,
		cfg:              o.Config,
		gatewaySelector:  o.GatewaySelector,
		registeredEvents: make(map[string]events.Unmarshaller),
		tp:               o.TraceProvider,
		tracer:           o.TraceProvider.Tracer("github.com/owncloud/ocis/services/clientlog/pkg/service"),
		publisher:        o.Stream,
		ch:               ch,
	}

	for _, e := range o.RegisteredEvents {
		typ := reflect.TypeOf(e)
		cl.registeredEvents[typ.String()] = e
	}

	return cl, nil
}

// Run runs the service
func (cl *ClientlogService) Run() error {
	for event := range cl.ch {
		cl.processEvent(event)
	}

	return nil
}

func (cl *ClientlogService) processEvent(event events.Event) {
	gwc, err := cl.gatewaySelector.Next()
	if err != nil {
		cl.log.Error().Err(err).Interface("event", event).Msg("error getting gatway client")
		return
	}

	ctx, err := utils.GetServiceUserContext(cl.cfg.ServiceAccount.ServiceAccountID, gwc, cl.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		cl.log.Error().Err(err).Interface("event", event).Msg("error authenticating service user")
		return
	}

	var (
		users  []string
		evType string
		data   interface{}
	)
	switch e := event.Event.(type) {
	default:
		err = errors.New("unhandled event")
	case events.UploadReady:
		evType = "postprocessing-finished"
		users, data, err = processFileEvent(ctx, e.FileRef, gwc)
	case events.ItemTrashed:
		evType = "item-trashed"

		var resp *provider.ListRecycleResponse
		resp, err = gwc.ListRecycle(ctx, &provider.ListRecycleRequest{
			Ref: e.Ref,
		})
		if err != nil || resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
			cl.log.Error().Err(err).Interface("event", event).Str("status code", resp.GetStatus().GetMessage()).Msg("error listing recycle")
			return
		}

		for _, item := range resp.GetRecycleItems() {
			if item.GetKey() == e.ID.GetOpaqueId() {

				data = FileEvent{
					ItemID: storagespace.FormatResourceID(*e.ID),
					// TODO: check with web if parentID is needed
					// ParentItemID: storagespace.FormatResourceID(*item.GetRef().GetResourceId()),
					SpaceID: storagespace.FormatStorageID(e.ID.GetStorageId(), e.ID.GetSpaceId()),
				}

				users, err = utils.GetSpaceMembers(ctx, e.ID.GetSpaceId(), gwc, utils.ViewerRole)
				break
			}
		}
	case events.ItemRestored:
		evType = "item-restored"
		users, data, err = processFileEvent(ctx, e.Ref, gwc)
	case events.ContainerCreated:
		evType = "folder-created"
		users, data, err = processFileEvent(ctx, e.Ref, gwc)
	case events.ItemMoved:
		// we are only interested in the rename case
		if !utils.ResourceIDEqual(e.OldReference.GetResourceId(), e.Ref.GetResourceId()) || e.Ref.GetPath() == e.OldReference.GetPath() {
			return
		}
		evType = "item-renamed"
		users, data, err = processFileEvent(ctx, e.Ref, gwc)
	case events.FileLocked:
		evType = "file-locked"
		users, data, err = processFileEvent(ctx, e.Ref, gwc)
	case events.FileUnlocked:
		evType = "file-unlocked"
		users, data, err = processFileEvent(ctx, e.Ref, gwc)
	}

	if err != nil {
		cl.log.Error().Err(err).Interface("event", event).Msg("error gathering members for event")
		return
	}

	// II) instruct sse service to send the information
	if err := cl.sendSSE(users, evType, data); err != nil {
		cl.log.Error().Err(err).Interface("userIDs", users).Str("eventid", event.ID).Msg("failed to store event for user")
		return
	}
}

func (cl *ClientlogService) sendSSE(userIDs []string, evType string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return events.Publish(context.Background(), cl.publisher, events.SendSSE{
		UserIDs: userIDs,
		Type:    evType,
		Message: b,
	})
}

func processFileEvent(ctx context.Context, ref *provider.Reference, gwc gateway.GatewayAPIClient) ([]string, FileEvent, error) {
	info, err := utils.GetResource(ctx, ref, gwc)
	if err != nil {
		return nil, FileEvent{}, err
	}

	data := FileEvent{
		ParentItemID: storagespace.FormatResourceID(*info.GetParentId()),
		ItemID:       storagespace.FormatResourceID(*info.GetId()),
		SpaceID:      storagespace.FormatStorageID(info.GetSpace().GetRoot().GetStorageId(), info.GetSpace().GetRoot().GetSpaceId()),
	}

	users, err := utils.GetSpaceMembers(ctx, info.GetSpace().GetId().GetOpaqueId(), gwc, utils.ViewerRole)
	return users, data, err
}
