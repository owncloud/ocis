package policy

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

func TestStaticSelector(t *testing.T) {
	ctx := context.Background()
	req := httptest.NewRequest("GET", "https://example.org/foo", nil)
	sel := NewStaticSelector(&config.StaticSelectorConf{Policy: "ocis"})

	want := "ocis"
	got, err := sel(ctx, req)
	if got != want {
		t.Errorf("Expected policy %v got %v", want, got)
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	sel = NewStaticSelector(&config.StaticSelectorConf{Policy: "foo"})

	want = "foo"
	got, err = sel(ctx, req)
	if got != want {
		t.Errorf("Expected policy %v got %v", want, got)
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

type testCase struct {
	AccSvcShouldReturnError bool
	Claims                  *oidc.StandardClaims
	Expected                string
}

func TestMigrationSelector(t *testing.T) {
	cfg := config.MigrationSelectorConf{
		AccFoundPolicy:        "found",
		AccNotFoundPolicy:     "not_found",
		UnauthenticatedPolicy: "unauth",
	}
	var tests = []testCase{
		{true, &oidc.StandardClaims{PreferredUsername: "Hans"}, "not_found"},
		{false, &oidc.StandardClaims{PreferredUsername: "Hans"}, "found"},
		{false, nil, "unauth"},
	}

	for _, tc := range tests {
		//t.Run(fmt.Sprintf("#%v", k), func(t *testing.T) {
		//	t.Parallel()
		tc := tc
		sut := NewMigrationSelector(&cfg, mockAccSvc(tc.AccSvcShouldReturnError))
		r := httptest.NewRequest("GET", "https://example.com", nil)
		ctx := oidc.NewContext(r.Context(), tc.Claims)
		nr := r.WithContext(ctx)

		got, err := sut(ctx, nr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got != tc.Expected {
			t.Errorf("Expected Policy %v got %v", tc.Expected, got)
		}
		//})
	}
}

func mockAccSvc(retErr bool) proto.AccountsService {
	if retErr {
		return &proto.MockAccountsService{
			GetFunc: func(ctx context.Context, in *proto.GetAccountRequest, opts ...client.CallOption) (record *proto.Account, err error) {
				return nil, fmt.Errorf("error returned by mockAccountsService GET")
			},
		}
	}

	return &proto.MockAccountsService{
		GetFunc: func(ctx context.Context, in *proto.GetAccountRequest, opts ...client.CallOption) (record *proto.Account, err error) {
			return &proto.Account{}, nil
		},
	}

}
