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
	"github.com/owncloud/ocis-pkg/v2/oidc"
	"github.com/owncloud/ocis-proxy/pkg/config"
)

// TODO testing the getAccount method should inject a cache
func TestGetAccountSuccess(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "success")
	if _, status := getAccount(log.NewLogger(), mockAccountUUIDMiddlewareAccSvc(false, true), "mail eq 'success'"); status != 0 {
		t.Errorf("expected an account")
	}
}
func TestGetAccountInternalError(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "failure")
	if _, status := getAccount(log.NewLogger(), mockAccountUUIDMiddlewareAccSvc(true, false), "mail eq 'failure'"); status != http.StatusInternalServerError {
		t.Errorf("expected an internal server error")
	}
}

func TestAccountUUIDMiddleware(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "success")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := AccountUUID(
		Logger(log.NewLogger()),
		TokenManagerConfig(config.TokenManager{JWTSecret: "secret"}),
		AccountsClient(mockAccountUUIDMiddlewareAccSvc(false, true)),
	)(next)

	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
	w := httptest.NewRecorder()
	ctx := oidc.NewContext(r.Context(), &oidc.StandardClaims{Email: "success"})
	r = r.WithContext(ctx)
	m.ServeHTTP(w, r)

	if r.Header.Get("x-access-token") == "" {
		t.Errorf("expected a token")
	}
}

func TestAccountUUIDMiddlewareWithDisabledAccount(t *testing.T) {
	svcCache.Invalidate(AccountsKey, "failure")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := AccountUUID(
		Logger(log.NewLogger()),
		TokenManagerConfig(config.TokenManager{JWTSecret: "secret"}),
		AccountsClient(mockAccountUUIDMiddlewareAccSvc(false, false)),
	)(next)

	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)
	w := httptest.NewRecorder()
	ctx := oidc.NewContext(r.Context(), &oidc.StandardClaims{Email: "failure"})
	r = r.WithContext(ctx)
	m.ServeHTTP(w, r)

	rsp := w.Result()
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected a disabled account to be unauthorized, got: %d", rsp.StatusCode)
	}
}

func mockAccountUUIDMiddlewareAccSvc(retErr, accEnabled bool) proto.AccountsService {
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
						Id:             "yay",
						AccountEnabled: accEnabled,
					},
				},
			}, nil
		},
	}

}
