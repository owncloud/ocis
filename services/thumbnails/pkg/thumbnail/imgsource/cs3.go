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
	"github.com/cs3org/reva/v2/pkg/bytesize"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const (
	// TokenTransportHeader holds the header key for the reva transfer token
	// "github.com/cs3org/reva/v2/internal/http/services/datagateway" is internal so we redeclare it here
	TokenTransportHeader = "X-Reva-Transfer"
)

// CS3 implements a CS3 image source
type CS3 struct {
	gatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	insecure         bool
	maxImageFileSize uint64
}

// NewCS3Source configures a new CS3 image source
func NewCS3Source(cfg config.Thumbnail, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], b bytesize.ByteSize) CS3 {
	return CS3{
		gatewaySelector:  gatewaySelector,
		insecure:         cfg.CS3AllowInsecure,
		maxImageFileSize: b.Bytes(),
	}
}

// Get downloads the file from a cs3 service
// The caller MUST make sure to close the returned ReadCloser
func (s CS3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	auth, ok := ContextGetAuthorization(ctx)
	if !ok {
		return nil, errors.ErrCS3AuthorizationMissing
	}
	ref, err := storagespace.ParseReference(path)
	if err != nil {
		// If the path is not a spaces reference try to handle it like a plain
		// path reference.
		ref = provider.Reference{
			Path: path,
		}
	}

	ctx = metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, auth)
	err = s.checkImageFileSize(ctx, ref)
	if err != nil {
		return nil, err
	}

	gwc, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	rsp, err := gwc.InitiateFileDownload(ctx, &provider.InitiateFileDownloadRequest{Ref: &ref})

	if err != nil {
		return nil, err
	}

	if rsp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not load image: %s", rsp.GetStatus().GetMessage())
	}
	var ep, tk string
	for _, p := range rsp.GetProtocols() {
		if p.GetProtocol() == "spaces" {
			ep, tk = p.GetDownloadEndpoint(), p.GetToken()
			break
		}
	}
	if (ep == "" || tk == "") && len(rsp.GetProtocols()) > 0 {
		ep, tk = rsp.GetProtocols()[0].GetDownloadEndpoint(), rsp.GetProtocols()[0].GetToken()
	}

	httpReq, err := rhttp.NewRequest(ctx, "GET", ep, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(revactx.TokenHeader, auth)
	httpReq.Header.Set(TokenTransportHeader, tk)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: s.insecure, //nolint:gosec
	}
	client := &http.Client{}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", path, resp.StatusCode)
	}

	return resp.Body, nil
}

func (s CS3) checkImageFileSize(ctx context.Context, ref provider.Reference) error {
	gwc, err := s.gatewaySelector.Next()
	if err != nil {
		return err
	}
	stat, err := gwc.Stat(ctx, &provider.StatRequest{Ref: &ref})
	if err != nil {
		return err
	}
	if stat.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return fmt.Errorf("could not stat image: %s", stat.GetStatus().GetMessage())
	}
	if stat.GetInfo().GetSize() > s.maxImageFileSize {
		return errors.ErrImageTooLarge
	}
	return nil
}
