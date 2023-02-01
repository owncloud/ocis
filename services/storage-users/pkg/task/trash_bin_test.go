package task_test

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"time"

	apiGateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	apiUser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	apiRpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	apiProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	apiTypes "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/task"
	"github.com/stretchr/testify/mock"
)

func MustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

var _ = Describe("trash", func() {
	var (
		gw            *cs3mocks.GatewayAPIClient
		ctx           context.Context
		personalSpace *apiProvider.StorageSpace
		projectSpace  *apiProvider.StorageSpace
		virtualSpace  *apiProvider.StorageSpace
		user          *apiUser.User
	)

	BeforeEach(func() {
		gw = &cs3mocks.GatewayAPIClient{}
		ctx = context.Background()
		personalSpace = &apiProvider.StorageSpace{
			SpaceType: "personal",
			Root: &apiProvider.ResourceId{
				OpaqueId: "personal",
			},
		}
		projectSpace = &apiProvider.StorageSpace{
			SpaceType: "project",
			Root: &apiProvider.ResourceId{
				OpaqueId: "project",
			},
			Opaque: &apiTypes.Opaque{},
		}
		// virtual is here as an example,
		// the task ignores all space types expect `project` and `personal`.
		virtualSpace = &apiProvider.StorageSpace{
			SpaceType: "virtual",
			Root: &apiProvider.ResourceId{
				OpaqueId: "virtual",
			},
		}
		user = &apiUser.User{
			Id: &apiUser.UserId{
				OpaqueId: "user",
			},
		}

	})

	Describe("PurgeTrashBin", func() {
		It("throws an error if the user cannot authenticate", func() {
			expectedErr := errors.New("expected")
			gw.On("Authenticate", mock.Anything, mock.Anything).Return(nil, expectedErr)

			err := task.PurgeTrashBin(gw, "", "", time.Now())
			Expect(err).To(MatchError(expectedErr))
		})
		It("only deletes items older than the specified period", func() {
			var (
				now          = time.Now()
				recycleItems = map[string][]*apiProvider.RecycleItem{
					"personal": []*apiProvider.RecycleItem{
						&apiProvider.RecycleItem{Key: "now", DeletionTime: utils.TimeToTS(now)},
						&apiProvider.RecycleItem{Key: "after", DeletionTime: utils.TimeToTS(now.Add(1 * time.Second))},
						&apiProvider.RecycleItem{Key: "before", DeletionTime: utils.TimeToTS(now.Add(-1 * time.Second))},
					},
					"project": []*apiProvider.RecycleItem{
						&apiProvider.RecycleItem{Key: "now", DeletionTime: utils.TimeToTS(now)},
						&apiProvider.RecycleItem{Key: "after", DeletionTime: utils.TimeToTS(now.Add(1 * time.Minute))},
						&apiProvider.RecycleItem{Key: "before", DeletionTime: utils.TimeToTS(now.Add(-1 * time.Minute))},
					},
					"virtual": []*apiProvider.RecycleItem{
						&apiProvider.RecycleItem{Key: "now", DeletionTime: utils.TimeToTS(now)},
						&apiProvider.RecycleItem{Key: "after", DeletionTime: utils.TimeToTS(now.Add(1 * time.Hour))},
						&apiProvider.RecycleItem{Key: "before", DeletionTime: utils.TimeToTS(now.Add(-1 * time.Hour))},
					},
				}
			)

			personalSpace.Owner = user
			projectSpace.Opaque.Map = map[string]*apiTypes.OpaqueEntry{
				"grants": &apiTypes.OpaqueEntry{
					Value: MustMarshal(map[string]*apiProvider.ResourcePermissions{
						"admin": &apiProvider.ResourcePermissions{
							Delete: true,
						},
					}),
				},
			}

			gw.On("Authenticate", mock.Anything, mock.Anything).Return(
				&apiGateway.AuthenticateResponse{
					Status: status.NewOK(ctx),
					Token:  "",
				}, nil,
			)

			gw.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(
				&apiProvider.ListStorageSpacesResponse{
					Status: status.NewOK(ctx),
					StorageSpaces: []*apiProvider.StorageSpace{
						personalSpace,
						projectSpace,
						virtualSpace,
					},
				}, nil,
			)

			gw.On("ListRecycle", mock.Anything, mock.Anything).Return(
				func(_ context.Context, req *apiProvider.ListRecycleRequest, _ ...grpc.CallOption) *apiProvider.ListRecycleResponse {
					return &apiProvider.ListRecycleResponse{
						RecycleItems: recycleItems[req.Ref.ResourceId.OpaqueId],
					}
				}, nil,
			)

			gw.On("PurgeRecycle", mock.Anything, mock.Anything).Return(
				func(_ context.Context, req *apiProvider.PurgeRecycleRequest, _ ...grpc.CallOption) *apiProvider.PurgeRecycleResponse {
					var items []*apiProvider.RecycleItem
					for _, item := range recycleItems[req.Ref.ResourceId.OpaqueId] {
						if req.Key == item.Key {
							continue
						}

						items = append(items, item)
					}

					recycleItems[req.Ref.ResourceId.OpaqueId] = items

					return &apiProvider.PurgeRecycleResponse{
						Status: &apiRpc.Status{
							Code: apiRpc.Code_CODE_OK,
						},
					}
				}, nil,
			)

			err := task.PurgeTrashBin(gw, "", "", now)
			Expect(err).To(BeNil())
			Expect(recycleItems["personal"]).To(HaveLen(2))
			Expect(recycleItems["project"]).To(HaveLen(2))
			// virtual spaces are ignored
			Expect(recycleItems["virtual"]).To(HaveLen(3))
		})
	})
})
