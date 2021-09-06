package svc

import (
	"fmt"
	"net/http"

	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
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
			//POST /drives/{drive-id}/items/{parent-item-id}/children
			// POST /drives/marketing // creates a space called Marketing
			r.Route("/drives", func(r chi.Router) {
				r.Post("/{drive-id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					us, ok := ctxpkg.ContextGetUser(r.Context())
					if !ok {
						errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "invalid user")
						return
					}

					// do request
					// prepare ms graph response (https://docs.microsoft.com/en-us/graph/api/resources/driveitem?view=graph-rest-1.0)

					// get reva client
					client, err := svc.GetClient()
					if err != nil {
						errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
						return
					}

					// prepare createspacerequest
					csr := provider.CreateStorageSpaceRequest{
						Owner: us,
						Type:  "share",
						Name:  chi.URLParam(r, "drive-id"),
						Quota: &provider.Quota{
							QuotaMaxBytes: 65536,
							QuotaMaxFiles: 20,
						},
					}

					resp, err := client.CreateStorageSpace(r.Context(), &csr)
					if err != nil {
						errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
						return
					}

					fmt.Println(resp)
				}))
			})
		})
	})

	return svc
}
