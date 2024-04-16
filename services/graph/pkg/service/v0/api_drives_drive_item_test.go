package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaborationv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	svc "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("DrivesDriveItemService", func() {
	var (
		drivesDriveItemService svc.DrivesDriveItemService
		gatewayClient          *cs3mocks.GatewayAPIClient
		gatewaySelector        *mocks.Selectable[gateway.GatewayAPIClient]
	)

	BeforeEach(func() {
		logger := log.NewLogger()
		gatewayClient = cs3mocks.NewGatewayAPIClient(GinkgoT())

		gatewaySelector = mocks.NewSelectable[gateway.GatewayAPIClient](GinkgoT())
		gatewaySelector.On("Next").Return(gatewayClient, nil)

		cache := identity.NewIdentityCache(identity.IdentityCacheWithGatewaySelector(gatewaySelector))

		service, err := svc.NewDrivesDriveItemService(logger, gatewaySelector, cache)
		Expect(err).ToNot(HaveOccurred())
		drivesDriveItemService = service
	})

	Describe("UnmountShare", func() {
		It("handles gateway selector related errors", func() {
			gatewaySelector.ExpectedCalls = nil

			expectedError := errors.New("obtaining next gatewayClient failed")
			gatewaySelector.On("Next").Return(gatewayClient, expectedError)

			_, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "")
			Expect(err).To(MatchError(expectedError))
		})

		Describe("gateway client share listing", func() {
			It("handles share listing errors", func() {
				expectedError := errors.New("listing shares failed")
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(&collaborationv1beta1.ListReceivedSharesResponse{}, expectedError)

				_, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "")
				Expect(err).To(MatchError(expectedError))
			})

			It("uses the correct filters to get the shares", func() {
				expectedResourceID := storageprovider.ResourceId{
					StorageId: "1",
					OpaqueId:  "2",
					SpaceId:   "3",
				}
				expectedShareID := collaborationv1beta1.ShareId{
					OpaqueId: "1:2:3",
				}
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						Expect(in.Filters).To(HaveLen(3))

						var shareStates []collaborationv1beta1.ShareState
						var resourceIDs []*storageprovider.ResourceId

						for _, filter := range in.Filters {
							switch filter.Term.(type) {
							case *collaborationv1beta1.Filter_State:
								shareStates = append(shareStates, filter.GetState())
							case *collaborationv1beta1.Filter_ResourceId:
								resourceIDs = append(resourceIDs, filter.GetResourceId())
							}
						}

						Expect(shareStates).To(HaveLen(2))
						Expect(shareStates).To(ContainElements(
							collaborationv1beta1.ShareState_SHARE_STATE_PENDING,
							collaborationv1beta1.ShareState_SHARE_STATE_REJECTED,
						))

						Expect(resourceIDs).To(HaveLen(1))
						Expect(resourceIDs[0]).To(Equal(&expectedResourceID))

						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{
									State: collaborationv1beta1.ShareState_SHARE_STATE_PENDING,
									Share: &collaborationv1beta1.Share{
										Id: &expectedShareID,
									},
								},
							},
						}, nil
					})
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						Expect(in.GetUpdateMask().GetPaths()).To(Equal([]string{"state"}))
						Expect(in.GetShare().GetState()).To(Equal(collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED))
						Expect(in.GetShare().GetShare().GetId().GetOpaqueId()).To(Equal(expectedShareID.GetOpaqueId()))
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id:         &expectedShareID,
									ResourceId: &expectedResourceID,
								},
							},
						}, nil
					})

				_, err := drivesDriveItemService.MountShare(context.Background(), expectedResourceID, "")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("gateway client share update", func() {
			It("updates the share state to be accepted", func() {
				expectedShareID := collaborationv1beta1.ShareId{
					OpaqueId: "1:2:3",
				}
				expectedResourceID := storageprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}

				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{
									State: collaborationv1beta1.ShareState_SHARE_STATE_PENDING,
									Share: &collaborationv1beta1.Share{
										Id: &expectedShareID,
									},
								},
							},
						}, nil
					})

				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						Expect(in.GetUpdateMask().GetPaths()).To(Equal([]string{"state"}))
						Expect(in.GetShare().GetState()).To(Equal(collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED))
						Expect(in.GetShare().GetShare().GetId().GetOpaqueId()).To(Equal(expectedShareID.GetOpaqueId()))
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id:         &expectedShareID,
									ResourceId: &expectedResourceID,
								},
							},
						}, nil
					})

				_, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "")
				Expect(err).ToNot(HaveOccurred())
			})

			It("updates the mountPoint", func() {
				expectedShareID := collaborationv1beta1.ShareId{
					OpaqueId: "1:2:3",
				}
				expectedResourceID := storageprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}

				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{},
							},
						}, nil
					})

				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						Expect(in.GetUpdateMask().GetPaths()).To(HaveLen(2))
						Expect(in.GetUpdateMask().GetPaths()).To(ContainElements("mount_point"))
						Expect(in.GetShare().GetMountPoint().GetPath()).To(Equal("new name"))
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id:         &expectedShareID,
									ResourceId: &expectedResourceID,
								},
								MountPoint: in.GetShare().GetMountPoint(),
							},
						}, nil
					})

				di, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "new name")
				Expect(err).ToNot(HaveOccurred())
				Expect(di[0].GetMountPoint().GetPath()).To(Equal("new name"))
			})

			It("succeeds when any of the shares was accepted", func() {
				expectedResourceID := storageprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}

				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{},
								{},
								{},
							},
						}, nil
					})

				var calls int
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						calls++
						Expect(calls).To(BeNumerically("<=", 3))

						if calls <= 2 {
							return nil, fmt.Errorf("error %d", calls)
						}

						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id: &collaborationv1beta1.ShareId{
										OpaqueId: strconv.Itoa(calls),
									},
									ResourceId: &expectedResourceID,
								},
							},
						}, nil
					})

				di, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "new name")
				Expect(err).To(BeNil())
				Expect(di[0].GetShare().GetId().GetOpaqueId()).To(Equal(expectedResourceID.OpaqueId))
			})
			It("errors when none of the shares can be accepted", func() {
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{},
								{},
								{},
							},
						}, nil
					})

				var calls int
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						calls++
						Expect(calls).To(BeNumerically("<=", 3))
						return nil, fmt.Errorf("error %d", calls)
					})

				_, err := drivesDriveItemService.MountShare(context.Background(), storageprovider.ResourceId{}, "new name")
				Expect(fmt.Sprint(err)).To(ContainSubstring("error 1"))
				Expect(fmt.Sprint(err)).To(ContainSubstring("error 2"))
				Expect(fmt.Sprint(err)).To(ContainSubstring("error 3"))
			})
		})
	})

	Describe("UnmountShare", func() {
		It("handles gateway selector related errors", func() {
			gatewaySelector.ExpectedCalls = nil

			expectedError := errors.New("obtaining next gatewayClient failed")
			gatewaySelector.On("Next").Return(gatewayClient, expectedError)

			err := drivesDriveItemService.UnmountShare(context.Background(), &collaborationv1beta1.ShareId{})
			Expect(err).To(MatchError(expectedError))
		})

		Describe("gateway client share listing", func() {
			It("handles share listing errors", func() {
				expectedError := errorcode.New(errorcode.GeneralException, "listing shares failed")
				gatewayClient.
					On("GetReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(&collaborationv1beta1.GetReceivedShareResponse{}, errors.New("listing shares failed"))

				err := drivesDriveItemService.UnmountShare(context.Background(), &collaborationv1beta1.ShareId{OpaqueId: "1"})
				Expect(err).To(MatchError(expectedError))
			})

			It("uses the correct filters to get the shares", func() {
				shareID := &collaborationv1beta1.ShareId{
					OpaqueId: "3:4:5",
				}
				expectedResourceID := storageprovider.ResourceId{
					StorageId: "3",
					SpaceId:   "4",
					OpaqueId:  "5",
				}
				expectedShareID := collaborationv1beta1.ShareId{
					OpaqueId: "3:4:5",
				}
				gatewayClient.
					On("GetReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.GetReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.GetReceivedShareResponse, error) {
						Expect(in.Ref.GetId().GetOpaqueId()).To(Equal(shareID.GetOpaqueId()))
						return &collaborationv1beta1.GetReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id: &collaborationv1beta1.ShareId{
										OpaqueId: shareID.GetOpaqueId(),
									},
									ResourceId: &expectedResourceID,
								},
							},
						}, nil
					})

				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						Expect(in.Filters).To(HaveLen(2))

						var shareStates []collaborationv1beta1.ShareState
						var resourceIDs []*storageprovider.ResourceId

						for _, filter := range in.Filters {
							switch filter.Term.(type) {
							case *collaborationv1beta1.Filter_State:
								shareStates = append(shareStates, filter.GetState())
							case *collaborationv1beta1.Filter_ResourceId:
								resourceIDs = append(resourceIDs, filter.GetResourceId())
							}
						}

						Expect(shareStates).To(HaveLen(1))
						Expect(shareStates).To(ContainElements(
							collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
						))

						Expect(resourceIDs).To(HaveLen(1))
						Expect(resourceIDs[0]).To(Equal(&expectedResourceID))

						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{
									State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
									Share: &collaborationv1beta1.Share{
										Id: &expectedShareID,
									},
								},
							},
						}, nil
					})
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						Expect(in.GetUpdateMask().GetPaths()).To(Equal([]string{"state"}))
						Expect(in.GetShare().GetState()).To(Equal(collaborationv1beta1.ShareState_SHARE_STATE_REJECTED))
						Expect(in.GetShare().GetShare().GetId().GetOpaqueId()).To(Equal(expectedShareID.GetOpaqueId()))
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id:         &expectedShareID,
									ResourceId: &expectedResourceID,
								},
							},
						}, nil
					})

				err := drivesDriveItemService.UnmountShare(context.Background(), shareID)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("gateway client share update", func() {
			It("updates the share state to be rejected", func() {
				expectedShareID := collaborationv1beta1.ShareId{
					OpaqueId: "1$2!3",
				}
				gatewayClient.
					On("GetReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.GetReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.GetReceivedShareResponse, error) {
						return &collaborationv1beta1.GetReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
								Share: &collaborationv1beta1.Share{
									Id: &expectedShareID,
								},
							},
						}, nil
					})
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{
									State: collaborationv1beta1.ShareState_SHARE_STATE_PENDING,
									Share: &collaborationv1beta1.Share{
										Id: &expectedShareID,
									},
								},
							},
						}, nil
					})

				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						Expect(in.GetUpdateMask().GetPaths()).To(Equal([]string{"state"}))
						Expect(in.GetShare().GetState()).To(Equal(collaborationv1beta1.ShareState_SHARE_STATE_REJECTED))
						Expect(in.GetShare().GetShare().GetId().GetOpaqueId()).To(Equal(expectedShareID.GetOpaqueId()))
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: &cs3rpc.Status{Code: cs3rpc.Code_CODE_OK},
						}, nil
					})

				err := drivesDriveItemService.UnmountShare(context.Background(), &collaborationv1beta1.ShareId{OpaqueId: "1"})
				Expect(err).ToNot(HaveOccurred())
			})
			It("succeeds when all shares could be rejected", func() {
				gatewayClient.
					On("GetReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.GetReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.GetReceivedShareResponse, error) {
						return &collaborationv1beta1.GetReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
							},
						}, nil
					})
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{},
								{},
								{},
							},
						}, nil
					})

				var calls int
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						calls++
						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
						}, nil
					})

				err := drivesDriveItemService.UnmountShare(context.Background(), &collaborationv1beta1.ShareId{OpaqueId: "1"})
				Expect(calls).To(Equal(3))
				Expect(err).ToNot(HaveOccurred())
			})

			It("bubbles errors when any share fails rejecting", func() {
				gatewayClient.
					On("GetReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.GetReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.GetReceivedShareResponse, error) {
						return &collaborationv1beta1.GetReceivedShareResponse{
							Status: status.NewOK(ctx),
							Share: &collaborationv1beta1.ReceivedShare{
								State: collaborationv1beta1.ShareState_SHARE_STATE_ACCEPTED,
							},
						}, nil
					})
				gatewayClient.
					On("ListReceivedShares", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.ListReceivedSharesRequest, opts ...grpc.CallOption) (*collaborationv1beta1.ListReceivedSharesResponse, error) {
						return &collaborationv1beta1.ListReceivedSharesResponse{
							Shares: []*collaborationv1beta1.ReceivedShare{
								{},
								{},
								{},
							},
						}, nil
					})

				var calls int
				gatewayClient.
					On("UpdateReceivedShare", mock.Anything, mock.Anything, mock.Anything).
					Return(func(ctx context.Context, in *collaborationv1beta1.UpdateReceivedShareRequest, opts ...grpc.CallOption) (*collaborationv1beta1.UpdateReceivedShareResponse, error) {
						calls++
						Expect(calls).To(BeNumerically("<=", 3))

						if calls <= 2 {
							return nil, fmt.Errorf("error %d", calls)
						}

						return &collaborationv1beta1.UpdateReceivedShareResponse{
							Status: status.NewOK(ctx),
						}, nil
					})

				err := drivesDriveItemService.UnmountShare(context.Background(), &collaborationv1beta1.ShareId{OpaqueId: "1"})
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				Expect(errs).To(HaveLen(2))

				var err1 errorcode.Error
				Expect(errors.As(errs[0], &err1)).To(BeTrue())
				Expect(err1).To(MatchError(errorcode.New(errorcode.GeneralException, "error 1")))

				var err2 errorcode.Error
				Expect(errors.As(errs[1], &err2)).To(BeTrue())
				Expect(err2).To(MatchError(errorcode.New(errorcode.GeneralException, "error 2")))
			})
		})
	})
})

