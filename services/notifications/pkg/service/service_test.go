package service_test

import (
	"context"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/service"
	"github.com/test-go/testify/mock"
)

var _ = Describe("Notifications", func() {
	var (
		gwc    *cs3mocks.GatewayAPIClient
		sharer = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharer",
			},
			DisplayName: "Dr. S. Harer",
		}
		sharee = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharee",
			},
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
		gwc.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gwc.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gwc.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev interface{}) {
			ch := make(chan interface{})
			evts := service.NewEventsNotifier(ch, tc, log.NewLogger(), gwc, "", "", "")
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
			expectedReceipients: map[string]bool{sharee.GetId().GetOpaqueId(): true},
			expectedSubject:     "Dr. S. Harer shared 'secrets of the board' with you",
			expectedSender:      sharer.GetDisplayName(),
			done:                make(chan struct{}),
		}, events.ShareCreated{
			Sharer:        sharer.GetId(),
			GranteeUserID: sharee.GetId(),
			CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
			ItemID:        resourceid,
		}),

		Entry("Share Expired", testChannel{
			expectedReceipients: map[string]bool{sharee.GetId().GetOpaqueId(): true},
			expectedSubject:     "Share to 'secrets of the board' expired at 2023-04-17 16:42:00",
			expectedSender:      sharer.GetDisplayName(),
			done:                make(chan struct{}),
		}, events.ShareExpired{
			ShareOwner:    sharer.GetId(),
			GranteeUserID: sharee.GetId(),
			ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			ItemID:        resourceid,
		}),

		Entry("Added to Space", testChannel{
			expectedReceipients: map[string]bool{sharee.GetId().GetOpaqueId(): true},
			expectedSubject:     "Dr. S. Harer invited you to join secret space",
			expectedSender:      sharer.GetDisplayName(),
			done:                make(chan struct{}),
		}, events.SpaceShared{
			Executant:     sharer.GetId(),
			Creator:       sharer.GetId(),
			GranteeUserID: sharee.GetId(),
			ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
		}),

		Entry("Removed from Space", testChannel{
			expectedReceipients: map[string]bool{sharee.GetId().GetOpaqueId(): true},
			expectedSubject:     "Dr. S. Harer removed you from secret space",
			expectedSender:      sharer.GetDisplayName(),
			done:                make(chan struct{}),
		}, events.SpaceUnshared{
			Executant:     sharer.GetId(),
			GranteeUserID: sharee.GetId(),
			ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
		}),

		Entry("Space Expired", testChannel{
			expectedReceipients: map[string]bool{sharee.GetId().GetOpaqueId(): true},
			expectedSubject:     "Membership of 'secret space' expired at 2023-04-17 16:42:00",
			expectedSender:      sharer.GetDisplayName(),
			done:                make(chan struct{}),
		}, events.SpaceMembershipExpired{
			SpaceOwner:    sharer.GetId(),
			GranteeUserID: sharee.GetId(),
			SpaceID:       &provider.StorageSpaceId{OpaqueId: "spaceid"},
			SpaceName:     "secret space",
			ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
		}),
	)
})

// NOTE: This is explictitly not testing the message itself. Should we?
type testChannel struct {
	expectedReceipients map[string]bool
	expectedSubject     string
	expectedSender      string
	done                chan struct{}
}

func (tc testChannel) SendMessage(ctx context.Context, userIDs []string, msg, subject, senderDisplayName string) error {
	defer GinkgoRecover()

	for _, u := range userIDs {
		Expect(tc.expectedReceipients[u]).To(Equal(true))
	}

	// TODO: test the message?
	//Expect(msg).To(Equal(tc.expectedMessage))
	Expect(subject).To(Equal(tc.expectedSubject))
	Expect(senderDisplayName).To(Equal(tc.expectedSender))
	tc.done <- struct{}{}
	return nil
}

func (tc testChannel) SendMessageToGroup(ctx context.Context, groupID *group.GroupId, msg, subject, senderDisplayName string) error {
	return tc.SendMessage(ctx, []string{groupID.GetOpaqueId()}, msg, subject, senderDisplayName)
}
