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
						Id:        &provider.StorageSpaceId{OpaqueId: "sameID"},
						SpaceType: "aspacetype",
						Root: &provider.ResourceId{
							StorageId: "sameID",
							OpaqueId:  "sameID",
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
						"id":"sameID",
						"name":"aspacename",
						"root":{
							"id":"sameID!sameID",
							"webDavUrl":"https://localhost:9200/dav/spaces/sameID"
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
						Id:        &provider.StorageSpaceId{OpaqueId: "bsameID"},
						SpaceType: "bspacetype",
						Root: &provider.ResourceId{
							StorageId: "bsameID",
							OpaqueId:  "bsameID",
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
						Id:        &provider.StorageSpaceId{OpaqueId: "asameID"},
						SpaceType: "aspacetype",
						Root: &provider.ResourceId{
							StorageId: "asameID",
							OpaqueId:  "asameID",
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
						"id":"asameID",
						"name":"aspacename",
						"root":{
							"eTag":"101112131415",
							"id":"asameID!asameID",
							"webDavUrl":"https://localhost:9200/dav/spaces/asameID"
						}
					},
					{
						"driveAlias":"bspacetype/bspacename",
						"driveType":"bspacetype",
						"id":"bsameID",
						"name":"bspacename",
						"root":{
							"eTag":"123456789",
							"id":"bsameID!bsameID",
							"webDavUrl":"https://localhost:9200/dav/spaces/bsameID"
						}
					}
				]
			}
			`))
		})
		It("can list a spaces type mountpoint", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "aID!differentID"},
						SpaceType: "mountpoint",
						Root: &provider.ResourceId{
							StorageId: "aID",
							OpaqueId:  "differentID",
						},
						Name: "New Folder",
						Opaque: &typesv1beta1.Opaque{
							Map: map[string]*typesv1beta1.OpaqueEntry{
								"spaceAlias":     {Decoder: "plain", Value: []byte("mountpoint/new-folder")},
								"etag":           {Decoder: "plain", Value: []byte("101112131415")},
								"grantStorageID": {Decoder: "plain", Value: []byte("ownerStorageID")},
								"grantOpaqueID":  {Decoder: "plain", Value: []byte("opaqueID")},
							},
						},
					},
				},
			}, nil)
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
				Status: status.NewOK(ctx),
				Info: &provider.ResourceInfo{
					Etag:  "123456789",
					Type:  provider.ResourceType_RESOURCE_TYPE_CONTAINER,
					Id:    &provider.ResourceId{StorageId: "ownerStorageID", OpaqueId: "opaqueID"},
					Path:  "New Folder",
					Mtime: &typesv1beta1.Timestamp{Seconds: 1648327606},
					Size:  uint64(1234),
				},
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
						"driveAlias":"mountpoint/new-folder",
						"driveType":"mountpoint",
						"id":"aID!differentID",
						"name":"New Folder",
						"root":{
							"eTag":"101112131415",
							"id":"aID!differentID",
							"remoteItem":{
								"eTag":"123456789",
								"folder":{},
								"id":"ownerStorageID!opaqueID",
								"lastModifiedDateTime":"2022-03-26T21:46:46+01:00",
								"name":"New Folder",
								"size":1234,
								"webDavUrl":"https://localhost:9200/dav/spaces/ownerStorageID!opaqueID"
							},
							"webDavUrl":"https://localhost:9200/dav/spaces/aID!differentID"
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
