package kiteworks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kwlib"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/reva/v2/pkg/utils"
)

func init() {
	registry.Register("kiteworks", New)
}

// Config holds the driver configuration.
type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	APIToken string `mapstructure:"api_token"`
	Insecure bool   `mapstructure:"insecure"`
	MountID  string `mapstructure:"mount_id"`
}

// Driver implements storage.FS against a Kiteworks box (read-only).
type Driver struct {
	factory   *kwlib.APIClientFactory
	apiToken  string
	storageID string
	log       zerolog.Logger
	// TODO: lock tokens are ephemeral; lost on restart and not shared across nodes.
	locks sync.Map // nodeID → *provider.Lock
}

// New returns a read-only Kiteworks storage driver.
func New(m map[string]interface{}, _ events.Stream, log *zerolog.Logger) (storage.FS, error) {
	c := &Config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	c.Endpoint = strings.TrimRight(c.Endpoint, "/")

	storageID := c.MountID
	if storageID == "" {
		storageID = "kiteworks"
	}

	l := zerolog.Nop()
	if log != nil {
		l = *log
	}

	return &Driver{
		factory:   kwlib.NewClientFactory(c.Endpoint, "reva-kiteworks/1.0", c.Insecure),
		apiToken:  c.APIToken,
		storageID: storageID,
		log:       l,
	}, nil
}

func (d *Driver) client(ctx context.Context) *kwlib.APIClient {
	token := d.apiToken
	if token == "" {
		token, _ = ctxpkg.ContextGetToken(ctx)
	}
	return d.factory.Build("", "", "", token, &d.log)
}

// toResourceInfo converts a kwlib.FileInfo to a CS3 ResourceInfo.
// spaceRootPath is the absolute KW path of the space root folder; it is stripped
// from fi.Path to produce a space-relative path with a leading "/".
func (d *Driver) toResourceInfo(fi *kwlib.FileInfo, spaceID, spaceRootPath string) *provider.ResourceInfo {
	relPath := strings.TrimPrefix(fi.Path, spaceRootPath)
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}
	spaceRoot := &provider.ResourceId{
		StorageId: d.storageID,
		SpaceId:   spaceID,
		OpaqueId:  spaceID,
	}
	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  fi.ID,
		},
		Space:         &provider.StorageSpace{Root: spaceRoot},
		Name:          fi.Name,
		Path:          relPath,
		Etag:          fi.ETag(),
		Mtime:         utils.TimeToTS(fi.MTime()),
		PermissionSet: permissionSet(fi),
	}
	if fi.ParentID != nil && *fi.ParentID != "" && *fi.ParentID != "0" {
		ri.ParentId = &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  *fi.ParentID,
		}
	}
	if fi.IsDir() {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_CONTAINER
		ri.MimeType = "httpd/unix-directory"
	} else {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_FILE
		if fi.Size != nil {
			ri.Size = uint64(*fi.Size)
		}
	}
	return ri
}

// Capabilities omits Trash, Sharing and ArbitraryMetadata: those methods still
// reject with NotSupported.
func (d *Driver) Capabilities(_ context.Context) storage.Capabilities {
	return storage.Capabilities{
		Upload:          true,
		CreateContainer: true,
		Delete:          true,
		Move:            true,
		Versioning:      true,
		Locking:         true,
	}
}

// --- Read methods ---

func (d *Driver) Shutdown(_ context.Context) error { return nil }

