package svc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/owncloud/ocis-graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-ldap/ldap/v3"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

// GroupCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func (g Graph) GroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupID := chi.URLParam(r, "groupID")
		if groupID == "" {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest)
			return
		}
		filter := fmt.Sprintf("(entryuuid=%s)", groupID)
		group, err := g.ldapGetSingleEntry(g.config.Ldap.BaseDNGroups, filter)
		if err != nil {
			g.logger.Info().Err(err).Msgf("Failed to read group %s", groupID)
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), groupIDKey, group)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetGroups implements the Service interface.
func (g Graph) GetGroups(w http.ResponseWriter, r *http.Request) {
	con, err := g.initLdap()
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to initialize ldap")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	result, err := g.ldapSearch(con, "(objectclass=*)", g.config.Ldap.BaseDNGroups)

	if err != nil {
		g.logger.Error().Err(err).Msg("Failed search ldap with filter: '(objectclass=*)'")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	var groups []*msgraph.Group

	for _, group := range result.Entries {
		groups = append(
			groups,
			createGroupModelFromLDAP(
				group,
			),
		)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: groups})
}

// GetGroup implements the Service interface.
func (g Graph) GetGroup(w http.ResponseWriter, r *http.Request) {
	group := r.Context().Value(groupIDKey).(*ldap.Entry)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, createGroupModelFromLDAP(group))
}
