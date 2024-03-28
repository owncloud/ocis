package svc_test

import (
	"context"
	"errors"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	linkv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"

	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	svc "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("ExtensionsOrgLibregraphInfoService", func() {
	var (
		service         svc.ExtensionsOrgLibregraphInfoService
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector *mocks.Selectable[gateway.GatewayAPIClient]
	)
	BeforeEach(func() {
		logger := log.NewLogger()
		gatewayClient = cs3mocks.NewGatewayAPIClient(GinkgoT())

		gatewaySelector = mocks.NewSelectable[gateway.GatewayAPIClient](GinkgoT())
		gatewaySelector.On("Next").Return(gatewayClient, nil)

		extensionsOrgLibregraphInfoService, err := svc.NewExtensionsOrgLibregraphInfoService(
			svc.ExtensionsOrgLibregraphInfoServiceOptions{}.
				WithLogger(logger).
				WithGatewaySelector(gatewaySelector),
		)
		Expect(err).ToNot(HaveOccurred())
		service = extensionsOrgLibregraphInfoService
	})

	Describe("TokenInfo", func() {
		It("exits if the next gatewayClient cannot be obtained", func() {
			gatewaySelector.ExpectedCalls = nil

			expectedError := errors.New("obtaining next gatewayClient failed")
			gatewaySelector.On("Next").Return(gatewayClient, expectedError)

			_, err := service.TokenInfo(context.Background(), "", "")
			Expect(err).To(MatchError(expectedError))
		})
		It("authorizes the request with correct parameters", func() {
			gatewayClient.
				On("Authenticate", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, r *gateway.AuthenticateRequest, _ ...grpc.CallOption) (*gateway.AuthenticateResponse, error) {
					Expect(r.Type).To(Equal("publicshares"))
					Expect(r.ClientId).To(Equal("123"))
					Expect(r.ClientSecret).To(Equal("password|456"))
					return nil, nil
				})

			_, _ = service.TokenInfo(context.Background(), "123", "456")
		})
		It("exits the Authorization fails", func() {
			expectedError := errorcode.New(errorcode.GeneralException, "authorization failed")
			gatewayClient.
				On("Authenticate", mock.Anything, mock.Anything).
				Return(&gateway.AuthenticateResponse{}, errors.New("authorization failed"))

			r, err := service.TokenInfo(context.Background(), "", "")
			Expect(err).To(MatchError(&expectedError))
			Expect(r.ID).To(Equal(""))
			Expect(r.IsInternal).To(BeFalse())
			Expect(r.HasPassword).To(BeFalse())
		})
		It("reports an info response with required password if the permission is denied", func() {
			gatewayClient.
				On("Authenticate", mock.Anything, mock.Anything).
				Return(&gateway.AuthenticateResponse{
					Status: &rpcv1beta1.Status{
						Code: rpcv1beta1.Code_CODE_PERMISSION_DENIED,
					},
				}, nil)

			r, err := service.TokenInfo(context.Background(), "", "")
			Expect(err).To(BeNil())
			Expect(r.ID).To(Equal(""))
			Expect(r.IsInternal).To(BeFalse())
			Expect(r.HasPassword).To(BeTrue())
		})
		It("uses proper authentication for the share lookup", func() {
			gatewayClient.
				On("Authenticate", mock.Anything, mock.Anything).
				Return(&gateway.AuthenticateResponse{
					Status: &rpcv1beta1.Status{
						Code: rpcv1beta1.Code_CODE_OK,
					},
					Token: "token",
					User: &userv1beta1.User{
						Username: "username",
					},
				}, nil)

			gatewayClient.
				On("GetPublicShare", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, r *linkv1beta1.GetPublicShareRequest, opts ...grpc.CallOption) (*linkv1beta1.GetPublicShareResponse, error) {
					u, _ := ctxpkg.ContextGetUser(ctx)
					Expect(u.GetUsername()).To(Equal("username"))
					t, _ := ctxpkg.ContextGetToken(ctx)
					Expect(t).To(Equal("token"))
					Expect(r.GetRef().GetToken()).To(Equal("token"))
					return nil, nil
				})

			_, _ = service.TokenInfo(context.Background(), "", "")
		})
		It("exits if the share lookup fails", func() {})
		It("returns the information for non internal shares", func() {})
		It("exits if the stat for internal shares fails", func() {})
		It("removes the info id field if the requested share is not found", func() {})
		It("retruns all the information for internal shares", func() {})
	})
})
