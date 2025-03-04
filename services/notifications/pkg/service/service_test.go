package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	ehmocks "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0/mocks"
	settingsmocks "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0/mocks"
	"github.com/owncloud/reva/v2/pkg/store"

	microstore "go-micro.dev/v4/store"

	"time"

	storemocks "go-micro.dev/v4/store/mocks"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc"
)

var _ = Describe("Notifications", func() {
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		vs              *settingssvc.MockValueService
		sharer          = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharer",
			},
			Mail:        "sharer@owncloud.com",
			DisplayName: "Dr. S. Harer",
		}
		sharee = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharee",
			},
			Mail:        "sharee@owncloud.com",
			DisplayName: "Eric Expireling",
		}
		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil).Once()
		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharee}, nil).Once()
		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)
		vs = &settingssvc.MockValueService{}
		vs.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return &settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_CollectionValue{
							CollectionValue: &settingsmsg.CollectionValue{
								Values: []*settingsmsg.CollectionOption{
									{
										Key:    "mail",
										Option: &settingsmsg.CollectionOption_BoolValue{BoolValue: true},
									},
								},
							},
						},
					},
				},
			}, nil
		}
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev events.Event) {
			ch := make(chan events.Event)
			evts := NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, vs, "",
				"", "", "", "", "", "",
				store.Create(), nil, nil)
			go evts.Run()

			ch <- ev
			select {
			case <-tc.done:
				// finished
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		Entry("Share Created", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer shared 'secrets of the board' with you",
			expectedTextBody: `Hello Eric Expireling

Dr. S. Harer has shared "secrets of the board" with you.

Click here to view it: files/shares/with-me


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:        sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}),

		Entry("Share Expired", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Share to 'secrets of the board' expired at 2023-04-17 16:42:00",
			expectedTextBody: `Hello Eric Expireling,

Your share to secrets of the board has expired at 2023-04-17 16:42:00

Even though this share has been revoked you still might have access through other shares and/or space memberships.


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.ShareExpired{
				ShareOwner:    sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:        resourceid,
			},
		}),

		Entry("Share Removed", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer unshared 'secrets of the board' with you",
			expectedTextBody: `Hello Eric Expireling,

Dr. S. Harer has unshared 'secrets of the board' with you.

Even though this share has been revoked you still might have access through other shares and/or space memberships.


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.ShareRemoved{
				Executant:     sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ItemID:        resourceid,
			},
		}),

		Entry("Added to Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer invited you to join secret space",
			expectedTextBody: `Hello Eric Expireling,

Dr. S. Harer has invited you to join "secret space".

Click here to view it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:     sharer.GetId(),
				Creator:       sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),

		Entry("Removed from Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer removed you from secret space",
			expectedTextBody: `Hello Eric Expireling,

Dr. S. Harer has removed you from "secret space".

You might still have access through your other groups or direct membership.

Click here to check it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceUnshared{
				Executant:     sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),

		Entry("Space Expired", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Membership of 'secret space' expired at 2023-04-17 16:42:00",
			expectedTextBody: `Hello Eric Expireling,

Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceMembershipExpired{
				SpaceOwner:    sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				SpaceID:       &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:     "secret space",
				ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		}),
	)
})

