package command

import (
	"bytes"
	"context"
	"errors"
	"testing"

	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/storage/utils/grants"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/test-go/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// fakeViewerTrashStorage is a storage provider that keeps grants in memory and
// applies the updates it receives, like the real one.
type fakeViewerTrashStorage struct {
	spaces       []*provider.StorageSpace
	grants       map[string][]*provider.Grant // by space root node id
	listGrantErr map[string]error             // by space root node id
	updateCode   map[string]rpc.Code          // by grantee id
	updates      []*provider.UpdateGrantRequest
}

func (f *fakeViewerTrashStorage) ListStorageSpaces(_ context.Context, req *provider.ListStorageSpacesRequest, _ ...grpc.CallOption) (*provider.ListStorageSpacesResponse, error) {
	var spaceType string
	for _, filter := range req.GetFilters() {
		if filter.GetType() == provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE {
			spaceType = filter.GetSpaceType()
		}
	}
	res := &provider.ListStorageSpacesResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}
	for _, s := range f.spaces {
		if spaceType == "" || s.GetSpaceType() == spaceType {
			res.StorageSpaces = append(res.StorageSpaces, s)
		}
	}
	return res, nil
}

func (f *fakeViewerTrashStorage) ListGrants(_ context.Context, req *provider.ListGrantsRequest, _ ...grpc.CallOption) (*provider.ListGrantsResponse, error) {
	root := req.GetRef().GetResourceId().GetOpaqueId()
	if err := f.listGrantErr[root]; err != nil {
		return nil, err
	}
	res := &provider.ListGrantsResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}
	for _, g := range f.grants[root] {
		res.Grants = append(res.Grants, proto.Clone(g).(*provider.Grant))
	}
	return res, nil
}

func (f *fakeViewerTrashStorage) UpdateGrant(_ context.Context, req *provider.UpdateGrantRequest, _ ...grpc.CallOption) (*provider.UpdateGrantResponse, error) {
	f.updates = append(f.updates, req)
	id := granteeID(req.GetGrant().GetGrantee())
	if code, ok := f.updateCode[id]; ok && code != rpc.Code_CODE_OK {
		return &provider.UpdateGrantResponse{Status: &rpc.Status{Code: code}}, nil
	}
	root := req.GetRef().GetResourceId().GetOpaqueId()
	for i, g := range f.grants[root] {
		if granteeID(g.GetGrantee()) == id {
			f.grants[root][i] = proto.Clone(req.GetGrant()).(*provider.Grant)
		}
	}
	return &provider.UpdateGrantResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil
}

func granteeID(g *provider.Grantee) string {
	if g.GetType() == provider.GranteeType_GRANTEE_TYPE_GROUP {
		return g.GetGroupId().GetOpaqueId()
	}
	return g.GetUserId().GetOpaqueId()
}

func projectSpace(id, name string) *provider.StorageSpace {
	return &provider.StorageSpace{
		Id:        &provider.StorageSpaceId{OpaqueId: "storage$" + id},
		Name:      name,
		SpaceType: _spaceTypeProject,
		Root:      &provider.ResourceId{StorageId: "storage", SpaceId: id, OpaqueId: id},
	}
}

func userGrant(id string, perms *provider.ResourcePermissions) *provider.Grant {
	return &provider.Grant{
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_USER,
			Id:   &provider.Grantee_UserId{UserId: &userpb.UserId{OpaqueId: id}},
		},
		Permissions: proto.Clone(perms).(*provider.ResourcePermissions),
		Creator:     &userpb.UserId{OpaqueId: "teacher"},
	}
}

func groupGrant(id string, perms *provider.ResourcePermissions) *provider.Grant {
	return &provider.Grant{
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_GROUP,
			Id:   &provider.Grantee_GroupId{GroupId: &grouppb.GroupId{OpaqueId: id}},
		},
		Permissions: proto.Clone(perms).(*provider.ResourcePermissions),
		Creator:     &userpb.UserId{OpaqueId: "teacher"},
	}
}

