package svc

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/constants"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/propfind"
)

const filterFilesBody = `<?xml version="1.0" encoding="utf-8" ?>
<oc:filter-files xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ns">
  <d:prop>
    <d:getlastmodified />
    <d:getetag />
    <d:getcontenttype />
    <d:getcontentlength />
    <oc:size />
    <d:resourcetype />
    <oc:fileid />
    <oc:favorite />
    <oc:permissions />
  </d:prop>
  <oc:filter-rules>
    <oc:favorite>1</oc:favorite>
  </oc:filter-rules>
</oc:filter-files>`

func setupTestWebdav(t *testing.T, gwClient *cs3mocks.GatewayAPIClient) Webdav {
	t.Helper()

	pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
	selector := pool.GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gwClient
		},
	)

	return Webdav{
		log:             log.NopLogger(),
		gatewaySelector: selector,
	}
}

func newFilterFilesRequest(username string) *http.Request {
	r := httptest.NewRequest("REPORT", "/dav/files/"+username, strings.NewReader(filterFilesBody))
	ctx := context.WithValue(r.Context(), constants.ContextKeyID, username)
	ctx = revactx.ContextSetToken(ctx, "test-token")
	ctx = revactx.ContextSetUser(ctx, &userpb.User{
		Id:       &userpb.UserId{OpaqueId: username},
		Username: username,
	})
	r.Header.Set(revactx.TokenHeader, "test-token")
	return r.WithContext(ctx)
}

// mustParseFilterFiles reads and decodes a filter-files REPORT body, failing the test
// immediately (rather than risking a nil-pointer panic later) if decoding errors or the
// filter-files element wasn't present.
func mustParseFilterFiles(t *testing.T, r *http.Request) *reportFilterFiles {
	t.Helper()

	rep, err := readReport(r.Body)
	if err != nil {
		t.Fatalf("readReport: %v", err)
	}
	if rep.FilterFiles == nil {
		t.Fatal("expected FilterFiles to be parsed")
	}
	return rep.FilterFiles
}

func okStatus() *rpcv1beta1.Status {
	return &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK}
}

func mockCheckPermission(gwClient *cs3mocks.GatewayAPIClient, allowed bool) {
	code := rpcv1beta1.Code_CODE_OK
	if !allowed {
		code = rpcv1beta1.Code_CODE_PERMISSION_DENIED
	}
	gwClient.On("CheckPermission", mock.Anything, mock.Anything).Return(
		&permissions.CheckPermissionResponse{
			Status: &rpcv1beta1.Status{Code: code},
		}, nil,
	)
}

func mockWhoAmI(gwClient *cs3mocks.GatewayAPIClient, username string) {
	gwClient.On("WhoAmI", mock.Anything, mock.Anything).Return(
		&gateway.WhoAmIResponse{
			Status: okStatus(),
			User: &userpb.User{
				Id:       &userpb.UserId{OpaqueId: username},
				Username: username,
			},
		}, nil,
	)
}

func mockListStorageSpaces(gwClient *cs3mocks.GatewayAPIClient, spaces []*provider.StorageSpace) {
	gwClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(
		&provider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: spaces,
		}, nil,
	)
}

func personalSpace(spaceID string) *provider.StorageSpace {
	return &provider.StorageSpace{
		Id:        &provider.StorageSpaceId{OpaqueId: spaceID},
		SpaceType: "personal",
		Name:      "alice",
		Root: &provider.ResourceId{
			StorageId: "storage1",
			SpaceId:   spaceID,
			OpaqueId:  spaceID,
		},
	}
}

func TestFilterFilesReturns207WithFavorites(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	// Root listing returns a favorited file and a non-favorited directory
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			{
				Id:       &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "file1"},
				Name:     "favorite-doc.txt",
				Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
				MimeType: "text/plain",
				Size:     1024,
				Etag:     "abc123",
				Mtime:    &typesv1beta1.Timestamp{Seconds: 1700000000},
				PermissionSet: &provider.ResourcePermissions{
					GetPath:              true,
					Stat:                 true,
					InitiateFileDownload: true,
				},
				ArbitraryMetadata: &provider.ArbitraryMetadata{
					Metadata: map[string]string{
						"http://owncloud.org/ns/favorite": "1",
					},
				},
			},
			{
				Id:    &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "dir1"},
				Name:  "not-favorite-dir",
				Type:  provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Size:  0,
				Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000},
				PermissionSet: &provider.ResourcePermissions{
					GetPath: true,
					Stat:    true,
				},
			},
		},
	}, nil)

	// Recursive listing of the non-favorite directory: empty
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "dir1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)

	// Re-create request since body was consumed
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Errorf("expected 207, got %d: %s", rr.Code, rr.Body.String())
	}

	// Verify XML response
	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(ms.Responses) != 1 {
		t.Fatalf("expected 1 favorite response, got %d", len(ms.Responses))
	}

	// Check href — personal space favorites use /dav/files/<user>/ format
	expectedHref := "/dav/files/alice/favorite-doc.txt"
	if ms.Responses[0].Href != expectedHref {
		t.Errorf("expected href %q, got %q", expectedHref, ms.Responses[0].Href)
	}
}

