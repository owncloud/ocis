package svc

import (
	"context"
	"errors"
	"fmt"
	revaContext "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/emersion/go-webdav/carddav"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	ocismiddleware "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/ccs/pkg/storage"
	"net/http"
	"os"
)

type userPrincipalBackend struct{}

func (u *userPrincipalBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	user, ok := revaContext.ContextGetUser(ctx)
	if !ok {
		return "", errors.New("no user in context")
	}

	// TODO: use user.Id.OpaqueId ????
	return fmt.Sprintf("/dav/principals/users/%s/", user.Username), nil
}

type groupwareHandler struct {
	upBackend userPrincipalBackend
	// authBackend    auth.AuthProvider
	caldavBackend caldav.Backend
}

func (u *groupwareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userPrincipalPath, err := u.upBackend.CurrentUserPrincipal(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	var homeSets []webdav.BackendSuppliedHomeSet
	if u.caldavBackend != nil {
		path, err := u.caldavBackend.CalendarHomeSetPath(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			homeSets = append(homeSets, caldav.NewCalendarHomeSet(path))
		}
	}

	if r.URL.Path == userPrincipalPath {
		opts := webdav.ServePrincipalOptions{
			CurrentUserPrincipalPath: userPrincipalPath,
			HomeSets:                 homeSets,
			Capabilities: []webdav.Capability{
				carddav.CapabilityAddressBook,
				caldav.CapabilityCalendar,
			},
		}

		webdav.ServePrincipal(w, r, &opts)
		return
	}
	opts := webdav.ServePrincipalOptions{
		CurrentUserPrincipalPath: userPrincipalPath,
		HomeSets:                 homeSets,
		Capabilities: []webdav.Capability{
			carddav.CapabilityAddressBook,
			caldav.CapabilityCalendar,
		},
	}
	webdav.ServePrincipal(w, r, &opts)
	return

	// TODO serve something on / that signals this being a DAV server?

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

type wellknownHandler struct {
}

func (h *wellknownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dav/", http.StatusMovedPermanently)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (*chi.Mux, error) {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)
	m.Use(ocismiddleware.ExtractAccountUUID(
		account.JWTSecret(options.Config.JWTSecret),
	),
	)
	// just some fake location for now
	storageURL := "/tmp/caldav/"
	os.Mkdir(storageURL, 0700)
	upBackend := &userPrincipalBackend{}

	caldavBackend, _, err := storage.NewFilesystem(storageURL, "/calendar/", "/contacts/", upBackend)
	if err != nil {
		return nil, err
	}

	handler := groupwareHandler{
		upBackend:     *upBackend,
		caldavBackend: caldavBackend,
	}
	caldavHandler := caldav.Handler{Backend: caldavBackend, Prefix: "/dav"}
	wellknownHandler := wellknownHandler{}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Mount("/.well-known/caldav", &wellknownHandler)
		r.Mount("/.well-known/carddav", &wellknownHandler)
		r.Mount("/dav/", &handler)
		r.Mount("/dav/principals/users/", &handler)
		r.Mount("/dav/calendars/{user}/", &caldavHandler)
	})

	_ = chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	return m, nil
}
