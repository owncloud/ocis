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

// GetApplication implements the Service interface.
func (g Graph) GetApplication(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling get application")

	// TODO make application id and name for this instance configurable
	applicationID := chi.URLParam(r, "applicationID")

	s := settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient())

	lbr, err := s.ListRoles(r.Context(), &settingssvc.ListBundlesRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("could not list roles: transport error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	roles := make([]libregraph.AppRole, 0, len(lbr.Bundles))
	for _, bundle := range lbr.GetBundles() {
		role := libregraph.NewAppRole(bundle.GetId())
		role.SetDisplayName(bundle.GetDisplayName())
		roles = append(roles, *role)
	}

	application := libregraph.NewApplication(applicationID)
	application.SetAppRoles(roles)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: application})
}
