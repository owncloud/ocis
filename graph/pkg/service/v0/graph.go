package svc

import (
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/go-chi/chi"
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/cs3"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Graph defines implements the business logic for Service.
type Graph struct {
	config *config.Config
	mux    *chi.Mux
	logger *log.Logger
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetClient() (gateway.GatewayAPIClient, error) {
	return cs3.GetGatewayServiceClient(g.config.Reva.Address)
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

const userIDKey key = 0
const groupIDKey key = 1

type listResponse struct {
	Value interface{} `json:"value,omitempty"`
}
