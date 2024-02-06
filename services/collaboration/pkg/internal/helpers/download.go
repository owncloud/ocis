package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

func DownloadFile(
	ctx context.Context,
	ref *providerv1beta1.Reference,
	gwc gatewayv1beta1.GatewayAPIClient,
	token string,
	insecure bool,
	logger log.Logger,
) (http.Response, error) {

	req := &providerv1beta1.InitiateFileDownloadRequest{
		Ref: ref,
	}

	resp, err := gwc.InitiateFileDownload(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Msg("DownloadHelper: InitiateFileDownload failed")
		return http.Response{}, err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("DownloadHelper: InitiateFileDownload failed with wrong status")
		return http.Response{}, errors.New("InitiateFileDownload failed with status " + resp.Status.Code.String())
	}

	downloadEndpoint := ""
	downloadToken := ""
	hasDownloadToken := false

	for _, proto := range resp.Protocols {
		if proto.Protocol == "simple" || proto.Protocol == "spaces" {
			downloadEndpoint = proto.DownloadEndpoint
			downloadToken = proto.Token
			hasDownloadToken = proto.Token != ""
		}
	}

	if downloadEndpoint == "" {
		logger.Error().
			Str("FileReference", ref.String()).
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("DownloadHelper: Download endpoint or token is missing")
		return http.Response{}, errors.New("download endpoint is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadEndpoint, bytes.NewReader([]byte("")))
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("DownloadHelper: Could not create the request to the endpoint")
		return http.Response{}, err
	}
	if downloadToken != "" {
		// public link downloads have the token in the download endpoint
		httpReq.Header.Add("X-Reva-Transfer", downloadToken)
	}
	// TODO: the access token shouldn't be needed
	httpReq.Header.Add("x-access-token", token)

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("FileReference", ref.String()).
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("DownloadHelper: Get request to the download endpoint failed")
		return http.Response{}, err
	}

	return *httpResp, nil
}