var _ = Describe("Notifications X-Site Scripting", func() {
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		vs              *settingssvc.MockValueService
		sharer          = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharer",
			},
			Mail:        "sharer@owncloud.com",
			DisplayName: "Dr. O'reilly",
		}
		sharee = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharee",
			},
			Mail:        "sharee@owncloud.com",
			DisplayName: "<script>alert('Eric Expireling');</script>",
		}
		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil).Once()
		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharee}, nil).Once()
		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
			Status: &rpc.Status{Code: rpc.Code_CODE_OK},
			Info: &provider.ResourceInfo{
				Name:  "<script>alert('secrets of the board');</script>",
				Space: &provider.StorageSpace{Name: "<script>alert('secret space');</script>"}},
		}, nil)
		vs = &settingssvc.MockValueService{}
		vs.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return &settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_CollectionValue{
							CollectionValue: &settingsmsg.CollectionValue{
								Values: []*settingsmsg.CollectionOption{
									{
										Key:    "mail",
										Option: &settingsmsg.CollectionOption_BoolValue{BoolValue: true},
									},
								},
							},
						},
					},
				},
			}, nil
		}
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev events.Event) {
			ch := make(chan events.Event)
			evts := NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, vs, "",
				"", "", "", "", "", "",
				store.Create(), nil, nil)
			go evts.Run()

			ch <- ev
			select {
			case <-tc.done:
				// finished
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		Entry("Share Created", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. O'reilly shared '<script>alert('secrets of the board');</script>' with you",
			expectedTextBody: `Hello <script>alert('Eric Expireling');</script>

Dr. O'reilly has shared "<script>alert('secrets of the board');</script>" with you.

Click here to view it: files/shares/with-me


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedHTMLBody: `<!DOCTYPE html>
<html>
<body>
<table cellspacing="0" cellpadding="0" border="0" width="100%">
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" width="600px">
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        Hello &lt;script&gt;alert(&#39;Eric Expireling&#39;);&lt;/script&gt;
                        <br><br>
                        Dr. O&#39;reilly has shared "&lt;script&gt;alert(&#39;secrets of the board&#39;);&lt;/script&gt;" with you.
                        <br><br>
                        Click here to view it: <a href="files/shares/with-me">files/shares/with-me</a>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        <footer>
                            <br>
                            <br>
                            --- <br>
                            ownCloud - Store. Share. Work.<br>
                            <a href="https://owncloud.com">https://owncloud.com</a>
                        </footer>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`,
			expectedSender: sharer.GetDisplayName(),

			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:        sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}),

		Entry("Added to Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. O'reilly invited you to join <script>alert('secret space');</script>",
			expectedTextBody: `Hello <script>alert('Eric Expireling');</script>,

Dr. O'reilly has invited you to join "<script>alert('secret space');</script>".

Click here to view it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			expectedHTMLBody: `<!DOCTYPE html>
<html>
<body>
<table cellspacing="0" cellpadding="0" border="0" width="100%">
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" width="600px">
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        Hello &lt;script&gt;alert(&#39;Eric Expireling&#39;);&lt;/script&gt;,
                        <br><br>
                        Dr. O&#39;reilly has invited you to join "&lt;script&gt;alert(&#39;secret space&#39;);&lt;/script&gt;".
                        <br><br>
                        Click here to view it: <a href="f/spaceid">f/spaceid</a>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        <footer>
                            <br>
                            <br>
                            --- <br>
                            ownCloud - Store. Share. Work.<br>
                            <a href="https://owncloud.com">https://owncloud.com</a>
                        </footer>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`,
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:     sharer.GetId(),
				Creator:       sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),
	)
})

