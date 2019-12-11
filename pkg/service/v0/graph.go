package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/owncloud/ocis-graph/pkg/service/v0/errorcode"

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

		userID := chi.URLParam(r, "userID")
		if userID == "" {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest)
			return
		}
		user, err = g.ldapGetSingleEntry(userID, g.config.Ldap.BaseDNUsers)
		if err != nil {
			g.logger.Info().Err(err).Msgf("Failed to read user %s", userID)
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
		group, err := g.ldapGetSingleEntry(groupID, g.config.Ldap.BaseDNGroups)
		if err != nil {
			g.logger.Info().Err(err).Msgf("Failed to read group %s", groupID)
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), groupIDKey, group)
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
		g.logger.Error().Err(err).Msgf("Failed to marshal object %s", me)
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
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

func (g Graph) ldapGetSingleEntry(resourceID string, baseDn string) (*ldap.Entry, error) {
	conn, err := g.initLdap()
	if err != nil {
		return nil, err
	}
	filter := fmt.Sprintf("(entryuuid=%s)", resourceID)
	result, err := g.ldapSearch(conn, filter, baseDn)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, errors.New("resource not found")
	}
	return result.Entries[0], nil
}

func (g Graph) initLdap() (*ldap.Conn, error) {
	g.logger.Info().Msgf("Dailing ldap %s://%s", g.config.Ldap.Network, g.config.Ldap.Address)
	con, err := ldap.Dial(g.config.Ldap.Network, g.config.Ldap.Address)

	if err != nil {
		return nil, err
	}

	if err := con.Bind(g.config.Ldap.UserName, g.config.Ldap.Password); err != nil {
		return nil, err
	}
	return con, nil
}

func (g Graph) ldapSearch(con *ldap.Conn, filter string, baseDN string) (*ldap.SearchResult, error) {
	search := ldap.NewSearchRequest(
		baseDN,
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
			"cn",
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

func createGroupModelFromLDAP(entry *ldap.Entry) *msgraph.Group {
	id := entry.GetAttributeValue("entryuuid")
	displayName := entry.GetAttributeValue("cn")

	return &msgraph.Group{
		DisplayName: &displayName,
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
const groupIDKey key = 1

type listResponse struct {
	Value interface{} `json:"value,omitempty"`
}
