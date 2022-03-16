package svc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/graph/mocks"
	"github.com/owncloud/ocis/graph/pkg/config/defaults"
	service "github.com/owncloud/ocis/graph/pkg/service/v0"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Graph", func() {
	var (
		svc           service.Service
		gatewayClient *mocks.GatewayClient
		httpClient    *mocks.HTTPClient
		ctx           context.Context
	)

	JustBeforeEach(func() {
		ctx = context.Background()
		gatewayClient = &mocks.GatewayClient{}
		httpClient = &mocks.HTTPClient{}
		svc = service.NewService(
			service.Config(defaults.DefaultConfig()),
			service.WithGatewayClient(gatewayClient),
			service.WithHTTPClient(httpClient),
		)
	})

	Describe("NewService", func() {
		It("returns a service", func() {
			Expect(svc).ToNot(BeNil())
		})
	})

	Describe("drive", func() {
		It("can list an empty list of spaces", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{},
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("can list a space without owner", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "aspaceid"},
						SpaceType: "aspacetype",
						Root: &provider.ResourceId{
							StorageId: "aspaceid",
							OpaqueId:  "anopaqueid",
						},
						Name: "aspacename",
					},
				},
			}, nil)
			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status: status.NewUnimplemented(ctx, fmt.Errorf("not supported"), "not supported"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			body, _ := io.ReadAll(rr.Body)
			Expect(body).To(MatchJSON(`
			{
				"value":[
					{
						"driveType":"aspacetype",
						"id":"aspaceid",
						"name":"aspacename",
						"root":{
							"id":"aspaceid!anopaqueid",
							"webDavUrl":"https://localhost:9200/dav/spaces/aspaceid"
						}
					}
				]
			}
			`))
		})
		It("can list a spaces with sort", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "bspaceid"},
						SpaceType: "bspacetype",
						Root: &provider.ResourceId{
							StorageId: "bspaceid",
							OpaqueId:  "bopaqueid",
						},
						Name: "bspacename",
						Opaque: &typesv1beta1.Opaque{
							Map: map[string]*typesv1beta1.OpaqueEntry{
								"spaceAlias": {Decoder: "plain", Value: []byte("bspacetype/bspacename")},
								"etag":       {Decoder: "plain", Value: []byte("123456789")},
							},
						},
					},
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "aspaceid"},
						SpaceType: "aspacetype",
						Root: &provider.ResourceId{
							StorageId: "aspaceid",
							OpaqueId:  "anopaqueid",
						},
						Name: "aspacename",
						Opaque: &typesv1beta1.Opaque{
							Map: map[string]*typesv1beta1.OpaqueEntry{
								"spaceAlias": {Decoder: "plain", Value: []byte("aspacetype/aspacename")},
								"etag":       {Decoder: "plain", Value: []byte("101112131415")},
							},
						},
					},
				},
			}, nil)
			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status: status.NewUnimplemented(ctx, fmt.Errorf("not supported"), "not supported"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives?$orderby=name%20asc", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			body, _ := io.ReadAll(rr.Body)
			Expect(body).To(MatchJSON(`
			{
				"value":[
					{
						"driveAlias":"aspacetype/aspacename",
						"driveType":"aspacetype",
						"id":"aspaceid",
						"name":"aspacename",
						"root":{
							"eTag":"101112131415",
							"id":"aspaceid!anopaqueid",
							"webDavUrl":"https://localhost:9200/dav/spaces/aspaceid"
						}
					},
					{
						"driveAlias":"bspacetype/bspacename",
						"driveType":"bspacetype",
						"id":"bspaceid",
						"name":"bspacename",
						"root":{
							"eTag":"123456789",
							"id":"bspaceid!bopaqueid",
							"webDavUrl":"https://localhost:9200/dav/spaces/bspaceid"
						}
					}
				]
			}
			`))
		})
		It("can not list spaces with wrong sort parameter", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{}}, nil)
			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status: status.NewUnimplemented(ctx, fmt.Errorf("not supported"), "not supported"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives?$orderby=owner%20asc", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			body, _ := io.ReadAll(rr.Body)
			var libreError libregraph.OdataError
			err := json.Unmarshal(body, &libreError)
			Expect(err).To(Not(HaveOccurred()))
			Expect(libreError.Error.Message).To(Equal("we do not support <owner> as a order parameter"))
			Expect(libreError.Error.Code).To(Equal(errorcode.InvalidRequest.String()))
		})
		It("can list a spaces with invalid query parameter", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{}}, nil)
			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status: status.NewUnimplemented(ctx, fmt.Errorf("not supported"), "not supported"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives?§orderby=owner%20asc", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			body, _ := io.ReadAll(rr.Body)
			var libreError libregraph.OdataError
			err := json.Unmarshal(body, &libreError)
			Expect(err).To(Not(HaveOccurred()))
			Expect(libreError.Error.Message).To(Equal("Query parameter '§orderby' is not supported. Cause: Query parameter '§orderby' is not supported"))
			Expect(libreError.Error.Code).To(Equal(errorcode.InvalidRequest.String()))
		})
	})
})