func (d *Driver) ListStorageSpaces(ctx context.Context, _ []*provider.ListStorageSpacesRequest_Filter, _ bool) ([]*provider.StorageSpace, error) {
	c := d.client(ctx)
	dirs, err := c.GetTopFolders()
	if err != nil {
		return nil, err
	}

	u, hasUser := ctxpkg.ContextGetUser(ctx)

	spaces := make([]*provider.StorageSpace, 0, len(dirs.Data))
	for i := range dirs.Data {
		fi := &dirs.Data[i]
		opaque := utils.AppendPlainToOpaque(nil, "spaceAlias", "project/"+fi.Name)
		if hasUser && u.GetId().GetOpaqueId() != "" {
			grants := map[string]*provider.ResourcePermissions{
				u.Id.OpaqueId: spaceRole(fi),
			}
			if b, err := json.Marshal(grants); err == nil {
				opaque.Map["grants"] = &types.OpaqueEntry{Decoder: "json", Value: b}
			}
		}
		spaces = append(spaces, &provider.StorageSpace{
			Id:        &provider.StorageSpaceId{OpaqueId: fi.ID},
			Name:      fi.Name,
			SpaceType: "project",
			Root: &provider.ResourceId{
				StorageId: d.storageID,
				SpaceId:   fi.ID,
				OpaqueId:  fi.ID,
			},
			RootInfo: d.toResourceInfo(fi, fi.ID, fi.Path),
			Mtime:    utils.TimeToTS(fi.MTime()),
			Opaque:   opaque,
		})
	}
	return spaces, nil
}

// resolveRef walks a CS3 reference to the target KW node ID and space ID.
// If ref.Path is non-empty it resolves each component through ListFolderContents.
func (d *Driver) resolveRef(ctx context.Context, ref *provider.Reference) (nodeID, spaceID string, err error) {
	spaceID = ref.GetResourceId().GetSpaceId()
	nodeID = ref.GetResourceId().GetOpaqueId()
	if nodeID == "" {
		nodeID = spaceID
	}
	relPath := strings.Trim(strings.TrimPrefix(ref.GetPath(), "./"), "/.")
	if relPath == "" {
		return nodeID, spaceID, nil
	}

	c := d.client(ctx)
	for _, part := range strings.Split(relPath, "/") {
		if part == "" {
			continue
		}
		children, err := c.ListFolderContents(nodeID)
		if err != nil {
			return "", "", err
		}
		var found bool
		for i := range children {
			if children[i].Name == part {
				nodeID = children[i].ID
				found = true
				break
			}
		}
		if !found {
			return "", "", errtypes.NotFound(part)
		}
	}
	return nodeID, spaceID, nil
}

// spaceRootPath fetches the absolute KW path of the space root folder.
// Used by callers that need to convert absolute KW paths to space-relative paths.
func (d *Driver) spaceRootPath(ctx context.Context, spaceID string) (string, error) {
	root, err := d.client(ctx).GetFolderByID(spaceID)
	if err != nil {
		return "", err
	}
	return root.Path, nil
}

func (d *Driver) GetMD(ctx context.Context, ref *provider.Reference, _, _ []string) (*provider.ResourceInfo, error) {
	if vr, ok := parseVersionRef(ref); ok {
		return d.getVersionMD(ctx, ref, vr)
	}
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return d.nodeMD(ctx, nodeID, spaceID, rootPath)
}

// nodeMD fetches metadata for a node, trying folder first then file.
func (d *Driver) nodeMD(ctx context.Context, nodeID, spaceID, spaceRootPath string) (*provider.ResourceInfo, error) {
	c := d.client(ctx)
	fi, err := c.GetFolderByID(nodeID)
	if err == nil {
		return d.toResourceInfo(fi, spaceID, spaceRootPath), nil
	}
	var ce *kwlib.ClientError
	if !errors.As(err, &ce) || (ce.StatusCode != http.StatusNotFound && ce.StatusCode != http.StatusForbidden) {
		return nil, err
	}
	fi, err = c.GetFileByID(nodeID)
	if err != nil {
		var fce *kwlib.ClientError
		if errors.As(err, &fce) && fce.StatusCode == http.StatusNotFound {
			return nil, errtypes.NotFound(nodeID)
		}
		return nil, err
	}
	return d.toResourceInfo(fi, spaceID, spaceRootPath), nil
}

func (d *Driver) ListFolder(ctx context.Context, ref *provider.Reference, _, _ []string) ([]*provider.ResourceInfo, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	items, err := d.client(ctx).ListFolderContents(nodeID)
	if err != nil {
		return nil, err
	}

	infos := make([]*provider.ResourceInfo, 0, len(items))
	for i := range items {
		ri := d.toResourceInfo(&items[i], spaceID, rootPath)
		// ocdav does path.Join(requestPath, info.Path) for ListFolder results,
		// so Path must be just the filename, not the space-root-relative path.
		ri.Path = ri.Name
		infos = append(infos, ri)
	}
	return infos, nil
}

