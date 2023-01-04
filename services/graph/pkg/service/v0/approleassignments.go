package svc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
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

	s := settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient())

	lrar, err := s.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("could not list role assginments: transport error")
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

	userID := chi.URLParam(r, "userID")

	appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
	err := json.NewDecoder(r.Body).Decode(appRoleAssignment)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not create user: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err.Error()))
		return
	}

	if appRoleAssignment.GetPrincipalId() != userID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("user id %s does not match principal id %v", userID, appRoleAssignment.GetPrincipalId()))
		return
	}
	if appRoleAssignment.GetResourceId() != g.config.Service.ApplicationID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("resource id %s does not match expected application id %v", userID, g.config.Service.ApplicationID))
		return
	}

	s := settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient())

	artur, err := s.AssignRoleToUser(r.Context(), &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      appRoleAssignment.AppRoleId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("could not assign role to user")
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
