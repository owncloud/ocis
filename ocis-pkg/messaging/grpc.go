package messaging

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"google.golang.org/grpc/metadata"
)

func authenticate(usr *user.User, gwc gateway.GatewayAPIClient, machineAuthAPIKey string) (context.Context, error) {
	ctx := revactx.ContextSetUser(context.Background(), usr)
	authRes, err := gwc.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + usr.GetId().GetOpaqueId(),
		ClientSecret: machineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error impersonating user: %s", authRes.Status.Message)
	}

	return metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, authRes.Token), nil
}

func getSpace(ctx context.Context, spaceID string, gwc gateway.GatewayAPIClient) (*storageprovider.StorageSpace, error) {
	res, err := gwc.ListStorageSpaces(ctx, listStorageSpaceRequest(spaceID))
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("Error while getting space: (%v) %s", res.GetStatus().GetCode(), res.GetStatus().GetMessage())
	}

	if len(res.StorageSpaces) == 0 {
		return nil, fmt.Errorf("error getting storage space %s: no space returned", spaceID)
	}

	return res.StorageSpaces[0], nil
}

func getUser(ctx context.Context, userid *user.UserId, gwc gateway.GatewayAPIClient) (*user.User, error) {
	getUserResponse, err := gwc.GetUser(context.Background(), &user.GetUserRequest{
		UserId: userid,
	})
	if err != nil {
		return nil, err
	}

	if getUserResponse.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	return getUserResponse.GetUser(), nil
}

func getGroup(ctx context.Context, groupid string, gwc gateway.GatewayAPIClient) (*group.Group, error) {
	r, err := gwc.GetGroup(ctx, &group.GetGroupRequest{GroupId: &group.GroupId{OpaqueId: groupid}})
	if err != nil {
		return nil, err
	}

	if r.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
	}

	return r.GetGroup(), nil
}

func getResource(ctx context.Context, resourceid *storageprovider.ResourceId, gwc gateway.GatewayAPIClient) (*storageprovider.ResourceInfo, error) {
	res, err := gwc.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: resourceid}})
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("Unexpected status code while getting space: %v", res.GetStatus().GetCode())
	}

	return res.GetInfo(), nil
}

func listStorageSpaceRequest(spaceID string) *storageprovider.ListStorageSpacesRequest {
	return &storageprovider.ListStorageSpacesRequest{
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: spaceID,
					},
				},
			},
		},
	}
}