func (d *Driver) Download(ctx context.Context, ref *provider.Reference, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	if vr, ok := parseVersionRef(ref); ok {
		return d.downloadVersion(ctx, ref, vr, openReaderFunc)
	}
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, nil, err
	}

	ri, err := d.nodeMD(ctx, nodeID, spaceID, rootPath)
	if err != nil {
		return nil, nil, err
	}

	if !openReaderFunc(ri) {
		return ri, nil, nil
	}

	resp, err := d.client(ctx).GetFileContents(nodeID, "")
	if err != nil {
		return nil, nil, err
	}
	return ri, resp.Body, nil
}

func (d *Driver) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	rootPath, err := d.spaceRootPath(ctx, id.GetSpaceId())
	if err != nil {
		return "", err
	}
	ri, err := d.nodeMD(ctx, id.GetOpaqueId(), id.GetSpaceId(), rootPath)
	if err != nil {
		return "", err
	}
	return ri.Path, nil
}

func (d *Driver) ListGrants(_ context.Context, _ *provider.Reference) ([]*provider.Grant, error) {
	return []*provider.Grant{}, nil
}

// TODO: this maps only one of Kiteworks' two quota models.
//
//  1. System quota: a per-folder limit on a top-level folder. Matches the oCIS per-space
//     quota.
//  2. Total folder quota: one pool shared by every folder the account owns. No oCIS
//     equivalent.
//
// The choice is per-folder and defaults to total folder quota. So out of the box every
// space reports the same total and used (the account pool), and an unlimited pool makes
// every space report unrestricted.
//
// Both come back under the same field names (storage_quota, storage_used,
// storage_available), so this endpoint cannot tell them apart. Read useFolderQuota from
// GET /rest/folders/{id} and treat the pooled case as "no real space quota".
//
// Enabling system quota:
//   - Admin > Users > Profiles > Standard > Collaboration > Enable Folders That Use System Quota
//   - Per folder, at creation time: Folder quota > Use system storage quota
//
// Total folder quota is set on the same profile page.
func (d *Driver) GetQuota(ctx context.Context, ref *provider.Reference) (uint64, uint64, uint64, error) {
	nodeID, _, err := d.resolveRef(ctx, ref)
	if err != nil || nodeID == "" {
		return 0, 0, math.MaxUint64, nil
	}
	q, err := d.client(ctx).GetFolderQuota(nodeID)
	if err != nil {
		// The endpoint requires file_add on the folder, so viewers get a 403.
		// Quota is informational, so report unlimited rather than fail the caller.
		return 0, 0, math.MaxUint64, nil
	}
	total := uint64(max(q.StorageQuota, 0))
	used := uint64(max(q.StorageUsed, 0))
	// total=0 means no quota is applied to the folder. Signal unlimited remaining
	// so the graph service doesn't compute 0/0 = NaN and report "exceeded".
	if total == 0 {
		return 0, used, math.MaxUint64, nil
	}
	return total, used, uint64(max(q.StorageAvailable, 0)), nil
}

func (d *Driver) GetLock(ctx context.Context, ref *provider.Reference) (*provider.Lock, error) {
	nodeID, _, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	fi, err := d.client(ctx).GetFileByID(nodeID)
	if err != nil {
		return nil, err
	}
	if !fi.Locked {
		return nil, nil
	}
	if v, ok := d.locks.Load(nodeID); ok {
		return v.(*provider.Lock), nil
	}
	return &provider.Lock{LockId: "kw-" + nodeID, Type: provider.LockType_LOCK_TYPE_EXCL}, nil
}

