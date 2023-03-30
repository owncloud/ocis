package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
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
		_, _ = w.Write([]byte("only json format is supported for personal data export"))
		return
	case ".json":
		marsh = json.Marshal
	}

	ref := &provider.Reference{
		ResourceId: &provider.ResourceId{SpaceId: u.GetId().GetOpaqueId(), OpaqueId: u.GetId().GetOpaqueId()},
		Path:       loc,
	}

	// touch file
	if err := mustTouchFile(ctx, ref, g.GetGatewayClient()); err != nil {
		g.logger.Error().Err(err).Msg("error touching file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// go start gathering
	go g.GatherPersonalData(u, ref, r.Header.Get(revactx.TokenHeader), marsh)

	w.WriteHeader(http.StatusCreated)
}

// GatherPersonalData will all gather all personal data of the user and save it to a file in the users personal space
func (g Graph) GatherPersonalData(usr *user.User, ref *provider.Reference, token string, marsh Marshaller) {
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
	var errmsg string
	if err := g.upload(usr, by, ref, token); err != nil {
		g.logger.Error().Err(err).Msg("failed uploading personal data export")
		errmsg = err.Error()
	}

	if err := events.Publish(g.eventsPublisher, events.PersonalDataExtracted{
		Executant: usr.GetId(),
		Timestamp: time.Now(),
		ErrorMsg:  errmsg,
	}); err != nil {
		g.logger.Error().Err(err).Msg("cannot publish event")
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
	ctx = revactx.ContextSetToken(ctx, th)
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

// touches the file, creating folders if necessary
func mustTouchFile(ctx context.Context, ref *provider.Reference, gwc gateway.GatewayAPIClient) error {
	if err := touchFile(ctx, ref, gwc); err == nil {
		return nil
	}

	if err := createFolders(ctx, ref, gwc); err != nil {
		return err
	}

	return touchFile(ctx, ref, gwc)
}

func touchFile(ctx context.Context, ref *provider.Reference, gwc gateway.GatewayAPIClient) error {
	resp, err := gwc.TouchFile(ctx, &provider.TouchFileRequest{
		Opaque: utils.AppendPlainToOpaque(nil, "markprocessing", "true"),
		Ref:    ref,
	})
	if err != nil {
		return err
	}
	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return fmt.Errorf("unexpected statuscode while touching file: %d %s", resp.GetStatus().GetCode(), resp.GetStatus().GetMessage())
	}
	return nil
}

func createFolders(ctx context.Context, ref *provider.Reference, gwc gateway.GatewayAPIClient) error {
	var paths []string
	p := filepath.Dir(ref.GetPath())
	for p != "." {
		paths = append([]string{p}, paths...)
		p = filepath.Dir(p)
	}

	for _, p := range paths {
		r := &provider.Reference{ResourceId: ref.GetResourceId(), Path: p}
		resp, err := gwc.CreateContainer(ctx, &provider.CreateContainerRequest{Ref: r})
		if err != nil {
			return err
		}

		code := resp.GetStatus().GetCode()
		if code != rpc.Code_CODE_OK && code != rpc.Code_CODE_ALREADY_EXISTS {
			return fmt.Errorf("unexpected statuscode while creating folder: %d %s", code, resp.GetStatus().GetMessage())
		}
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