func TestFilterFilesEmptyFavoritesReturns207(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	// Root listing: no files at all
	gwClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)

	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Errorf("expected 207, got %d: %s", rr.Code, rr.Body.String())
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(ms.Responses) != 0 {
		t.Errorf("expected 0 responses, got %d", len(ms.Responses))
	}
}

func TestFilterFilesReturnsBothFileAndFolderFavorites(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	favMeta := &provider.ArbitraryMetadata{
		Metadata: map[string]string{propOcFavorite: "1"},
	}

	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			{
				Id:                &provider.ResourceId{StorageId: "s1", SpaceId: "space1", OpaqueId: "file1"},
				Name:              "my-file.pdf",
				Type:              provider.ResourceType_RESOURCE_TYPE_FILE,
				Size:              2048,
				MimeType:          "application/pdf",
				Mtime:             &typesv1beta1.Timestamp{Seconds: 1700000000},
				ArbitraryMetadata: favMeta,
			},
			{
				Id:                &provider.ResourceId{StorageId: "s1", SpaceId: "space1", OpaqueId: "dir1"},
				Name:              "my-folder",
				Type:              provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Size:              4096,
				Mtime:             &typesv1beta1.Timestamp{Seconds: 1700000000},
				ArbitraryMetadata: favMeta,
			},
		},
	}, nil)

	// Recursive listing of the favorite folder: empty
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "dir1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207, got %d", rr.Code)
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(ms.Responses) != 2 {
		t.Fatalf("expected 2 favorites (file + folder), got %d", len(ms.Responses))
	}

	// Verify one is a file (no collection) and one is a folder (collection)
	var hasFile, hasFolder bool
	for _, resp := range ms.Responses {
		if strings.Contains(resp.Href, "my-file.pdf") {
			hasFile = true
		}
		if strings.Contains(resp.Href, "my-folder") {
			hasFolder = true
		}
	}
	if !hasFile {
		t.Error("expected file favorite in response")
	}
	if !hasFolder {
		t.Error("expected folder favorite in response")
	}
}

func TestFilterFilesHrefsUseFilesPrefix(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "bob")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			{
				Id:    &provider.ResourceId{StorageId: "s1", SpaceId: "space1", OpaqueId: "nested"},
				Name:  "Documents",
				Type:  provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000},
			},
		},
	}, nil)

	// Nested file that is a favorite
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "nested"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			{
				Id:       &provider.ResourceId{StorageId: "s1", SpaceId: "space1", OpaqueId: "deepfile"},
				Name:     "notes.md",
				Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
				Size:     512,
				MimeType: "text/markdown",
				Mtime:    &typesv1beta1.Timestamp{Seconds: 1700000000},
				ArbitraryMetadata: &provider.ArbitraryMetadata{
					Metadata: map[string]string{propOcFavorite: "1"},
				},
			},
		},
	}, nil)

	r := newFilterFilesRequest("bob")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("bob")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207, got %d", rr.Code)
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(ms.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(ms.Responses))
	}

	// Personal space favorites use /dav/files/<user>/ format
	expectedHref := "/dav/files/bob/Documents/notes.md"
	if ms.Responses[0].Href != expectedHref {
		t.Errorf("expected href %q, got %q", expectedHref, ms.Responses[0].Href)
	}
}

func TestFilterFilesSkipsProjectSpaces(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")

	projectSpace := &provider.StorageSpace{
		Id:        &provider.StorageSpaceId{OpaqueId: "proj-space"},
		SpaceType: "project",
		Name:      "Engineering",
		Root: &provider.ResourceId{
			StorageId: "storage1",
			SpaceId:   "proj-space",
			OpaqueId:  "proj-space",
		},
	}
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{projectSpace})

	// No ListContainer mock needed — project spaces should be skipped entirely

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207, got %d: %s", rr.Code, rr.Body.String())
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Project spaces are not addressable under /dav/files/<user>/,
	// so they should be skipped — no favorites returned.
	if len(ms.Responses) != 0 {
		t.Fatalf("expected 0 responses (project spaces skipped), got %d", len(ms.Responses))
	}
}