var _ = Describe("Notifications grouped store", func() {
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		valueService    *settingsmocks.ValueService
		store           *storemocks.Store
		executant       = &user.User{
			Id: &user.UserId{
				OpaqueId: "executantId",
			},
			Mail:        "executant@owncloud.com",
			DisplayName: "executant",
		}
		receiverGrouped = &user.User{
			Id: &user.UserId{
				OpaqueId: "receiverGroupedId",
			},
			Mail:        "receiverGrouped@owncloud.com",
			DisplayName: "receiverGrouped",
		}
		// receiverInstant is used to trigger tc.done
		receiverInstant = &user.User{
			Id: &user.UserId{
				OpaqueId: "receiverInstantId",
			},
			Mail:        "receiverInstant@owncloud.com",
			DisplayName: "receiverInstant",
		}
		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	DescribeTable("Sending notifications",
		func(tc testChannelGroupedStore, ev events.Event, interval string) {
			// setup mocks
			// Note: This is done here and not inside a BeforeEach, because some mocks need variables from entries.

			pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
			gatewayClient = &cs3mocks.GatewayAPIClient{}
			gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
				"GatewaySelector",
				"com.owncloud.api.gateway",
				func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
					return gatewayClient
				},
			)

			gatewayClient.EXPECT().Authenticate(mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: executant}, nil)
			gatewayClient.EXPECT().Stat(mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)

			gatewayClient.EXPECT().GetGroup(mock.Anything, mock.Anything).
				Return(&group.GetGroupResponse{
					Status: &rpc.Status{Code: rpc.Code_CODE_OK},
					Group: &group.Group{
						Members: []*user.UserId{receiverGrouped.GetId(), receiverInstant.GetId()},
					},
				}, nil)
			gatewayClient.EXPECT().GetUser(mock.Anything, mock.Anything).
				Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: executant}, nil).Once()
			gatewayClient.EXPECT().GetUser(mock.Anything, mock.Anything).
				Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: receiverGrouped}, nil).Once()
			gatewayClient.EXPECT().GetUser(mock.Anything, mock.Anything).
				Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: receiverInstant}, nil).Once()

			valueService = settingsmocks.NewValueService(GinkgoT())

			// disabled email
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).
				Return(&settingssvc.GetValueResponse{Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_BoolValue{
							BoolValue: false,
						},
					},
				}}, nil).Twice()

			// filtering
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).Return(&settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_CollectionValue{
							CollectionValue: &settingsmsg.CollectionValue{
								Values: []*settingsmsg.CollectionOption{
									{
										Key:    "mail",
										Option: &settingsmsg.CollectionOption_BoolValue{BoolValue: true},
									},
								},
							},
						},
					},
				},
			}, nil).Twice()

			// interval
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).
				Return(&settingssvc.GetValueResponse{Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_StringValue{
							StringValue: interval,
						},
					},
				}}, nil).Once()
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).
				Return(&settingssvc.GetValueResponse{Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_StringValue{
							StringValue: "instant",
						},
					},
				}}, nil).Once()

			// locale
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).Return(&settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_ListValue{
							ListValue: &settingsmsg.ListValue{
								Values: []*settingsmsg.ListOptionValue{
									{
										Option: &settingsmsg.ListOptionValue_StringValue{StringValue: "en"},
									},
								},
							},
						},
					},
				},
			}, nil).Once()

			// fail if mocked functions are not called
			store = storemocks.NewStore(GinkgoT())
			store.EXPECT().Read(mock.Anything).Return([]*microstore.Record{}, nil).Once()
			store.EXPECT().Write(mock.Anything).Return(nil).Once()

			// setup EventsNotifier
			ch := make(chan events.Event)
			evts := NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, valueService, "",
				"", "", "", "", "", "",
				store, nil, nil)
			go evts.Run()

			// test
			ch <- ev
			select {
			case <-tc.done:
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		Entry("Share Created daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:         executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				CTime:          utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:         resourceid,
			},
		},
			"daily"),

		Entry("Share Expired daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareExpired{
				ShareOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:         resourceid,
			},
		},
			"daily"),

		Entry("Share Removed daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareRemoved{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ItemID:         resourceid,
			},
		},
			"daily"),

		Entry("Added to Space daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:      executant.GetId(),
				Creator:        executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		},
			"daily"),

		Entry("Removed from Space daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceUnshared{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		},
			"daily"),

		Entry("Space Expired daily", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceMembershipExpired{
				SpaceOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:      "secret space",
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		},
			"daily"),

		Entry("Share Created weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:         executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				CTime:          utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:         resourceid,
			},
		},
			"weekly"),

		Entry("Share Expired weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareExpired{
				ShareOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:         resourceid,
			},
		},
			"weekly"),

		Entry("Share Removed weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareRemoved{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ItemID:         resourceid,
			},
		},
			"weekly"),

		Entry("Added to Space weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:      executant.GetId(),
				Creator:        executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		},
			"weekly"),

		Entry("Removed from Space weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceUnshared{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		},
			"weekly"),

		Entry("Space Expired weekly", testChannelGroupedStore{
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceMembershipExpired{
				SpaceOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:      "secret space",
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		},
			"weekly"),
	)
})

