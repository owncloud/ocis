package service

import (
	"context"
	"encoding/json"
	"reflect"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	event "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0/mocks"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/stretchr/testify/mock"

	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/store"
	storemocks "go-micro.dev/v4/store/mocks"
)

func Test_userEventStore_persist(t *testing.T) {
	tEventId := "event id"
	tUser := &user.User{Id: &user.UserId{OpaqueId: "userid"}}
	tInterval := "interval"
	tKey := tInterval + "_" + tUser.Id.OpaqueId
	tErr := errors.New("some error")
	tValueSingleEventId, err := json.Marshal(userEventIds{
		User:     tUser,
		EventIds: []string{tEventId},
	})
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		log           log.Logger
		store         store.Store
		historyClient ehsvc.EventHistoryService
	}
	type args struct {
		interval string
		eventId  string
		users    []*user.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*user.User
	}{
		{
			name: "new record",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{}, nil).Once()

					s.EXPECT().Write(&store.Record{Key: tKey, Value: tValueSingleEventId}).Return(nil).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				interval: tInterval,
				eventId:  tEventId,
				users:    []*user.User{tUser},
			},
			want: []*user.User(nil),
		},
		{
			name: "append to record",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{{Key: tKey, Value: tValueSingleEventId}}, nil).Once()

					b, err := json.Marshal(userEventIds{
						User:     tUser,
						EventIds: []string{tEventId, tEventId},
					})
					if err != nil {
						t.Fatal(err)
					}
					s.EXPECT().Write(&store.Record{Key: tKey, Value: b}).Return(nil).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				interval: tInterval,
				eventId:  tEventId,
				users:    []*user.User{tUser},
			},
			want: []*user.User(nil),
		},
		{
			name: "error on read",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return(nil, tErr).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				interval: tInterval,
				eventId:  tEventId,
				users:    []*user.User{tUser},
			},
			want: []*user.User{tUser},
		},
		{
			name: "error on write",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{}, nil).Once()

					b, err := json.Marshal(userEventIds{
						User:     tUser,
						EventIds: []string{tEventId},
					})
					if err != nil {
						t.Fatal(err)
					}
					s.EXPECT().Write(&store.Record{Key: tKey, Value: b}).Return(tErr).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				interval: tInterval,
				eventId:  tEventId,
				users:    []*user.User{tUser},
			},
			want: []*user.User{tUser},
		},
		{
			name: "multiple users",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{}, nil).Once()
					s.EXPECT().Write(&store.Record{Key: tKey, Value: tValueSingleEventId}).Return(nil).Once()

					tUser2 := &user.User{Id: &user.UserId{OpaqueId: "userid2"}}
					tKey2 := tInterval + "_" + tUser2.Id.OpaqueId
					tValueSingleEventId2, err := json.Marshal(userEventIds{
						User:     tUser2,
						EventIds: []string{tEventId},
					})
					if err != nil {
						t.Fatal(err)
					}
					s.EXPECT().Read(tKey2).Return([]*store.Record{}, nil).Once()
					s.EXPECT().Write(&store.Record{Key: tKey2, Value: tValueSingleEventId2}).Return(nil).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				interval: tInterval,
				eventId:  tEventId,
				users:    []*user.User{tUser, {Id: &user.UserId{OpaqueId: "userid2"}}},
			},
			want: []*user.User(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userEventStore{
				log:           tt.fields.log,
				store:         tt.fields.store,
				historyClient: tt.fields.historyClient,
			}
			assert.Equalf(t, tt.want, s.persist(tt.args.interval, tt.args.eventId, tt.args.users), "persist(%v, %v, %v)", tt.args.interval, tt.args.eventId, tt.args.users)
		})
	}
}

