package permissions

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
)

type ListPermissionsRequest struct {
	UserID string `json:"user_id"`
}

type ListPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

type ListPermissionsHandler interface {
	ListPermissions(context.Context, *ListPermissionsRequest, *ListPermissionsResponse) error
}

type listPermissionsHandler struct {
	r chi.Router
	h ListPermissionsHandler
}

func RegisterListPermissionsHandler(r chi.Router, i ListPermissionsHandler, middlewares ...func(http.Handler) http.Handler) {
	handler := &listPermissionsHandler{
		r: r,
		h: i,
	}

	r.MethodFunc("POST", "/api/v0/settings/permissions-list", handler.ListPermissions)
}

func (h listPermissionsHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	req := &ListPermissionsRequest{}
	resp := &ListPermissionsResponse{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	if err := h.h.ListPermissions(
		r.Context(),
		req,
		resp,
	); err != nil {
		if errors.Is(err, settings.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
