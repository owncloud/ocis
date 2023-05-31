package utils

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"google.golang.org/grpc/metadata"
)

// Impersonate returns an authenticated reva context and the user it represents
func Impersonate(userID *user.UserId, selector pool.Selectable[gateway.GatewayAPIClient], machineAuthAPIKey string) (context.Context, *user.User, error) {
	usr, err := GetUser(userID, selector, machineAuthAPIKey)
	if err != nil {
		return nil, nil, err
	}

	ctx, err := ImpersonateUser(usr, selector, machineAuthAPIKey)
	return ctx, usr, err
}

// GetUser gets the specified user
func GetUser(userID *user.UserId, selector pool.Selectable[gateway.GatewayAPIClient], machineAuthAPIKey string) (*user.User, error) {
	gwc, err := selector.Next()
	if err != nil {
		return nil, err
	}
	getUserResponse, err := gwc.GetUser(context.Background(), &user.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if getUserResponse.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	return getUserResponse.GetUser(), nil
}

// ImpersonateUser impersonates the given user
func ImpersonateUser(usr *user.User, selector pool.Selectable[gateway.GatewayAPIClient], machineAuthAPIKey string) (context.Context, error) {
	ctx := revactx.ContextSetUser(context.Background(), usr)
	gwc, err := selector.Next()
	if err != nil {
		return nil, err
	}
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
