package svc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CiscoM31/godata"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/extensions/graph/pkg/service/v0/errorcode"
)

// ChangeOwnPassword implements the Service interface. It allows the user to change
// its own password
func (g Graph) ChangeOwnPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := revactx.ContextGetUser(ctx)
	if !ok {
		g.logger.Error().Msg("user not in context")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	_, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		g.logger.Err(err).Interface("query", r.URL.Query()).Msg("query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	cpw := libregraph.NewPasswordChange()
	err = json.NewDecoder(r.Body).Decode(cpw)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	currentPw := cpw.GetCurrentPassword()
	if currentPw == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "current password cannot be empty")
		return
	}

	newPw := cpw.GetNewPassword()
	if newPw == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "new password cannot be empty")
		return
	}

	if newPw == currentPw {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "new password must be differnt from current password")
		return
	}

	authReq := &gateway.AuthenticateRequest{
		Type:         "basic",
		ClientId:     u.Username,
		ClientSecret: currentPw,
	}
	authRes, err := g.gatewayClient.Authenticate(r.Context(), authReq)
	if err != nil {
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if authRes.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "password change failed")
		return
	}

	newPwProfile := libregraph.NewPasswordProfile()
	newPwProfile.SetPassword(newPw)
	changes := libregraph.NewUser()
	changes.SetPasswordProfile(*newPwProfile)
	_, err = g.identityBackend.UpdateUser(ctx, u.Id.OpaqueId, *changes)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "password change failed")
		g.logger.Debug().Err(err).Str("userid", u.Id.OpaqueId).Msg("failed to update user password")
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