func (d *Driver) ListRecycle(_ context.Context, _ *provider.Reference, _, _ string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

// --- Write methods: all return NotSupported ---

func (d *Driver) CreateReference(_ context.Context, _ string, _ *url.URL) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CreateDir(ctx context.Context, ref *provider.Reference) (*storage.CreateDirResult, error) {
	parentRef := &provider.Reference{
		ResourceId: ref.GetResourceId(),
		Path:       path.Dir(ref.GetPath()),
	}
	parentID, spaceID, err := d.resolveRef(ctx, parentRef)
	if err != nil {
		return nil, err
	}

	name := path.Base(ref.GetPath())
	if name == "" || name == "." {
		return nil, errtypes.BadRequest("kiteworks: CreateDir requires a folder name")
	}

	folderID, err := d.client(ctx).CreateFolder(parentID, kwlib.CreateDirRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &storage.CreateDirResult{
		SpaceID: spaceID,
		ResourceID: &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  folderID,
		},
	}, nil
}

func (d *Driver) TouchFile(ctx context.Context, ref *provider.Reference, _ bool, _ string) (*storage.TouchFileResult, error) {
	parentRef := &provider.Reference{
		ResourceId: ref.GetResourceId(),
		Path:       path.Dir(ref.GetPath()),
	}
	parentID, spaceID, err := d.resolveRef(ctx, parentRef)
	if err != nil {
		return nil, err
	}

	name := path.Base(ref.GetPath())
	if name == "" || name == "." {
		return nil, errtypes.BadRequest("kiteworks: TouchFile requires a filename")
	}

	c := d.client(ctx)
	upload, err := c.InitializeUpload(parentID, name, 0, 1)
	if err != nil {
		return nil, err
	}
	fi, err := c.UploadChunk(ctx, upload.URI, name, bytes.NewReader(nil), 0, 0, true)
	if err != nil {
		return nil, err
	}
	return &storage.TouchFileResult{
		SpaceID: spaceID,
		ResourceID: &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  fi.ID,
		},
	}, nil
}

func (d *Driver) Delete(ctx context.Context, ref *provider.Reference) (*storage.DeleteResult, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}

	c := d.client(ctx)
	err = c.DeleteFolder(nodeID)
	if err != nil {
		var ce *kwlib.ClientError
		if !errors.As(err, &ce) || ce.StatusCode != http.StatusNotFound {
			return nil, err
		}
		if err = c.DeleteFile(nodeID); err != nil {
			return nil, err
		}
	}
	d.locks.Delete(nodeID)
	return &storage.DeleteResult{
		ResourceId: &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  nodeID,
		},
	}, nil
}

// rawNodeInfo fetches the raw KW FileInfo for a node, trying folder first.
// Used by write methods that need ParentID or type without a full toResourceInfo conversion.
func (d *Driver) rawNodeInfo(ctx context.Context, nodeID string) (*kwlib.FileInfo, error) {
	c := d.client(ctx)
	fi, err := c.GetFolderByID(nodeID)
	if err == nil {
		return fi, nil
	}
	var ce *kwlib.ClientError
	if !errors.As(err, &ce) || (ce.StatusCode != http.StatusNotFound && ce.StatusCode != http.StatusForbidden) {
		return nil, err
	}
	return c.GetFileByID(nodeID)
}

func moveNode(c *kwlib.APIClient, fi *kwlib.FileInfo, parentID string) error {
	if fi.IsDir() {
		return c.MoveFolder(fi.ID, parentID)
	}
	_, err := c.Move(fi, &kwlib.FileInfo{ID: parentID}, false)
	return err
}

func renameNode(c *kwlib.APIClient, fi *kwlib.FileInfo, name string) error {
	if fi.IsDir() {
		_, err := c.RenameFolder(fi, name)
		return err
	}
	_, err := c.RenameFile(fi, name, false)
	return err
}

