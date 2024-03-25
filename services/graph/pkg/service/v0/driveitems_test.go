package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	roleconversions "github.com/cs3org/reva/v2/pkg/conversions"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/linktype"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

type itemsList struct {
	Value []*libregraph.DriveItem
}

var _ = Describe("Driveitems", func() {
	var (
		svc                    service.Service
		ctx                    context.Context
		cfg                    *config.Config
		gatewayClient          *cs3mocks.GatewayAPIClient
		gatewaySelector        pool.Selectable[gateway.GatewayAPIClient]
		eventsPublisher        mocks.Publisher
		identityBackend        *identitymocks.Backend
		getPublicShareResponse *link.GetPublicShareResponse
		getShareResponse       *collaboration.GetShareResponse
		listSpacesResponse     *provider.ListStorageSpacesResponse

		rr *httptest.ResponseRecorder

		newGroup *libregraph.Group

		currentUser = &userpb.User{
			Id: &userpb.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)
		getPublicShareResponse = &link.GetPublicShareResponse{
			Status: status.NewNotFound(ctx, "not found"),
		}
		getShareResponse = &collaboration.GetShareResponse{
			Status: status.NewNotFound(ctx, "not found"),
		}

		grantMapJSON, _ := json.Marshal(
			map[string]*provider.ResourcePermissions{
				"userid": roleconversions.NewSpaceViewerRole().CS3ResourcePermissions(),
			},
		)
		spaceOpaque := &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"grants": {
					Decoder: "json",
					Value:   grantMapJSON,
				},
			},
		}
		listSpacesResponse = &provider.ListStorageSpacesResponse{
			Status: status.NewOK(ctx),
			StorageSpaces: []*provider.StorageSpace{
				{
					Id: &provider.StorageSpaceId{
						OpaqueId: "2",
					},
					Opaque: spaceOpaque,
				},
			},
		}

		identityBackend = &identitymocks.Backend{}
		newGroup = libregraph.NewGroup()
		newGroup.SetMembersodataBind([]string{"/users/user1"})
		newGroup.SetId("group1")

		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewaySelector(gatewaySelector),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("DeletePermission", func() {
		It("deletes a user permission as expected", func() {
			gatewayClient.On("GetPublicShare", mock.Anything, mock.Anything).Return(getPublicShareResponse, nil)

			getShareResponse.Status = status.NewOK(ctx)
			getShareResponse.Share = &collaboration.Share{
				Id: &collaboration.ShareId{
					OpaqueId: "permissionid",
				},
				ResourceId: &provider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				},
			}
			gatewayClient.On("GetShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.GetShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			).Return(getShareResponse, nil)

			rmShareMock := gatewayClient.On("RemoveShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.RemoveShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			rmShareMockResponse := &collaboration.RemoveShareResponse{
				Status: status.NewOK(ctx),
			}
			rmShareMock.Return(rmShareMockResponse, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "permissionid")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.DeletePermission(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
		It("deletes a link permission as expected", func() {
			getPublicShareMock := gatewayClient.On("GetPublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.GetPublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "linkpermissionid"
				}),
			)
			getPublicShareMock.Return(&link.GetPublicShareResponse{
				Status: status.NewOK(ctx),
				Share: &link.PublicShare{
					Id: &link.PublicShareId{
						OpaqueId: "permissionid",
					},
					ResourceId: &provider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					},
				},
			}, nil)

			rmPublicShareMock := gatewayClient.On("RemovePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.RemovePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "linkpermissionid"
				}),
			)
			rmPublicShareMockResponse := &link.RemovePublicShareResponse{
				Status: status.NewOK(ctx),
			}
			rmPublicShareMock.Return(rmPublicShareMockResponse, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "linkpermissionid")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.DeletePermission(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
		It("deletes a space permission as expected", func() {
			gatewayClient.On("GetPublicShare", mock.Anything, mock.Anything).Return(getPublicShareResponse, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!2")
			rctx.URLParams.Add("permissionID", "u:userid")

			gatewayClient.On("RemoveShare",
				mock.Anything,
				mock.Anything,
			).Return(func(ctx context.Context, in *collaboration.RemoveShareRequest, opts ...grpc.CallOption) (*collaboration.RemoveShareResponse, error) {
				Expect(in.Ref.GetKey()).ToNot(BeNil())
				Expect(in.Ref.GetKey().GetGrantee().GetUserId().GetOpaqueId()).To(Equal("userid"))
				return &collaboration.RemoveShareResponse{Status: status.NewOK(ctx)}, nil
			})

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.DeletePermission(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})

		It("fails to delete permission when the item id does not match the shared resource's id", func() {
			gatewayClient.On("GetPublicShare", mock.Anything, mock.Anything).Return(getPublicShareResponse, nil)
			getShareResponse.Status = status.NewOK(ctx)
			getShareResponse.Share = &collaboration.Share{
				Id: &collaboration.ShareId{
					OpaqueId: "permissionid",
				},
				ResourceId: &provider.ResourceId{
					StorageId: "3",
					SpaceId:   "4",
					OpaqueId:  "5",
				},
			}
			gatewayClient.On("GetShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.GetShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			).Return(getShareResponse, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "permissionid")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.DeletePermission(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("UpdatePermission", func() {
		var (
			driveItemPermission           *libregraph.Permission
			getShareMockResponse          *collaboration.GetShareResponse
			getPublicShareMockResponse    *link.GetPublicShareResponse
			getUserMockResponse           *userpb.GetUserResponse
			updateShareMockResponse       *collaboration.UpdateShareResponse
			updatePublicShareMockResponse *link.UpdatePublicShareResponse
		)
		const TestLinkName = "Test Link"
		BeforeEach(func() {
			rr = httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "permissionid")

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)

			driveItemPermission = &libregraph.Permission{}

			getUserMock := gatewayClient.On("GetUser", mock.Anything, mock.Anything)
			getUserMockResponse = &userpb.GetUserResponse{
				Status: status.NewOK(ctx),
				User: &userpb.User{
					Id:          &userpb.UserId{OpaqueId: "useri"},
					DisplayName: "Test User",
				},
			}
			getUserMock.Return(getUserMockResponse, nil)

			getShareMock := gatewayClient.On("GetShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.GetShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			share := &collaboration.Share{
				Id: &collaboration.ShareId{
					OpaqueId: "permissionid",
				},
				ResourceId: &provider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				},
				Grantee: &provider.Grantee{
					Type: provider.GranteeType_GRANTEE_TYPE_USER,
					Id: &provider.Grantee_UserId{
						UserId: &userpb.UserId{
							OpaqueId: "userid",
						},
					},
				},
				Permissions: &collaboration.SharePermissions{
					Permissions: roleconversions.NewViewerRole(true).CS3ResourcePermissions(),
				},
			}
			getShareMockResponse = &collaboration.GetShareResponse{
				Status: status.NewOK(ctx),
				Share:  share,
			}
			getShareMock.Return(getShareMockResponse, nil)

			updateShareMockResponse = &collaboration.UpdateShareResponse{
				Status: status.NewOK(ctx),
				Share:  share,
			}

			updatePublicShareMockResponse = &link.UpdatePublicShareResponse{
				Status: status.NewOK(ctx),
				Share:  &link.PublicShare{DisplayName: TestLinkName},
			}

			getPublicShareMock := gatewayClient.On("GetPublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.GetPublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			getPublicShareMockResponse = &link.GetPublicShareResponse{
				Status: status.NewOK(ctx),
				Share: &link.PublicShare{
					Id: &link.PublicShareId{
						OpaqueId: "permissionid",
					},
					ResourceId: &provider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					},
					Permissions: &link.PublicSharePermissions{
						Permissions: linktype.NewViewLinkPermissionSet().GetPermissions(),
					},
					Token: "token",
				},
			}
			getPublicShareMock.Return(getPublicShareMockResponse, nil)

			statMock := gatewayClient.On("Stat",
				mock.Anything,
				mock.MatchedBy(func(req *provider.StatRequest) bool {
					return utils.ResourceIDEqual(
						req.GetRef().GetResourceId(),
						&provider.ResourceId{
							StorageId: "1",
							SpaceId:   "2",
							OpaqueId:  "3",
						},
					) && req.GetRef().GetPath() == "."
				}))

			statResponse := &provider.StatResponse{
				Status: status.NewOK(ctx),
				Info: &provider.ResourceInfo{
					Id: &provider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					},
					Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				},
			}
			statMock.Return(statResponse, nil)
			spaceRootStatMock := gatewayClient.On("Stat",
				mock.Anything,
				mock.MatchedBy(func(req *provider.StatRequest) bool {
					return utils.ResourceIDEqual(
						req.GetRef().GetResourceId(),
						&provider.ResourceId{
							StorageId: "1",
							SpaceId:   "2",
							OpaqueId:  "2",
						},
					)
				}))

			spaceRootStatMock.Return(
				&provider.StatResponse{
					Status: status.NewOK(ctx),
					Info: &provider.ResourceInfo{
						Id: &provider.ResourceId{
							StorageId: "1",
							SpaceId:   "2",
							OpaqueId:  "2",
						},
						Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER,
					},
				}, nil)

		})
		It("fails when no share is found", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})
		It("fails to update permission when no request body is sent", func() {
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("fails to update password when no request body is sent", func() {
			svc.SetLinkPassword(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("fails to update password when itemID mismatches with the driveID", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$4!3")
			rctx.URLParams.Add("permissionID", "permissionid")

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.SetLinkPassword(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Error.Code).To(Equal(errorcode.ItemNotFound.String()))
			Expect(res.Error.Message).To(Equal("driveID and itemID do not match"))
		})
		It("fails to update permission when itemID mismatches with the driveID", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$4!3")
			rctx.URLParams.Add("permissionID", "permissionid")

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Error.Code).To(Equal(errorcode.ItemNotFound.String()))
			Expect(res.Error.Message).To(Equal("driveID and itemID do not match"))
		})
		It("fails to update permission when the resourceID mismatches with the shared resource's id", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			getPublicShareMockResponse.Share.ResourceId = &provider.ResourceId{
				StorageId: "1",
				SpaceId:   "2",
				OpaqueId:  "4",
			}

			driveItemPermission.SetExpirationDateTime(time.Now().Add(time.Hour))
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
		})
		It("fails to update public link password when the permissionID is not parseable", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "permi%ssionid")

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.SetLinkPassword(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Error.Code).To(Equal(errorcode.InvalidRequest.String()))
			Expect(res.Error.Message).To(Equal("invalid permissionID"))
		})
		It("fails to update permission when the permissionID is not parseable", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			rctx.URLParams.Add("itemID", "1$2!3")
			rctx.URLParams.Add("permissionID", "permi%ssionid")

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Error.Code).To(Equal(errorcode.InvalidRequest.String()))
			Expect(res.Error.Message).To(Equal("invalid permissionID"))
		})
		It("succeeds when trying to update a link permission with displayname", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					if req.GetRef().GetId().GetOpaqueId() == "permissionid" {
						return req.Update.GetDisplayName() == TestLinkName
					}
					return false
				}),
			)

			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			link := libregraph.NewSharingLink()
			link.SetLibreGraphDisplayName(TestLinkName)

			driveItemPermission.SetLink(*link)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Link).ToNot(BeNil())
			Expect(res.Link.GetLibreGraphDisplayName() == TestLinkName)
		})
		It("fails updating the id", func() {
			driveItemPermission.SetId("permissionid")
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("fails updating the password flag", func() {
			driveItemPermission.SetHasPassword(true)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("updates the expiration date", func() {
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			expiration := time.Now().Add(time.Hour)
			updateShareMock := gatewayClient.On("UpdateShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
					if req.GetShare().GetId().GetOpaqueId() == "permissionid" {
						return expiration.Equal(utils.TSToTime(req.GetShare().GetExpiration()))
					}
					return false
				}),
			)
			updateShareMockResponse.Share.Expiration = utils.TimeToTS(expiration)
			updateShareMock.Return(updateShareMockResponse, nil)

			driveItemPermission.SetExpirationDateTime(expiration)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.GetExpirationDateTime().Equal(expiration)).To(BeTrue())
		})
		It("deletes the expiration date", func() {
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			updateShareMock := gatewayClient.On("UpdateShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
					if req.GetShare().GetId().GetOpaqueId() == "permissionid" {
						return true
					}
					return false
				}),
			)
			updateShareMock.Return(updateShareMockResponse, nil)

			driveItemPermission.SetExpirationDateTimeNil()
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			_, ok := res.GetExpirationDateTimeOk()
			Expect(ok).To(BeFalse())
		})
		// that is resharing test. Please delete after disable resharing feature
		
		// It("updates the share permissions with changing the role", func() {
		// 	getPublicShareMockResponse.Share = nil
		// 	getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
		// 	updateShareMock := gatewayClient.On("UpdateShare",
		// 		mock.Anything,
		// 		mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
		// 			return req.GetShare().GetId().GetOpaqueId() == "permissionid"
		// 		}),
		// 	)
		// 	updateShareMock.Return(updateShareMockResponse, nil)
		// 	driveItemPermission.SetRoles([]string{unifiedrole.NewViewerUnifiedRole(false).GetId()})
		// 	body, err := driveItemPermission.MarshalJSON()
		// 	Expect(err).To(BeNil())
		// 	svc.UpdatePermission(
		// 		rr,
		// 		httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
		// 			WithContext(ctx),
		// 	)
		// 	Expect(rr.Code).To(Equal(http.StatusOK))
		// 	data, err := io.ReadAll(rr.Body)
		// 	Expect(err).ToNot(HaveOccurred())

		// 	res := libregraph.Permission{}

		// 	err = json.Unmarshal(data, &res)
		// 	Expect(err).ToNot(HaveOccurred())
		// 	_, ok := res.GetRolesOk()
		// 	Expect(ok).To(BeTrue())
		// })
		It("fails to update the share permissions for a file share when setting a space specific role", func() {
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			updateShareMock := gatewayClient.On("UpdateShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
					return req.GetShare().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			updateShareMock.Return(updateShareMockResponse, nil)

			driveItemPermission.SetRoles([]string{unifiedrole.NewSpaceViewerUnifiedRole().GetId()})
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("fails to update the space permissions for a space share when setting a file specific role", func() {
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			gatewayClient.On("GetPublicShare",
				mock.Anything,
				mock.Anything,
			).Return(getPublicShareMockResponse, nil)
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(listSpacesResponse, nil)

			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			updateShareMock := gatewayClient.On("UpdateShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
					return req.GetShare().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			updateShareMock.Return(updateShareMockResponse, nil)

			driveItemPermission.SetRoles([]string{unifiedrole.NewFileEditorUnifiedRole(false).GetId()})
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2")
			// This is a space root
			rctx.URLParams.Add("itemID", "1$2!2")
			rctx.URLParams.Add("permissionID", "u:userid")
			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("updates the share permissions when changing the resource permission actions", func() {
			getPublicShareMockResponse.Share = nil
			getPublicShareMockResponse.Status = status.NewNotFound(ctx, "not found")
			updateShareMock := gatewayClient.On("UpdateShare",
				mock.Anything,
				mock.MatchedBy(func(req *collaboration.UpdateShareRequest) bool {
					return req.GetShare().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			updateShareMockResponse.Share.Permissions = &collaboration.SharePermissions{
				Permissions: &provider.ResourcePermissions{
					GetPath: true,
				},
			}

			updateShareMock.Return(updateShareMockResponse, nil)

			driveItemPermission.SetLibreGraphPermissionsActions([]string{unifiedrole.DriveItemPathRead})
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			_, ok := res.GetRolesOk()
			Expect(ok).To(BeFalse())
			_, ok = res.GetLibreGraphPermissionsActionsOk()
			Expect(ok).To(BeTrue())
		})
		It("updates the expiration date on a public share", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			expiration := time.Now().UTC().Add(time.Hour)
			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)

			updatePublicShareMockResponse.Share.Expiration = utils.TimeToTS(expiration)
			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			driveItemPermission.SetExpirationDateTime(expiration)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.GetExpirationDateTime().Equal(expiration)).To(BeTrue())
		})
		It("updates the permissions on a public share", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			newLink := libregraph.NewSharingLink()
			newLinkType, err := libregraph.NewSharingLinkTypeFromValue("edit")
			Expect(err).ToNot(HaveOccurred())
			newLink.SetType(*newLinkType)

			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)

			updatePublicShareMockResponse.Share.Permissions = &link.PublicSharePermissions{
				Permissions: linktype.NewFolderEditLinkPermissionSet().Permissions,
			}
			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			driveItemPermission.SetLink(*newLink)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			linkType := res.Link.GetType()
			Expect(string(linkType)).To(Equal("edit"))
		})
		It("updates the public share to internal link", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			getPublicShareMock := gatewayClient.On("GetPublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.GetPublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)
			getPublicShareMockResponse = &link.GetPublicShareResponse{
				Status: status.NewOK(ctx),
				Share: &link.PublicShare{
					Id: &link.PublicShareId{
						OpaqueId: "permissionid",
					},
					ResourceId: &provider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					},
					PasswordProtected: true,
					Permissions: &link.PublicSharePermissions{
						Permissions: linktype.NewFileEditLinkPermissionSet().GetPermissions(),
					},
					Token: "token",
				},
			}
			getPublicShareMock.Return(getPublicShareMockResponse, nil)

			newLink := libregraph.NewSharingLink()
			newLinkType, err := libregraph.NewSharingLinkTypeFromValue("internal")
			Expect(err).ToNot(HaveOccurred())
			newLink.SetType(*newLinkType)

			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)

			updatePublicShareMockResponse.Share.Permissions = &link.PublicSharePermissions{
				Permissions: linktype.NewInternalLinkPermissionSet().Permissions,
			}
			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			driveItemPermission.SetLink(*newLink)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			linkType := res.Link.GetType()
			Expect(string(linkType)).To(Equal("internal"))
			pp, hasPP := res.GetHasPasswordOk()
			Expect(hasPP).To(Equal(true))
			Expect(*pp).To(Equal(false))
		})
		It("updates the password on a public share", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			newLinkPassword := libregraph.NewSharingLinkPassword()
			newLinkPassword.SetPassword("OC123!")

			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)

			updatePublicShareMockResponse.Share.Permissions = &link.PublicSharePermissions{
				Permissions: linktype.NewViewLinkPermissionSet().Permissions,
			}
			updatePublicShareMockResponse.Share.PasswordProtected = true
			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			body, err := newLinkPassword.MarshalJSON()
			Expect(err).To(BeNil())
			svc.SetLinkPassword(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.Permission{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			linkType := res.Link.GetType()
			Expect(string(linkType)).To(Equal("view"))
			Expect(*res.HasPassword).To(BeTrue())
		})
		It("fails when updating the expiration date on a public share", func() {
			getShareMockResponse.Share = nil
			getShareMockResponse.Status = status.NewNotFound(ctx, "not found")

			expiration := time.Now().UTC().AddDate(0, 0, -1)
			updatePublicShareMock := gatewayClient.On("UpdatePublicShare",
				mock.Anything,
				mock.MatchedBy(func(req *link.UpdatePublicShareRequest) bool {
					return req.GetRef().GetId().GetOpaqueId() == "permissionid"
				}),
			)

			updatePublicShareMockResponse.Share = nil
			updatePublicShareMockResponse.Status = status.NewFailedPrecondition(ctx, nil, "expiration date is in the past")
			updatePublicShareMock.Return(updatePublicShareMockResponse, nil)

			driveItemPermission.SetExpirationDateTime(expiration)
			body, err := driveItemPermission.MarshalJSON()
			Expect(err).To(BeNil())
			svc.UpdatePermission(
				rr,
				httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(body))).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := libregraph.OdataError{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.GetError().Code).To(Equal(errorcode.InvalidRequest.String()))
			Expect(res.GetError().Message).To(Equal("expiration date is in the past"))
		})
	})

	Describe("GetRootDriveChildren", func() {
		It("handles ListStorageSpaces not found", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListStorageSpaces error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewInternal(ctx, "internal error"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("handles ListContainer not found", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("succeeds", func() {
			mtime := time.Now()
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewOK(ctx),
				Infos: []*provider.ResourceInfo{
					{
						Type:  provider.ResourceType_RESOURCE_TYPE_FILE,
						Id:    &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
						Etag:  "etag",
						Mtime: utils.TimeToTS(mtime),
					},
				},
			}, nil)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
			Expect(res.Value[0].GetETag()).To(Equal("etag"))
			Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
		})
	})

	Describe("GetDriveItemChildren", func() {
		It("handles ListContainer not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied as not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		Context("it succeeds", func() {
			var (
				r     *http.Request
				mtime = time.Now()
			)

			BeforeEach(func() {
				r = httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("driveID", "storageid$spaceid")
				rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			})

			assertItemsList := func(length int) itemsList {
				svc.GetDriveItemChildren(rr, r)
				Expect(rr.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(rr.Body)
				Expect(err).ToNot(HaveOccurred())

				res := itemsList{}

				err = json.Unmarshal(data, &res)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(res.Value)).To(Equal(1))
				Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
				Expect(res.Value[0].GetETag()).To(Equal("etag"))
				Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
				Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))

				return res
			}

			It("returns a generic file", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:              provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:                &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:              "etag",
							Mtime:             utils.TimeToTS(mtime),
							ArbitraryMetadata: nil,
						},
					},
				}, nil)

				res := assertItemsList(1)
				Expect(res.Value[0].Audio).To(BeNil())
				Expect(res.Value[0].Location).To(BeNil())
			})

			It("returns the audio facet if metadata is available", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:       &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:     "etag",
							Mtime:    utils.TimeToTS(mtime),
							MimeType: "audio/mpeg",
							ArbitraryMetadata: &provider.ArbitraryMetadata{
								Metadata: map[string]string{
									"libre.graph.audio.album":             "Some Album",
									"libre.graph.audio.albumArtist":       "Some AlbumArtist",
									"libre.graph.audio.artist":            "Some Artist",
									"libre.graph.audio.bitrate":           "192",
									"libre.graph.audio.composers":         "Some Composers",
									"libre.graph.audio.copyright":         "Some Copyright",
									"libre.graph.audio.disc":              "2",
									"libre.graph.audio.discCount":         "5",
									"libre.graph.audio.duration":          "225000",
									"libre.graph.audio.genre":             "Some Genre",
									"libre.graph.audio.hasDrm":            "false",
									"libre.graph.audio.isVariableBitrate": "true",
									"libre.graph.audio.title":             "Some Title",
									"libre.graph.audio.track":             "6",
									"libre.graph.audio.trackCount":        "9",
									"libre.graph.audio.year":              "1994",
								},
							},
						},
					},
				}, nil)

				res := assertItemsList(1)
				audio := res.Value[0].Audio

				Expect(audio).ToNot(BeNil())
				Expect(audio.Album).To(Equal(libregraph.PtrString("Some Album")))
				Expect(audio.AlbumArtist).To(Equal(libregraph.PtrString("Some AlbumArtist")))
				Expect(audio.Artist).To(Equal(libregraph.PtrString("Some Artist")))
				Expect(audio.Bitrate).To(Equal(libregraph.PtrInt64(192)))
				Expect(audio.Composers).To(Equal(libregraph.PtrString("Some Composers")))
				Expect(audio.Copyright).To(Equal(libregraph.PtrString("Some Copyright")))
				Expect(audio.Disc).To(Equal(libregraph.PtrInt32(2)))
				Expect(audio.DiscCount).To(Equal(libregraph.PtrInt32(5)))
				Expect(audio.Duration).To(Equal(libregraph.PtrInt64(225000)))
				Expect(audio.Genre).To(Equal(libregraph.PtrString("Some Genre")))
				Expect(audio.HasDrm).To(Equal(libregraph.PtrBool(false)))
				Expect(audio.IsVariableBitrate).To(Equal(libregraph.PtrBool(true)))
				Expect(audio.Title).To(Equal(libregraph.PtrString("Some Title")))
				Expect(audio.Track).To(Equal(libregraph.PtrInt32(6)))
				Expect(audio.TrackCount).To(Equal(libregraph.PtrInt32(9)))
				Expect(audio.Year).To(Equal(libregraph.PtrInt32(1994)))
			})

			It("returns the location facet if metadata is available", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:       &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:     "etag",
							Mtime:    utils.TimeToTS(mtime),
							MimeType: "image/jpeg",
							ArbitraryMetadata: &provider.ArbitraryMetadata{
								Metadata: map[string]string{
									"libre.graph.location.altitude":  "1047.7",
									"libre.graph.location.latitude":  "49.48675890884328",
									"libre.graph.location.longitude": "11.103870357204285",
								},
							},
						},
					},
				}, nil)

				res := assertItemsList(1)
				location := res.Value[0].Location

				Expect(location).ToNot(BeNil())
				Expect(location.Altitude).To(Equal(libregraph.PtrFloat64(1047.7)))
				Expect(location.Latitude).To(Equal(libregraph.PtrFloat64(49.48675890884328)))
				Expect(location.Longitude).To(Equal(libregraph.PtrFloat64(11.103870357204285)))
			})
		})
	})
})
