package svc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

const principalTypeUser = "User"

// ListAppRoleAssignments implements the Service interface.
func (g Graph) ListAppRoleAssignments(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling list appRoleAssignments")

	userID := chi.URLParam(r, "userID")

	lrar, err := g.roleService.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})
	if err != nil {
		// TODO check the error type and return proper error code
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	values := make([]libregraph.AppRoleAssignment, 0, len(lrar.GetAssignments()))
	for _, assignment := range lrar.GetAssignments() {
		values = append(values, g.assignmentToAppRoleAssignment(assignment))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: values})
}

// CreateAppRoleAssignment implements the Service interface.
func (g Graph) CreateAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling create appRoleAssignment")

	appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
	err := json.NewDecoder(r.Body).Decode(appRoleAssignment)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err.Error()))
		return
	}

	userID := chi.URLParam(r, "userID")

	if appRoleAssignment.GetPrincipalId() != userID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("user id %s does not match principal id %v", userID, appRoleAssignment.GetPrincipalId()))
		return
	}
	if appRoleAssignment.GetResourceId() != g.config.Service.ApplicationID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("resource id %s does not match expected application id %v", userID, g.config.Service.ApplicationID))
		return
	}

	artur, err := g.roleService.AssignRoleToUser(r.Context(), &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      appRoleAssignment.AppRoleId,
	})
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, g.assignmentToAppRoleAssignment(artur.GetAssignment()))
}

// DeleteAppRoleAssignment implements the Service interface.
func (g Graph) DeleteAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("body", r.Body).Msg("calling delete appRoleAssignment")

	userID := chi.URLParam(r, "userID")

	// check assignment belongs to the user
	lrar, err := g.roleService.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	appRoleAssignmentID := chi.URLParam(r, "appRoleAssignmentID")

	assignmentFound := false
	for _, roleAssignment := range lrar.GetAssignments() {
		if roleAssignment.Id == appRoleAssignmentID {
			assignmentFound = true
		}
	}
	if !assignmentFound {
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, fmt.Sprintf("appRoleAssignment %v not found for user %v", appRoleAssignmentID, userID))
		return
	}

	_, err = g.roleService.RemoveRoleFromUser(r.Context(), &settingssvc.RemoveRoleFromUserRequest{
		Id: appRoleAssignmentID,
	})
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.NoContent(w, r)
}

func (g Graph) assignmentToAppRoleAssignment(assignment *settingsmsg.UserRoleAssignment) libregraph.AppRoleAssignment {
	appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
	appRoleAssignment.SetId(assignment.Id)
	appRoleAssignment.SetAppRoleId(assignment.RoleId)
	appRoleAssignment.SetPrincipalType(principalTypeUser) // currently always assigned to the user
	appRoleAssignment.SetResourceId(g.config.Service.ApplicationID)
	appRoleAssignment.SetResourceDisplayName(g.config.Service.ApplicationDisplayName)
	appRoleAssignment.SetPrincipalId(assignment.AccountUuid)
	// appRoleAssignment.SetPrincipalDisplayName() // TODO fetch and cache
	return *appRoleAssignment
}
