package connector_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/middleware"

	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
)

var _ = Describe("ContentConnector", func() {
	var (
		cc            *connector.ContentConnector
		gatewayClient *cs3mocks.GatewayAPIClient
		cfg           *config.Config
		wopiCtx       middleware.WopiContext

		srv           *httptest.Server
		srvReqHeader  http.Header
		randomContent string
	)

	BeforeEach(func() {
		// contentConnector only uses "cfg.CS3Api.DataGateway.Insecure", which is irrelevant for the tests
		cfg = &config.Config{}
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		cc = connector.NewContentConnector(gatewayClient, cfg)

		wopiCtx = middleware.WopiContext{
			AccessToken: "abcdef123456",
			FileReference: providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "abc",
					OpaqueId:  "12345",
					SpaceId:   "zzz",
				},
				Path: ".",
			},
			User:       &userv1beta1.User{}, // Not used for now
			ViewMode:   appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
			EditAppUrl: "http://test.ex.prv/edit",
			ViewAppUrl: "http://test.ex.prv/view",
		}

		randomContent = "This is the content of the test.txt file"
		srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			srvReqHeader = req.Header // save the request header to check later
			switch req.URL.Path {
			case "/download/failed.png":
				w.WriteHeader(404)
			case "/download/test.txt":
				w.Write([]byte(randomContent))
			case "/upload/failed.png":
				w.WriteHeader(404)
			case "/upload/test.txt":
				w.WriteHeader(200)
			}
		}))
	})

	AfterEach(func() {
		srv.Close()
	})

	Describe("GetFile", func() {
		It("No valid context", func() {
			sb := &strings.Builder{}
			ctx := context.Background()
			err := cc.GetFile(ctx, sb)
			Expect(err).To(HaveOccurred())
		})

		It("Initiate download failed", func() {
			sb := &strings.Builder{}
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, targetErr)

			err := cc.GetFile(ctx, sb)
			Expect(err).To(Equal(targetErr))
		})

		It("Initiate download status not ok", func() {
			sb := &strings.Builder{}
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, nil)

			err := cc.GetFile(ctx, sb)
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
		})

		It("Missing download endpoint", func() {
			sb := &strings.Builder{}
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewOK(ctx),
			}, nil)

			err := cc.GetFile(ctx, sb)
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
		})

		It("Download request failed", func() {
			sb := &strings.Builder{}
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileDownloadProtocol{
					{
						Protocol:         "simple",
						DownloadEndpoint: srv.URL + "/download/failed.png",
						Token:            "MyDownloadToken",
					},
				},
			}, nil)

			err := cc.GetFile(ctx, sb)
			Expect(srvReqHeader.Get("X-Access-Token")).To(Equal(wopiCtx.AccessToken))
			Expect(srvReqHeader.Get("X-Reva-Transfer")).To(Equal("MyDownloadToken"))
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
		})

		It("Download request success", func() {
			sb := &strings.Builder{}
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("InitiateFileDownload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileDownloadProtocol{
					{
						Protocol:         "simple",
						DownloadEndpoint: srv.URL + "/download/test.txt",
						Token:            "MyDownloadToken",
					},
				},
			}, nil)

			err := cc.GetFile(ctx, sb)
			Expect(srvReqHeader.Get("X-Access-Token")).To(Equal(wopiCtx.AccessToken))
			Expect(srvReqHeader.Get("X-Reva-Transfer")).To(Equal("MyDownloadToken"))
			Expect(err).To(Succeed())
			Expect(sb.String()).To(Equal(randomContent))
		})

		It("ViewOnlyMode Download request success", func() {
			sb := &strings.Builder{}

			wopiCtx = middleware.WopiContext{
				AccessToken:   "abcdef123456",
				ViewOnlyToken: "view.only.123456",
				FileReference: providerv1beta1.Reference{
					ResourceId: &providerv1beta1.ResourceId{
						StorageId: "abc",
						OpaqueId:  "12345",
						SpaceId:   "zzz",
					},
					Path: ".",
				},
				User:       &userv1beta1.User{}, // Not used for now
				ViewMode:   appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY,
				EditAppUrl: "http://test.ex.prv/edit",
				ViewAppUrl: "http://test.ex.prv/view",
			}

			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("InitiateFileDownload",
				mock.MatchedBy(func(ctx context.Context) bool {
					return revactx.ContextMustGetToken(ctx) == "view.only.123456"
				}), mock.Anything).Times(1).Return(&gateway.InitiateFileDownloadResponse{
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileDownloadProtocol{
					{
						Protocol:         "simple",
						DownloadEndpoint: srv.URL + "/download/test.txt",
						Token:            "MyDownloadToken",
					},
				},
			}, nil)

			err := cc.GetFile(ctx, sb)
			Expect(srvReqHeader.Get("X-Access-Token")).To(Equal(wopiCtx.ViewOnlyToken))
			Expect(srvReqHeader.Get("X-Reva-Transfer")).To(Equal("MyDownloadToken"))
			Expect(err).To(Succeed())
			Expect(sb.String()).To(Equal(randomContent))
		})
	})

	Describe("PutFile", func() {
		It("No valid context", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := context.Background()
			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "notARandomLockId")
			Expect(err).To(HaveOccurred())
			Expect(newLockId).To(Equal(""))
		})

		It("Stat call failed", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, targetErr)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "notARandomLockId")
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Stat call status not ok", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "notARandomLockId")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("Mismatched lockId", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
				},
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "notARandomLockId")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal("goodAndValidLock"))
		})

		It("Upload without lockId but on a non empty file", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: nil,
					Size: uint64(123456789),
				},
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal(""))
		})

		It("Initiate upload fails", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, targetErr)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(err).To(HaveOccurred())
			Expect(newLockId).To(Equal(""))
		})

		It("Initiate upload status not ok", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("Empty upload successful", func() {
			reader := strings.NewReader("")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(err).To(Succeed())
			Expect(newLockId).To(Equal(""))
		})

		It("Missing upload endpoint", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("upload request failed", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileUploadProtocol{
					{
						Protocol:       "simple",
						UploadEndpoint: srv.URL + "/upload/failed.png",
					},
				},
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(srvReqHeader.Get("X-Access-Token")).To(Equal(wopiCtx.AccessToken))
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("upload request success", func() {
			reader := strings.NewReader("Content to upload is here!")
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Lock: &providerv1beta1.Lock{
						LockId: "goodAndValidLock",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
					Size: uint64(123456789),
				},
			}, nil)

			gatewayClient.On("InitiateFileUpload", mock.Anything, mock.Anything).Times(1).Return(&gateway.InitiateFileUploadResponse{
				Status: status.NewOK(ctx),
				Protocols: []*gateway.FileUploadProtocol{
					{
						Protocol:       "simple",
						UploadEndpoint: srv.URL + "/upload/test.txt",
					},
				},
			}, nil)

			newLockId, err := cc.PutFile(ctx, reader, reader.Size(), "goodAndValidLock")
			Expect(srvReqHeader.Get("X-Access-Token")).To(Equal(wopiCtx.AccessToken))
			Expect(err).To(Succeed())
			Expect(newLockId).To(Equal(""))
		})
	})
})
