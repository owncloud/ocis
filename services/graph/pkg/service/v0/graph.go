package svc

import (
	"context"
	"net/http"
	"net/url"
	"path"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/jellydator/ttlcache/v2"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"go-micro.dev/v4/client"
	mevents "go-micro.dev/v4/events"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate make -C ../../.. generate

// GatewayClient is the subset of the gateway.GatewayAPIClient that is being used to interact with the gateway
type GatewayClient interface {
	//gateway.GatewayAPIClient

	// Authenticates a user.
	Authenticate(ctx context.Context, in *gateway.AuthenticateRequest, opts ...grpc.CallOption) (*gateway.AuthenticateResponse, error)
	// Returns the home path for the given authenticated user.
	// When a user has access to multiple storage providers, one of them is the home.
	GetHome(ctx context.Context, in *provider.GetHomeRequest, opts ...grpc.CallOption) (*provider.GetHomeResponse, error)
	// GetPath does a path lookup for a resource by ID
	GetPath(ctx context.Context, in *provider.GetPathRequest, opts ...grpc.CallOption) (*provider.GetPathResponse, error)
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
	SetArbitraryMetadata(ctx context.Context, request *provider.SetArbitraryMetadataRequest, opts ...grpc.CallOption) (*provider.SetArbitraryMetadataResponse, error)
}

// Publisher is the interface for events publisher
type Publisher interface {
	Publish(string, interface{}, ...mevents.PublishOption) error
}

// Permissions is the interface used to access the permissions service
type Permissions interface {
	GetPermissionByID(ctx context.Context, request *settingssvc.GetPermissionByIDRequest, opts ...client.CallOption) (*settingssvc.GetPermissionByIDResponse, error)
	ListPermissionsByResource(ctx context.Context, in *settingssvc.ListPermissionsByResourceRequest, opts ...client.CallOption) (*settingssvc.ListPermissionsByResourceResponse, error)
}

// HTTPClient is the subset of the http.Client that is being used to interact with the download gateway
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetGatewayServiceClientFunc is a callback used to pass in a mock during testing
type GetGatewayServiceClientFunc func() (GatewayClient, error)

// RoleService is the interface used to access the role service
type RoleService interface {
	ListRoles(ctx context.Context, in *settingssvc.ListBundlesRequest, opts ...client.CallOption) (*settingssvc.ListBundlesResponse, error)
	ListRoleAssignments(ctx context.Context, in *settingssvc.ListRoleAssignmentsRequest, opts ...client.CallOption) (*settingssvc.ListRoleAssignmentsResponse, error)
	AssignRoleToUser(ctx context.Context, in *settingssvc.AssignRoleToUserRequest, opts ...client.CallOption) (*settingssvc.AssignRoleToUserResponse, error)
	RemoveRoleFromUser(ctx context.Context, in *settingssvc.RemoveRoleFromUserRequest, opts ...client.CallOption) (*emptypb.Empty, error)
}

// Graph defines implements the business logic for Service.
type Graph struct {
	config               *config.Config
	mux                  *chi.Mux
	logger               *log.Logger
	identityBackend      identity.Backend
	gatewayClient        GatewayClient
	roleService          RoleService
	permissionsService   Permissions
	spacePropertiesCache *ttlcache.Cache
	eventsPublisher      events.Publisher
	searchService        searchsvc.SearchProviderService
}

// ServeHTTP implements the Service interface.
func (g Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// GetClient returns a gateway client to talk to reva
func (g Graph) GetGatewayClient() GatewayClient {
	return g.gatewayClient
}

func (g Graph) publishEvent(ev interface{}) {
	if g.eventsPublisher != nil {
		if err := events.Publish(g.eventsPublisher, ev); err != nil {
			g.logger.Error().
				Err(err).
				Msg("could not publish user created event")
		}
	}
}

func (g Graph) getWebDavBaseURL() (*url.URL, error) {
	webDavBaseURL, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		return nil, err
	}
	webDavBaseURL.Path = path.Join(webDavBaseURL.Path, g.config.Spaces.WebDavPath)
	return webDavBaseURL, nil
}

type ListResponse struct {
	Value interface{} `json:"value,omitempty"`
}

const (
	ReadmeSpecialFolderName     = "readme"
	SpaceImageSpecialFolderName = "image"
)
