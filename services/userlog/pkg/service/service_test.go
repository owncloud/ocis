package service_test

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/store"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/mocks"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/service"
	"github.com/test-go/testify/mock"
	microevents "go-micro.dev/v4/events"
	microstore "go-micro.dev/v4/store"
)

var _ = Describe("UserlogService", func() {
	var (
		cfg = &config.Config{}

		ul  *service.UserlogService
		bus testBus
		sto microstore.Store

		gwc cs3mocks.GatewayAPIClient
		ehc mocks.EventHistoryService
	)

	BeforeEach(func() {
		var err error
		sto = store.Create()
		bus = testBus(make(chan events.Event))
		o := utils.AppendJSONToOpaque(nil, "grants", map[string]*provider.ResourcePermissions{"userid": {Stat: true}})
		gwc.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{StorageSpaces: []*provider.StorageSpace{
			{
				Opaque:    o,
				SpaceType: "project",
			},
		}, Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil)
		gwc.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{User: &user.User{Id: &user.UserId{OpaqueId: "userid"}}, Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil)
		gwc.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil)

		ul, err = service.NewUserlogService(
			service.Config(cfg),
			service.Consumer(bus),
			service.Store(sto),
			service.Logger(log.NewLogger()),
			service.Mux(chi.NewMux()),
			service.GatewayClient(&gwc),
			service.HistoryClient(&ehc),
			service.RegisteredEvents([]events.Unmarshaller{
				events.UploadReady{},
			}),
		)
		Expect(err).ToNot(HaveOccurred())

	})

	It("it stores, returns and deletes a couple of events", func() {
		ids := make(map[string]struct{})
		ids[bus.Publish(events.SpaceCreated{Executant: &user.UserId{OpaqueId: "userid"}})] = struct{}{}
		ids[bus.Publish(events.UploadReady{SpaceOwner: &user.UserId{OpaqueId: "userid"}})] = struct{}{}
		ids[bus.Publish(events.ContainerCreated{SpaceOwner: &user.UserId{OpaqueId: "userid"}})] = struct{}{}

		time.Sleep(500 * time.Millisecond)

		var events []*ehmsg.Event
		for id := range ids {
			events = append(events, &ehmsg.Event{Id: id})
		}

		ehc.On("GetEvents", mock.Anything, mock.Anything).Return(&ehsvc.GetEventsResponse{Events: events}, nil)

		evs, err := ul.GetEvents(context.Background(), "userid")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(evs)).To(Equal(len(ids)))

		var evids []string
		for _, e := range evs {
			_, exists := ids[e.Id]
			Expect(exists).To(BeTrue())
			delete(ids, e.Id)
			evids = append(evids, e.Id)
		}

		Expect(len(ids)).To(Equal(0))
		err = ul.DeleteEvents("userid", evids)
		Expect(err).ToNot(HaveOccurred())

		evs, err = ul.GetEvents(context.Background(), "userid")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(evs)).To(Equal(0))
	})

	AfterEach(func() {
		close(bus)
	})
})

type testBus chan events.Event

func (tb testBus) Consume(_ string, _ ...microevents.ConsumeOption) (<-chan microevents.Event, error) {
	ch := make(chan microevents.Event)
	go func() {
		for ev := range tb {
			b, _ := json.Marshal(ev.Event)
			ch <- microevents.Event{
				Payload: b,
				Metadata: map[string]string{
					events.MetadatakeyEventID:   ev.ID,
					events.MetadatakeyEventType: ev.Type,
				},
			}
		}
	}()
	return ch, nil
}

func (tb testBus) Publish(e interface{}) string {
	ev := events.Event{
		ID:    uuid.New().String(),
		Type:  reflect.TypeOf(e).String(),
		Event: e,
	}

	tb <- ev
	return ev.ID
}
