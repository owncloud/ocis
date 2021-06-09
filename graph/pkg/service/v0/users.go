package svc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-ldap/ldap/v3"

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
		// TODO make filter configurable
		filter := fmt.Sprintf("(&(objectClass=posixAccount)(ownCloudUUID=%s))", userID)
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

	u := user.ContextMustGetUser(r.Context())

	me := createUserModelFromCS3User(u)

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

	// TODO make filter configurable
	result, err := g.ldapSearch(con, "(objectClass=posixAccount)", g.config.Ldap.BaseDNUsers)

	if err != nil {
		g.logger.Error().Err(err).Msg("Failed search ldap with filter: '(objectClass=posixAccount)'")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	users := make([]*msgraph.User, 0, len(result.Entries))

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

func createUserModelFromCS3User(u *userpb.User) *msgraph.User {
	return &msgraph.User{
		DisplayName:   &u.DisplayName,
		Mail:          &u.Mail,
		PreferredName: &u.Username,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &u.Id.OpaqueId,
			},
		},
	}
}
