package middleware_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/gomega"
	pMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
	policiesPG "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc"
)

func TestPolicies_NoQuery_PassThrough(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, _, _ := prepare("")

	responseRecorder := httptest.NewRecorder()
	policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/policies", nil))

	g.Expect(responseRecorder.Code).To(Equal(http.StatusOK))
}

func TestPolicies_ErrorsOnEvaluationError(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, policiesProviderService, _ := prepare("any")
	policiesProviderService.On("Evaluate", mock.Anything, mock.Anything).Return(
		nil,
		errors.New("any"),
	)

	responseRecorder := httptest.NewRecorder()
	policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/policies", nil))

	g.Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
}

func TestPolicies_ErrorsOnDeny(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, policiesProviderService, _ := prepare("any")
	policiesProviderService.On("Evaluate", mock.Anything, mock.Anything).Return(
		&policiesPG.EvaluateResponse{},
		nil,
	)

	responseRecorder := httptest.NewRecorder()
	policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/policies", nil))

	result := responseRecorder.Result()
	defer func() {
		g.Expect(result.Body.Close()).ToNot(HaveOccurred())
	}()

	data, err := io.ReadAll(result.Body)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(data).To(ContainSubstring(middleware.DeniedMessage))
	g.Expect(responseRecorder.Code).To(Equal(http.StatusForbidden))
}

func TestPolicies_EvaluationEnvironment_HTTPStage(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, policiesProviderService, _ := prepare("any")
	policiesProviderService.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(
		func(_ context.Context, in *policiesPG.EvaluateRequest, _ ...client.CallOption) (*policiesPG.EvaluateResponse, error) {
			g.Expect(in.Environment.Stage).To(Equal(pMessage.Stage_STAGE_HTTP))

			return &policiesPG.EvaluateResponse{Result: false}, nil
		},
	)

	responseRecorder := httptest.NewRecorder()
	policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/policies", nil))
}

func TestPolicies_EvaluationEnvironment_Request(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, policiesProviderService, _ := prepare("any")
	policiesProviderService.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(
		func(_ context.Context, in *policiesPG.EvaluateRequest, _ ...client.CallOption) (*policiesPG.EvaluateResponse, error) {
			g.Expect(in.Environment.Request.Method).To(Equal(http.MethodDelete))
			g.Expect(in.Environment.Request.Path).To(Equal("/whatever"))

			return &policiesPG.EvaluateResponse{Result: false}, nil
		},
	)

	responseRecorder := httptest.NewRecorder()
	policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodDelete, "/whatever", nil))
}

func TestPolicies_EvaluationEnvironment_Resource(t *testing.T) {
	var g = NewWithT(t)

	policiesMiddleware, policiesProviderService, _ := prepare("any")

	// tus metadata
	{
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/remote.php/dav/spaces", nil)
		request.Header.Set(net.HeaderUploadMetadata, fmt.Sprintf("filename %v", base64.StdEncoding.EncodeToString([]byte("tus-file-name.png"))))
		policiesProviderService.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(
			func(_ context.Context, in *policiesPG.EvaluateRequest, _ ...client.CallOption) (*policiesPG.EvaluateResponse, error) {
				g.Expect(in.Environment.Resource.Name).To(Equal("tus-file-name.png"))

				return &policiesPG.EvaluateResponse{Result: false}, nil
			},
		).Once()
		policiesMiddleware.ServeHTTP(responseRecorder, request)
	}

	// url path
	{
		responseRecorder := httptest.NewRecorder()
		policiesProviderService.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(
			func(_ context.Context, in *policiesPG.EvaluateRequest, _ ...client.CallOption) (*policiesPG.EvaluateResponse, error) {
				g.Expect(in.Environment.Resource.Name).To(Equal("simple-file-name.png"))

				return &policiesPG.EvaluateResponse{Result: false}, nil
			},
		).Once()
		policiesMiddleware.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodPut, "/remote.php/dav/spaces/simple-file-name.png", nil))
	}
}

func prepare(q string) (http.Handler, *mocks.PoliciesProviderService, *cs3mocks.GatewayAPIClient) {

	// mocked gatewaySelector
	gatewayClient := &cs3mocks.GatewayAPIClient{}
	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)
	defer pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")

	// mocked policiesProviderService
	policiesProviderService := &mocks.PoliciesProviderService{}

	// spin up middleware
	policiesMiddleware := middleware.Policies(
		q,
		middleware.WithRevaGatewaySelector(gatewaySelector),
		middleware.PoliciesProviderService(policiesProviderService),
	)(mockHandler{})

	return policiesMiddleware, policiesProviderService, gatewayClient
}

type mockHandler struct{}

func (m mockHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {}