func TestFilterFilesPermissionDenied(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockWhoAmI(gwClient, "alice")
	mockCheckPermission(gwClient, false)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d: %s", rr.Code, rr.Body.String())
	}
}

// favoritedFile builds a minimal favorited ResourceInfo for pagination tests, where the
// exact prop values don't matter beyond being present and distinguishable by name.
func favoritedFile(storageID, spaceID, opaqueID, name string) *provider.ResourceInfo {
	return &provider.ResourceInfo{
		Id:       &provider.ResourceId{StorageId: storageID, SpaceId: spaceID, OpaqueId: opaqueID},
		Name:     name,
		Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
		MimeType: "text/plain",
		Size:     1,
		Mtime:    &typesv1beta1.Timestamp{Seconds: 1700000000},
		ArbitraryMetadata: &provider.ArbitraryMetadata{
			Metadata: map[string]string{propOcFavorite: "1"},
		},
	}
}

// TestFilterFilesListStorageSpacesPagination regression-tests that ListStorageSpaces'
// pagination is followed to completion: without it, favorites in any space beyond the
// first page would be silently omitted.
func TestFilterFilesListStorageSpacesPagination(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")

	gwClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *provider.ListStorageSpacesRequest) bool {
		return req.PageToken == ""
	})).Return(&provider.ListStorageSpacesResponse{
		Status:        okStatus(),
		StorageSpaces: []*provider.StorageSpace{personalSpace("space1")},
		NextPageToken: "page2",
	}, nil)
	gwClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *provider.ListStorageSpacesRequest) bool {
		return req.PageToken == "page2"
	})).Return(&provider.ListStorageSpacesResponse{
		Status: okStatus(),
		StorageSpaces: []*provider.StorageSpace{{
			Id:        &provider.StorageSpaceId{OpaqueId: "space2"},
			SpaceType: "personal",
			Name:      "alice-second-space",
			Root:      &provider.ResourceId{StorageId: "storage1", SpaceId: "space2", OpaqueId: "space2"},
		}},
		NextPageToken: "",
	}, nil)

	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{favoritedFile("storage1", "space1", "file-in-space1", "from-space-1.txt")},
	}, nil)
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space2"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{favoritedFile("storage1", "space2", "file-in-space2", "from-space-2.txt")},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207, got %d: %s", rr.Code, rr.Body.String())
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// A favorite from the second page of storage spaces must be present, proving
	// ListStorageSpaces' NextPageToken was followed rather than stopping at page one.
	if len(ms.Responses) != 2 {
		t.Fatalf("expected 2 favorites (one per space, across two ListStorageSpaces pages), got %d: %+v", len(ms.Responses), ms.Responses)
	}
	var hasSpace1, hasSpace2 bool
	for _, resp := range ms.Responses {
		if strings.Contains(resp.Href, "from-space-1.txt") {
			hasSpace1 = true
		}
		if strings.Contains(resp.Href, "from-space-2.txt") {
			hasSpace2 = true
		}
	}
	if !hasSpace1 || !hasSpace2 {
		t.Fatalf("expected favorites from both spaces, got hasSpace1=%v hasSpace2=%v (responses: %+v)", hasSpace1, hasSpace2, ms.Responses)
	}
}

