package svc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
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

	values := make([]*libregraph.AppRoleAssignment, 0, len(lrar.GetAssignments()))
	for _, assignment := range lrar.GetAssignments() {
		appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
		appRoleAssignment.SetId(assignment.Id)
		appRoleAssignment.SetAppRoleId(assignment.RoleId)
		appRoleAssignment.SetPrincipalType(principalTypeUser)  // currently always assigned to the user
		appRoleAssignment.SetResourceId("todo-application-id") // TODO read from config
		// appRoleAssignment.SetResourceDisplayName() // TODO read from config?
		appRoleAssignment.SetPrincipalId(assignment.AccountUuid)
		// appRoleAssignment.SetPrincipalDisplayName() // TODO fetch and cache
		values = append(values, appRoleAssignment)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: values})
}

// CreateAppRoleAssignment implements the Service interface.
func (g Graph) CreateAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling create appRoleAssignment")

	userID := chi.URLParam(r, "userID")

	ara := libregraph.NewAppRoleAssignmentWithDefaults()
	ara.SetAppRoleId("new-appRoleID")
	ara.SetPrincipalId(userID)
	ara.SetResourceId("some-application-id")

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ara)
}

// DeleteAppRoleAssignment implements the Service interface.
func (g Graph) DeleteAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("body", r.Body).Msg("calling delete appRoleAssignment")

	render.NoContent(w, r)
}
