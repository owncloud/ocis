package service_test

import (
	"context"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/service"
	"github.com/test-go/testify/mock"
	"go-micro.dev/v4/client"
)

var _ = Describe("Notifications", func() {
	var (
		gwc    *cs3mocks.GatewayAPIClient
		vs     *settingssvc.MockValueService
		sharer = &user.User{
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
		gwc = &cs3mocks.GatewayAPIClient{}
		gwc.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil).Once()
		gwc.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharee}, nil).Once()
		gwc.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gwc.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)
		vs = &settingssvc.MockValueService{}
		vs.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return nil, nil
		}
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev events.Event) {
			cfg := defaults.FullDefaultConfig()
			cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
			_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			ch := make(chan events.Event)
			evts := service.NewEventsNotifier(ch, tc, log.NewLogger(), gwc, vs, "", "", "")
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
			expectedMessage: `Hello Eric Expireling

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
			expectedMessage: `Hello Eric Expireling,

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

		Entry("Added to Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer invited you to join secret space",
			expectedMessage: `Hello Eric Expireling,

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
			expectedMessage: `Hello Eric Expireling,

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
			expectedMessage: `Hello Eric Expireling,

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

// NOTE: This is explictitly not testing the message itself. Should we?
type testChannel struct {
	expectedReceipients []string
	expectedSubject     string
	expectedMessage     string
	expectedSender      string
	done                chan struct{}
}

func (tc testChannel) SendMessage(ctx context.Context, m *channels.Message) error {
	defer GinkgoRecover()

	Expect(m.Recipient).To(Equal(tc.expectedReceipients))
	Expect(m.Subject).To(Equal(tc.expectedSubject))
	Expect(m.TextBody).To(Equal(tc.expectedMessage))
	Expect(m.Sender).To(Equal(tc.expectedSender))
	tc.done <- struct{}{}
	return nil
}