func Test_userEventStore_pop(t *testing.T) {
	tCtx := context.TODO()
	tUser := &user.User{Id: &user.UserId{OpaqueId: "userid"}}
	tInterval := "interval"
	tKey := tInterval + "_" + tUser.Id.OpaqueId
	tErr := errors.New("some error")

	// event 1
	tEventId := "event id"
	tValueSingleEventId, err := json.Marshal(userEventIds{
		User:     tUser,
		EventIds: []string{tEventId},
	})
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(events.ShareCreated{
		ShareID: &collaboration.ShareId{OpaqueId: "shareid"},
	})
	tEvent := &event.Event{
		Type:  reflect.TypeOf(&events.ShareCreated{}).String(),
		Id:    tEventId,
		Event: b,
	}
	tUserEvents := &userEvents{
		User:   tUser,
		Events: []*event.Event{tEvent},
	}

	// event 2
	tEventId2 := "event id2"
	if err != nil {
		t.Fatal(err)
	}
	b, err = json.Marshal(events.ShareRemoved{
		ShareID: &collaboration.ShareId{OpaqueId: "shareid"},
	})
	if err != nil {
		t.Fatal(err)
	}
	tEvent2 := &event.Event{
		Type:  reflect.TypeOf(&events.ShareRemoved{}).String(),
		Id:    tEventId,
		Event: b,
	}

	tValueMultipleEventIds, err := json.Marshal(userEventIds{
		User:     tUser,
		EventIds: []string{tEventId, tEventId2},
	})
	tUserMultipleEvents := &userEvents{
		User:   tUser,
		Events: []*event.Event{tEvent, tEvent2},
	}

	type fields struct {
		log           log.Logger
		store         store.Store
		historyClient ehsvc.EventHistoryService
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *userEvents
		wantErr bool
	}{
		{
			name: "single event",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{{Key: tKey, Value: tValueSingleEventId}}, nil).Once()

					s.EXPECT().Delete(tKey).Return(nil).Once()

					return s
				}(),
				historyClient: func() ehsvc.EventHistoryService {
					hc := mocks.NewEventHistoryService(t)

					hc.EXPECT().GetEvents(mock.Anything, &ehsvc.GetEventsRequest{Ids: []string{tEventId}}).
						Return(&ehsvc.GetEventsResponse{Events: []*event.Event{tEvent}}, nil).Once()

					return hc
				}(),
			},
			args: args{
				ctx: tCtx,
				key: tKey,
			},
			want:    tUserEvents,
			wantErr: false,
		},
		{
			name: "multiple events",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)
					s.EXPECT().Read(tKey).Return([]*store.Record{{Key: tKey, Value: tValueMultipleEventIds}}, nil).Once()

					s.EXPECT().Delete(tKey).Return(nil).Once()

					return s
				}(),
				historyClient: func() ehsvc.EventHistoryService {
					hc := mocks.NewEventHistoryService(t)

					hc.EXPECT().GetEvents(tCtx, &ehsvc.GetEventsRequest{Ids: []string{tEventId, tEventId2}}).
						Return(&ehsvc.GetEventsResponse{Events: []*event.Event{tEvent, tEvent2}}, nil).Once()

					return hc
				}(),
			},
			args: args{
				ctx: tCtx,
				key: tKey,
			},
			want:    tUserMultipleEvents,
			wantErr: false,
		},
		{
			name: "error on read",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)
					s.EXPECT().Read(tKey).Return(nil, tErr).Once()

					return s
				}(),
				historyClient: nil,
			},
			args: args{
				ctx: tCtx,
				key: tKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error on get events",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)
					s.EXPECT().Read(tKey).Return([]*store.Record{{Key: tKey, Value: tValueSingleEventId}}, nil).Once()
					return s
				}(),
				historyClient: func() ehsvc.EventHistoryService {
					hc := mocks.NewEventHistoryService(t)

					hc.EXPECT().GetEvents(mock.Anything, &ehsvc.GetEventsRequest{Ids: []string{tEventId}}).
						Return(nil, tErr).Once()

					return hc
				}(),
			},
			args: args{
				ctx: tCtx,
				key: tKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error on delete",
			fields: fields{
				log: testLogger,
				store: func() store.Store {
					s := storemocks.NewStore(t)

					s.EXPECT().Read(tKey).Return([]*store.Record{{Key: tKey, Value: tValueSingleEventId}}, nil).Once()

					s.EXPECT().Delete(tKey).Return(tErr).Once()

					return s
				}(),
				historyClient: func() ehsvc.EventHistoryService {
					hc := mocks.NewEventHistoryService(t)

					hc.EXPECT().GetEvents(mock.Anything, &ehsvc.GetEventsRequest{Ids: []string{tEventId}}).
						Return(&ehsvc.GetEventsResponse{Events: []*event.Event{tEvent}}, nil).Once()

					return hc
				}(),
			},
			args: args{
				ctx: tCtx,
				key: tKey,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userEventStore{
				log:           tt.fields.log,
				store:         tt.fields.store,
				historyClient: tt.fields.historyClient,
			}
			got, err := s.pop(tt.args.ctx, tt.args.key)
			if tt.wantErr {
				assert.NotNil(t, err)
			}
			assert.Equalf(t, tt.want, got, "pop(%v, %v)", tt.args.ctx, tt.args.key)
		})
	}
}
