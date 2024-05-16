package connector

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strconv"
	"time"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/middleware"
	"github.com/rs/zerolog"
)

// ContentConnectorService is the interface to implement the "File contents"
// endpoint. Basically upload and download contents.
// All operations need a context containing a WOPI context and, optionally,
// a zerolog logger.
// Target file is within the WOPI context
type ContentConnectorService interface {
	// GetFile downloads the file and write its contents in the provider writer
	GetFile(ctx context.Context, writer io.Writer) error
	// PutFile uploads the stream up to the stream length. The file should be
	// locked beforehand, so the lockID needs to be provided.
	// The current lockID will be returned ONLY if a conflict happens (the file is
	// locked with a different lockID)
	PutFile(ctx context.Context, stream io.Reader, streamLength int64, lockID string) (string, error)
}

// ContentConnector implements the "File contents" endpoint.
// Basically, the ContentConnector handles downloads (GetFile) and
// uploads (PutFile)
// Note that operations might return any kind of error, not just ConnectorError
type ContentConnector struct {
	gwc gatewayv1beta1.GatewayAPIClient
	cfg *config.Config
}

// NewContentConnector creates a new content connector
func NewContentConnector(gwc gatewayv1beta1.GatewayAPIClient, cfg *config.Config) *ContentConnector {
	return &ContentConnector{
		gwc: gwc,
		cfg: cfg,
	}
}

// GetFile downloads the file from the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getfile
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// The contents of the file will be written directly into the writer passed as
// parameter.
func (c *ContentConnector) GetFile(ctx context.Context, writer io.Writer) error {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return err
	}

	logger := zerolog.Ctx(ctx)

	// Initiate download request
	req := &providerv1beta1.InitiateFileDownloadRequest{
		Ref: &wopiContext.FileReference,
	}

	if wopiContext.ViewMode == appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY && wopiContext.ViewOnlyToken != "" {
		ctx = revactx.ContextSetToken(ctx, wopiContext.ViewOnlyToken)
	}
	resp, err := c.gwc.InitiateFileDownload(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("GetFile: InitiateFileDownload failed")
		return err
	}

	if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("GetFile: InitiateFileDownload failed with wrong status")
		return NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
	}

	// Figure out the download endpoint and download token
	downloadEndpoint := ""
	downloadToken := ""
	hasDownloadToken := false

	for _, proto := range resp.GetProtocols() {
		if proto.GetProtocol() == "simple" || proto.GetProtocol() == "spaces" {
			downloadEndpoint = proto.GetDownloadEndpoint()
			downloadToken = proto.GetToken()
			hasDownloadToken = proto.GetToken() != ""
			break
		}
	}

	if downloadEndpoint == "" {
		logger.Error().
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("GetFile: Download endpoint or token is missing")
		return NewConnectorError(500, "GetFile: Download endpoint is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.cfg.CS3Api.DataGateway.Insecure,
			},
		},
	}

	// Prepare the request to download the file
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadEndpoint, bytes.NewReader([]byte("")))
	if err != nil {
		logger.Error().
			Err(err).
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("GetFile: Could not create the request to the endpoint")
		return err
	}
	if downloadToken != "" {
		// public link downloads have the token in the download endpoint
		httpReq.Header.Add("X-Reva-Transfer", downloadToken)
	}
	if wopiContext.ViewMode == appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY && wopiContext.ViewOnlyToken != "" {
		httpReq.Header.Add("X-Access-Token", wopiContext.ViewOnlyToken)
	} else {
		httpReq.Header.Add("X-Access-Token", wopiContext.AccessToken)
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("Endpoint", downloadEndpoint).
			Bool("HasDownloadToken", hasDownloadToken).
			Msg("GetFile: Get request to the download endpoint failed")
		return err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logger.Error().
			Err(err).
			Int("HttpCode", httpResp.StatusCode).
			Msg("GetFile: downloading the file failed")
		return NewConnectorError(500, "GetFile: Downloading the file failed")
	}

	// Copy the download into the writer
	_, err = io.Copy(writer, httpResp.Body)
	if err != nil {
		logger.Error().Msg("GetFile: copying the file content to the response body failed")
		return err
	}

	logger.Debug().Msg("GetFile: success")
	return nil
}