func (d *Driver) Move(ctx context.Context, src, dst *provider.Reference) (*storage.MoveResult, error) {
	dstName := path.Base(dst.GetPath())
	if dstName == "" || dstName == "." || dstName == "/" {
		return nil, errtypes.BadRequest("kiteworks: Move requires a destination name")
	}

	srcNodeID, _, err := d.resolveRef(ctx, src)
	if err != nil {
		return nil, err
	}

	dstParentID, _, err := d.resolveRef(ctx, &provider.Reference{
		ResourceId: dst.GetResourceId(),
		Path:       path.Dir(dst.GetPath()),
	})
	if err != nil {
		return nil, err
	}
	if dstParentID == srcNodeID {
		return nil, errtypes.BadRequest("kiteworks: cannot move a resource into itself")
	}

	srcFI, err := d.rawNodeInfo(ctx, srcNodeID)
	if err != nil {
		return nil, err
	}

	srcParentID := ""
	if srcFI.ParentID != nil {
		srcParentID = *srcFI.ParentID
	}

	// KW has no combined move+rename. Rename first: it is the undoable step, and
	// the new name cannot collide in a parent the node has not left yet.
	c := d.client(ctx)
	renamed := dstName != srcFI.Name
	if renamed {
		if err := renameNode(c, srcFI, dstName); err != nil {
			return nil, err
		}
	}

	if srcParentID != dstParentID {
		if err := moveNode(c, srcFI, dstParentID); err != nil {
			if renamed {
				if rbErr := renameNode(c, srcFI, srcFI.Name); rbErr != nil {
					d.log.Error().Err(rbErr).Str("nodeID", srcNodeID).Msg("could not restore the original name after a failed move")
					return nil, fmt.Errorf("kiteworks: move failed and %q is left renamed to %q: %w", srcFI.Name, dstName, err)
				}
			}
			return nil, err
		}
	}

	return &storage.MoveResult{
		OldReference: src,
		NewReference: dst,
	}, nil
}

