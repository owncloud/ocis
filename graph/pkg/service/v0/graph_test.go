package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"

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
		svc    service.Service
		client *mocks.GatewayClient
		ctx    context.Context
	)

	JustBeforeEach(func() {
		ctx = context.Background()
		client = &mocks.GatewayClient{}
		svc = service.NewService(
			service.Config(config.DefaultConfig()),
			service.GatewayServiceClient(client),
		)
	})

	Describe("NewService", func() {
		It("returns a service", func() {
			Expect(svc).ToNot(BeNil())
		})
	})
	Describe("drive", func() {
		It("can list an empty list of spaces", func() {
			client.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{},
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives", nil)
			rr := httptest.NewRecorder()
			svc.GetDrives(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})
})
