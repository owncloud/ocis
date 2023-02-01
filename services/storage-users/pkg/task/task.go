package task

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
)

func getToken(gw gateway.GatewayAPIClient, clientSecret, userID string) (string, error) {
	ctx := ctxpkg.ContextSetUser(context.Background(), &user.User{
		Id: &user.UserId{
			OpaqueId: userID,
		},
	})
	authenticateResponse, err := gw.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + userID,
		ClientSecret: clientSecret,
	})

	if err != nil {
		return "", err
	}

	if authenticateResponse.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return "", errtypes.NewErrtypeFromStatus(authenticateResponse.Status)
	}

	return authenticateResponse.Token, nil
}
