package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/clientlog/pkg/config"
	"go.opentelemetry.io/otel/trace"
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
		info, err := utils.GetResource(ctx, e.FileRef, gwc)
		if err != nil {
			cl.log.Error().Err(err).Interface("event", event).Msg("error getting resource")
			return
		}

		evType = "postprocessing-finished"
		data = FileReadyEvent{
			ParentItemID: storagespace.FormatResourceID(*info.GetParentId()),
			ItemID:       storagespace.FormatResourceID(*info.GetId()),
		}

		users, err = utils.GetSpaceMembers(ctx, info.GetSpace().GetId().GetOpaqueId(), gwc, utils.ViewerRole)
	}

	if err != nil {
		cl.log.Info().Err(err).Interface("event", event).Msg("error gathering members for event")
		return
	}

	// II) instruct sse service to send the information
	for _, id := range users {
		if err := cl.sendSSE(id, evType, data); err != nil {
			cl.log.Error().Err(err).Str("userID", id).Str("eventid", event.ID).Msg("failed to store event for user")
			return
		}
	}
}

func (cl *ClientlogService) sendSSE(userid string, evType string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return events.Publish(context.Background(), cl.publisher, events.SendSSE{
		UserID:  userid,
		Type:    evType,
		Message: b,
	})
}
