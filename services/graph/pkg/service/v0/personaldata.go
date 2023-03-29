package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revaCtx "github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/utils"
)

var (
	_backupFileName = "personal_data_export.json"

	// TokenTransportHeader holds the header key for the reva transfer token
	TokenTransportHeader = "X-Reva-Transfer"
)

// Marshaller is the common interface for a marshaller
type Marshaller func(any) ([]byte, error)

// ExportPersonalDataRequest is the body of the request
type ExportPersonalDataRequest struct {
	StorageLocation string `json:"storageLocation"`
}

// ExportPersonalData exports all personal data the system holds
func (g Graph) ExportPersonalData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := revactx.ContextMustGetUser(ctx)
	// Get location from request
	loc := getLocation(r)

	// prepare marshaller
	var marsh Marshaller
	switch filepath.Ext(loc) {
	default:
		g.logger.Info().Str("path", loc).Msg("invalid location")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("only json format is supported for personal data export"))
		return
	case ".json":
		marsh = json.Marshal
	}

	ref := &provider.Reference{
		ResourceId: &provider.ResourceId{SpaceId: u.GetId().GetOpaqueId(), OpaqueId: u.GetId().GetOpaqueId()},
		Path:       loc,
	}

	// touch file
	gwc := g.GetGatewayClient()
	resp, err := gwc.TouchFile(ctx, &provider.TouchFileRequest{
		Opaque: utils.AppendPlainToOpaque(nil, "markprocessing", "true"),
		Ref:    ref,
	})
	if err != nil || resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Str("status", resp.GetStatus().GetMessage()).Msg("error touching file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// go start gathering
	go g.GatherPersonalData(u, ref, r.Header.Get(revaCtx.TokenHeader), marsh)

	w.WriteHeader(http.StatusOK)
}

// GatherPersonalData will all gather all personal data of the user and save it to a file in the users personal space
func (g Graph) GatherPersonalData(usr *user.User, ref *provider.Reference, token string, marsh Marshaller) {
	// TMP - Delay processing - comment if you read this on PR
	time.Sleep(10 * time.Second)

	// create data
	data := make(map[string]interface{})

	// reva user
	data["user"] = usr

	// marshal
	by, err := marsh(data)
	if err != nil {
		g.logger.Error().Err(err).Msg("cannot marshal personal user data")
		return
	}

	// upload
	if err := g.upload(usr, by, ref, token); err != nil {
		g.logger.Error().Err(err).Msg("failed uploading personal data export")
	}
}

func (g Graph) upload(u *user.User, data []byte, ref *provider.Reference, th string) error {
	uReq := &provider.InitiateFileUploadRequest{
		Ref:    ref,
		Opaque: utils.AppendPlainToOpaque(nil, "Upload-Length", strconv.FormatUint(uint64(len(data)), 10)),
	}

	gwc := g.GetGatewayClient()
	ctx, err := utils.ImpersonateUser(u, gwc, g.config.MachineAuthAPIKey)
	if err != nil {
		return err
	}
	ctx = revaCtx.ContextSetToken(ctx, th)
	uRes, err := gwc.InitiateFileUpload(ctx, uReq)
	if err != nil {
		return err
	}

	if uRes.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("wrong status code while initiating upload: %s", uRes.GetStatus().GetMessage())
	}

	var uploadEP, uploadToken string
	for _, p := range uRes.Protocols {
		if p.Protocol == "simple" {
			uploadEP, uploadToken = p.UploadEndpoint, p.Token
		}
	}

	httpUploadReq, err := rhttp.NewRequest(ctx, "PUT", uploadEP, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	httpUploadReq.Header.Set(TokenTransportHeader, uploadToken)

	httpUploadRes, err := rhttp.GetHTTPClient(rhttp.Insecure(true)).Do(httpUploadReq)
	if err != nil {
		return err
	}
	defer httpUploadRes.Body.Close()
	if httpUploadRes.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status uploading file: %d", httpUploadRes.StatusCode)
	}

	return nil
}

func getLocation(r *http.Request) string {
	// from body
	var req ExportPersonalDataRequest
	if b, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(b, &req); err == nil && req.StorageLocation != "" {
			return req.StorageLocation
		}
	}

	// from header?

	return _backupFileName
}
