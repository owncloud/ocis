package svc

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis-graph/pkg/config"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	ldap "gopkg.in/ldap.v3"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Me(http.ResponseWriter, *http.Request)
	Users(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Graph{
		config: options.Config,
		mux:    m,
	}

	m.HandleFunc("/v1.0/me", svc.Me)
	m.HandleFunc("/v1.0/users", svc.Users)

	return svc
}

// Graph defines implements the business logic for Service.
type Graph struct {
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Me implements the Service interface.
func (g Graph) Me(w http.ResponseWriter, r *http.Request) {
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

// Users implements the Service interface.
func (g Graph) Users(w http.ResponseWriter, r *http.Request) {
	con, err := ldap.Dial("tcp", "localhost:10389")

	if err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := con.Bind("cn=admin,dc=example,dc=org", "admin"); err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	search := ldap.NewSearchRequest(
		"ou=groups,dc=example,dc=org",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectclass=*)",
		[]string{
			"dn",
			"uuid",
			"uid",
			"givenName",
			"mail",
		},
		nil,
	)

	result, err := con.Search(search)

	if err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := make([]*msgraph.User, len(result.Entries))

	for _, user := range result.Entries {
		users = append(
			users,
			createUserModel(
				user.DN,
				"1234-5678-9000-000",
			),
		)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
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