// PutFile uploads the file to the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putfile
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// The contents of the file will be read from the stream. The full stream
// length must be provided in order to upload the file.
//
// A lock ID must be provided for the upload (which must match the lock in the
// file). The only case where an empty lock ID can be used is if the target
// file has 0 size.
//
// This method will return the lock ID that should be returned in case of a
// conflict, otherwise it will return an empty string. This means that if the
// method returns a ConnectorError with code 409, the returned string is the
// lock ID that should be used in the X-WOPI-Lock header. In other error
// cases or if the method is successful, an empty string will be returned
// (check for err != nil to know if something went wrong)
func (c *ContentConnector) PutFile(ctx context.Context, stream io.Reader, streamLength int64, lockID string) (string, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return "", err
	}

	logger := zerolog.Ctx(ctx).With().
		Str("RequestedLockID", lockID).
		Int64("UploadLength", streamLength).
		Logger()

	// We need a stat call on the target file in order to get both the lock
	// (if any) and the current size of the file
	statRes, err := c.gwc.Stat(ctx, &providerv1beta1.StatRequest{
		Ref: &wopiContext.FileReference,
	})
	if err != nil {
		logger.Error().Err(err).Msg("PutFile: stat failed")
		return "", err
	}

	if statRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("StatusCode", statRes.GetStatus().GetCode().String()).
			Str("StatusMsg", statRes.GetStatus().GetMessage()).
			Msg("PutFile: stat failed with unexpected status")
		return "", NewConnectorError(500, statRes.GetStatus().GetCode().String()+" "+statRes.GetStatus().GetMessage())
	}

	// If there is a lock and it mismatches, return 409
	if statRes.GetInfo().GetLock() != nil && statRes.GetInfo().GetLock().GetLockId() != lockID {
		logger.Error().
			Str("LockID", statRes.GetInfo().GetLock().GetLockId()).
			Msg("PutFile: wrong lock")
		// onlyoffice says it's required to send the current lockId, MS doesn't say anything
		return statRes.GetInfo().GetLock().GetLockId(), NewConnectorError(409, "Wrong lock")
	}

	// only unlocked uploads can go through if the target file is empty,
	// otherwise the X-WOPI-Lock header is required even if there is no lock on the file
	// This is part of the onlyoffice documentation (https://api.onlyoffice.com/editors/wopi/restapi/putfile)
	// Wopivalidator fails some tests if we don't also check for the X-WOPI-Lock header.
	if lockID == "" && statRes.GetInfo().GetLock() == nil && statRes.GetInfo().GetSize() > 0 {
		logger.Error().Msg("PutFile: file must be locked first")
		// onlyoffice says to send an empty string if the file is unlocked, MS doesn't say anything
		return "", NewConnectorError(409, "File must be locked first")
	}

	// Prepare the data to initiate the upload
	opaque := &types.Opaque{
		Map: make(map[string]*types.OpaqueEntry),
	}

	if streamLength >= 0 {
		opaque.Map["Upload-Length"] = &types.OpaqueEntry{
			Decoder: "plain",
			Value:   []byte(strconv.FormatInt(streamLength, 10)),
		}
	}

	req := &providerv1beta1.InitiateFileUploadRequest{
		Opaque: opaque,
		Ref:    &wopiContext.FileReference,
		LockId: lockID,
		Options: &providerv1beta1.InitiateFileUploadRequest_IfMatch{
			IfMatch: statRes.GetInfo().GetEtag(),
		},
	}

	// Initiate the upload request
	resp, err := c.gwc.InitiateFileUpload(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("UploadHelper: InitiateFileUpload failed")
		return "", err
	}

	if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("UploadHelper: InitiateFileUpload failed with wrong status")
		return "", NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
	}

	// if the content length is greater than 0, we need to upload the content to the
	// target endpoint, otherwise we're done
	if streamLength > 0 {

		uploadEndpoint := ""
		uploadToken := ""
		hasUploadToken := false

		for _, proto := range resp.GetProtocols() {
			if proto.GetProtocol() == "simple" || proto.GetProtocol() == "spaces" {
				uploadEndpoint = proto.GetUploadEndpoint()
				uploadToken = proto.GetToken()
				hasUploadToken = proto.GetToken() != ""
				break
			}
		}

		if uploadEndpoint == "" {
			logger.Error().
				Str("Endpoint", uploadEndpoint).
				Bool("HasUploadToken", hasUploadToken).
				Msg("UploadHelper: Upload endpoint or token is missing")
			return "", NewConnectorError(500, "upload endpoint or token is missing")
		}

		httpClient := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.cfg.CS3Api.DataGateway.Insecure,
				},
			},
			Timeout: 10 * time.Second,
		}

		// prepare the request to upload the contents to the upload endpoint
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadEndpoint, stream)
		if err != nil {
			logger.Error().
				Err(err).
				Str("Endpoint", uploadEndpoint).
				Bool("HasUploadToken", hasUploadToken).
				Msg("UploadHelper: Could not create the request to the endpoint")
			return "", err
		}
		// "stream" is an *http.body and doesn't fill the httpReq.ContentLength automatically
		// we need to fill the ContentLength ourselves, and must match the stream length in order
		// to prevent issues
		httpReq.ContentLength = streamLength

		if uploadToken != "" {
			// public link uploads have the token in the upload endpoint
			httpReq.Header.Add("X-Reva-Transfer", uploadToken)
		}
		httpReq.Header.Add("X-Access-Token", wopiContext.AccessToken)

		httpReq.Header.Add("X-Lock-Id", lockID)
		// TODO: better mechanism for the upload while locked, relies on patch in REVA
		//if lockID, ok := ctxpkg.ContextGetLockID(ctx); ok {
		//	httpReq.Header.Add("X-Lock-Id", lockID)
		//}

		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			logger.Error().
				Err(err).
				Str("Endpoint", uploadEndpoint).
				Bool("HasUploadToken", hasUploadToken).
				Msg("UploadHelper: Put request to the upload endpoint failed")
			return "", err
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			logger.Error().
				Str("Endpoint", uploadEndpoint).
				Bool("HasUploadToken", hasUploadToken).
				Int("HttpCode", httpResp.StatusCode).
				Msg("UploadHelper: Put request to the upload endpoint failed with unexpected status")
			return "", NewConnectorError(500, "PutFile: Uploading the file failed")
		}
	}

	logger.Debug().Msg("PutFile: success")
	return "", nil
}
