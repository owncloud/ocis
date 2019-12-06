package svc

import (
	"encoding/json"
	"fmt"
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

	m.Route("/v1.0", func(r chi.Router) {
		r.Get("/me", svc.Me)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", svc.Users)
			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", svc.Users)
			})
		})
	})

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
	userID := chi.URLParam(r, "userId")
	filter := "(objectclass=*)"
	if userID != "" {
		filter = fmt.Sprintf("(entryuuid=%s)", userID)
	}

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

	result, err := con.Search(search)

	if err != nil {
		// TODO: we should not give this error out to users
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID != "" {
		if len(result.Entries) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user := createUserModelFromLDAP(result.Entries[0])
		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
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
