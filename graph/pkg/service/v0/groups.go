package svc

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

// DeleteGroup implements the Service interface.
func (g Graph) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	groupID, err := url.PathUnescape(groupID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping group id failed")
		return
	}

	if groupID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
		return
	}

	err = g.identityBackend.DeleteGroup(r.Context(), groupID)

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func (g Graph) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	groupID, err := url.PathUnescape(groupID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping group id failed")
		return
	}

	if groupID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
		return
	}

	members, err := g.identityBackend.GetGroupMembers(r.Context(), groupID)
	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, members)
}

// PostGroupMember implements the Service interface.
func (g Graph) PostGroupMember(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling PostGroupMember")

	groupID := chi.URLParam(r, "groupID")
	groupID, err := url.PathUnescape(groupID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping group id failed")
		return
	}

	if groupID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
		return
	}
	memberRef := libregraph.NewMemberReference()
	err = json.NewDecoder(r.Body).Decode(memberRef)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	memberRefURL, ok := memberRef.GetOdataIdOk()
	if !ok {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "@odata.id refernce is missing")
		return
	}
	memberURL, err := url.ParseRequestURI(*memberRefURL)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Error parsing @odata.id url")
		return
	}
	segments := strings.Split(memberURL.Path, "/")
	if len(segments) < 2 {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Error parsing @odata.id url path")
		return
	}
	id := segments[len(segments)-1]
	memberType := segments[len(segments)-2]
	// The MS Graph spec allows "directoryObject", "user", "group" and "organizational Contact"
	// we restrict this to users for now. Might add Groups as members later
	if memberType != "users" {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Only user are allowed as group members")
		return
	}

	g.logger.Debug().Str("memberType", memberType).Str("id", id).Msg("Add Member")
	err = g.identityBackend.AddMemberToGroup(r.Context(), groupID, id)

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
