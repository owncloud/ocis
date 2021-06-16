package svc

import (
	"net/http"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis/graph/pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/account"
	ocism "github.com/owncloud/ocis/ocis-pkg/middleware"
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
		r.Use(chim.StripSlashes)
		r.Use(ocism.ExtractAccountUUID(
			account.JWTSecret(options.Config.JWTSecret)),
		)
		r.Use(middleware.ForwardToken())
		r.Use(middleware.ExtractRelativePath())
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Get("/", svc.GetMe)
				r.Get("/drives", svc.GetDrives)
				r.Get("/drive/root/children", svc.GetPersonalDriveChildren)
			})
			r.Route("/drives", func(r chi.Router) {
				r.Get("/", svc.GetDrives)
				r.Route("/{drive-id}", func(r chi.Router) {
					r.Get("/", svc.GetDrive)
					r.Get("/root*", svc.RootRouter().ServeHTTP)
					//r.Get("/root", svc.GetDriveItem)         // /me/drive/root:/path/to/file
					//r.Get("/root/children", svc.GetChildren) // /me/drive/root:/path/to/folder:/children
					// r.Get("/items/{item-id}:{path:[^:]*}", svc.GetDriveItem)

				})
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
		})
	})

	return svc
}
