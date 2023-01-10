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
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// GetEducationSchools implements the Service interface.
func (g Graph) GetEducationSchools(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling get schools")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get schools: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	schools, err := g.identityEducationBackend.GetEducationSchools(r.Context(), r.URL.Query())
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

	schools, err = sortEducationSchools(odataReq, schools)
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("cannot get schools: could not sort schools according to query")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: schools})
}

// PostEducationSchool implements the Service interface.
func (g Graph) PostEducationSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling post school")
	school := libregraph.NewEducationSchool()
	err := json.NewDecoder(r.Body).Decode(school)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not create school: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
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

	if _, ok := school.GetDisplayNameOk(); !ok {
		logger.Debug().Interface("school", school).Msg("could not create school: missing required attribute")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Missing Required Attribute")
		return
	}

	if _, ok := school.GetSchoolNumberOk(); !ok {
		logger.Debug().Interface("school", school).Msg("could not create school: missing required attribute")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Missing Required Attribute")
		return
	}

	if school, err = g.identityEducationBackend.CreateEducationSchool(r.Context(), *school); err != nil {
		logger.Debug().Err(err).Interface("school", school).Msg("could not create school: backend error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	/* TODO requires reva changes
	if school != nil && school.Id != nil {
		e := events.SchoolCreated{SchoolID: *school.Id}
		if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
			e.Executant = currentUser.GetId()
		}
		g.publishEvent(e)
	}
	*/

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, school)
}

