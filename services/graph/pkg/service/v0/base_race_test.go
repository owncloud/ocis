package svc

import (
	"context"
	"fmt"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// numRaceTestResources needs to be comfortably larger than the number of
// workers so that some workers are still reading from the driveItems map while
// the collecting loop has already started writing to it.
const numRaceTestResources = 200

// newRaceTestService builds a BaseGraphService with a mocked gateway that can
// stat any resource and resolve any user.
func newRaceTestService(t *testing.T, maxConcurrency int) BaseGraphService {
	t.Helper()

	selectorName := fmt.Sprintf("GatewaySelectorRace%d", maxConcurrency)
	pool.RemoveSelector(selectorName + "com.owncloud.api.gateway")

	gatewayClient := &cs3mocks.GatewayAPIClient{}
	gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(
		func(_ context.Context, req *storageprovider.StatRequest, _ ...grpc.CallOption) *storageprovider.StatResponse {
			return &storageprovider.StatResponse{
				Status: status.NewOK(context.Background()),
				Info: &storageprovider.ResourceInfo{
					Id:   req.GetRef().GetResourceId(),
					Type: storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
				},
			}
		},
		nil,
	)
	gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(
		&userpb.GetUserResponse{
			Status: status.NewOK(context.Background()),
			User: &userpb.User{
				Id:          &userpb.UserId{Idp: "idp", OpaqueId: "user-id"},
				DisplayName: "User Name",
			},
		},
		nil,
	)
	gatewayClient.On("GetPublicShare", mock.Anything, mock.Anything).Return(
		&link.GetPublicShareResponse{Status: status.NewOK(context.Background())},
		nil,
	)

	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		selectorName,
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)

	cfg := defaults.FullDefaultConfig()
	cfg.Identity.LDAP.CACert = ""
	cfg.TokenManager.JWTSecret = "loremipsum"
	cfg.MaxConcurrency = maxConcurrency

	logger := log.NewLogger()

	return BaseGraphService{
		logger:          &logger,
		gatewaySelector: gatewaySelector,
		identityCache:   identity.NewIdentityCache(identity.IdentityCacheWithGatewaySelector(gatewaySelector)),
		config:          cfg,
	}
}

// raceTestShares returns one user share per resource, all distinct resources so
// that every share becomes its own unit of work.
func raceTestShares() []*collaboration.Share {
	editorResourcePermissions := conversions.NewEditorRole().CS3ResourcePermissions()

	shares := make([]*collaboration.Share, 0, numRaceTestResources)
	for i := 0; i < numRaceTestResources; i++ {
		shares = append(shares, &collaboration.Share{
			Id: &collaboration.ShareId{OpaqueId: fmt.Sprintf("share-id-%d", i)},
			ResourceId: &storageprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  fmt.Sprintf("opaqueid-%d", i),
			},
			Grantee: &storageprovider.Grantee{
				Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
				Id: &storageprovider.Grantee_UserId{
					UserId: &userpb.UserId{OpaqueId: "user-id"},
				},
			},
			Permissions: &collaboration.SharePermissions{Permissions: editorResourcePermissions},
		})
	}
	return shares
}

// TestCS3UserSharesToDriveItemsConcurrentMapAccess reproduces OCISDEV-1088: the
// errgroup workers read the driveItems map while the collecting loop writes to
// it, which makes the Go runtime abort the process with
// "fatal error: concurrent map read and map write".
//
// Run with -race to catch it deterministically.
func TestCS3UserSharesToDriveItemsConcurrentMapAccess(t *testing.T) {
	// A non-empty map is required: the workers only read the map when looking
	// for an already known drive item, and an empty map read is not enough for
	// the race detector to flag a conflicting access pattern on every run.
	driveItems := make(driveItemsByResourceID, numRaceTestResources)
	for i := 0; i < numRaceTestResources; i++ {
		driveItems[fmt.Sprintf("storageid$spaceid!seed-%d", i)] = libregraph.DriveItem{}
	}

	g := newRaceTestService(t, defaults.FullDefaultConfig().MaxConcurrency)

	items, err := g.cs3UserSharesToDriveItems(context.Background(), raceTestShares(), driveItems)
	require.NoError(t, err)
	require.Len(t, items, numRaceTestResources+numRaceTestResources)
}

// TestCS3UserSharesToDriveItemsConcurrentMapAccessMaxConcurrencyOne covers the
// single worker path.
//
// This case exists because lowering OCIS_MAX_CONCURRENCY /
// GRAPH_MAX_CONCURRENCY to 1 looks like an obvious operational workaround for
// OCISDEV-1088, but it never was one: even with a single worker the map was
// still read concurrently with the collecting loop, because the results channel
// is buffered and the worker never blocks handing off an item. Before the fix
// this test failed under -race exactly like the one above.
func TestCS3UserSharesToDriveItemsConcurrentMapAccessMaxConcurrencyOne(t *testing.T) {
	driveItems := make(driveItemsByResourceID, numRaceTestResources)
	for i := 0; i < numRaceTestResources; i++ {
		driveItems[fmt.Sprintf("storageid$spaceid!seed-%d", i)] = libregraph.DriveItem{}
	}

	g := newRaceTestService(t, 1)

	items, err := g.cs3UserSharesToDriveItems(context.Background(), raceTestShares(), driveItems)
	require.NoError(t, err)
	require.Len(t, items, numRaceTestResources+numRaceTestResources)
}
