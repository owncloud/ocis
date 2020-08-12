package svc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/owncloud/ocis-graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-ldap/ldap/v3"
	"github.com/owncloud/ocis-pkg/v2/oidc"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

// UserCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func (g Graph) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *ldap.Entry
		var err error

		userID := chi.URLParam(r, "userID")
		if userID == "" {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest)
			return
		}
		filter := fmt.Sprintf("(entryuuid=%s)", userID)
		user, err = g.ldapGetSingleEntry(g.config.Ldap.BaseDNUsers, filter)
		if err != nil {
			g.logger.Info().Err(err).Msgf("Failed to read user %s", userID)
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetMe implements the Service interface.
func (g Graph) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := oidc.FromContext(r.Context())
	g.logger.Info().Interface("Claims", claims).Msg("Claims in /me")

	filter := fmt.Sprintf("(uid=%s)", claims.PreferredUsername)
	user, err := g.ldapGetSingleEntry(g.config.Ldap.BaseDNUsers, filter)
	if err != nil {
		g.logger.Info().Err(err).Msgf("Failed to read user %s", claims.PreferredUsername)
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound)
		return
	}

	me := createUserModelFromLDAP(user)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, me)
}

// GetUsers implements the Service interface.
func (g Graph) GetUsers(w http.ResponseWriter, r *http.Request) {
	con, err := g.initLdap()
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to initialize ldap")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	result, err := g.ldapSearch(con, "(objectclass=*)", g.config.Ldap.BaseDNUsers)

	if err != nil {
		g.logger.Error().Err(err).Msg("Failed search ldap with filter: '(objectclass=*)'")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	var users []*msgraph.User

	for _, user := range result.Entries {
		users = append(
			users,
			createUserModelFromLDAP(
				user,
			),
		)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: users})
}

// GetUser implements the Service interface.
func (g Graph) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userIDKey).(*ldap.Entry)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, createUserModelFromLDAP(user))
}
