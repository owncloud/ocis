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
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/oklog/run"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	pMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
	policiesPG "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/middleware"
)

func TestPolicies_ErrorsOnEvaluationError(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)

	pmt.eval = func(request *policiesPG.EvaluateRequest, response *policiesPG.EvaluateResponse) error {
		return errors.New("any error")
	}

	pmt.run(func(g *WithT, w *httptest.ResponseRecorder, _ *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(w.Code).To(Equal(http.StatusInternalServerError))
	})
}

func TestPolicies_ErrorsOnDeny(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)

	pmt.run(func(g *WithT, w *httptest.ResponseRecorder, _ *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		res := w.Result()
		defer func() {
			g.Expect(res.Body.Close()).ToNot(HaveOccurred())
		}()

		data, err := io.ReadAll(res.Body)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(data).To(ContainSubstring(middleware.DeniedMessage))
		g.Expect(w.Code).To(Equal(http.StatusForbidden))
	})
}

func TestPolicies_EvaluationEnvironment_HTTPStage(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)

	pmt.run(func(g *WithT, _ *httptest.ResponseRecorder, r *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(r.Environment.Stage).To(Equal(pMessage.Stage_STAGE_HTTP))
	})
}

func TestPolicies_EvaluationEnvironment_Request(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)
	pmt.httpRequest = httptest.NewRequest(http.MethodDelete, "/whatever", nil)

	pmt.run(func(g *WithT, _ *httptest.ResponseRecorder, r *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(r.Environment.Request.Method).To(Equal(http.MethodDelete))
		g.Expect(r.Environment.Request.Path).To(Equal("/whatever"))
	})
}

func TestPolicies_EvaluationEnvironment_Resource(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)

	// tus metadata
	pmt.httpRequest = httptest.NewRequest(http.MethodPost, "/remote.php/dav/spaces", nil)
	pmt.httpRequest.Header.Set(net.HeaderUploadMetadata, fmt.Sprintf("filename %v", base64.StdEncoding.EncodeToString([]byte("tus-file-name.png"))))

	pmt.run(func(g *WithT, _ *httptest.ResponseRecorder, r *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(r.Environment.Resource.Name).To(Equal("tus-file-name.png"))
	})

	// url path
	pmt.httpRequest = httptest.NewRequest(http.MethodPut, "/remote.php/dav/spaces/simple-file-name.png", nil)
	pmt.run(func(g *WithT, _ *httptest.ResponseRecorder, r *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(r.Environment.Resource.Name).To(Equal("simple-file-name.png"))
	})

	// shared-resource put
	pmt.httpRequest = httptest.NewRequest(http.MethodPut, "/remote.php/dav/spaces/897987fd978dffdfds9f78dsf97fd", nil)
	pmt.gwClientStatResponse.Info.Name = "shared-file-name.png"
	pmt.run(func(g *WithT, _ *httptest.ResponseRecorder, r *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(r.Environment.Resource.Name).To(Equal("shared-file-name.png"))
	})
}

func TestPolicies_NoQuery_PassThrough(t *testing.T) {
	pmt := newPoliciesMiddlewareTester(t)
	pmt.regoQuery = ""

	pmt.run(func(g *WithT, w *httptest.ResponseRecorder, _ *policiesPG.EvaluateRequest, _ *policiesPG.EvaluateResponse) {
		g.Expect(w.Code).To(Equal(http.StatusOK))
	})
}

func newPoliciesMiddlewareTester(t *testing.T) policiesMiddlewareTester {
	return policiesMiddlewareTester{
		g:           NewWithT(t),
		regoQuery:   "any",
		httpRequest: httptest.NewRequest(http.MethodGet, "/policies", nil),
		eval: func(request *policiesPG.EvaluateRequest, response *policiesPG.EvaluateResponse) error {
			return nil
		},
		gwClientStatResponse: &sprovider.StatResponse{
			Status: status.NewOK(context.Background()),
			Info:   &sprovider.ResourceInfo{},
		},
	}
}

type policiesMiddlewareTester struct {
	g                    *WithT
	regoQuery            string
	httpRequest          *http.Request
	grpcEvaluateRequest  *policiesPG.EvaluateRequest
	grpcEvaluateResponse *policiesPG.EvaluateResponse
	gwClientStatResponse *sprovider.StatResponse
	eval                 func(request *policiesPG.EvaluateRequest, response *policiesPG.EvaluateResponse) error
}

func (pmt *policiesMiddlewareTester) run(e func(g *WithT, w *httptest.ResponseRecorder, grpcRequest *policiesPG.EvaluateRequest, grpcResponse *policiesPG.EvaluateResponse)) {
	var (
		polServiceAddress   = fmt.Sprintf("127.0.0.1:%d", freeport.GetPort()) //non-ephemeral port didn't work
		polServiceName      = "policies"
		polServiceNamespace = "com.owncloud.api"
		ctx                 = context.Background()
		rg                  = run.Group{}
	)

	registry.Configure("memory")
	reg := registry.GetRegistry()

	polService := registry.BuildGRPCService(polServiceNamespace+"."+polServiceName, "", polServiceAddress, "")

	err := reg.Register(polService)
	pmt.g.Expect(err).ToNot(HaveOccurred())

	defer func() {
		err := reg.Deregister(polService)
		pmt.g.Expect(err).ToNot(HaveOccurred())
	}()

	gwService := registry.BuildGRPCService("com.owncloud.api.gateway", "", "any", "")

	err = reg.Register(gwService)
	pmt.g.Expect(err).ToNot(HaveOccurred())

	defer func() {
		err := reg.Deregister(gwService)
		pmt.g.Expect(err).ToNot(HaveOccurred())
	}()

	err = ogrpc.Configure()
	pmt.g.Expect(err).ToNot(HaveOccurred())

	srv, err := ogrpc.NewService(
		ogrpc.Name(polServiceName),
		ogrpc.Address(polServiceAddress),
		ogrpc.Namespace(polServiceNamespace),
		ogrpc.Context(ctx),
	)
	pmt.g.Expect(err).ToNot(HaveOccurred())

	defer func() {
		err := srv.Server().Stop()
		pmt.g.Expect(err).ToNot(HaveOccurred())
	}()

	err = policiesPG.RegisterPoliciesProviderHandler(srv.Server(), pmt)
	pmt.g.Expect(err).ToNot(HaveOccurred())

	rg.Add(srv.Run, func(err error) {
		pmt.g.Expect(err).ToNot(HaveOccurred())
	})

	go func() {
		err := rg.Run()
		pmt.g.Expect(err).ToNot(HaveOccurred())
	}()

	gatewayClient := &cs3mocks.GatewayAPIClient{}
	gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(pmt.gwClientStatResponse, nil)

	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		"com.owncloud.api.gateway",
		func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)
	defer pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")

	w := httptest.NewRecorder()

	p := middleware.Policies(
		pmt.regoQuery,
		middleware.WithRevaGatewaySelector(gatewaySelector),
	)(pmt)
	p.ServeHTTP(w, pmt.httpRequest)

	e(pmt.g, w, pmt.grpcEvaluateRequest, pmt.grpcEvaluateResponse)
}

func (pmt *policiesMiddlewareTester) ServeHTTP(writer http.ResponseWriter, request *http.Request) {}

func (pmt *policiesMiddlewareTester) Evaluate(ctx context.Context, request *policiesPG.EvaluateRequest, response *policiesPG.EvaluateResponse) error {
	err := pmt.eval(request, response)
	pmt.grpcEvaluateRequest = request
	pmt.grpcEvaluateResponse = response

	return err
}
