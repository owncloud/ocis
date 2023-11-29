package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("Utils", func() {
	var (
		svc service.Graph
	)

	BeforeEach(func() {
		cfg := defaults.FullDefaultConfig()
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		identityBackend := &identitymocks.Backend{}
		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithIdentityBackend(identityBackend),
		)
	})

	DescribeTable("GetDriveAndItemIDParam",
		func(driveID, itemID string, shouldPass bool) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", driveID)
			rctx.URLParams.Add("itemID", itemID)

			extractedDriveID, extractedItemID, err := svc.GetDriveAndItemIDParam(
				httptest.NewRequest(http.MethodGet, "/", nil).
					WithContext(
						context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
					),
			)

			switch shouldPass {
			case true:
				Expect(err).To(BeNil())
				parsedItemID, _ := storagespace.ParseID(itemID)
				Expect(extractedItemID).To(Equal(parsedItemID))

				parsedDriveID, _ := storagespace.ParseID(driveID)
				Expect(extractedDriveID).To(Equal(parsedDriveID))
			default:
				Expect(err).ToNot(BeNil())
			}
		},
		Entry("fails: invalid driveID", "", "1$2!3", false),
		Entry("fails: invalid itemID", "1$2", "", false),
		Entry("fails: incompatible driveID and itemID", "1$2", "3$4!5", false),
		Entry("fails: no itemID opaqueId", "1$2", "3$4", false),
		Entry("pass: valid driveID and itemID", "1$2", "1$2!5", true),
	)

	DescribeTable("IsSpaceRoot",
		func(resourceID *provider.ResourceId, isRoot bool) {
			Expect(service.IsSpaceRoot(resourceID)).To(Equal(isRoot))
		},
		Entry("spaceId and opaqueID equal", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
			SpaceId:   "2",
		}, true),
		Entry("nil", nil, false),
		Entry("spaceID empty", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
		}, false),
		Entry("opaqueID empty", &provider.ResourceId{
			StorageId: "1",
			SpaceId:   "3",
		}, false),
		Entry("spaceID and opaqueID unequal", &provider.ResourceId{
			OpaqueId: "2",
			SpaceId:  "3",
		}, false),
	)
})
