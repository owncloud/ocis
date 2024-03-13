package helpers

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func UploadFile(
	ctx context.Context,
	content io.Reader, // content won't be closed inside the method
	contentLength int64,
	ref *providerv1beta1.Reference,
	gwc gatewayv1beta1.GatewayAPIClient,
	token string,
	lockID string,
	insecure bool,
	logger log.Logger,
) error {
	opaque := &types.Opaque{
		Map: make(map[string]*types.OpaqueEntry),
	}

	strContentLength := strconv.FormatInt(contentLength, 10)
	if contentLength >= 0 {
		opaque.Map["Upload-Length"] = &types.OpaqueEntry{
			Decoder: "plain",
			Value:   []byte(strContentLength),
		}
	}

	req := &providerv1beta1.InitiateFileUploadRequest{
		Opaque: opaque,
		Ref:    ref,
		LockId: lockID,
		// TODO: if-match
		//Options: &providerv1beta1.InitiateFileUploadRequest_IfMatch{
		//	IfMatch: "",
		//},
	}

	resp, err := gwc.InitiateFileUpload(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Msg("UploadHelper: InitiateFileUpload failed")
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("UploadHelper: InitiateFileUpload failed with wrong status")
		return errors.New("InitiateFileUpload failed with status " + resp.Status.Code.String())
	}

	// if content length is 0, we're done. We don't upload anything to the target endpoint
	if contentLength == 0 {
		return nil
	}

	uploadEndpoint := ""
	uploadToken := ""
	hasUploadToken := false

	for _, proto := range resp.Protocols {
		if proto.Protocol == "simple" || proto.Protocol == "spaces" {
			uploadEndpoint = proto.UploadEndpoint
			uploadToken = proto.Token
			hasUploadToken = proto.Token != ""
			break
		}
	}

	if uploadEndpoint == "" {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Msg("UploadHelper: Upload endpoint or token is missing")
		return errors.New("upload endpoint or token is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
		Timeout: 10 * time.Second,
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadEndpoint, content)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Msg("UploadHelper: Could not create the request to the endpoint")
		return err
	}
	// "content" is an *http.body and doesn't fill the httpReq.ContentLength automatically
	// we need to fill the ContentLength ourselves, and must match the stream length in order
	// to prevent issues
	httpReq.ContentLength = contentLength

	if uploadToken != "" {
		// public link uploads have the token in the upload endpoint
		httpReq.Header.Add("X-Reva-Transfer", uploadToken)
	}
	// TODO: the access token shouldn't be needed
	httpReq.Header.Add("X-Access-Token", token)

	httpReq.Header.Add("X-Lock-Id", lockID)
	// TODO: better mechanism for the upload while locked, relies on patch in REVA
	//if lockID, ok := ctxpkg.ContextGetLockID(ctx); ok {
	//	httpReq.Header.Add("X-Lock-Id", lockID)
	//}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Msg("UploadHelper: Put request to the upload endpoint failed")
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("UploadLength", strContentLength).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Int("HttpCode", httpResp.StatusCode).
			Msg("UploadHelper: Put request to the upload endpoint failed with unexpected status")
		return errors.New("Put request failed with status " + strconv.Itoa(httpResp.StatusCode))
	}

	return nil
}
