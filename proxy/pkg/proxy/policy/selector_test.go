package policy

import (
	"context"
	"fmt"
	"testing"

	"github.com/asim/go-micro/v3/client"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revauser "github.com/cs3org/reva/pkg/user"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

func TestLoadSelector(t *testing.T) {
	type test struct {
		cfg         *config.PolicySelector
		expectedErr error
	}
	sCfg := &config.StaticSelectorConf{Policy: "reva"}
	mcfg := &config.MigrationSelectorConf{
		AccFoundPolicy:        "found",
		AccNotFoundPolicy:     "not_found",
		UnauthenticatedPolicy: "unauth",
	}
	ccfg := &config.ClaimsSelectorConf{}
	rcfg := &config.RegexSelectorConf{}

	table := []test{
		{cfg: &config.PolicySelector{Static: sCfg, Migration: mcfg}, expectedErr: ErrMultipleSelectors},
		{cfg: &config.PolicySelector{Static: sCfg, Claims: ccfg, Regex: rcfg}, expectedErr: ErrMultipleSelectors},
		{cfg: &config.PolicySelector{}, expectedErr: ErrSelectorConfigIncomplete},
		{cfg: &config.PolicySelector{Static: sCfg}, expectedErr: nil},
		{cfg: &config.PolicySelector{Migration: mcfg}, expectedErr: nil},
		{cfg: &config.PolicySelector{Claims: ccfg}, expectedErr: nil},
		{cfg: &config.PolicySelector{Regex: rcfg}, expectedErr: nil},
	}

	for _, test := range table {
		_, err := LoadSelector(test.cfg)
		if err != test.expectedErr {
			t.Errorf("Unexpected error %v", err)
		}
	}
}

func TestStaticSelector(t *testing.T) {
	sel := NewStaticSelector(&config.StaticSelectorConf{Policy: "ocis"})
	ctx := context.Background()
	want := "ocis"
	got, err := sel(ctx)
	if got != want {
		t.Errorf("Expected policy %v got %v", want, got)
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	sel = NewStaticSelector(&config.StaticSelectorConf{Policy: "foo"})

	want = "foo"
	got, err = sel(ctx)
	if got != want {
		t.Errorf("Expected policy %v got %v", want, got)
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

type migrationTestCase struct {
	AccSvcShouldReturnError bool
	Claims                  map[string]interface{}
	Expected                string
}

func TestMigrationSelector(t *testing.T) {
	cfg := config.MigrationSelectorConf{
		AccFoundPolicy:        "found",
		AccNotFoundPolicy:     "not_found",
		UnauthenticatedPolicy: "unauth",
	}
	var tests = []migrationTestCase{
		{true, map[string]interface{}{oidc.PreferredUsername: "Hans"}, "not_found"},
		{true, map[string]interface{}{oidc.Email: "hans@example.test"}, "not_found"},
		{false, map[string]interface{}{oidc.PreferredUsername: "Hans"}, "found"},
		{false, nil, "unauth"},
	}

	for _, tc := range tests {
		tc := tc
		sut := NewMigrationSelector(&cfg, mockAccSvc(tc.AccSvcShouldReturnError))
		ctx := oidc.NewContext(context.Background(), tc.Claims)

		got, err := sut(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got != tc.Expected {
			t.Errorf("Expected Policy %v got %v", tc.Expected, got)
		}
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

type testCase struct {
	Name     string
	Context  context.Context
	Expected string
}

func TestClaimsSelector(t *testing.T) {
	sel := NewClaimsSelector(&config.ClaimsSelectorConf{
		DefaultPolicy:         "default",
		UnauthenticatedPolicy: "unauthenticated",
	})

	var tests = []testCase{
		{"unatuhenticated", context.Background(), "unauthenticated"},
		{"default", oidc.NewContext(context.Background(), map[string]interface{}{oidc.OcisRoutingPolicy: ""}), "default"},
		{"claim-value", oidc.NewContext(context.Background(), map[string]interface{}{oidc.OcisRoutingPolicy: "ocis.routing.policy-value"}), "ocis.routing.policy-value"},
	}
	for _, tc := range tests {
		got, err := sel(tc.Context)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got != tc.Expected {
			t.Errorf("Expected Policy %v got %v", tc.Expected, got)
		}
	}
}

func TestRegexSelector(t *testing.T) {
	sel := NewRegexSelector(&config.RegexSelectorConf{
		DefaultPolicy: "default",
		MatchesPolicies: []config.RegexRuleConf{
			{Priority: 10, Property: "mail", Match: "marie@example.org", Policy: "ocis"},
			{Priority: 20, Property: "mail", Match: "[^@]+@example.org", Policy: "oc10"},
			{Priority: 30, Property: "username", Match: "(einstein|feynman)", Policy: "ocis"},
			{Priority: 40, Property: "username", Match: ".+", Policy: "oc10"},
			{Priority: 50, Property: "id", Match: "4c510ada-c86b-4815-8820-42cdf82c3d51", Policy: "ocis"},
			{Priority: 60, Property: "id", Match: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c", Policy: "oc10"},
		},
		UnauthenticatedPolicy: "unauthenticated",
	})

	var tests = []testCase{
		{"unauthenticated", context.Background(), "unauthenticated"},
		{"default", revauser.ContextSetUser(context.Background(), &userv1beta1.User{}), "default"},
		{"mail-ocis", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Mail: "marie@example.org"}), "ocis"},
		{"mail-oc10", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Mail: "einstein@example.org"}), "oc10"},
		{"username-einstein", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Username: "einstein"}), "ocis"},
		{"username-feynman", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Username: "feynman"}), "ocis"},
		{"username-marie", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Username: "marie"}), "oc10"},
		{"id-nil", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Id: &userv1beta1.UserId{}}), "default"},
		{"id-1", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Id: &userv1beta1.UserId{OpaqueId: "4c510ada-c86b-4815-8820-42cdf82c3d51"}}), "ocis"},
		{"id-2", revauser.ContextSetUser(context.Background(), &userv1beta1.User{Id: &userv1beta1.UserId{OpaqueId: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}}), "oc10"},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			got, err := sel(tc.Context)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got != tc.Expected {
				t.Errorf("Expected Policy %v got %v", tc.Expected, got)
			}
		})
	}
}