// newFakeStorage returns one project space with an old "Can view" user and
// group, a new "Can view" user, a manager, a "Can edit" member and a read-only
// member with trash access that is not an old "Can view" member.
func newFakeStorage() *fakeViewerTrashStorage {
	oldViewer := userGrant("old-viewer", oldSpaceViewerPermissions())
	oldViewer.Expiration = &typesv1beta1.Timestamp{Seconds: 2000000000}

	return &fakeViewerTrashStorage{
		spaces: []*provider.StorageSpace{
			projectSpace("space-1", "Space 1"),
			{
				Id:        &provider.StorageSpaceId{OpaqueId: "storage$personal-1"},
				Name:      "Personal",
				SpaceType: _spaceTypePersonal,
				Root:      &provider.ResourceId{StorageId: "storage", SpaceId: "personal-1", OpaqueId: "personal-1"},
			},
		},
		grants: map[string][]*provider.Grant{
			"space-1": {
				oldViewer,
				groupGrant("old-viewer-group", oldSpaceViewerPermissions()),
				userGrant("new-viewer", fixedSpaceViewerPermissions()),
				userGrant("manager", conversions.NewManagerRole().CS3ResourcePermissions()),
				userGrant("editor", conversions.NewSpaceEditorWithoutVersionsWithoutTrashbinRole().CS3ResourcePermissions()),
				userGrant("file-viewer", conversions.NewViewerRole().CS3ResourcePermissions()),
			},
			"personal-1": {
				userGrant("personal-viewer", oldSpaceViewerPermissions()),
			},
		},
	}
}

func resultsByMember(report viewerTrashReport) map[string]viewerTrashResult {
	results := map[string]viewerTrashResult{}
	for _, e := range report.entries {
		results[e.member] = e.result
	}
	return results
}

func TestViewerPermissionSets(t *testing.T) {
	oldPerms := oldSpaceViewerPermissions()
	fixedPerms := fixedSpaceViewerPermissions()

	require.True(t, oldPerms.GetListRecycle())
	require.False(t, fixedPerms.GetListRecycle())

	// the two sets differ only in the trash bin access
	fixedPerms.ListRecycle = true
	require.True(t, grants.PermissionsEqual(oldPerms, fixedPerms))

	// building the sets must not change the role itself
	require.Equal(t, oldPerms.GetListRecycle(), oldSpaceViewerPermissions().GetListRecycle())
}

func TestFixViewerTrashDryRun(t *testing.T) {
	storage := newFakeStorage()

	report, err := fixViewerTrash(context.Background(), storage, "", true)
	require.NoError(t, err)

	require.Empty(t, storage.updates, "a dry run must not change anything")
	for _, e := range report.entries {
		require.Equal(t, "storage$space-1", e.spaceID, "the report shows the space ID in the form --space-id takes")
	}
	require.Equal(t, 1, report.spacesChecked, "only project spaces are checked")
	require.Equal(t, map[string]viewerTrashResult{
		"user old-viewer":        _viewerTrashWouldFix,
		"group old-viewer-group": _viewerTrashWouldFix,
		"user file-viewer":       _viewerTrashCheck,
	}, resultsByMember(report))
}

