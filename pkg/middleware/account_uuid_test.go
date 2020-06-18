package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/oidc"
)

// TODO testing the getAccount method should inject a cache
func TestGetAccountSuccess(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "success")
	if _, status := getAccount(log.NewLogger(), oidc.StandardClaims{Email: "success"}, mockAccSvc(false)); status != 0 {
		t.Errorf("expected an account")
	}
}
func TestGetAccountInternalError(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "failure")
	if _, status := getAccount(log.NewLogger(), oidc.StandardClaims{Email: "failure"}, mockAccSvc(true)); status != http.StatusInternalServerError {
		t.Errorf("expected an internal server error")
	}
}

func mockAccSvc(retErr bool) proto.AccountsService {
	if retErr {
		return &proto.MockAccountsService{
			ListFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
				return nil, fmt.Errorf("error returned by mockAccountsService LIST")
			},
		}
	}

	return &proto.MockAccountsService{
		ListFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
			return &proto.ListAccountsResponse{
				Accounts: []*proto.Account{
					{
						Id: "yay",
					},
				},
			}, nil
		},
	}

}
