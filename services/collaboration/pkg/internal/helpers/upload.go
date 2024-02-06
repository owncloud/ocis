package helpers

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strconv"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func UploadFile(ctx context.Context, content io.ReadCloser, ref *providerv1beta1.Reference, gwc gatewayv1beta1.GatewayAPIClient, token string, lockID string, insecure bool, logger log.Logger) error {

	req := &providerv1beta1.InitiateFileUploadRequest{
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
			Msg("UploadHelper: InitiateFileUpload failed")
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("UploadHelper: InitiateFileUpload failed with wrong status")
		return errors.New("InitiateFileUpload failed with status " + resp.Status.Code.String())
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
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadEndpoint, content)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Msg("UploadHelper: Could not create the request to the endpoint")
		return err
	}

	if uploadToken != "" {
		// public link uploads have the token in the upload endpoint
		httpReq.Header.Add("X-Reva-Transfer", uploadToken)
	}
	// TODO: the access token shouldn't be needed
	httpReq.Header.Add("x-access-token", token)

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
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Msg("UploadHelper: Put request to the upload endpoint failed")
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("RequestedLockID", lockID).
			Str("Endpoint", uploadEndpoint).
			Bool("HasUploadToken", hasUploadToken).
			Int("HttpCode", httpResp.StatusCode).
			Msg("UploadHelper: Put request to the upload endpoint failed with unexpected status")
		return errors.New("Put request failed with status " + strconv.Itoa(httpResp.StatusCode))
	}

	return nil
}
