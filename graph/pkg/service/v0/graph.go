package svc

import (
	"context"
	"net/http"

	"github.com/ReneKroon/ttlcache/v2"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/identity"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"google.golang.org/grpc"
)

//go:generate make generate

// GatewayClient is the subset of the gateway.GatewayAPIClient that is being used to interact with the gateway
type GatewayClient interface {
	//gateway.GatewayAPIClient

	// Returns the home path for the given authenticated user.
	// When a user has access to multiple storage providers, one of them is the home.
	GetHome(ctx context.Context, in *provider.GetHomeRequest, opts ...grpc.CallOption) (*provider.GetHomeResponse, error)
	// Returns a list of resource information
	// for the provided reference.
	// MUST return CODE_NOT_FOUND if the reference does not exists.
	ListContainer(ctx context.Context, in *provider.ListContainerRequest, opts ...grpc.CallOption) (*provider.ListContainerResponse, error)
	// Returns the resource information at the provided reference.
	// MUST return CODE_NOT_FOUND if the reference does not exist.
	Stat(ctx context.Context, in *provider.StatRequest, opts ...grpc.CallOption) (*provider.StatResponse, error)
	// Initiates the download of a file using an
	// out-of-band data transfer mechanism.
	InitiateFileDownload(ctx context.Context, in *provider.InitiateFileDownloadRequest, opts ...grpc.CallOption) (*gateway.InitiateFileDownloadResponse, error)
	// Creates a storage space.
	CreateStorageSpace(ctx context.Context, in *provider.CreateStorageSpaceRequest, opts ...grpc.CallOption) (*provider.CreateStorageSpaceResponse, error)
	// Lists storage spaces.
	ListStorageSpaces(ctx context.Context, in *provider.ListStorageSpacesRequest, opts ...grpc.CallOption) (*provider.ListStorageSpacesResponse, error)
	// Updates a storage space.
	UpdateStorageSpace(ctx context.Context, in *provider.UpdateStorageSpaceRequest, opts ...grpc.CallOption) (*provider.UpdateStorageSpaceResponse, error)
	// Deletes a storage space.
	DeleteStorageSpace(ctx context.Context, in *provider.DeleteStorageSpaceRequest, opts ...grpc.CallOption) (*provider.DeleteStorageSpaceResponse, error)
	// Returns the quota available under the provided
	// reference.
	// MUST return CODE_NOT_FOUND if the reference does not exist
	// MUST return CODE_RESOURCE_EXHAUSTED on exceeded quota limits.
	GetQuota(ctx context.Context, in *gateway.GetQuotaRequest, opts ...grpc.CallOption) (*provider.GetQuotaResponse, error)
}

// HTTPClient is the subset of the http.Client that is being used to interact with the download gateway
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetGatewayServiceClientFunc is a callback used to pass in a mock during testing
type GetGatewayServiceClientFunc func() (GatewayClient, error)

// Graph defines implements the business logic for Service.
type Graph struct {
	config               *config.Config
	mux                  *chi.Mux
	logger               *log.Logger
	identityBackend      identity.Backend
	gatewayClient        GatewayClient
	httpClient           HTTPClient
	spacePropertiesCache *ttlcache.Cache
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetGatewayClient() GatewayClient {
	return g.gatewayClient
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetHTTPClient() HTTPClient {
	return g.httpClient
}

type listResponse struct {
	Value interface{} `json:"value,omitempty"`
}

const (
	NoSpaceFoundMessage           = "space with id `%s` not found"
	ListStorageSpacesTransportErr = "transport error sending list storage spaces grpc request"
	ListStorageSpacesReturnsErr   = "list storage spaces grpc request returns an errorcode in the response"
)
