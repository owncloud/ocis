package utils

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"google.golang.org/grpc/metadata"
)

// GetUser gets the specified user
func GetUser(userID *user.UserId, gwc gateway.GatewayAPIClient) (*user.User, error) {
	getUserResponse, err := gwc.GetUser(context.Background(), &user.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if getUserResponse.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	return getUserResponse.GetUser(), nil
}

// GetServiceUserContext returns an authenticated context of the given service user
func GetServiceUserContext(serviceUserID string, gwc gateway.GatewayAPIClient, serviceUserSecret string) (context.Context, error) {
	ctx := context.Background()
	authRes, err := gwc.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "serviceaccounts",
		ClientId:     serviceUserID,
		ClientSecret: serviceUserSecret,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error authenticating service user: %s", authRes.Status.Message)
	}

	return metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, authRes.Token), nil
}
