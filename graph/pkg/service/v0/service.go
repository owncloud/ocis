package svc

import (
	"net/http"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/account"
	opkgm "github.com/owncloud/ocis/ocis-pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetMe(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Graph{
		config: options.Config,
		mux:    m,
		logger: &options.Logger,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Get("/", svc.GetMe)
				r.Get("/drives", svc.GetDrives)
				r.Get("/drive/root/children", svc.GetRootDriveChildren)
			})
			r.Route("/users", func(r chi.Router) {
				r.Get("/", svc.GetUsers)
				r.Route("/{userID}", func(r chi.Router) {
					r.Use(svc.UserCtx)
					r.Get("/", svc.GetUser)
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Get("/", svc.GetGroups)
				r.Route("/{groupID}", func(r chi.Router) {
					r.Use(svc.GroupCtx)
					r.Get("/", svc.GetGroup)
				})
			})
			r.Route("/drives", func(r chi.Router) {
				r.Use(opkgm.ExtractAccountUUID(
					account.Logger(options.Logger),
					account.JWTSecret(options.Config.TokenManager.JWTSecret)),
				)
				r.Post("/", svc.CreateDrive)
			})
			r.Route("/Drive({id})", func(r chi.Router) {
				r.Use(opkgm.ExtractAccountUUID(
					account.Logger(options.Logger),
					account.JWTSecret(options.Config.TokenManager.JWTSecret)),
				)
				r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
					d := strings.ReplaceAll(chi.URLParam(r, "id"), `"`, "")
					/*
						1. get storage space by id
						2. prepare UpdateStorageSpaceRequest
						3. get a reva client
						4. do UpdateStorageSpace request

						Known loose ends:
						1. Reva's FS interface does not yet contain UpdateStorageSpace. Needs to be added.
						2. There are many ways to select a resource on OData. Because spaces names are not unique, we will support
						unique updates and not bulk updates. Supported URLs look like:

						https://localhost:9200/graph/v1.0/DriveById(id=1284d238-aa92-42ce-bdc4-0b0000009157!cdf8d353-dd02-46ed-b06a-3bb66f29743c)

						3. How are uploading images to the space being handled? Since an image is not a property of the Drive (speaking OData)
						it can be handled directly by doing an upload to the storage itself.
						4. Ditto for descriptions. We want to persist a space's description on a file inside the `.space` reserved folder.
					*/
					_, _ = w.Write([]byte(d))
					w.WriteHeader(http.StatusOK)
				})
			})
		})
	})

	return svc
}
