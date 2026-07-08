package backend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/selector"
	"google.golang.org/grpc"
)

// TestGroupsClaimPresent verifies only a present, non-null groups claim triggers reconciliation, so an absent claim can't strip a user's memberships.
func TestGroupsClaimPresent(t *testing.T) {
	const claim = "groups"

	tests := []struct {
		name   string
		claims map[string]interface{}
		want   bool
	}{
		{"claim absent", map[string]interface{}{"other": "x"}, false},
		{"claim null", map[string]interface{}{claim: nil}, false},
		{"empty claims", map[string]interface{}{}, false},
		{"present but empty", map[string]interface{}{claim: []interface{}{}}, true},
		{"present with groups", map[string]interface{}{claim: []interface{}{"admins"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := groupsClaimPresent(tt.claims, claim); got != tt.want {
				t.Errorf("groupsClaimPresent(%v) = %v, want %v", tt.claims, got, tt.want)
			}
		})
	}
}

// With Groups unconfigured, SyncGroupMemberships must return before touching the gateway; the nil selector panics if it doesn't.
func TestSyncGroupMemberships_GroupsClaimUnconfigured(t *testing.T) {
	b := &cs3backend{
		Options: Options{
			autoProvisionClaims: config.AutoProvisionClaims{Groups: ""},
		},
	}
	claims := map[string]interface{}{
		"groups": []interface{}{"admins"},
	}
	if err := b.SyncGroupMemberships(context.Background(), &cs3.User{}, claims); err != nil {
		t.Fatal(err)
	}
}

// graphCalls records the mutating libregraph requests so tests assert on side effects, not control flow.
type graphCalls struct {
	mu          sync.Mutex
	addMember   []string // group ids the user was added to
	createGroup []string // display names of created groups
	delMember   []string // group ids the user was removed from
}

func (g *graphCalls) snapshot() ([]string, []string, []string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string(nil), g.addMember...),
		append([]string(nil), g.createGroup...),
		append([]string(nil), g.delMember...)
}

// newGraphServer emulates the libregraph endpoints SyncGroupMemberships calls; existingGroup exists without the user as a member, currentGroups are the user's memberships.
func newGraphServer(t *testing.T, existingGroup string, currentGroups []string, calls *graphCalls) *httptest.Server {
	t.Helper()

	// Names GetGroup resolves (else 404); CreateGroup adds to it so a create-then-reread stays consistent.
	resolvable := make(map[string]bool, len(currentGroups)+1)
	for _, name := range currentGroups {
		resolvable[name] = true
	}
	if existingGroup != "" {
		resolvable[existingGroup] = true
	}

	mux := http.NewServeMux()

	// GetUser with Expand=memberOf: return the user's current memberships.
	mux.HandleFunc("/graph/v1.0/users/", func(w http.ResponseWriter, r *http.Request) {
		trimmed := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/users/")
		u := libregraph.NewUserWithDefaults()
		u.SetId(trimmed)
		groups := make([]libregraph.Group, 0, len(currentGroups))
		for _, name := range currentGroups {
			gr := libregraph.NewGroupWithDefaults()
			gr.SetId(name)
			gr.SetDisplayName(name)
			groups = append(groups, *gr)
		}
		u.SetMemberOf(groups)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	})

	mux.HandleFunc("/graph/v1.0/groups/", func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/groups/")

		switch {
		case strings.HasSuffix(rest, "/members/$ref") && r.Method == http.MethodPost:
			gid := strings.TrimSuffix(rest, "/members/$ref")
			calls.mu.Lock()
			calls.addMember = append(calls.addMember, gid)
			calls.mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		case strings.Contains(rest, "/members/") && r.Method == http.MethodDelete:
			gid := strings.SplitN(rest, "/members/", 2)[0]
			calls.mu.Lock()
			calls.delMember = append(calls.delMember, gid)
			calls.mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// GetGroup: rest is the group id/name.
		gid, _ := url.PathUnescape(rest)
		calls.mu.Lock()
		known := resolvable[gid]
		calls.mu.Unlock()
		if known {
			gr := libregraph.NewGroupWithDefaults()
			gr.SetId(gid)
			gr.SetDisplayName(gid)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gr)
			return
		}
		http.Error(w, `{"error":{"code":"itemNotFound"}}`, http.StatusNotFound)
	})

	// CreateGroup.
	mux.HandleFunc("/graph/v1.0/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		var gr libregraph.Group
		_ = json.NewDecoder(r.Body).Decode(&gr)
		gr.SetId("created-" + gr.GetDisplayName())
		calls.mu.Lock()
		calls.createGroup = append(calls.createGroup, gr.GetDisplayName())
		resolvable[gr.GetDisplayName()] = true
		resolvable[gr.GetId()] = true
		calls.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(gr)
	})

	return httptest.NewServer(mux)
}

