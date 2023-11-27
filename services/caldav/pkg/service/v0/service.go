package svc

import (
	"context"
	"github.com/emersion/go-webdav/caldav"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/services/caldav/pkg/storage"
	"net/http"
)

type userPrincipalBackend struct{}

func (u *userPrincipalBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	// TODO: ask ocis for the user ....
	return "/principals/einstein", nil
	/*
		authCtx, ok := auth.FromContext(ctx)
		if !ok {
			panic("Invalid data in auth context!")
		}
		if authCtx == nil {
			return "", fmt.Errorf("unauthenticated requests are not supported")
		}

		userDir := base64.RawStdEncoding.EncodeToString([]byte(authCtx.UserName))
		return "/" + userDir + "/", nil
	*/
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (*chi.Mux, error) {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	// just some fake localtion for now
	storageURL := "file:///temp"
	upBackend := &userPrincipalBackend{}

	caldavBackend, _, err := storage.NewFilesystem(storageURL, "/calendar/", "/contacts/", upBackend)
	if err != nil {
		return nil, err
	}

	caldavHandler := caldav.Handler{Backend: caldavBackend}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Mount("/.well-known/caldav", &caldavHandler)
	})

	_ = chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	return m, nil
}
