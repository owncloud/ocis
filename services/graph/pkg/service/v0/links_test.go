package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/linktype"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var _ = Describe("createLinkTests", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend
		currentUser     = &userpb.User{
			Id: &userpb.UserId{
				OpaqueId: "user",
			},
		}

		rr *httptest.ResponseRecorder
	)
	const (
		ViewerLinkString = "Viewer Link"
		ItemID           = "f0042750-23c5-441c-9f2c-ff7c53e5bd2a$cd621428-dfbe-44c1-9393-65bf0dd440a6!1177add3-b4eb-434e-a2e8-1859b31b17bf"
		DriveId          = "f0042750-23c5-441c-9f2c-ff7c53e5bd2a$cd621428-dfbe-44c1-9393-65bf0dd440a6!cd621428-dfbe-44c1-9393-65bf0dd440a6"
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		rr = httptest.NewRecorder()

		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		identityBackend = &identitymocks.Backend{}

		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
		cfg.FilesSharing.EnableResharing = true

		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewaySelector(gatewaySelector),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("CreateLink", func() {
		var (
			itemID              string
			driveItemCreateLink *libregraph.DriveItemCreateLink
			statMock            *mock.Call
			statResponse        *provider.StatResponse
			createLinkResponse  *link.CreatePublicShareResponse
			createLinkMock      *mock.Call
		)

		BeforeEach(func() {
			itemID = ItemID
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", DriveId)
			rctx.URLParams.Add("itemID", itemID)

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)

			driveItemCreateLink = &libregraph.DriveItemCreateLink{
				Type:                nil,
				ExpirationDateTime:  nil,
				Password:            nil,
				DisplayName:         nil,
				LibreGraphQuickLink: nil,
			}

			statMock = gatewayClient.On("Stat", mock.Anything, mock.Anything)
			statResponse = &provider.StatResponse{
				Status: status.NewOK(ctx),
			}
			statMock.Return(statResponse, nil)

			createLinkMock = gatewayClient.On("CreatePublicShare", mock.Anything, mock.Anything)
			createLinkResponse = &link.CreatePublicShareResponse{
				Status: status.NewOK(ctx),
			}
			createLinkMock.Return(createLinkResponse, nil)

			linkType, err := libregraph.NewSharingLinkTypeFromValue("view")
			Expect(err).ToNot(HaveOccurred())
			driveItemCreateLink.Type = linkType
			driveItemCreateLink.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			permissions, err := linktype.CS3ResourcePermissionsFromSharingLink(*driveItemCreateLink, provider.ResourceType_RESOURCE_TYPE_CONTAINER)
			Expect(err).ToNot(HaveOccurred())
			createLinkResponse.Share = &link.PublicShare{
				Id:                &link.PublicShareId{OpaqueId: "123"},
				Expiration:        utils.TimeToTS(*driveItemCreateLink.ExpirationDateTime),
				PasswordProtected: false,
				DisplayName:       ViewerLinkString,
				Token:             "SomeGOODCoffee",
				Permissions:       &link.PublicSharePermissions{Permissions: permissions},
			}

			statResponse.Info = &provider.ResourceInfo{Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER}
		})

		toJSONReader := func(v any) *strings.Reader {
			driveItemInviteBytes, err := json.Marshal(v)
			Expect(err).ToNot(HaveOccurred())

			return strings.NewReader(string(driveItemInviteBytes))
		}

		// Public Shares / "links" in graph terms
		It("creates a public link as expected (happy path)", func() {
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemCreateLink)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var createLinkResponseBody *libregraph.Permission
			err := json.Unmarshal(rr.Body.Bytes(), &createLinkResponseBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(createLinkResponseBody.GetId()).To(Equal("123"))
			Expect(createLinkResponseBody.GetExpirationDateTime().Unix()).To(Equal(driveItemCreateLink.ExpirationDateTime.Unix()))
			Expect(createLinkResponseBody.GetHasPassword()).To(Equal(false))
			Expect(createLinkResponseBody.GetLink().LibreGraphDisplayName).To(Equal(libregraph.PtrString(ViewerLinkString)))
			link := createLinkResponseBody.GetLink()
			respLinkType := link.GetType()
			expected, err := libregraph.NewSharingLinkTypeFromValue("view")
			Expect(err).ToNot(HaveOccurred())
			Expect(&respLinkType).To(Equal(expected))
		})

		It("handles a failing CreateLink", func() {
			linkType, err := libregraph.NewSharingLinkTypeFromValue("edit")
			Expect(err).ToNot(HaveOccurred())
			driveItemCreateLink.Type = linkType
			permissions, err := linktype.CS3ResourcePermissionsFromSharingLink(*driveItemCreateLink, provider.ResourceType_RESOURCE_TYPE_CONTAINER)
			Expect(err).ToNot(HaveOccurred())
			createLinkResponse.Status = status.NewInternal(ctx, "transport error")
			createLinkResponse.Share = &link.PublicShare{
				Id:          &link.PublicShareId{OpaqueId: "123"},
				Permissions: &link.PublicSharePermissions{Permissions: permissions},
			}

			statResponse.Info = &provider.ResourceInfo{Type: provider.ResourceType_RESOURCE_TYPE_FILE}

			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemCreateLink)).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			var odataError libregraph.OdataError
			err = json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("generalException"))
			Expect(getError.GetMessage()).To(Equal("transport error"))
		})

		It("fails due to an invalid spaceID", func() {
			driveItemCreateLink = libregraph.NewDriveItemCreateLink()

			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/graph/v1beta1/drives/space-id/items/item-id/createShare", nil),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			var odataError libregraph.OdataError
			err := json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("invalidRequest"))
			Expect(getError.GetMessage()).To(Equal("invalid driveID"))
		})

		It("fails due to an empty itemID", func() {
			driveItemCreateLink = libregraph.NewDriveItemCreateLink()

			itemID = ""
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", DriveId)
			rctx.URLParams.Add("itemID", itemID)

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			var odataError libregraph.OdataError
			err := json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("invalidRequest"))
			Expect(getError.GetMessage()).To(Equal("invalid itemID"))
		})

		It("fails due to an itemID on a different storage", func() {
			driveItemCreateLink = libregraph.NewDriveItemCreateLink()

			// use wrong storageID within itemID
			itemID = "f0042750-23c5-441c-9f2c-ff7c53e5bd2b$cd621428-dfbe-44c1-9393-65bf0dd440a6!1177add3-b4eb-434e-a2e8-1859b31b17bf"
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", DriveId)
			rctx.URLParams.Add("itemID", itemID)

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).
					WithContext(ctx),
			)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
			var odataError libregraph.OdataError
			err := json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("itemNotFound"))
			Expect(getError.GetMessage()).To(Equal("driveID and itemID do not match"))
		})

		// Public Shares / "links" in graph terms
		It("fails when creating a public link with an empty request body", func() {
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			var odataError libregraph.OdataError
			err := json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("invalidRequest"))
			Expect(getError.GetMessage()).To(Equal("invalid body schema definition"))

		})

		It("fails when the stat returns access denied", func() {
			err := errors.New("no permission to stat the file")
			statResponse.Status = status.NewPermissionDenied(ctx, err, err.Error())
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemCreateLink)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusForbidden))

			var odataError libregraph.OdataError
			err = json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("accessDenied"))
			Expect(getError.GetMessage()).To(Equal("no permission to stat the file"))
		})

		It("fails when the stat returns resource is locked", func() {
			err := errors.New("the resource is locked")
			statResponse.Status = status.NewLocked(ctx, err.Error())
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemCreateLink)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusLocked))

			var odataError libregraph.OdataError
			err = json.Unmarshal(rr.Body.Bytes(), &odataError)
			Expect(err).ToNot(HaveOccurred())
			getError := odataError.GetError()
			Expect(getError.GetCode()).To(Equal("itemIsLocked"))
			Expect(getError.GetMessage()).To(Equal("the resource is locked"))
		})

		It("succeeds when the link type mapping is not successful", func() {
			// we need to send a valid link type
			linkType, err := libregraph.NewSharingLinkTypeFromValue("edit")
			Expect(err).ToNot(HaveOccurred())
			driveItemCreateLink.Type = linkType
			permissions := &provider.ResourcePermissions{
				CreateContainer:    true,
				InitiateFileUpload: true,
				Move:               true,
			}
			// return different permissions which do not match a link type
			createLinkResponse.Share = &link.PublicShare{
				Id:                &link.PublicShareId{OpaqueId: "123"},
				Expiration:        utils.TimeToTS(*driveItemCreateLink.ExpirationDateTime),
				PasswordProtected: false,
				DisplayName:       ViewerLinkString,
				Token:             "SomeGOODCoffee",
				Permissions:       &link.PublicSharePermissions{Permissions: permissions},
			}
			svc.CreateLink(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemCreateLink)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var createLinkResponseBody *libregraph.Permission
			err = json.Unmarshal(rr.Body.Bytes(), &createLinkResponseBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(createLinkResponseBody.GetId()).To(Equal("123"))
			Expect(createLinkResponseBody.GetExpirationDateTime().Unix()).To(Equal(driveItemCreateLink.ExpirationDateTime.Unix()))
			Expect(createLinkResponseBody.GetHasPassword()).To(Equal(false))
			Expect(createLinkResponseBody.GetLink().LibreGraphDisplayName).To(Equal(libregraph.PtrString(ViewerLinkString)))
			respLink := createLinkResponseBody.GetLink()
			// some conversion gymnastics
			respLinkType := respLink.GetType()
			Expect(err).ToNot(HaveOccurred())
			mockLink := libregraph.SharingLink{}
			lt, _ := linktype.SharingLinkTypeFromCS3Permissions(&link.PublicSharePermissions{Permissions: permissions})
			mockLink.Type = lt
			expectedType := mockLink.GetType()
			Expect(&respLinkType).To(Equal(&expectedType))
			libreGraphActions := createLinkResponseBody.LibreGraphPermissionsActions
			Expect(libreGraphActions[0]).To(Equal("libre.graph/driveItem/children/create"))
			Expect(libreGraphActions[1]).To(Equal("libre.graph/driveItem/upload/create"))
			Expect(libreGraphActions[2]).To(Equal("libre.graph/driveItem/path/update"))
		})
	})
})