func (d *Driver) InitiateUpload(_ context.Context, _ *provider.Reference, _ int64, _ map[string]string) (map[string]string, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) Upload(_ context.Context, _ storage.UploadRequest, _ storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

// MarkProcessing is a no-op: KW controls its own file metadata so there is no
// reliable way to mark a node as "in-flight" without abusing the checkout-lock
// API (which would block concurrent uploads and leave orphaned locks on crash).
// Trade-off: a file touched by TouchFile is downloadable as a zero-byte stub
// until CommitUpload replaces it.
func (d *Driver) MarkProcessing(_ context.Context, _ *provider.Reference, _ bool, _ string) error {
	return nil
}

// TODO: large uploads still fail for the client in the sync path.
//
// There are two legs, and their speeds are unrelated:
//
//	client -> oCIS   TUS PATCHes, staged to a local bin file
//	oCIS   -> KW     this method, 8 MiB chunks over a KW upload session
//
// async: no problem. CommitUpload runs after the client is already done, so it
// is bound to no request.
//
// sync: problem. The whole oCIS -> KW transfer happens inside the client's last
// chunk request. That transfer can take minutes; the request is not meant to
// live that long.
//
// Ideas:
//  1. Always run the KW driver with async enabled. asyncfileuploads is not
//     wired into revaconfig, so the commit currently always runs inline.
//  2. Make sync partly async: answer the client before the file is fully in KW.
//     Sensible, but needs verifying against real clients.
//  3. Let the last chunk request run arbitrarily long for big files. Doubtful,
//     timeouts exist for a reason.
//  4. Commit chunk by chunk instead of at the end. Far from the current
//     architecture, and TUS resume does not map onto KW's sequential chunks.
//  5. Upload client -> KW directly and only register the file here. Drops
//     antivirus and checksums.
func (d *Driver) CommitUpload(ctx context.Context, ref *provider.Reference, _ string, source storage.UploadSource) error {
	if source.Body == nil {
		return errtypes.BadRequest("kiteworks: CommitUpload requires a non-nil body")
	}
	// the client connection's lifetime must not abort a commit already in flight
	ctx = context.WithoutCancel(ctx)
	nodeID, _, err := d.resolveRef(ctx, ref)
	if err != nil {
		return err
	}
	c := d.client(ctx)
	fi, err := c.GetFileByID(nodeID)
	if err != nil {
		return err
	}
	if err = c.UploadFileVersion(ctx, nodeID, fi.Name, source.Body, source.Length); err != nil {
		return err
	}
	if !source.NodeExisted {
		d.deletePlaceholderVersion(c, nodeID)
	}
	return nil
}

// deletePlaceholderVersion removes the zero-byte stub version created by TouchFile,
// leaving only the real content at version 1. Best-effort: all errors are logged,
// never returned; the upload already succeeded before this is called.
func (d *Driver) deletePlaceholderVersion(c *kwlib.APIClient, fileID string) {
	log := d.log.With().Str("fileID", fileID).Logger()

	versions, err := c.GetFileVersions(fileID)
	if err != nil {
		log.Warn().Err(err).Msg("placeholder cleanup: could not list versions")
		return
	}
	for _, v := range versions {
		if v.VersionNumber > 0 && v.Size == 0 {
			if err := c.DeleteFileVersion(fileID, v.ID); err != nil {
				ce := kwlib.AsClientError(err)
				switch ce.StatusCode {
				case http.StatusNotFound: // already gone, fine
				case http.StatusUnprocessableEntity: // last version: real bytes not uploaded, leave it
					log.Warn().Msg("placeholder cleanup: cannot delete last version, real upload may have failed")
				default:
					log.Warn().Err(err).Msg("placeholder cleanup: delete version failed")
				}
			}
			return
		}
	}
}

func (d *Driver) PrepareUpload(_ context.Context, _ *provider.Reference, _ string, info storage.UploadInfo) (*storage.PrepareUploadResult, error) {
	return &storage.PrepareUploadResult{VersionCreated: info.NodeExisted}, nil
}

func (d *Driver) RollbackUpload(ctx context.Context, _ *provider.Reference, _ string, info storage.RollbackInfo) error {
	// the prior version stays current, so there is nothing to undo
	if info.NodeExisted {
		return nil
	}
	if err := d.client(ctx).DeleteFile(info.NodeID); err != nil {
		if kwlib.AsClientError(err).StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (d *Driver) RestoreRecycleItem(_ context.Context, _ *provider.Reference, _, _ string, _ *provider.Reference) (*storage.RestoreRecycleItemResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) PurgeRecycleItem(_ context.Context, _ *provider.Reference, _, _ string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) EmptyRecycle(_ context.Context, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) AddGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) DenyGrant(_ context.Context, _ *provider.Reference, _ *provider.Grantee) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) RemoveGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UpdateGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) SetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UnsetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ []string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) SetLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) (*storage.SetLockResult, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	if err := d.client(ctx).LockFile(nodeID); err != nil {
		var ce *kwlib.ClientError
		if errors.As(err, &ce) && ce.StatusCode == http.StatusForbidden {
			return nil, errtypes.PreconditionFailed("kiteworks: file already locked")
		}
		return nil, err
	}
	d.locks.Store(nodeID, lock)
	return &storage.SetLockResult{SpaceID: spaceID}, nil
}

// RefreshLock is a no-op. KW has no lock-refresh endpoint, locks carry no
// expiry, and re-calling the lock endpoint on an already-locked file returns
// 403 even when the caller owns the lock.
func (d *Driver) RefreshLock(_ context.Context, _ *provider.Reference, _ *provider.Lock, _ string) error {
	return nil
}

func (d *Driver) Unlock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) (*storage.UnlockResult, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	if err := d.client(ctx).UnlockFile(nodeID); err != nil {
		var ce *kwlib.ClientError
		if errors.As(err, &ce) && ce.StatusCode == http.StatusForbidden {
			return nil, errtypes.PreconditionFailed("kiteworks: file not locked or locked by another user")
		}
		return nil, err
	}
	d.locks.Delete(nodeID)
	return &storage.UnlockResult{SpaceID: spaceID}, nil
}

func (d *Driver) CreateStorageSpace(_ context.Context, _ *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UpdateStorageSpace(_ context.Context, _ *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) DeleteStorageSpace(_ context.Context, _ *provider.DeleteStorageSpaceRequest) (*storage.DeleteStorageSpaceResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CreateHome(_ context.Context) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) GetHome(_ context.Context) (string, error) {
	return "", errtypes.NotSupported("kiteworks: read-only driver")
}