var _ = Describe("Notifications grouped send", func() {
	const sender = "foo@bar.com"
	const subject = "Report"
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		valueService    *settingsmocks.ValueService
		store           *storemocks.Store
		historyClient   *ehmocks.EventHistoryService

		executant = &user.User{
			Id: &user.UserId{
				OpaqueId: "executantId",
			},
			Mail:        "executant@owncloud.com",
			DisplayName: "executant",
		}
		receiver = &user.User{
			Id: &user.UserId{
				OpaqueId: "receiver",
			},
			Mail:        "receiver@owncloud.com",
			DisplayName: "receiver",
		}

		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	DescribeTable("Sending notifications",
		func(tc testChannel, storedEvents []events.Event, interval string) {
			// setup mocks
			// Note: This is done here and not inside a BeforeEach, because some mocks need variables from entries.

			b, err := json.Marshal(userEventIds{
				User:     receiver,
				EventIds: []string{"event id"},
			})
			if err != nil {
				Fail(err.Error())
			}

			store = &storemocks.Store{}
			store.EXPECT().List(mock.Anything).Return([]string{"key"}, nil)
			store.EXPECT().Read(mock.Anything).Return([]*microstore.Record{{Value: b}}, nil).Once()
			store.EXPECT().Delete(mock.Anything).Return(nil).Once()

			var hcEvs []*ehmsg.Event
			for _, e := range storedEvents {
				b, err = json.Marshal(e.Event)
				if err != nil {
					Fail(err.Error())
				}
				hcEvs = append(hcEvs, &ehmsg.Event{Type: e.Type, Event: b})
			}

			historyClient = &ehmocks.EventHistoryService{}
			historyClient.EXPECT().GetEvents(mock.Anything, mock.Anything).Return(&v0.GetEventsResponse{Events: hcEvs}, nil)

			// locale
			valueService = &settingsmocks.ValueService{}
			valueService.EXPECT().GetValueByUniqueIdentifiers(mock.Anything, mock.Anything).Return(&settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_ListValue{
							ListValue: &settingsmsg.ListValue{
								Values: []*settingsmsg.ListOptionValue{
									{
										Option: &settingsmsg.ListOptionValue_StringValue{StringValue: "en"},
									},
								},
							},
						},
					},
				},
			}, nil).Once()

			pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
			gatewayClient = &cs3mocks.GatewayAPIClient{}
			gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
				"GatewaySelector",
				"com.owncloud.api.gateway",
				func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
					return gatewayClient
				},
			)

			gatewayClient.EXPECT().Authenticate(mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: executant}, nil)
			gatewayClient.EXPECT().Stat(mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)
			gatewayClient.EXPECT().GetUser(mock.Anything, mock.Anything).
				Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: executant}, nil).Once()

			// setup EventsNotifier
			evs := []events.Unmarshaller{
				events.ShareCreated{},
				events.ShareExpired{},
				events.ShareRemoved{},
				events.SpaceShared{},
				events.SpaceUnshared{},
				events.SpaceMembershipExpired{},
				events.SendEmailsEvent{},
			}
			registeredEvents := make(map[string]events.Unmarshaller)
			for _, e := range evs {
				typ := reflect.TypeOf(e)
				registeredEvents[typ.String()] = e
			}
			ch := make(chan events.Event)
			evts := NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, valueService, "",
				"", "", "", "", "",
				sender, store, historyClient, registeredEvents)
			go evts.Run()

			// trigger sending
			ch <- events.Event{
				Event: events.SendEmailsEvent{
					Interval: interval,
				},
			}
			select {
			case <-tc.done:
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		// daily
		Entry("Share Created daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody:    buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.`}),
			expectedHTMLBody:    buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.`}),
			done:                make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareCreated{}).String(),
			Event: events.ShareCreated{
				Sharer:        executant.GetId(),
				GranteeUserID: receiver.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}}, "daily"),

		Entry("Share Expired daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`Your share to secrets of the board has expired at 2023-04-17 16:42:00

Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`Your share to secrets of the board has expired at 2023-04-17 16:42:00<br><br>Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareExpired{}).String(),
			Event: events.ShareExpired{
				ShareOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:         resourceid,
			},
		}}, "daily"),

		Entry("Share Removed daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has unshared 'secrets of the board' with you.

Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has unshared 'secrets of the board' with you.<br><br>Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareRemoved{}).String(),
			Event: events.ShareRemoved{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ItemID:         resourceid,
			},
		}}, "daily"),

		Entry("Added to Space daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody:    buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has invited you to join "secret space".`}),
			expectedHTMLBody:    buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has invited you to join "secret space".`}),
			done:                make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceShared{}).String(),
			Event: events.SpaceShared{
				Executant:      executant.GetId(),
				Creator:        executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}}, "daily"),

		Entry("Removed from Space daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has removed you from "secret space".

You might still have access through your other groups or direct membership.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has removed you from "secret space".<br><br>You might still have access through your other groups or direct membership.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceUnshared{}).String(),
			Event: events.SpaceUnshared{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}}, "daily"),

		Entry("Space Expired daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`Your membership of space secret space has expired at 2023-04-17 16:42:00<br><br>Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceMembershipExpired{}).String(),
			Event: events.SpaceMembershipExpired{
				SpaceOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:      "secret space",
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		}}, "daily"),

		Entry("Share Created and Space Expired grouped daily", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.


Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.<br><br><br>Your membership of space secret space has expired at 2023-04-17 16:42:00<br><br>Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			done:             make(chan struct{}),
		}, []events.Event{
			{
				Type: reflect.TypeOf(events.ShareCreated{}).String(),
				Event: events.ShareCreated{
					Sharer:        executant.GetId(),
					GranteeUserID: receiver.GetId(),
					CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
					ItemID:        resourceid,
				},
			},
			{
				Type: reflect.TypeOf(events.SpaceMembershipExpired{}).String(),
				Event: events.SpaceMembershipExpired{
					SpaceOwner:     executant.GetId(),
					GranteeGroupID: &group.GroupId{},
					SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
					SpaceName:      "secret space",
					ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				},
			},
		}, "daily"),

		// weekly
		Entry("Share Created weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody:    buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.`}),
			expectedHTMLBody:    buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.`}),
			done:                make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareCreated{}).String(),
			Event: events.ShareCreated{
				Sharer:        executant.GetId(),
				GranteeUserID: receiver.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}}, "weekly"),

		Entry("Share Expired weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`Your share to secrets of the board has expired at 2023-04-17 16:42:00

Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`Your share to secrets of the board has expired at 2023-04-17 16:42:00<br><br>Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareExpired{}).String(),
			Event: events.ShareExpired{
				ShareOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:         resourceid,
			},
		}}, "weekly"),

		Entry("Share Removed weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has unshared 'secrets of the board' with you.

Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has unshared 'secrets of the board' with you.<br><br>Even though this share has been revoked you still might have access through other shares and/or space memberships.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.ShareRemoved{}).String(),
			Event: events.ShareRemoved{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ItemID:         resourceid,
			},
		}}, "weekly"),

		Entry("Added to Space weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody:    buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has invited you to join "secret space".`}),
			expectedHTMLBody:    buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has invited you to join "secret space".`}),
			done:                make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceShared{}).String(),
			Event: events.SpaceShared{
				Executant:      executant.GetId(),
				Creator:        executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}}, "weekly"),

		Entry("Removed from Space weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has removed you from "secret space".

You might still have access through your other groups or direct membership.`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has removed you from "secret space".<br><br>You might still have access through your other groups or direct membership.`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceUnshared{}).String(),
			Event: events.SpaceUnshared{
				Executant:      executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				ID:             &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}}, "weekly"),

		Entry("Space Expired weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`Your membership of space secret space has expired at 2023-04-17 16:42:00<br><br>Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			done:             make(chan struct{}),
		}, []events.Event{{
			Type: reflect.TypeOf(events.SpaceMembershipExpired{}).String(),
			Event: events.SpaceMembershipExpired{
				SpaceOwner:     executant.GetId(),
				GranteeGroupID: &group.GroupId{},
				SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:      "secret space",
				ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		}}, "weekly"),

		Entry("Share Created and Space Expired grouped weekly", testChannel{
			expectedReceipients: []string{receiver.GetMail()},
			expectedSender:      sender,
			expectedSubject:     subject,
			expectedTextBody: buildExpectedTextBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.


Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			expectedHTMLBody: buildExpectedHTMLBodyGrouped(receiver.DisplayName, []string{`executant has shared "secrets of the board" with you.<br><br><br>Your membership of space secret space has expired at 2023-04-17 16:42:00<br><br>Even though this membership has expired you still might have access through other shares and/or space memberships`}),
			done:             make(chan struct{}),
		}, []events.Event{
			{
				Type: reflect.TypeOf(events.ShareCreated{}).String(),
				Event: events.ShareCreated{
					Sharer:        executant.GetId(),
					GranteeUserID: receiver.GetId(),
					CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
					ItemID:        resourceid,
				},
			},
			{
				Type: reflect.TypeOf(events.SpaceMembershipExpired{}).String(),
				Event: events.SpaceMembershipExpired{
					SpaceOwner:     executant.GetId(),
					GranteeGroupID: &group.GroupId{},
					SpaceID:        &provider.StorageSpaceId{OpaqueId: "spaceid"},
					SpaceName:      "secret space",
					ExpiredAt:      time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				},
			},
		}, "weekly"),
	)
})

