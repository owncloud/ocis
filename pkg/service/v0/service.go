package svc

import (
	"net/http"

	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	//"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/owncloud/ocis-ocs/pkg/config"
	ocsm "github.com/owncloud/ocis-ocs/pkg/middleware"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis-pkg/v2/log"
	storepb "github.com/owncloud/ocis-store/pkg/proto/v0"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetConfig(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Ocs{
		config: options.Config,
		mux:    m,
		logger: options.Logger,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.NotFound(svc.NotFound)
		r.Use(middleware.StripSlashes)
		r.Use(ocsm.OCSFormatCtx) // updates request Accept header according to format=(json|xml) query parameter
		r.Route("/v{version:(1|2)}.php", func(r chi.Router) {
			r.Use(svc.VersionCtx) // stores version in context
			r.Route("/apps/files_sharing/api/v1", func(r chi.Router) {})
			r.Route("/apps/notifications/api/v1", func(r chi.Router) {})
			r.Route("/cloud", func(r chi.Router) {
				r.Route("/capabilities", func(r chi.Router) {})
				r.Route("/user", func(r chi.Router) {
					r.Get("/", svc.GetUser)
					r.Get("/signing-key", svc.GetSigningKey)
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", svc.ListUsers)
				})
			})
			r.Route("/config", func(r chi.Router) {
				r.Get("/", svc.GetConfig)
			})
		})
	})

	return svc
}

// Ocs defines implements the business logic for Service.
type Ocs struct {
	config *config.Config
	logger log.Logger
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (o Ocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.mux.ServeHTTP(w, r)
}

// NotFound uses ErrRender to always return a proper OCS payload
func (o Ocs) NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, ErrRender(MetaUnknownError.StatusCode, "please check the syntax. API specifications are here: http://www.freedesktop.org/wiki/Specifications/open-collaboration-services"))
}

// GetUser returns the currently logged in user
func (o Ocs) GetUser(w http.ResponseWriter, r *http.Request) {

	// TODO move token marshaling to ocis-proxy
	u, ok := user.ContextGetUser(r.Context())
	if !ok {
		render.Render(w, r, ErrRender(MetaBadRequest.StatusCode, "missing user in context"))
		return
	}

	render.Render(w, r, DataRender(&data.User{
		ID:          u.Username, // TODO userid vs username! implications for clients if we return the userid here? -> implement graph ASAP?
		DisplayName: u.DisplayName,
		Email:       u.Mail,
	}))
}

// GetSigningKey returns the signing key for the current user. It will create it on the fly if it does not exist
// The signing key is part of the user settings and is used by the proxy to authenticate requests
// TODO middleware for the proxy
func (o Ocs) GetSigningKey(w http.ResponseWriter, r *http.Request) {

	// TODO move token marshaling to ocis-proxy
	_, ok := user.ContextGetUser(r.Context())
	if !ok {
		//	render.Render(w, r, ErrRender(MetaBadRequest.StatusCode, "missing user in context"))
		//	return
	}
	c := storepb.NewStoreService("com.owncloud.api.store", grpc.NewClient())
	res, err := c.Read(r.Context(), &storepb.ReadRequest{
		Key: "TODO replace with user from ctx",
		Options: &storepb.ReadOptions{
			Database: "ocs",
			Table:    "signing-keys",
		}})
	if err != nil {
		// TODO check return code, if 404 / not found error continue and try to create it
		o.logger.Error().Err(err).Msg("error reading key")
	}
	o.logger.Info().Interface("release", res).Msg("read key")
	// TODO check if signing key empty
	signingKey := string(res.Records[0].Value)
	/* TODO create key if it is missing
	if ($signingKey === null) {
			$signingKey = \OC::$server->getSecureRandom()->generate(64);
			\OC::$server->getConfig()->setUserValue($userId, 'core', 'signing-key', $signingKey, null);
	}
	*/

	render.Render(w, r, DataRender(&data.SigningKey{
		//	User:       u.Username, // TODO userid vs username?
		SigningKey: signingKey,
	}))
}

// ListUsers lists the users
func (o Ocs) ListUsers(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, ErrRender(MetaUnknownError.StatusCode, "please check the syntax. API specifications are here: http://www.freedesktop.org/wiki/Specifications/open-collaboration-services"))
}
