package svc

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	libregraph "github.com/owncloud/libre-graph-api-go"
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

// PostGroup implements the Service interface.
func (g Graph) PostGroup(w http.ResponseWriter, r *http.Request) {
	grp := libregraph.NewGroup()
	err := json.NewDecoder(r.Body).Decode(grp)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if isNilOrEmpty(grp.DisplayName) {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Missing Required Attribute")
		return
	}

	// Disallow user-supplied IDs. It's supposed to be readonly. We're either
	// generating them in the backend ourselves or rely on the Backend's
	// storage (e.g. LDAP) to provide a unique ID.
	if !isNilOrEmpty(grp.Id) {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "group id is a read-only attribute")
		return
	}

	if grp, err = g.identityBackend.CreateGroup(r.Context(), *grp); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, grp)
}

// GetGroup implements the Service interface.
func (g Graph) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	groupID, err := url.PathUnescape(groupID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping group id failed")
	}

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