// TestFilterFilesListContainerPagination regression-tests that a single container's own
// ListContainer pagination is followed to completion, and that a subdirectory only present
// on a later page is still recursed into — without it, favorites (and whole subtrees) beyond
// a container's first page would be silently omitted.
func TestFilterFilesListContainerPagination(t *testing.T) {
	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	// Root container, page 1: one favorite, plus a NextPageToken.
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1" && req.PageToken == ""
	})).Return(&provider.ListContainerResponse{
		Status:        okStatus(),
		Infos:         []*provider.ResourceInfo{favoritedFile("storage1", "space1", "file-a", "file-a.txt")},
		NextPageToken: "root-page2",
	}, nil)
	// Root container, page 2: another favorite, plus a subdirectory that only appears here.
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1" && req.PageToken == "root-page2"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			favoritedFile("storage1", "space1", "file-b", "file-b.txt"),
			{
				Id:    &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "subdir"},
				Name:  "subdir-from-page-2",
				Type:  provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000},
			},
		},
	}, nil)
	// The subdirectory discovered only on the root's second page.
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "subdir"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{favoritedFile("storage1", "space1", "file-c", "file-c.txt")},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207, got %d: %s", rr.Code, rr.Body.String())
	}

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// file-a.txt (page 1), file-b.txt (page 2), and file-c.txt (inside the subdir only
	// discoverable via page 2) must all be present.
	if len(ms.Responses) != 3 {
		t.Fatalf("expected 3 favorites (2 from paginated root + 1 nested), got %d: %+v", len(ms.Responses), ms.Responses)
	}
	var hasA, hasB, hasC bool
	for _, resp := range ms.Responses {
		hasA = hasA || strings.Contains(resp.Href, "file-a.txt")
		hasB = hasB || strings.Contains(resp.Href, "file-b.txt")
		hasC = hasC || strings.Contains(resp.Href, "file-c.txt")
	}
	if !hasA || !hasB || !hasC {
		t.Fatalf("expected file-a.txt, file-b.txt and file-c.txt all present, got hasA=%v hasB=%v hasC=%v (responses: %+v)", hasA, hasB, hasC, ms.Responses)
	}
}

// TestFilterFilesRespectsContainerVisitLimit regression-tests maxFavoriteContainers: once the
// budget of ListContainer calls is exhausted, collection must stop cleanly (partial results,
// no error) rather than continuing to recurse without bound. The test lowers the limit to 2
// so it can prove the cutoff with a small, fast fixture instead of a real 5000-container tree.
func TestFilterFilesRespectsContainerVisitLimit(t *testing.T) {
	original := maxFavoriteContainers
	maxFavoriteContainers = 2
	t.Cleanup(func() { maxFavoriteContainers = original })

	gwClient := cs3mocks.NewGatewayAPIClient(t)
	svc := setupTestWebdav(t, gwClient)

	mockCheckPermission(gwClient, true)
	mockWhoAmI(gwClient, "alice")
	mockListStorageSpaces(gwClient, []*provider.StorageSpace{personalSpace("space1")})

	// Root container: three non-favorited subdirectories, none of which are favorited
	// themselves — every favorite in this fixture lives one level deeper.
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "space1"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos: []*provider.ResourceInfo{
			{Id: &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "dirA"}, Name: "dirA", Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER, Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000}},
			{Id: &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "dirB"}, Name: "dirB", Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER, Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000}},
			{Id: &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "dirC"}, Name: "dirC", Type: provider.ResourceType_RESOURCE_TYPE_CONTAINER, Mtime: &typesv1beta1.Timestamp{Seconds: 1700000000}},
		},
	}, nil)
	// Only dirA should ever be queried: the root call plus dirA's call already consume the
	// budget of 2, so dirB/dirC must never be requested. If the code regresses and queries
	// them anyway, this test fails via the mock's unexpected-call panic (no .On(...) is
	// registered for dirB/dirC), not just via the assertion below.
	gwClient.On("ListContainer", mock.Anything, mock.MatchedBy(func(req *provider.ListContainerRequest) bool {
		return req.Ref.ResourceId.OpaqueId == "dirA"
	})).Return(&provider.ListContainerResponse{
		Status: okStatus(),
		Infos:  []*provider.ResourceInfo{favoritedFile("storage1", "space1", "deep-file", "deep-file.txt")},
	}, nil)

	r := newFilterFilesRequest("alice")
	filterFiles := mustParseFilterFiles(t, r)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, filterFiles)

	// The request must still succeed (partial results), not error out.
	if rr.Code != http.StatusMultiStatus {
		t.Fatalf("expected 207 even when the container-visit limit is hit, got %d: %s", rr.Code, rr.Body.String())
	}

	gwClient.AssertNumberOfCalls(t, "ListContainer", 2)

	var ms propfind.MultiStatusResponseUnmarshalXML
	if err := xml.Unmarshal(rr.Body.Bytes(), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(ms.Responses) != 1 {
		t.Fatalf("expected exactly 1 favorite (from dirA, before the budget ran out), got %d: %+v", len(ms.Responses), ms.Responses)
	}
	if !strings.Contains(ms.Responses[0].Href, "deep-file.txt") {
		t.Fatalf("expected the favorite found in dirA, got %+v", ms.Responses[0])
	}
}
