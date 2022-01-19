package svc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/graph/mocks"
	"github.com/owncloud/ocis/graph/pkg/config"
	service "github.com/owncloud/ocis/graph/pkg/service/v0"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var spaceYaml = strings.NewReader(`---
version: "1.0"
description: read from yaml
special:
  readme: readme2.md
  image: .img/space.png
`)

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
						"id":"aspaceid!anopaqueid",
						"name":"aspacename",
						"root":{
							"id":"aspaceid!anopaqueid",
							"webDavUrl":"https://localhost:9200/dav/spaces/aspaceid!anopaqueid"
						}
					}
				]
			}		
			`))
		})

		It("can list a space with extended properties from a space.yaml", func() {
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
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileDownloadProtocol{
					{
						Protocol:         "spaces",
						DownloadEndpoint: "ignored",
					},
				},
			}, nil)
			// mock space.yaml
			httpClient.On("Do", mock.Anything, mock.Anything).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(spaceYaml),
			}, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status: status.NewUnimplemented(ctx, fmt.Errorf("not supported"), "not supported"),
			}, nil)
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(
				func(_ context.Context, req *provider.StatRequest, _ ...grpc.CallOption) *provider.StatResponse {
					switch req.Ref.GetPath() {
					case "./readme2.md":
						return &provider.StatResponse{
							Status: status.NewOK(ctx),
							Info: &provider.ResourceInfo{
								Type: provider.ResourceType_RESOURCE_TYPE_FILE,
								Path: "readme2.md",
								Id: &provider.ResourceId{
									StorageId: "aspaceid",
									OpaqueId:  "readmeid",
								},
								PermissionSet: &provider.ResourcePermissions{
									Stat: true,
								},
								Size: 10,
							},
						}
					case "./.img/space.png":
						return &provider.StatResponse{
							Status: status.NewOK(ctx),
							Info: &provider.ResourceInfo{
								Type: provider.ResourceType_RESOURCE_TYPE_FILE,
								Path: "space.png",
								Id: &provider.ResourceId{
									StorageId: "aspaceid",
									OpaqueId:  "imageid",
								},
								PermissionSet: &provider.ResourcePermissions{
									Stat: true,
								},
								Size: 20,
							},
						}
					default:
						return &provider.StatResponse{
							Status: status.NewNotFound(ctx, "not found"),
						}
					}
				},
				nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			type expectedJSON struct {
				Value []struct {
					DriveType   string `json:"driveType"`
					ID          string
					Name        string
					Description string
					Root        struct {
						ID        string
						WebDavURL string
					}
					Special []struct {
						ID            string
						Name          string
						Size          int
						SpecialFolder struct {
							Name string
						}
						WebDavURL string
					}
				}
			}
			body, _ := io.ReadAll(rr.Body)
			e := expectedJSON{}
			err := json.Unmarshal(body, &e)
			Expect(err).To(BeNil())

			Expect(len(e.Value)).To(Equal(1))

			v := e.Value[0]
			eq := func(actual string, expected string) {
				Expect(actual).To(Equal(expected))
			}
			eq(v.DriveType, "aspacetype")
			eq(v.ID, "aspaceid!anopaqueid")
			eq(v.Name, "aspacename")
			eq(v.Description, "read from yaml")
			eq(v.Root.ID, "aspaceid!anopaqueid")
			eq(v.Root.WebDavURL, "https://localhost:9200/dav/spaces/aspaceid!anopaqueid")

			Expect(len(v.Special)).To(Equal(2))
			for _, s := range v.Special {
				switch s.ID {
				default:
					Fail("unexpected id" + s.ID)
				case "readmeid":
					eq(s.Name, "readme2.md")
					Expect(s.Size).To(Equal(10))
					eq(s.SpecialFolder.Name, "readme")
					eq(s.WebDavURL, "https://localhost:9200/dav/spaces/aspaceid/readme2.md")
				case "imageid":
					eq(s.Name, "space.png")
					Expect(s.Size).To(Equal(20))
					eq(s.SpecialFolder.Name, "image")
					eq(s.WebDavURL, "https://localhost:9200/dav/spaces/aspaceid/.img/space.png")

				}
			}
		})
	})
})
