package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis-graph/pkg/config"
	"github.com/owncloud/ocis-pkg/log"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	ldap "gopkg.in/ldap.v3"
)

// Graph defines implements the business logic for Service.
type Graph struct {
	config *config.Config
	mux    *chi.Mux
	logger *log.Logger
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// UserCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func (g Graph) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *ldap.Entry
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			user, err = g.ldapGetUser(userID)
		} else {
			// TODO: we should not give this error out to users
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			render.Status(r, http.StatusNotFound)
			return
		}
		if err != nil {
			g.logger.Info().Msgf("error reading user: %s", err.Error())
			render.Status(r, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetMe implements the Service interface.
func (g Graph) GetMe(w http.ResponseWriter, r *http.Request) {
	me := createUserModel(
		"Alice",
		"1234-5678-9000-000",
	)

	resp, err := json.Marshal(me)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// GetUsers implements the Service interface.
func (g Graph) GetUsers(w http.ResponseWriter, r *http.Request) {
	con, err := g.initLdap()
	if err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := g.ldapSearch(con, "(objectclass=*)")

	if err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	render.JSON(w, r, users)
}

// GetUser implements the Service interface.
func (g Graph) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userIDKey).(*ldap.Entry)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, createUserModelFromLDAP(user))
}

func (g Graph) ldapGetUser(userID string) (*ldap.Entry, error) {
	conn, err := g.initLdap()
	if err != nil {
		return nil, err
	}
	filter := fmt.Sprintf("(entryuuid=%s)", userID)
	result, err := g.ldapSearch(conn, filter)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, errors.New("user not found")
	}
	return result.Entries[0], nil
}

func (g Graph) initLdap() (*ldap.Conn, error) {
	g.logger.Info().Msg("Dailing ldap.... ")
	con, err := ldap.Dial("tcp", "localhost:10389")

	if err != nil {
		return nil, err
	}

	if err := con.Bind("cn=admin,dc=example,dc=org", "admin"); err != nil {
		return nil, err
	}
	return con, nil
}

func (g Graph) ldapSearch(con *ldap.Conn, filter string) (*ldap.SearchResult, error) {
	search := ldap.NewSearchRequest(
		"ou=users,dc=example,dc=org",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"dn",
			"uid",
			"givenname",
			"mail",
			"displayname",
			"entryuuid",
			"sn",
		},
		nil,
	)

	return con.Search(search)
}

func createUserModel(displayName string, id string) *msgraph.User {
	return &msgraph.User{
		DisplayName: &displayName,
		GivenName:   &displayName,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &id,
			},
		},
	}
}

func createUserModelFromLDAP(entry *ldap.Entry) *msgraph.User {
	displayName := entry.GetAttributeValue("displayname")
	givenName := entry.GetAttributeValue("givenname")
	mail := entry.GetAttributeValue("mail")
	surName := entry.GetAttributeValue("sn")
	id := entry.GetAttributeValue("entryuuid")
	return &msgraph.User{
		DisplayName: &displayName,
		GivenName:   &givenName,
		Surname:     &surName,
		Mail:        &mail,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &id,
			},
		},
	}
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

const userIDKey key = 0
