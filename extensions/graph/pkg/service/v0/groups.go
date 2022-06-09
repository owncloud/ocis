package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/CiscoM31/godata"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/extensions/graph/pkg/service/v0/errorcode"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const memberRefsLimit = 20

// GetGroups implements the Service interface.
func (g Graph) GetGroups(w http.ResponseWriter, r *http.Request) {
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		g.logger.Err(err).Interface("query", r.URL.Query()).Msg("query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	groups, err := g.identityBackend.GetGroups(r.Context(), r.URL.Query())

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	groups, err = sortGroups(odataReq, groups)
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

	if _, ok := grp.GetDisplayNameOk(); !ok {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Missing Required Attribute")
		return
	}

	// Disallow user-supplied IDs. It's supposed to be readonly. We're either
	// generating them in the backend ourselves or rely on the Backend's
	// storage (e.g. LDAP) to provide a unique ID.
	if _, ok := grp.GetIdOk(); ok {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "group id is a read-only attribute")
		return
	}

	if grp, err = g.identityBackend.CreateGroup(r.Context(), *grp); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if grp != nil && grp.Id != nil {
		g.publishEvent(events.GroupCreated{GroupID: *grp.Id})
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, grp)
}

// PatchGroup implements the Service interface.
func (g Graph) PatchGroup(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msg("Calling PatchGroup")
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
	changes := libregraph.NewGroup()
	err = json.NewDecoder(r.Body).Decode(changes)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if memberRefs, ok := changes.GetMembersodataBindOk(); ok {
		// The spec defines a limit of 20 members maxium per Request
		if len(memberRefs) > memberRefsLimit {
			errorcode.NotAllowed.Render(w, r, http.StatusInternalServerError,
				fmt.Sprintf("Request is limited to %d members", memberRefsLimit))
			return
		}
		memberIDs := make([]string, 0, len(memberRefs))
		for _, memberRef := range memberRefs {
			memberType, id, err := g.parseMemberRef(memberRef)
			if err != nil {
				errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Error parsing member@odata.bind values")
				return
			}
			g.logger.Debug().Str("memberType", memberType).Str("memberid", id).Msg("Add Member")
			// The MS Graph spec allows "directoryObject", "user", "group" and "organizational Contact"
			// we restrict this to users for now. Might add Groups as members later
			if memberType != "users" {
				errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Only user are allowed as group members")
				return
			}
			memberIDs = append(memberIDs, id)
		}
		err = g.identityBackend.AddMembersToGroup(r.Context(), groupID, memberIDs)
	}

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

	group, err := g.identityBackend.GetGroup(r.Context(), groupID, r.URL.Query())
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

	g.publishEvent(events.GroupDeleted{GroupID: groupID})
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
	memberType, id, err := g.parseMemberRef(*memberRefURL)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Error parsing @odata.id url")
		return
	}
	// The MS Graph spec allows "directoryObject", "user", "group" and "organizational Contact"
	// we restrict this to users for now. Might add Groups as members later
	if memberType != "users" {
		errorcode.InvalidRequest.Render(w, r, http.StatusInternalServerError, "Only user are allowed as group members")
		return
	}

	g.logger.Debug().Str("memberType", memberType).Str("id", id).Msg("Add Member")
	err = g.identityBackend.AddMembersToGroup(r.Context(), groupID, []string{id})

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	g.publishEvent(events.GroupMemberAdded{GroupID: groupID, UserID: id})
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// DeleteGroupMember implements the Service interface.
func (g Graph) DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling DeleteGroupMember")

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

	memberID := chi.URLParam(r, "memberID")
	memberID, err = url.PathUnescape(memberID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping group id failed")
		return
	}

	if memberID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
		return
	}
	g.logger.Debug().Str("groupID", groupID).Str("memberID", memberID).Msg("DeleteGroupMember")
	err = g.identityBackend.RemoveMemberFromGroup(r.Context(), groupID, memberID)

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	g.publishEvent(events.GroupMemberRemoved{GroupID: groupID, UserID: memberID})
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func (g Graph) parseMemberRef(ref string) (string, string, error) {
	memberURL, err := url.ParseRequestURI(ref)
	if err != nil {
		return "", "", err
	}
	segments := strings.Split(memberURL.Path, "/")
	if len(segments) < 2 {
		return "", "", errors.New("invalid member reference")
	}
	id := segments[len(segments)-1]
	memberType := segments[len(segments)-2]
	return memberType, id, nil
}

func sortGroups(req *godata.GoDataRequest, groups []*libregraph.Group) ([]*libregraph.Group, error) {
	var sorter sort.Interface
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return groups, nil
	}
	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case "displayName":
		sorter = groupsByDisplayName{groups}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == "desc" {
		sorter = sort.Reverse(sorter)
	}
	sort.Sort(sorter)
	return groups, nil
}
