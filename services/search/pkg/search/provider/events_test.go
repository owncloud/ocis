package provider_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/search/mocks"
	provider "github.com/owncloud/ocis/v2/services/search/pkg/search/provider"
)

var _ = Describe("Searchprovider", func() {
	var (
		p           *provider.Provider
		gwClient    *cs3mocks.GatewayAPIClient
		indexClient *mocks.IndexClient

		ctx        context.Context
		eventsChan chan interface{}

		logger = log.NewLogger()
		user   = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}

		ref = &sprovider.Reference{
			ResourceId: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "rootopaqueid",
			},
			Path: "./foo.pdf",
		}
		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			Path: "foo.pdf",
			Size: 12345,
		}
	)

	BeforeEach(func() {
		ctx = context.Background()
		eventsChan = make(chan interface{})
		gwClient = &cs3mocks.GatewayAPIClient{}
		indexClient = &mocks.IndexClient{}

		p = provider.New(gwClient, indexClient, "", eventsChan, logger)

		gwClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
			Status: status.NewOK(ctx),
			Token:  "authtoken",
		}, nil)
		gwClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
			Status: status.NewOK(context.Background()),
			Info:   ri,
		}, nil)
		indexClient.On("DocCount").Return(uint64(1), nil)
	})

	Describe("New", func() {
		It("returns a new instance", func() {
			p = provider.New(gwClient, indexClient, "", eventsChan, logger)
			Expect(p).ToNot(BeNil())
		})
	})

	Describe("events", func() {
		It("triggers an index update when a file has been uploaded", func() {
			called := false
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.FileUploaded{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("triggers an index update when a file has been touched", func() {
			called := false
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.FileTouched{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("removes an entry from the index when the file has been deleted", func() {
			called := false
			gwClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
				Status: status.NewNotFound(context.Background(), ""),
			}, nil)
			indexClient.On("Delete", mock.MatchedBy(func(id *sprovider.ResourceId) bool {
				return id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.ItemTrashed{
				Ref:       ref,
				ID:        ri.Id,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when they are being restored", func() {
			called := false
			indexClient.On("Restore", mock.MatchedBy(func(id *sprovider.ResourceId) bool {
				return id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.ItemRestored{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when a version has been restored", func() {
			called := false
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.FileVersionRestored{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when they are being moved", func() {
			called := false
			gwClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
				Status: status.NewOK(ctx),
				Path:   "./new/path.pdf",
			}, nil)
			indexClient.On("Move", mock.MatchedBy(func(id *sprovider.ResourceId) bool {
				return id.OpaqueId == ri.Id.OpaqueId
			}), "./new/path.pdf").Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			ref.Path = "./new/path.pdf"
			eventsChan <- events.ItemMoved{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})
	})
})
