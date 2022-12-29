package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
	"github.com/r3labs/sse/v2"
	"google.golang.org/grpc/metadata"
)

// SSE provides server sent events functionality
type SSE struct {
	cfg    *config.Config
	gwc    gateway.GatewayAPIClient
	server *sse.Server
}

// NewSSE returns a new SSE instance
func NewSSE(cfg *config.Config) (*SSE, error) {
	server := sse.New()

	gwc, err := pool.GetGatewayServiceClient(cfg.Reva.Address)
	if err != nil {
		return nil, err
	}

	return &SSE{
		cfg:    cfg,
		gwc:    gwc,
		server: server,
	}, nil
}

// ServeHTTP allows clients to subscribe to sse
func (s *SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := ctx.ContextGetUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uid := u.GetId().GetOpaqueId()
	if uid == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream := s.server.CreateStream(uid)
	stream.AutoReplay = false

	// add stream to URL
	q := r.URL.Query()
	q.Set("stream", uid)
	r.URL.RawQuery = q.Encode()

	s.server.ServeHTTP(w, r)

}

// ListenForEvents listens for events to inform clients about changes. Blocking. Start in seperate go routine.
func (s *SSE) ListenForEvents(evts <-chan events.Event) {
	for e := range evts {
		rcps, ev := s.extractDetails(e.Event)

		for r := range rcps {
			s.server.Publish(r, &sse.Event{Data: ev})
		}
	}
}

// extracts recipients and builds event to send to client
func (s *SSE) extractDetails(e interface{}) (<-chan string, []byte) {

	// determining recipients can take longer. We spawn a seperate go routine to do it
	ch := make(chan string)
	go s.determineRecipients(e, ch)

	var event interface{}

	switch ev := e.(type) {
	case events.UploadReady:
		t := ev.Timestamp.Format("2006-01-02 15:04:05")
		id, _ := storagespace.FormatReference(ev.FileRef)
		event = UploadReady{
			FileID:    id,
			SpaceID:   ev.FileRef.GetResourceId().GetSpaceId(),
			Filename:  ev.Filename,
			Timestamp: t,
			Message:   fmt.Sprintf("[%s] Hello! The file %s is ready to work with", t, ev.Filename),
		}

	}

	b, err := json.Marshal(event)
	if err != nil {
		log.Println("ERROR:", err)
	}

	return ch, b
}

func (s *SSE) determineRecipients(e interface{}, ch chan<- string) {
	defer close(ch)

	var (
		ref  *provider.Reference
		user *user.User
	)
	switch ev := e.(type) {
	case events.UploadReady:
		ref = ev.FileRef
		user = ev.ExecutingUser
	}

	// impersonate executing user to stat the resource
	// FIXME: What to do if executing user is not member of the space?
	ctx, err := s.impersonate(user.GetId())
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	space, err := s.getStorageSpace(ctx, ref.GetResourceId().GetSpaceId())
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	informed := make(map[string]struct{})
	inform := func(users ...string) {
		for _, id := range users {
			if _, ok := informed[id]; ok {
				continue
			}
			ch <- id
			informed[id] = struct{}{}
		}
	}

	// inform executing user first and foremost
	inform(user.GetId().GetOpaqueId())

	// inform space members next
	var grants map[string]*provider.ResourcePermissions
	if err := utils.ReadJSONFromOpaque(space.GetOpaque(), "grants", &grants); err == nil {
		groups := s.resolveGroups(ctx, space)

		// FIXME: Which space permissions allow me to get this event?
		for id := range grants {
			users, ok := groups[id]
			if !ok {
				users = []string{id}
			}

			inform(users...)
		}
	}

	// inform share recipients also
	// TODO: How to get all shares pointing to the resource?
	resp, err := s.gwc.ListShares(ctx, listSharesRequest(ref.GetResourceId()))
	if err != nil || resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		log.Println("ERROR:", err, resp.GetStatus().GetCode())
		return
	}

	for _, share := range resp.GetShares() {
		users := []string{share.GetGrantee().GetUserId().GetOpaqueId()}
		if users[0] == "" {
			users = s.getGroupMembers(ctx, share.GetGrantee().GetGroupId().GetOpaqueId())
		}

		inform(users...)
	}
}

// resolve group members for all groups of a space
func (s *SSE) resolveGroups(ctx context.Context, space *provider.StorageSpace) map[string][]string {
	groups := make(map[string][]string)

	var grpids map[string]struct{}
	if err := utils.ReadJSONFromOpaque(space.GetOpaque(), "groups", &grpids); err != nil {
		return groups
	}

	for g := range grpids {
		groups[g] = append(groups[g], s.getGroupMembers(ctx, g)...)
	}

	return groups
}

// returns authenticated context of `userID`
func (s *SSE) impersonate(userID *user.UserId) (context.Context, error) {
	getUserResponse, err := s.gwc.GetUser(context.Background(), &user.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	if getUserResponse.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	// Get auth context
	ownerCtx := revactx.ContextSetUser(context.Background(), getUserResponse.User)
	authRes, err := s.gwc.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + userID.OpaqueId,
		ClientSecret: s.cfg.MachineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error impersonating user: %s", authRes.Status.Message)
	}

	return metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, authRes.Token), nil
}

// calls gateway for storage space
func (s *SSE) getStorageSpace(ctx context.Context, id string) (*provider.StorageSpace, error) {
	resp, err := s.gwc.ListStorageSpaces(ctx, listStorageSpaceRequest(id))
	if err != nil {
		return nil, err
	}

	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK || len(resp.GetStorageSpaces()) != 1 {
		return nil, fmt.Errorf("can't fetch storage space: %s", resp.GetStatus().GetCode())
	}

	return resp.GetStorageSpaces()[0], nil
}

// calls gateway for group members
func (s *SSE) getGroupMembers(ctx context.Context, id string) []string {
	var members []string
	r, err := s.gwc.GetGroup(ctx, getGroupRequest(id))
	if err != nil || r.GetStatus().GetCode() != rpc.Code_CODE_OK {
		log.Println("ERROR", err, r.GetStatus().GetCode())
	}

	for _, uid := range r.GetGroup().GetMembers() {
		members = append(members, uid.GetOpaqueId())
	}

	return members
}

func listStorageSpaceRequest(id string) *provider.ListStorageSpacesRequest {
	return &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &provider.ListStorageSpacesRequest_Filter_Id{
					Id: &provider.StorageSpaceId{
						OpaqueId: id,
					},
				},
			},
		},
	}
}

func getGroupRequest(id string) *group.GetGroupRequest {
	return &group.GetGroupRequest{
		GroupId: &group.GroupId{
			OpaqueId: id,
		},
	}
}

func listSharesRequest(id *provider.ResourceId) *collaboration.ListSharesRequest {
	return &collaboration.ListSharesRequest{
		Filters: []*collaboration.Filter{
			{
				Type: collaboration.Filter_TYPE_RESOURCE_ID,
				Term: &collaboration.Filter_ResourceId{
					ResourceId: id,
				},
			},
		},
	}
}
