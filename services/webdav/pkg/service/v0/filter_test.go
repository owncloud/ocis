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
					GetPath:  true,
					Stat:     true,
					InitiateFileDownload: true,
				},
				ArbitraryMetadata: &provider.ArbitraryMetadata{
					Metadata: map[string]string{
						"http://owncloud.org/ns/favorite": "1",
					},
				},
			},
			{
				Id:       &provider.ResourceId{StorageId: "storage1", SpaceId: "space1", OpaqueId: "dir1"},
				Name:     "not-favorite-dir",
				Type:     provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Size:     0,
				Mtime:    &typesv1beta1.Timestamp{Seconds: 1700000000},
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
	rep, err := readReport(r.Body)
	if err != nil {
		t.Fatalf("readReport: %v", err)
	}
	if rep.FilterFiles == nil {
		t.Fatal("expected FilterFiles to be parsed")
	}

	// Re-create request since body was consumed
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

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
	rep, _ := readReport(r.Body)

	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

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
	rep, _ := readReport(r.Body)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

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
				Id:       &provider.ResourceId{StorageId: "s1", SpaceId: "space1", OpaqueId: "nested"},
				Name:     "Documents",
				Type:     provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Mtime:    &typesv1beta1.Timestamp{Seconds: 1700000000},
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
	rep, _ := readReport(r.Body)
	r = newFilterFilesRequest("bob")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

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
	rep, _ := readReport(r.Body)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

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
	rep, _ := readReport(r.Body)
	r = newFilterFilesRequest("alice")
	rr := httptest.NewRecorder()
	svc.handleFilterFiles(rr, r, rep.FilterFiles)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d: %s", rr.Code, rr.Body.String())
	}
}