func TestFixViewerTrash(t *testing.T) {
	storage := newFakeStorage()
	before := newFakeStorage().grants["space-1"]

	report, err := fixViewerTrash(context.Background(), storage, "", false)
	require.NoError(t, err)

	require.Equal(t, map[string]viewerTrashResult{
		"user old-viewer":        _viewerTrashFixed,
		"group old-viewer-group": _viewerTrashFixed,
		"user file-viewer":       _viewerTrashCheck,
	}, resultsByMember(report))
	require.Len(t, storage.updates, 2)

	for _, req := range storage.updates {
		require.True(t, utils.ExistsInOpaque(req.GetOpaque(), "spacegrant"), "must be marked as a space member change")
		require.Equal(t, _spaceTypeProject, utils.ReadPlainFromOpaque(req.GetOpaque(), "spacetype"))
		require.Equal(t, "space-1", req.GetRef().GetResourceId().GetOpaqueId(), "must change the member on the space root")
		require.True(t, grants.PermissionsEqual(fixedSpaceViewerPermissions(), req.GetGrant().GetPermissions()))
		require.Equal(t, "teacher", req.GetGrant().GetCreator().GetOpaqueId(), "must keep who added the member")
	}
	require.Equal(t, uint64(2000000000), storage.updates[0].GetGrant().GetExpiration().GetSeconds(), "must keep the expiration date")

	// all other members are unchanged
	for i, g := range storage.grants["space-1"] {
		switch granteeID(g.GetGrantee()) {
		case "old-viewer", "old-viewer-group":
			require.False(t, g.GetPermissions().GetListRecycle())
		default:
			require.True(t, proto.Equal(before[i], g), "member %s must not change", granteeID(g.GetGrantee()))
		}
	}
	// personal spaces are not touched
	require.True(t, storage.grants["personal-1"][0].GetPermissions().GetListRecycle())

	// a second run finds nothing left to change
	storage.updates = nil
	report, err = fixViewerTrash(context.Background(), storage, "", false)
	require.NoError(t, err)
	require.Empty(t, storage.updates)
	require.Equal(t, 0, report.countMatchingResults(_viewerTrashFixed))
}

func TestFixViewerTrashOneSpace(t *testing.T) {
	storage := newFakeStorage()
	storage.spaces = append(storage.spaces, projectSpace("space-2", "Space 2"))
	storage.grants["space-2"] = []*provider.Grant{userGrant("other-viewer", oldSpaceViewerPermissions())}

	for _, id := range []string{"space-2", "storage$space-2", "storage$space-2!space-2"} {
		t.Run(id, func(t *testing.T) {
			report, err := fixViewerTrash(context.Background(), storage, id, true)
			require.NoError(t, err)
			require.Equal(t, 1, report.spacesChecked)
			require.Equal(t, map[string]viewerTrashResult{"user other-viewer": _viewerTrashWouldFix}, resultsByMember(report))
		})
	}

	_, err := fixViewerTrash(context.Background(), storage, "does-not-exist", true)
	require.Error(t, err)
}

func TestFixViewerTrashFailures(t *testing.T) {
	storage := newFakeStorage()
	storage.updateCode = map[string]rpc.Code{"old-viewer": rpc.Code_CODE_PERMISSION_DENIED}
	storage.spaces = append(storage.spaces, projectSpace("broken", "Broken"))
	storage.listGrantErr = map[string]error{"broken": errors.New("connection refused")}

	report, err := fixViewerTrash(context.Background(), storage, "", false)
	require.NoError(t, err, "one failing member or space must not stop the run")

	results := resultsByMember(report)
	require.Contains(t, results["user old-viewer"], _viewerTrashFailed)
	require.Equal(t, _viewerTrashFixed, results["group old-viewer-group"], "the other members are still changed")
	require.Equal(t, 2, report.spacesChecked)
	require.Equal(t, 2, report.countMatchingResults(_viewerTrashFailed), "the failing member and the space whose members could not be listed")
	for _, e := range report.entries {
		if e.member == "group old-viewer-group" {
			require.Empty(t, e.reason, "a failure must not carry over to the next member")
		}
	}

	// the summary reports the failures
	var out bytes.Buffer
	require.NoError(t, printViewerTrashReport(&out, report, false))
	require.Contains(t, out.String(), "failed: CODE_PERMISSION_DENIED")
	require.Contains(t, out.String(), "failed: could not list members: connection refused")
	require.Contains(t, out.String(), "spaces checked: 2, members changed: 1, failed: 2, check manually: 1")
}

// failingWriter fails every write, like a full disk or a closed pipe.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("disk full")
}

func TestPrintViewerTrashReportWriteError(t *testing.T) {
	report, err := fixViewerTrash(context.Background(), newFakeStorage(), "", true)
	require.NoError(t, err)
	require.NotEmpty(t, report.entries)

	err = printViewerTrashReport(failingWriter{}, report, true)
	require.Error(t, err, "a report that could not be written must not look like a success")
	require.Contains(t, err.Error(), "disk full")
}
