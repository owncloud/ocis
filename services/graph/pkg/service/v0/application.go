package svc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// GetApplication implements the Service interface.
func (g Graph) GetApplication(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling get application")

	applicationID := chi.URLParam(r, "applicationID")

	role1 := libregraph.NewAppRole("uuid-for-employee-role")
	role1.SetDisplayName("Employee")
	role2 := libregraph.NewAppRole("uuid-for-managemer-role")
	role2.SetDisplayName("Manager")
	role3 := libregraph.NewAppRole("uuid-for-staff-role")
	role3.SetDisplayName("Staff")
	role4 := libregraph.NewAppRole("uuid-for-student-role")
	role4.SetDisplayName("Student")
	role5 := libregraph.NewAppRole("uuid-for-admin-role")
	role5.SetDisplayName("Administrator")
	role5.SetDescription("Can administrate all aspects of an application")
	role6 := libregraph.NewAppRole("uuid-for-guest-role")
	role6.SetDisplayName("Guest")
	role5.SetDescription("Can access shared resources, but has no personal drive")
	role7 := libregraph.NewAppRole("uuid-for-configurator-role")
	role7.SetDisplayName("Configurator")

	application := libregraph.NewApplication(applicationID)
	application.SetAppRoles([]libregraph.AppRole{
		*role1, *role2, *role3, *role4, *role5, *role6, *role7,
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: application})
}
