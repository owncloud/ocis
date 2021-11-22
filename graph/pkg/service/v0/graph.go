package svc

import (
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/identity"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Graph defines implements the business logic for Service.
type Graph struct {
	config          *config.Config
	mux             *chi.Mux
	logger          *log.Logger
	identityBackend identity.Backend
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetClient() (gateway.GatewayAPIClient, error) {
	return pool.GetGatewayServiceClient(g.config.Reva.Address)
}

type listResponse struct {
	Value interface{} `json:"value,omitempty"`
}
