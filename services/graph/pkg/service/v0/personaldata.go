package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
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

// ExportPersonalDataRequest is the body of the request
type ExportPersonalDataRequest struct {
	StorageLocation string `json:"storageLocation"`
}

// ExportPersonalData exports all personal data the system holds
func (g Graph) ExportPersonalData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := revactx.ContextMustGetUser(ctx)
	// Get location from request
	loc := ""
	if loc == "" {
		loc = _backupFileName
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
	go func() {
		time.Sleep(10 * time.Second)
		by, _ := json.Marshal(map[string]string{u.GetId().GetOpaqueId(): "no data stored"})
		b := bytes.NewBuffer(by)
		th := r.Header.Get(revaCtx.TokenHeader)
		err := g.upload(u, b, ref, th)
		fmt.Println("Upload error", err)
	}()

	w.WriteHeader(http.StatusOK)
	return
}

func (g Graph) upload(u *user.User, data io.Reader, ref *provider.Reference, th string) error {
	uReq := &provider.InitiateFileUploadRequest{
		Ref: ref,
		//Opaque: &typespb.Opaque{
		//Map: map[string]*typespb.OpaqueEntry{
		//"Upload-Length": {
		//Decoder: "plain",
		//// TODO: handle case where size is not known in advance
		//Value: []byte(strconv.FormatUint(cp.sourceInfo.GetSize(), 10)),
		//},
		//},
		//},
	}

	gwc := g.GetGatewayClient()
	ctx, _, err := utils.Impersonate(u.GetId(), gwc.(gateway.GatewayAPIClient), g.config.MachineAuthAPIKey)
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

	httpUploadReq, err := rhttp.NewRequest(ctx, "PUT", uploadEP, data)
	if err != nil {
		return err
	}
	httpUploadReq.Header.Set(TokenTransportHeader, uploadToken)

	httpUploadRes, err := rhttp.GetHTTPClient(
		// rhttp.Timeout(time.Duration(conf.Timeout*int64(time.Second))),
		rhttp.Insecure(true),
	).Do(httpUploadReq)
	if err != nil {
		return err
	}
	defer httpUploadRes.Body.Close()
	if httpUploadRes.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status uploading file: %d", httpUploadRes.StatusCode)
	}

	return nil
}