func buildExpectedTextBodyGrouped(receiver string, events []string) string {
	body := strings.Join(events, "\n\n\n")
	return fmt.Sprintf(
		`Hi %s,

%s


---
ownCloud - Store. Share. Work.
https://owncloud.com
`, receiver, body)
}

func buildExpectedHTMLBodyGrouped(receiver string, events []string) string {
	body := strings.Join(events, "<br><br><br>")
	return fmt.Sprintf(
		`<!DOCTYPE html>
<html>
<body>
<table cellspacing="0" cellpadding="0" border="0" width="100%%">
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" width="600px">
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        Hi %s,
                        <br><br>
                        %s
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        <footer>
                            <br>
                            <br>
                            --- <br>
                            ownCloud - Store. Share. Work.<br>
                            <a href="https://owncloud.com">https://owncloud.com</a>
                        </footer>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`, receiver, body)
}

// NOTE: This is explicitly not testing the message itself. Should we?
type testChannel struct {
	expectedReceipients []string
	expectedSubject     string
	expectedTextBody    string
	expectedHTMLBody    string
	expectedSender      string
	done                chan struct{}
}

func (tc testChannel) SendMessage(_ context.Context, m *channels.Message) error {
	defer GinkgoRecover()

	Expect(tc.expectedReceipients).To(Equal(m.Recipient))
	Expect(tc.expectedSubject).To(Equal(m.Subject))
	Expect(tc.expectedTextBody).To(Equal(m.TextBody))
	Expect(tc.expectedSender).To(Equal(m.Sender))
	if tc.expectedHTMLBody != "" {
		Expect(tc.expectedHTMLBody).To(Equal(m.HTMLBody))
	}
	tc.done <- struct{}{}
	return nil
}

type testChannelGroupedStore struct {
	done chan struct{}
}

func (tc testChannelGroupedStore) SendMessage(_ context.Context, _ *channels.Message) error {
	defer GinkgoRecover()

	tc.done <- struct{}{}
	return nil
}
