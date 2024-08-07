package service

import (
	"fmt"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
)

// AuthAppService defines the service interface.
type AuthAppService struct {
	gws pool.Selectable[gateway.GatewayAPIClient]
	m   *chi.Mux
}

// NewAuthAppService initializes a new AuthAppService.
func NewAuthAppService(opts ...Option) (*AuthAppService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	a := &AuthAppService{
		gws: o.GatewaySelector,
		m:   o.Mux,
	}

	a.m.Route("/auth-app/tokens", func(r chi.Router) {
		r.Post("/", a.HandleCreate)
	})

	return a, nil
}

// ServeHTTP implements the http.Handler interface.
func (a *AuthAppService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.m.ServeHTTP(w, r)
}

// HandleCreate handles the creation of a new auth-token
func (a *AuthAppService) HandleCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ALIVE")
}
