package kiteworks

import (
	"context"
	"io"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kwlib"
)

// versionRef identifies a specific KW file version. Version keys are encoded
// as "<fileID>@<versionID>" so that GetMD and Download can distinguish them
// from regular file references and look up the right KW resource.
type versionRef struct {
	fileID    string
	versionID string
}

// parseVersionRef returns a versionRef if ref.OpaqueId encodes a version
// ("fileID@versionID"), and false otherwise.
func parseVersionRef(ref *provider.Reference) (versionRef, bool) {
	fileID, versionID, ok := strings.Cut(ref.GetResourceId().GetOpaqueId(), "@")
	if !ok {
		return versionRef{}, false
	}
	return versionRef{fileID: fileID, versionID: versionID}, true
}

// getVersionMD returns a minimal ResourceInfo for a version reference.
// Calls GetFileByID on the parent file to verify access; the returned info
// carries only what handleGet requires (Type=FILE, no "processing" status).
func (d *Driver) getVersionMD(ctx context.Context, ref *provider.Reference, vr versionRef) (*provider.ResourceInfo, error) {
	if _, err := d.client(ctx).GetFileByID(vr.fileID); err != nil {
		return nil, err
	}
	return &provider.ResourceInfo{
		Id:   ref.GetResourceId(),
		Type: provider.ResourceType_RESOURCE_TYPE_FILE,
	}, nil
}

// downloadVersion streams the content of a specific version, used by Download
// when the ref carries a version OpaqueId.
func (d *Driver) downloadVersion(ctx context.Context, ref *provider.Reference, vr versionRef, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	fi, err := d.client(ctx).GetFileByID(vr.fileID)
	if err != nil {
		return nil, nil, err
	}
	if !fi.HasPermission(kwlib.PermDownload) {
		return nil, nil, errtypes.PermissionDenied(vr.fileID)
	}
	// Fetch content first to read Content-Length from the response headers.
	// This is required because the dataprovider uses ri.Size to set Content-Length,
	// and an empty size causes browsers to receive a zero-byte file.
	resp, err := d.client(ctx).GetVersionContents(vr.fileID, vr.versionID)
	if err != nil {
		return nil, nil, err
	}
	ri := &provider.ResourceInfo{
		Id:   ref.GetResourceId(),
		Type: provider.ResourceType_RESOURCE_TYPE_FILE,
		Name: fi.Name,
		Path: fi.Path,
		Size: uint64(max(resp.ContentLength, 0)),
	}
	if !openReaderFunc(ri) {
		resp.Body.Close()
		return ri, nil, nil
	}
	return ri, resp.Body, nil
}

// ListRevisions returns file versions as CS3 FileVersion entries.
// Keys are encoded as "<fileID>@<versionID>" so that GetMD and Download
// can route version download requests to the right KW endpoint.
func (d *Driver) ListRevisions(ctx context.Context, ref *provider.Reference) ([]*provider.FileVersion, error) {
	nodeID, _, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	fi, err := d.client(ctx).GetFileByID(nodeID)
	if err != nil {
		return nil, err
	}
	if !fi.HasPermission(kwlib.PermVersionView) {
		return nil, errtypes.PermissionDenied(nodeID)
	}
	versions, err := d.client(ctx).GetFileVersions(nodeID)
	if err != nil {
		return nil, err
	}
	revs := make([]*provider.FileVersion, 0, len(versions))
	for _, v := range versions {
		revs = append(revs, &provider.FileVersion{
			Key:   nodeID + "@" + v.ID,
			Mtime: uint64(time.Time(v.Created).Unix()),
			Size:  uint64(v.Size),
		})
	}
	return revs, nil
}

// RestoreRevision promotes revisionKey to be the current version of the file.
// revisionKey is "<fileID>@<versionID>" as returned by ListRevisions.
// ocdav passes the space-root ResourceId as ref (opaque_id == space_id), so we
// use the fileID embedded in the key rather than resolveRef.
func (d *Driver) RestoreRevision(ctx context.Context, _ *provider.Reference, revisionKey string) (*storage.RestoreRevisionResult, error) {
	fileID, versionID, ok := strings.Cut(revisionKey, "@")
	if !ok {
		return nil, errtypes.BadRequest("kiteworks: invalid revision key: " + revisionKey)
	}
	fi, err := d.client(ctx).GetFileByID(fileID)
	if err != nil {
		return nil, err
	}
	if !fi.HasPermission(kwlib.PermVersionPromote) {
		return nil, errtypes.PermissionDenied(fileID)
	}
	if err := d.client(ctx).PromoteFileVersion(fileID, versionID); err != nil {
		return nil, err
	}
	return &storage.RestoreRevisionResult{}, nil
}

// DownloadRevision streams the content of the revision identified by revisionKey.
// revisionKey is "<fileID>@<versionID>" as returned by ListRevisions.
func (d *Driver) DownloadRevision(ctx context.Context, ref *provider.Reference, revisionKey string, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	nodeID, _, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, nil, err
	}
	fi, err := d.client(ctx).GetFileByID(nodeID)
	if err != nil {
		return nil, nil, err
	}
	if !fi.HasPermission(kwlib.PermDownload) {
		return nil, nil, errtypes.PermissionDenied(nodeID)
	}
	_, versionID, _ := strings.Cut(revisionKey, "@")
	resp, err := d.client(ctx).GetVersionContents(nodeID, versionID)
	if err != nil {
		return nil, nil, err
	}
	ri := &provider.ResourceInfo{
		Id:   ref.GetResourceId(),
		Type: provider.ResourceType_RESOURCE_TYPE_FILE,
		Name: fi.Name,
		Path: fi.Path,
		Size: uint64(max(resp.ContentLength, 0)),
	}
	if !openReaderFunc(ri) {
		resp.Body.Close()
		return ri, nil, nil
	}
	return ri, resp.Body, nil
}
