package svc_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/graph/mocks"
	"github.com/owncloud/ocis/graph/pkg/config"
	service "github.com/owncloud/ocis/graph/pkg/service/v0"
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
			service.Config(config.DefaultConfig()),
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
					},
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

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives?$orderby=name%20asc", nil)
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
					},
					{
						"driveType":"bspacetype",
						"id":"bspaceid",
						"name":"bspacename",
						"root":{
							"id":"bspaceid!bopaqueid",
							"webDavUrl":"https://localhost:9200/dav/spaces/bspaceid"
						}
					}
				]
			}
			`))
		})
	})
})
