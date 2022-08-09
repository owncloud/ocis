package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/coreos/go-oidc/v3/oidc"
	. "github.com/onsi/ginkgo/v2"
	"golang.org/x/oauth2"
)

var _ = Describe("Test OIDC Authenticator", func() {
	It("should authenticate requests", func() {
		m := OIDCAuthenticator{
			ProviderFunc: func() (OIDCProvider, error) { return mockOP(false), nil },
		}

		r := httptest.NewRequest(http.MethodGet, "https://idp.example.com", nil)
		r.Header.Set("Authorization", "Bearer sometoken")

		_, ok := m.Authenticate(r)
		if ok {
			Fail("expected an internal server error")
		}
	})
})

type mockOIDCProvider struct {
	UserInfoFunc func(ctx context.Context, ts oauth2.TokenSource) (*oidc.UserInfo, error)
}

// UserInfo will panic if the function has been called, but not mocked
func (m mockOIDCProvider) UserInfo(ctx context.Context, ts oauth2.TokenSource) (*oidc.UserInfo, error) {
	if m.UserInfoFunc != nil {
		return m.UserInfoFunc(ctx, ts)
	}

	panic("UserInfo was called in test but not mocked")
}

func mockOP(retErr bool) OIDCProvider {
	if retErr {
		return &mockOIDCProvider{
			UserInfoFunc: func(ctx context.Context, ts oauth2.TokenSource) (*oidc.UserInfo, error) {
				return nil, fmt.Errorf("error returned by mockOIDCProvider UserInfo")
			},
		}

	}
	return &mockOIDCProvider{
		UserInfoFunc: func(ctx context.Context, ts oauth2.TokenSource) (*oidc.UserInfo, error) {
			ui := &oidc.UserInfo{
				// claims: private ...
			}
			return ui, nil
		},
	}

}