// PatchEducationSchool implements the Service interface.
func (g Graph) PatchEducationSchool(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling patch school")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Str("id", schoolID).Msg("could not update school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not update school: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	school := libregraph.NewEducationSchool()
	err = json.NewDecoder(r.Body).Decode(school)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not update school: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	if school, err = g.identityEducationBackend.UpdateEducationSchool(r.Context(), schoolID, *school); err != nil {
		logger.Debug().Err(err).Interface("school", school).Msg("could not update school: backend error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	/* TODO requires reva changes
	e := events.SchoolFeatureChanged{SchoolID: schoolID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)
	*/

	render.Status(r, http.StatusOK)
	render.JSON(w, r, school)
}

// GetEducationSchool implements the Service interface.
func (g Graph) GetEducationSchool(w http.ResponseWriter, r *http.Request) {
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
	school, err := g.identityEducationBackend.GetEducationSchool(r.Context(), schoolID, r.URL.Query())
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

// DeleteEducationSchool implements the Service interface.
func (g Graph) DeleteEducationSchool(w http.ResponseWriter, r *http.Request) {
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

	logger.Debug().Str("schoolID", schoolID).Msg("Getting users of school")
	users, err := g.identityEducationBackend.GetEducationSchoolUsers(r.Context(), schoolID)
	if err != nil {
		logger.Debug().Err(err).Msg("could not get school users: backend error")
		renderInternalServerError(w, r, err)
		return
	}

	for _, user := range users {
		logger.Debug().Str("schoolID", schoolID).Str("userID", *user.Id).Msg("calling delete member on backend")
		if err := g.identityEducationBackend.RemoveUserFromEducationSchool(r.Context(), schoolID, *user.Id); err != nil {
			if errors.Is(err, identity.ErrNotFound) {
				logger.Debug().Str("schoolID", schoolID).Str("userID", *user.Id).Msg("user not found")
				continue
			}
			logger.Debug().Err(err).Msg("could not delete school member: backend error")
			renderInternalServerError(w, r, err)
		}
	}

	logger.Debug().Str("id", schoolID).Msg("calling delete school on backend")
	err = g.identityEducationBackend.DeleteEducationSchool(r.Context(), schoolID)

	if err != nil {
		logger.Debug().Err(err).Msg("could not delete school: backend error")
		renderInternalServerError(w, r, err)
		return
	}

	/* TODO requires reva changes
	e := events.SchoolDeleted{SchoolID: schoolID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)
	*/

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func renderInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	var errcode errorcode.Error
	if errors.As(err, &errcode) {
		errcode.Render(w, r)
	} else {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
	}
}

// GetEducationSchoolUsers implements the Service interface.
func (g Graph) GetEducationSchoolUsers(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling get school users")
	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().Str("id", schoolID).Msg("could not get school users: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not get school users: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}

	logger.Debug().Str("id", schoolID).Msg("calling get school users on backend")
	users, err := g.identityEducationBackend.GetEducationSchoolUsers(r.Context(), schoolID)
	if err != nil {
		logger.Debug().Err(err).Msg("could not get school users: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}

// PostEducationSchoolUser implements the Service interface.
func (g Graph) PostEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("Calling post school user")

	schoolID := chi.URLParam(r, "schoolID")
	schoolID, err := url.PathUnescape(schoolID)
	if err != nil {
		logger.Debug().
			Err(err).
			Str("id", schoolID).
			Msg("could not add user to school: unescaping school id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping school id failed")
		return
	}

	if schoolID == "" {
		logger.Debug().Msg("could not add school user: missing school id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing school id")
		return
	}
	memberRef := libregraph.NewMemberReference()
	err = json.NewDecoder(r.Body).Decode(memberRef)
	if err != nil {
		logger.Debug().
			Err(err).
			Interface("body", r.Body).
			Msg("could not add school user: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	memberRefURL, ok := memberRef.GetOdataIdOk()
	if !ok {
		logger.Debug().Msg("could not add school user: @odata.id reference is missing")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "@odata.id reference is missing")
		return
	}
	memberType, id, err := g.parseMemberRef(*memberRefURL)
	if err != nil {
		logger.Debug().Err(err).Msg("could not add school user: error parsing @odata.id url")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Error parsing @odata.id url")
		return
	}
	// The MS Graph spec allows "directoryObject", "user", "school" and "organizational Contact"
	// we restrict this to users for now. Might add Schools as members later
	if memberType != "users" {
		logger.Debug().Str("type", memberType).Msg("could not add school user: Only users are allowed as school members")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Only users are allowed as school members")
		return
	}

	logger.Debug().Str("memberType", memberType).Str("id", id).Msg("calling add user on backend")
	err = g.identityEducationBackend.AddUsersToEducationSchool(r.Context(), schoolID, []string{id})

	if err != nil {
		logger.Debug().Err(err).Msg("could not add school user: backend error")
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	/* TODO requires reva changes
	e := events.SchoolMemberAdded{SchoolID: schoolID, UserID: id}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)
	*/

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// DeleteEducationSchoolUser implements the Service interface.
func (g Graph) DeleteEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
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

	userID := chi.URLParam(r, "userID")
	userID, err = url.PathUnescape(userID)
	if err != nil {
		logger.Debug().Err(err).Str("id", userID).Msg("could not delete school member: unescaping member id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping member id failed")
		return
	}

	if userID == "" {
		logger.Debug().Msg("could not delete school member: missing member id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing member id")
		return
	}
	logger.Debug().Str("schoolID", schoolID).Str("userID", userID).Msg("calling delete member on backend")
	err = g.identityEducationBackend.RemoveUserFromEducationSchool(r.Context(), schoolID, userID)

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

	/* TODO requires reva changes
	e := events.SchoolMemberRemoved{SchoolID: schoolID, UserID: userID}
	if currentUser, ok := ctxpkg.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(e)
	*/

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func sortEducationSchools(req *godata.GoDataRequest, schools []*libregraph.EducationSchool) ([]*libregraph.EducationSchool, error) {
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return schools, nil
	}
	var less func(i, j int) bool

	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case displayNameAttr:
		less = func(i, j int) bool {
			return strings.ToLower(schools[i].GetDisplayName()) < strings.ToLower(schools[j].GetDisplayName())
		}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == _sortDescending {
		sort.Slice(schools, reverse(less))
	} else {
		sort.Slice(schools, less)
	}

	return schools, nil
}
