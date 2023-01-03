package svc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// ListAppRoleAssignments implements the Service interface.
func (g Graph) ListAppRoleAssignments(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling list appRoleAssignments")

	userID := chi.URLParam(r, "userID")

	ara1 := libregraph.NewAppRoleAssignmentWithDefaults()
	ara1.SetAppRoleId("appRoleID-1")
	ara1.SetPrincipalId(userID)
	ara1.SetResourceId("some-application-id")
	ara2 := libregraph.NewAppRoleAssignmentWithDefaults()
	ara2.SetAppRoleId("appRoleID-2")
	ara2.SetPrincipalId(userID)
	ara2.SetResourceId("some-application-id")

	values := []*libregraph.AppRoleAssignment{ara1, ara2}

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