// newSyncBackend wires a cs3backend against srv and a token-issuing mock gateway to drive the real SyncGroupMemberships.
func newSyncBackend(t *testing.T, srv *httptest.Server, claims config.AutoProvisionClaims) *cs3backend {
	t.Helper()

	gatewayClient := &cs3mocks.GatewayAPIClient{}
	gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(
		&gateway.AuthenticateResponse{
			Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
			Token:  "service-token",
		}, nil)

	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)
	t.Cleanup(func() { pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway") })

	// Register the graph service so the libregraph client resolves to srv.
	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	reg := registry.GetRegistry(registry.Inmemory())
	graphService := registry.BuildHTTPService("com.owncloud.web.graph", u.Host, "")
	if err := reg.Register(graphService); err != nil {
		t.Fatal(err)
	}
	// Register a gateway node so gatewaySelector.Next() succeeds; address is unused as the selector returns the mock.
	gatewayService := registry.BuildGRPCService("com.owncloud.api.gateway", "", "any", "")
	if err := reg.Register(gatewayService); err != nil {
		t.Fatal(err)
	}
	// Deregister from the process-global registry so stale nodes don't leak into later tests.
	t.Cleanup(func() {
		_ = reg.Deregister(graphService)
		_ = reg.Deregister(gatewayService)
	})
	graphSelector := selector.NewSelector(selector.Registry(reg))

	return &cs3backend{
		graphSelector: graphSelector,
		Options: Options{
			logger:              log.NewLogger(),
			gatewaySelector:     gatewaySelector,
			selector:            graphSelector,
			autoProvisionClaims: claims,
		},
	}
}

func testUser() *cs3.User {
	return &cs3.User{Id: &cs3.UserId{OpaqueId: "user-1"}}
}

// TestSyncGroupMemberships_ClaimSet_CreatesAndJoins verifies that with the claim configured, a non-existing claimed group is created and joined.
func TestSyncGroupMemberships_ClaimSet_CreatesAndJoins(t *testing.T) {
	calls := &graphCalls{}
	srv := newGraphServer(t, "", nil, calls)
	defer srv.Close()

	b := newSyncBackend(t, srv, config.AutoProvisionClaims{Groups: "groups"})

	claims := map[string]interface{}{"groups": []interface{}{"new-group"}}
	if err := b.SyncGroupMemberships(context.Background(), testUser(), claims); err != nil {
		t.Fatal(err)
	}

	add, create, _ := calls.snapshot()
	if len(create) != 1 || create[0] != "new-group" {
		t.Errorf("expected group 'new-group' to be created, got %v", create)
	}
	if len(add) != 1 {
		t.Errorf("expected user to be added to the created group exactly once, got %v", add)
	}
}

// TestSyncGroupMemberships_ExistingGroupJoined verifies a claim naming an existing
// group joins the user without re-creating it.
func TestSyncGroupMemberships_ExistingGroupJoined(t *testing.T) {
	calls := &graphCalls{}
	// "admin" already exists as a local group; user is not yet a member.
	srv := newGraphServer(t, "admin", nil, calls)
	defer srv.Close()

	b := newSyncBackend(t, srv, config.AutoProvisionClaims{Groups: "groups"})

	claims := map[string]interface{}{"groups": []interface{}{"admin"}}
	if err := b.SyncGroupMemberships(context.Background(), testUser(), claims); err != nil {
		t.Fatal(err)
	}

	add, create, _ := calls.snapshot()
	if len(add) != 1 || add[0] != "admin" {
		t.Errorf("user should be joined to existing group 'admin', joined: %v", add)
	}
	if len(create) != 0 {
		t.Errorf("an existing group must not be (re)created: %v", create)
	}
}

// TestSyncGroupMemberships_Deprovision verifies a user is removed from a group absent from their claim.
func TestSyncGroupMemberships_Deprovision(t *testing.T) {
	calls := &graphCalls{}
	// User is currently a member of "editors"; their claim no longer carries it.
	srv := newGraphServer(t, "", []string{"editors"}, calls)
	defer srv.Close()

	b := newSyncBackend(t, srv, config.AutoProvisionClaims{Groups: "groups"})

	claims := map[string]interface{}{"groups": []interface{}{}}
	if err := b.SyncGroupMemberships(context.Background(), testUser(), claims); err != nil {
		t.Fatal(err)
	}

	_, _, del := calls.snapshot()
	if len(del) != 1 || del[0] != "editors" {
		t.Errorf("user should be removed from 'editors' (not in claim), removals: %v", del)
	}
}
