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
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"

	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// GetSchools implements the Service interface.
func (g Graph) GetSchools(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling get schools")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get schools: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	schools, err := g.identityEducationBackend.GetSchools(r.Context(), r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Msg("could not get schools: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	schools, err = sortSchools(odataReq, schools)
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("cannot get schools: could not sort schools according to query")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: schools})
}

// PostSchool implements the Service interface.
func (g Graph) PostSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling post school")
	school := libregraph.NewEducationSchool()
	err := json.NewDecoder(r.Body).Decode(school)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not create school: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	if _, ok := school.GetDisplayNameOk(); !ok {
		logger.Debug().Err(err).Interface("school", school).Msg("could not create school: missing required attribute")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Missing Required Attribute")
		return
	}

	// Disallow user-supplied IDs. It's supposed to be readonly. We're either
	// generating them in the backend ourselves or rely on the Backend's
	// storage (e.g. LDAP) to provide a unique ID.
	if _, ok := school.GetIdOk(); ok {
		logger.Debug().Msg("could not create school: id is a read-only attribute")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "school id is a read-only attribute")
		return
	}

	if school, err = g.identityEducationBackend.CreateSchool(r.Context(), *school); err != nil {
		logger.Debug().Interface("school", school).Msg("could not create school: backend error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if school != nil && school.Id != nil {
		e := events.SchoolCreated{SchoolID: *school.Id}
		if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
			e.Executant = currentUser.GetId()
		}
		g.publishEvent(e)
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, school)
}

// PatchSchool implements the Service interface.
func (g Graph) PatchSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling patch school")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Str("id", schoolID).Msg("could not change school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not change school: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}
	changes := libregraph.NewEducationSchool()
	err = json.NewDecoder(r.Body).Decode(changes)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not change school: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	e := events.SchoolFeatureChanged{SchoolID: schoolID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// GetSchool implements the Service interface.
func (g Graph) GetSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling get school")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Str("id", schoolID).Msg("could not get school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
	}

	if schoolID == "" {
		logger.Debug().Msg("could not get school: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	logger.Debug().
		Str("id", schoolID).
		Interface("query", r.URL.Query()).
		Msg("calling get school on backend")
	school, err := g.identityEducationBackend.GetSchool(r.Context(), schoolID, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Msg("could not get school: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, school)
}

// DeleteSchool implements the Service interface.
func (g Graph) DeleteSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling delete school")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Err(err).Str("id", schoolID).Msg("could not delete school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not delete school: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	logger.Debug().Str("id", schoolID).Msg("calling delete school on backend")
	err = g.identityEducationBackend.DeleteSchool(r.Context(), schoolID)

	if err != nil {
		logger.Debug().Err(err).Msg("could not delete school: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	e := events.SchoolDeleted{SchoolID: schoolID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// GetSchoolMembers implements the Service interface.
func (g Graph) GetSchoolMembers(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling get school members")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Str("id", schoolID).Msg("could not get school members: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not get school members: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	logger.Debug().Str("id", schoolID).Msg("calling get school members on backend")
	members, err := g.identityEducationBackend.GetSchoolMembers(r.Context(), schoolID)
	if err != nil {
		logger.Debug().Err(err).Msg("could not get school members: backend error")
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

// PostSchoolMember implements the Service interface.
func (g Graph) PostSchoolMember(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("Calling post school member")

	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().
			Err(err).
			Str("id", schoolID).
			Msg("could not add member to school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not add school member: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}
	memberRef := libregraph.NewMemberReference()
	err = json.NewDecoder(r.Body).Decode(memberRef)
	if err != nil {
		logger.Debug().
			Err(err).
			Interface("body", r.Body).
			Msg("could not add school member: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	memberRefURL, ok := memberRef.GetOdataIdOk()
	if !ok {
		logger.Debug().Msg("could not add school member: @odata.id reference is missing")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "@odata.id reference is missing")
		return
	}
	memberType, id, err := g.parseMemberRef(*memberRefURL)
	if err != nil {
		logger.Debug().Err(err).Msg("could not add school member: error parsing @odata.id url")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Error parsing @odata.id url")
		return
	}
	// The MS Graph spec allows "directoryObject", "user", "school" and "organizational Contact"
	// we restrict this to users for now. Might add Schools as members later
	if memberType != "users" {
		logger.Debug().Str("type", memberType).Msg("could not add school member: Only users are allowed as school members")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Only users are allowed as school members")
		return
	}

	logger.Debug().Str("memberType", memberType).Str("id", id).Msg("calling add member on backend")
	err = g.identityEducationBackend.AddMembersToSchool(r.Context(), schoolID, []string{id})

	if err != nil {
		logger.Debug().Err(err).Msg("could not add school member: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	e := events.SchoolMemberAdded{SchoolID: schoolID, UserID: id}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// DeleteSchoolMember implements the Service interface.
func (g Graph) DeleteSchoolMember(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling delete school member")

	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Err(err).Str("id", schoolID).Msg("could not delete school member: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not delete school member: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	memberID := chi.URLParam(r, "memberID")
	memberID, err = url.PathUnescape(memberID)
	if err != nil {
		logger.Debug().Err(err).Str("id", memberID).Msg("could not delete school member: unescaping member id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping member id failed")
		return
	}

	if memberID == "" {
		logger.Debug().Msg("could not delete school member: missing member id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing member id")
		return
	}
	logger.Debug().Str("schoolID", schoolID).Str("memberID", memberID).Msg("calling delete member on backend")
	err = g.identityEducationBackend.RemoveMemberFromSchool(r.Context(), schoolID, memberID)

	if err != nil {
		logger.Debug().Err(err).Msg("could not delete school member: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	e := events.SchoolMemberRemoved{SchoolID: schoolID, UserID: memberID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func sortSchools(req *godata.GoDataRequest, schools []*libregraph.EducationSchool) ([]*libregraph.EducationSchool, error) {
	var sorter sort.Interface
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return schools, nil
	}
	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case "displayName":
		sorter = schoolsByDisplayName{schools}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == "desc" {
		sorter = sort.Reverse(sorter)
	}
	sort.Sort(sorter)
	return schools, nil
}
