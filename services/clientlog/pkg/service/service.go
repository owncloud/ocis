package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
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
		cl.log.Error().Err(err).Interface("event", event).Msg("error getting gateway client")
		return
	}

	ctx, err := utils.GetServiceUserContextWithContext(context.Background(), gwc, cl.cfg.ServiceAccount.ServiceAccountID, cl.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		cl.log.Error().Err(err).Interface("event", event).Msg("error authenticating service user")
		return
	}

	var (
		users  []string
		evType string
		data   interface{}
	)

	p := func(typ string, ref *provider.Reference) {
		evType = typ
		users, data, err = processFileEvent(ctx, ref, gwc, event.InitiatorID)
	}

	switch e := event.Event.(type) {
	default:
		err = errors.New("unhandled event")
	case events.UploadReady:
		p("postprocessing-finished", e.FileRef)
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
					SpaceID:     storagespace.FormatStorageID(e.ID.GetStorageId(), e.ID.GetSpaceId()),
					InitiatorID: event.InitiatorID,
				}

				gwc, err = cl.gatewaySelector.Next()
				if err != nil {
					cl.log.Error().Err(err).Interface("event", event).Msg("error getting gateway client")
					return
				}
				users, err = utils.GetSpaceMembers(ctx, e.ID.GetSpaceId(), gwc, utils.ViewerRole)
				break
			}
		}
	case events.ItemRestored:
		p("item-restored", e.Ref)
	case events.ContainerCreated:
		p("folder-created", e.Ref)
	case events.ItemMoved:
		// we send a dedicated event in case the item was only renamed
		if utils.ResourceIDEqual(e.OldReference.GetResourceId(), e.Ref.GetResourceId()) || e.Ref.GetPath() == e.OldReference.GetPath() {
			p("item-renamed", e.Ref)
		} else {
			p("item-moved", e.Ref)
		}
	case events.FileLocked:
		p("file-locked", e.Ref)
	case events.FileUnlocked:
		p("file-unlocked", e.Ref)
	case events.FileTouched:
		p("file-touched", e.Ref)
	case events.SpaceShared:
		r, _ := storagespace.ParseReference(e.ID.GetOpaqueId())
		p("space-member-added", &r)
	case events.SpaceUnshared:
		r, _ := storagespace.ParseReference(e.ID.GetOpaqueId())
		p("space-member-removed", &r)
		users, err = addSharees(ctx, users, gwc, e.GranteeUserID, e.GranteeGroupID)
	case events.ShareCreated:
		p("share-created", &provider.Reference{ResourceId: e.ItemID})
		users, err = addSharees(ctx, users, gwc, e.GranteeUserID, e.GranteeGroupID)
	case events.ShareUpdated:
		p("share-updated", &provider.Reference{ResourceId: e.ItemID})
		users, err = addSharees(ctx, users, gwc, e.GranteeUserID, e.GranteeGroupID)
	case events.ShareRemoved:
		p("share-removed", &provider.Reference{ResourceId: e.ItemID})
		users, err = addSharees(ctx, users, gwc, e.GranteeUserID, e.GranteeGroupID)
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

func processFileEvent(ctx context.Context, ref *provider.Reference, gwc gateway.GatewayAPIClient, initiatorid string) ([]string, FileEvent, error) {
	info, err := utils.GetResource(ctx, ref, gwc)
	if err != nil {
		return nil, FileEvent{}, err
	}

	data := FileEvent{
		ParentItemID: storagespace.FormatResourceID(*info.GetParentId()),
		ItemID:       storagespace.FormatResourceID(*info.GetId()),
		SpaceID:      storagespace.FormatStorageID(info.GetSpace().GetRoot().GetStorageId(), info.GetSpace().GetRoot().GetSpaceId()),
		InitiatorID:  initiatorid,
		Etag:         info.GetEtag(),
	}

	users, err := utils.GetSpaceMembers(ctx, info.GetSpace().GetId().GetOpaqueId(), gwc, utils.ViewerRole)
	return users, data, err
}

// adds userid to users slice or gets members of groupid and adds them to users slice
func addSharees(ctx context.Context, users []string, gwc gateway.GatewayAPIClient, uid *user.UserId, gid *group.GroupId) ([]string, error) {
	if uid != nil {
		return append(users, uid.GetOpaqueId()), nil
	}
	us, err := utils.GetGroupMembers(ctx, gid.GetOpaqueId(), gwc)
	return append(users, us...), err
}
