package indexer_test

import (
	"context"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
	engineMocks "github.com/owncloud/ocis/v2/services/search/pkg/engine/mocks"
	"github.com/owncloud/ocis/v2/services/search/pkg/indexer"
	"github.com/test-go/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer", func() {
	var (
		i           indexer.Indexer
		extractor   *contentMocks.Extractor
		gw          *cs3mocks.GatewayAPIClient
		indexClient *engineMocks.Engine
		ctx         context.Context

		logger = log.NewLogger()
		user   = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "parentopaqueid",
			},
			Path:  "foo.pdf",
			Size:  12345,
			Mtime: utils.TimeToTS(time.Now().Add(-time.Hour)),
		}
	)

	BeforeEach(func() {
		ctx = context.Background()
		gw = &cs3mocks.GatewayAPIClient{}
		indexClient = &engineMocks.Engine{}
		extractor = &contentMocks.Extractor{}

		i = indexer.NewIndexer(gw, indexClient, extractor, logger, "")

		gw.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
			Status: status.NewOK(ctx),
			Token:  "authtoken",
		}, nil)
		gw.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
			Status: status.NewOK(context.Background()),
			Info:   ri,
		}, nil)
	})

	Describe("IndexSpace", func() {
		It("walks the space", func() {
			extractor.Mock.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			indexClient.On("Upsert", mock.Anything, mock.Anything).Return(nil)
			indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{}, nil)
			gw.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)

			err := i.IndexSpace(ctx, &provider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"}, &userv1beta1.UserId{OpaqueId: "user"})
			Expect(err).ToNot(HaveOccurred())

			indexClient.AssertCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
		})

		It("doesn't reindex if there's nothing to do", func() {
			extractor.Mock.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
				Matches: []*v0.Match{{}},
			}, nil)
			gw.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)

			err := i.IndexSpace(ctx, &provider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"}, &userv1beta1.UserId{OpaqueId: "user"})
			Expect(err).ToNot(HaveOccurred())

			indexClient.AssertNotCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
		})
	})
})
