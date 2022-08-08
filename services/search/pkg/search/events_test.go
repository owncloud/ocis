package search_test

import (
	"context"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
	engineMocks "github.com/owncloud/ocis/v2/services/search/pkg/engine/mocks"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	"github.com/stretchr/testify/mock"
	mEvents "go-micro.dev/v4/events"
)

var _ = Describe("Events", func() {
	var (
		gw        *cs3mocks.GatewayAPIClient
		engine    *engineMocks.Engine
		extractor *contentMocks.Extractor
		bus       events.Stream
		ctx       context.Context

		logger = log.NewLogger()
		user   = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}

		ref = &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "rootopaqueid",
			},
			Path: "./foo.pdf",
		}
		ri = &provider.ResourceInfo{
			Id: &provider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			Path: "foo.pdf",
			Size: 12345,
		}
	)

	BeforeEach(func() {
		ctx = context.Background()
		gw = &cs3mocks.GatewayAPIClient{}
		engine = &engineMocks.Engine{}
		bus, _ = mEvents.NewStream()
		extractor = &contentMocks.Extractor{}

		_ = search.HandleEvents(engine, extractor, gw, bus, logger, &config.Config{})

		gw.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
			Status: status.NewOK(ctx),
			Token:  "authtoken",
		}, nil)
		gw.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
			Status: status.NewOK(ctx),
			Info:   ri,
		}, nil)
		gw.On("GetPath", mock.Anything, mock.Anything).Return(&provider.GetPathResponse{
			Status: status.NewOK(ctx),
			Path:   "",
		}, nil)
		engine.On("DocCount").Return(uint64(1), nil)
	})
	Describe("events", func() {
		It("triggers an index update when a file has been uploaded", func() {
			called := false
			extractor.Mock.On("Extract", mock.Anything, mock.Anything).Return(content.Document{}, nil)
			engine.On("Upsert", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})

			_ = events.Publish(bus, events.FileUploaded{
				Ref:       ref,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("triggers an index update when a file has been touched", func() {
			called := false
			extractor.Mock.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			engine.On("Upsert", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			_ = events.Publish(bus, events.FileTouched{
				Ref:       ref,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("removes an entry from the index when the file has been deleted", func() {
			called := false
			gw.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
				Status: status.NewNotFound(context.Background(), ""),
			}, nil)
			engine.On("Delete", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})

			_ = events.Publish(bus, events.ItemTrashed{
				Ref:       ref,
				ID:        ri.Id,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when they are being restored", func() {
			called := false
			engine.On("Restore", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})

			_ = events.Publish(bus, events.ItemRestored{
				Ref:       ref,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when a version has been restored", func() {
			called := false
			extractor.Mock.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			engine.On("Upsert", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})

			_ = events.Publish(bus, events.FileVersionRestored{
				Ref:       ref,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})

		It("indexes items when they are being moved", func() {
			called := false
			engine.On("Move", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})

			_ = events.Publish(bus, events.ItemMoved{
				Ref:       ref,
				Executant: user.Id,
			})

			Eventually(func() bool {
				return called
			}, "2s").Should(BeTrue())
		})
	})
})
