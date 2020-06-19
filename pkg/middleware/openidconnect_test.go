package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-pkg/v2/log"
	"golang.org/x/oauth2"
)

func TestOpenIDConnectMiddleware(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "success")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	m := OpenIDConnect(
		Logger(log.NewLogger()),
		OIDCProviderFunc(func() (OIDCProvider, error) {
			return mockOP(false, mockUI(false)), nil
		}),
	)(next)

	r := httptest.NewRequest(http.MethodGet, "https://idp.example.com", nil)
	r.Header.Set("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	if r.Header.Get("x-access-token") == "" {
		t.Errorf("expected a token")
	}
}

type mockOIDCProvider struct {
	UserInfoFunc func(ctx context.Context, ts oauth2.TokenSource) (OIDCUserInfo, error)
}

// UserInfo will panic if the function has been called, but not mocked
func (m mockOIDCProvider) UserInfo(ctx context.Context, ts oauth2.TokenSource) (OIDCUserInfo, error) {
	if m.UserInfoFunc != nil {
		return m.UserInfoFunc(ctx, ts)
	}

	panic("UserInfo was called in test but not mocked")
}

func mockOP(retErr bool, ui OIDCUserInfo) OIDCProvider {
	if retErr {
		return &mockOIDCProvider{
			UserInfoFunc: func(ctx context.Context, ts oauth2.TokenSource) (OIDCUserInfo, error) {
				return nil, fmt.Errorf("error returned by mockOIDCProvider UserInfo")
			},
		}

	}
	return &mockOIDCProvider{
		UserInfoFunc: func(ctx context.Context, ts oauth2.TokenSource) (OIDCUserInfo, error) {
			return ui, nil
		},
	}

}

type mockOIDCUserInfo struct {
	ClaimsFunc func(v interface{}) error
}

// UserInfo will panic if the function has been called, but not mocked
func (m mockOIDCUserInfo) Claims(v interface{}) error {
	if m.ClaimsFunc != nil {
		return m.ClaimsFunc(v)
	}

	panic("ClaimsFunc was called in test but not mocked")
}
func mockUI(retErr bool) OIDCUserInfo {
	if retErr {
		return &mockOIDCUserInfo{
			ClaimsFunc: func(v interface{}) error {
				return fmt.Errorf("error returned by mockOIDCProvider UserInfo")
			},
		}

	}
	return &mockOIDCUserInfo{
		ClaimsFunc: func(v interface{}) error {
			// TODO fill in claims
			return nil
		},
	}

}

func mockOpenidConnectMiddlewareAccSvc(retErr bool) proto.AccountsService {
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
