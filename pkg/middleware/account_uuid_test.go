package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
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
		return &mockAccountsService{
			listFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
				return nil, fmt.Errorf("error returned by mockAccountsService LIST")
			},
		}
	}

	return &mockAccountsService{
		listFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
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

type mockAccountsService struct {
	listFunc   func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (*proto.ListAccountsResponse, error)
	getFunc    func(ctx context.Context, in *proto.GetAccountRequest, opts ...client.CallOption) (*proto.Account, error)
	createFunc func(ctx context.Context, in *proto.CreateAccountRequest, opts ...client.CallOption) (*proto.Account, error)
	updateFunc func(ctx context.Context, in *proto.UpdateAccountRequest, opts ...client.CallOption) (*proto.Account, error)
	deleteFunc func(ctx context.Context, in *proto.DeleteAccountRequest, opts ...client.CallOption) (*empty.Empty, error)
}

func (m mockAccountsService) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (*proto.ListAccountsResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, in, opts...)
	}

	panic("listFunc was called in test but not mocked")
}
func (m mockAccountsService) GetAccount(ctx context.Context, in *proto.GetAccountRequest, opts ...client.CallOption) (*proto.Account, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, in, opts...)
	}

	panic("getFunc was called in test but not mocked")
}

func (m mockAccountsService) CreateAccount(ctx context.Context, in *proto.CreateAccountRequest, opts ...client.CallOption) (*proto.Account, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, in, opts...)
	}

	panic("createFunc was called in test but not mocked")
}
func (m mockAccountsService) UpdateAccount(ctx context.Context, in *proto.UpdateAccountRequest, opts ...client.CallOption) (*proto.Account, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, in, opts...)
	}

	panic("updateFunc was called in test but not mocked")
}
func (m mockAccountsService) DeleteAccount(ctx context.Context, in *proto.DeleteAccountRequest, opts ...client.CallOption) (*empty.Empty, error) {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, in, opts...)
	}

	panic("deleteFunc was called in test but not mocked")
}
