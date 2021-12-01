package svc

import (
	"errors"
	"net/http"

	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	//msgraph "github.com/owncloud/open-graph-api-go" // FIXME add groups to open graph, needs OnPremisesSamAccountName and OnPremisesDomainName
)

// GetGroups implements the Service interface.
func (g Graph) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := g.identityBackend.GetGroups(r.Context(), r.URL.Query())

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: groups})
}

// GetGroup implements the Service interface.
func (g Graph) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
		return
	}

	group, err := g.identityBackend.GetGroup(r.Context(), groupID)
	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, group)
}
