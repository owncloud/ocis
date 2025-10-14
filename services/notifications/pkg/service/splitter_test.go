package service

import (
	"context"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settings "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/client"
	"strings"
	"testing"
)

func Test_intervalSplitter_execute(t *testing.T) {
	type fields struct {
		log         log.Logger
		valueClient v0.ValueService
	}
	type args struct {
		ctx       context.Context
		users     []*user.User
		settingId string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantInstant []*user.User
		wantDaily   []*user.User
		wantWeekly  []*user.User
	}{
		{"no connection to ValueService",
			fields{
				log: testLogger,
				valueClient: settings.MockValueService{
					GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
						return nil, errors.New("no connection to ValueService")
					}}}, args{
				ctx:       context.TODO(),
				users:     newUsers("foo"),
				settingId: "",
			},
			newUsers("foo"), []*user.User(nil), []*user.User(nil),
		},
		{"no setting in ValueService response",
			fields{
				log: testLogger,
				valueClient: settings.MockValueService{
					GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
						return &settings.GetValueResponse{}, nil
					}}},
			args{
				ctx:       context.TODO(),
				users:     newUsers("foo"),
				settingId: "",
			},
			newUsers("foo"), []*user.User(nil), []*user.User(nil),
		},
		{"ValueService nil response",
			fields{
				log: testLogger,
				valueClient: settings.MockValueService{
					GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
						return nil, nil
					}}},
			args{
				ctx:       context.TODO(),
				users:     newUsers("foo"),
				settingId: "",
			},
			newUsers("foo"), []*user.User(nil), []*user.User(nil),
		},
		{"input users nil",
			fields{
				log: testLogger,
				valueClient: settings.MockValueService{
					GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
						return nil, nil
					}},
			},
			args{
				ctx:   context.TODO(),
				users: nil,
			},
			[]*user.User(nil), []*user.User(nil), []*user.User(nil),
		},
		{"interval never",
			fields{
				log:         testLogger,
				valueClient: newStringValueMockValueService("never"),
			},
			args{
				ctx:   context.TODO(),
				users: newUsers("foo"),
			},
			[]*user.User(nil), []*user.User(nil), []*user.User(nil),
		},
		{"interval instant",
			fields{
				log:         testLogger,
				valueClient: newStringValueMockValueService("instant"),
			},
			args{
				ctx:   context.TODO(),
				users: newUsers("foo"),
			},
			newUsers("foo"), []*user.User(nil), []*user.User(nil),
		},
		{"interval daily",
			fields{
				log:         testLogger,
				valueClient: newStringValueMockValueService("daily"),
			},
			args{
				ctx:   context.TODO(),
				users: newUsers("foo"),
			},
			[]*user.User(nil), newUsers("foo"), []*user.User(nil),
		},
		{"interval weekly",
			fields{
				log:         testLogger,
				valueClient: newStringValueMockValueService("weekly"),
			},
			args{
				ctx:   context.TODO(),
				users: newUsers("foo"),
			},
			[]*user.User(nil), []*user.User(nil), newUsers("foo"),
		},
		{"multiple users and intervals",
			fields{
				log: testLogger,
				valueClient: settings.MockValueService{
					GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
						if strings.Contains(req.AccountUuid, "never") {
							return newGetValueResponseStringValue("never"), nil
						} else if strings.Contains(req.AccountUuid, "instant") {
							return newGetValueResponseStringValue("instant"), nil
						} else if strings.Contains(req.AccountUuid, "daily") {
							return newGetValueResponseStringValue("daily"), nil
						} else if strings.Contains(req.AccountUuid, "weekly") {
							return newGetValueResponseStringValue("weekly"), nil
						}
						return nil, nil
					}},
			},
			args{
				ctx:   context.TODO(),
				users: newUsers("never1", "instant1", "daily1", "weekly1", "never2", "instant2", "daily2", "weekly2"),
			},
			newUsers("instant1", "instant2"), newUsers("daily1", "daily2"), newUsers("weekly1", "weekly2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := intervalSplitter{
				log:         tt.fields.log,
				valueClient: tt.fields.valueClient,
			}
			gotInstant, gotDaily, gotWeekly := s.execute(tt.args.ctx, tt.args.users)
			assert.Equalf(t, tt.wantInstant, gotInstant, "execute(%v, %v, %v)", tt.args.ctx, tt.args.users)
			assert.Equalf(t, tt.wantDaily, gotDaily, "execute(%v, %v, %v)", tt.args.ctx, tt.args.users)
			assert.Equalf(t, tt.wantWeekly, gotWeekly, "execute(%v, %v, %v)", tt.args.ctx, tt.args.users)
		})
	}
}

func newStringValueMockValueService(strVal string) settings.ValueService {
	return settings.MockValueService{
		GetValueByUniqueIdentifiersFunc: func(ctx context.Context, req *settings.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settings.GetValueResponse, error) {
			return newGetValueResponseStringValue(strVal), nil
		},
	}
}

func newGetValueResponseStringValue(strVal string) *settings.GetValueResponse {
	return &settings.GetValueResponse{Value: &settingsmsg.ValueWithIdentifier{
		Value: &settingsmsg.Value{
			Value: &settingsmsg.Value_StringValue{
				StringValue: strVal,
			},
		},
	}}
}

func newUsers(ids ...string) []*user.User {
	var users []*user.User
	for _, s := range ids {
		users = append(users, &user.User{Id: &user.UserId{OpaqueId: s}})
	}
	return users
}
