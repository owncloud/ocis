package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rhttp"
	"github.com/cs3org/reva/pkg/token"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type CS3 struct {
	client gateway.GatewayAPIClient
}

func NewCS3Source(c gateway.GatewayAPIClient) CS3 {
	return CS3{
		client: c,
	}
}

// Get downloads the file from a cs3 service
// The caller MUST make sure to close the returned ReadCloser
func (s CS3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	auth, ok := ContextGetAuthorization(ctx)
	if !ok {
		return nil, errors.New("cs3source: authorization missing")
	}
	ctx = metadata.AppendToOutgoingContext(context.Background(), token.TokenHeader, auth)
	rsp, err := s.client.InitiateFileDownload(ctx, &provider.InitiateFileDownloadRequest{
		Ref: &provider.Reference{Path: path},
	})

	if err != nil {
		return nil, err
	}

	if rsp.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not load image: %s", rsp.Status.Message)
	}
	var ep, tk string
	for _, p := range rsp.Protocols {
		if p.Protocol == "simple" {
			ep, tk = p.DownloadEndpoint, p.Token
		}
	}

	httpReq, err := rhttp.NewRequest(ctx, "GET", ep, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(token.TokenHeader, auth)
	httpReq.Header.Set("X-REVA-TRANSFER", tk)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	client := &http.Client{}

	resp, err := client.Do(httpReq) // nolint:bodyclose
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", path, resp.StatusCode)
	}

	return resp.Body, nil
}