var _ = Describe("DrivesDriveItemApi", func() {
	var (
		mockDrivesDriveItemService *mocks.DrivesDriveItemProvider
		mockBaseGraphService       *mocks.BaseGraphProvider
		httpAPI                    svc.DrivesDriveItemApi
		rCTX                       *chi.Context
	)

	BeforeEach(func() {
		logger := log.NewLogger()

		mockBaseGraphService = mocks.NewBaseGraphProvider(GinkgoT())

		mockDrivesDriveItemService = mocks.NewDrivesDriveItemProvider(GinkgoT())
		api, err := svc.NewDrivesDriveItemApi(mockDrivesDriveItemService, mockBaseGraphService, logger)
		Expect(err).ToNot(HaveOccurred())

		httpAPI = api

		rCTX = chi.NewRouteContext()
		rCTX.URLParams.Add("driveID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668")
	})

	checkDriveIDAndItemIDValidation := func(handler http.HandlerFunc) {
		rCTX.URLParams.Add("driveID", "1$2")
		rCTX.URLParams.Add("itemID", "3$4!5")

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", nil).
			WithContext(
				context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
			)

		handler(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

		jsonData := gjson.Get(responseRecorder.Body.String(), "error")
		Expect(jsonData.Get("message").String()).To(Equal("invalid driveID or itemID"))
	}

	Describe("DeleteDriveItem", func() {
		It("validates the driveID and itemID url param", func() {
			checkDriveIDAndItemIDValidation(httpAPI.DeleteDriveItem)
		})

		It("uses the UnmountShare provider implementation", func() {
			rCTX.URLParams.Add("itemID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!a0ca6a90-a365-4782-871e-d44447bbc668")
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodDelete, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onUnmountShare := mockDrivesDriveItemService.On("UnmountShare", mock.Anything, mock.Anything)
			onUnmountShare.
				Return(func(ctx context.Context, shareID *collaborationv1beta1.ShareId) error {
					return errors.New("any")
				}).Once()

			httpAPI.DeleteDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusFailedDependency))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("unmounting share failed"))

			// happy path
			responseRecorder = httptest.NewRecorder()

			onUnmountShare.
				Return(func(ctx context.Context, shareID *collaborationv1beta1.ShareId) error {
					Expect(shareID.GetOpaqueId()).To(Equal("a0ca6a90-a365-4782-871e-d44447bbc668"))
					return nil
				}).Once()

			httpAPI.DeleteDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNoContent))
		})
	})

	Describe("CreateDriveItem", func() {
		It("checks if the idemID and driveID is in share jail", func() {
			rCTX.URLParams.Add("driveID", "1$2")

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal(svc.ErrNotAShareJail.Error()))
		})

		It("checks that the request body is valid", func() {
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("invalid request body"))

			// valid drive item, but invalid remote item id
			driveItem := libregraph.DriveItem{}

			driveItemJson, err := json.Marshal(driveItem)
			Expect(err).ToNot(HaveOccurred())

			responseRecorder = httptest.NewRecorder()

			request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))

			jsonData = gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("invalid remote item id"))
		})

		It("uses the MountShare provider implementation", func() {
			driveItemName := "a name"
			remoteItemID := "d66d28d8-3558-4f0f-ba2a-34a7185b806d$831997cf-a531-491b-ae72-9037739f04e9!c131a84c-7506-46b4-8e5e-60c56382da3b"
			driveItem := libregraph.DriveItem{
				Name: &driveItemName,
				RemoteItem: &libregraph.RemoteItem{
					Id: &remoteItemID,
				},
			}

			driveItemJson, err := json.Marshal(driveItem)
			Expect(err).ToNot(HaveOccurred())

			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onMountShare := mockDrivesDriveItemService.On("MountShare", mock.Anything, mock.Anything, mock.Anything)
			onMountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId, name string) ([]*collaborationv1beta1.ReceivedShare, error) {
					return nil, errors.New("any")
				}).Once()

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("mounting share failed"))

			// happy path
			responseRecorder = httptest.NewRecorder()

			request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onMountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId, name string) ([]*collaborationv1beta1.ReceivedShare, error) {
					Expect(storagespace.FormatResourceID(resourceID)).To(Equal(remoteItemID))
					Expect(driveItemName).To(Equal(name))
					return nil, nil
				}).Once()

			mockBaseGraphService.On("CS3ReceivedSharesToDriveItems", mock.Anything, mock.Anything).
				Return([]libregraph.DriveItem{{}}, nil).Once()

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
	})

	Describe("UpdateDriveItem", func() {
		It("fails without a valid driveID || itemID", func() {
			checkDriveIDAndItemIDValidation(httpAPI.UpdateDriveItem)
		})

		It("fails without a valid shareJail", func() {
			rCTX.URLParams.Add("driveID", "1$2!3")
			rCTX.URLParams.Add("itemID", "1$2!3")
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPut, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.UpdateDriveItem(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("fails without a valid request body", func() {
			rCTX.URLParams.Add("driveID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668")
			rCTX.URLParams.Add("itemID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!1")
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPut, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.UpdateDriveItem(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("invalid request body"))
		})

		It("updates the share", func() {
			rCTX.URLParams.Add("driveID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668")
			rCTX.URLParams.Add("itemID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!1")

			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"@UI.Hidden":true}`)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onUnmountShare := mockDrivesDriveItemService.On("UpdateShare", mock.Anything, mock.Anything, mock.Anything)
			onUnmountShare.
				Return(func(context.Context, *collaborationv1beta1.ShareId, *svc.ShareUpdateInstruction) (*collaborationv1beta1.ReceivedShare, error) {
					return nil, errors.New("any")
				}).Once()

			httpAPI.UpdateDriveItem(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusFailedDependency))

			onUnmountShare.
				Return(func(context.Context, *collaborationv1beta1.ShareId, *svc.ShareUpdateInstruction) (*collaborationv1beta1.ReceivedShare, error) {
					return &collaborationv1beta1.ReceivedShare{
						Hidden: true,
					}, nil
				}).Once()

			responseRecorder = httptest.NewRecorder()

			request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"@UI.Hidden":true}`)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.UpdateDriveItem(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
	})
})
