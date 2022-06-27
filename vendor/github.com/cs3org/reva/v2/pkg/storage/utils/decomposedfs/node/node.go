// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package node

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/grpc/services/storageprovider"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/storage/utils/ace"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pkg/xattr"
)

// Define keys and values used in the node metadata
const (
	LockdiscoveryKey = "DAV:lockdiscovery"
	FavoriteKey      = "http://owncloud.org/ns/favorite"
	ShareTypesKey    = "http://owncloud.org/ns/share-types"
	ChecksumsKey     = "http://owncloud.org/ns/checksums"
	UserShareType    = "0"
	QuotaKey         = "quota"

	QuotaUnlimited    = "0"
	QuotaUncalculated = "-1"
	QuotaUnknown      = "-2"

	// TrashIDDelimiter represents the characters used to separate the nodeid and the deletion time.
	TrashIDDelimiter    = ".T."
	RevisionIDDelimiter = ".REV."

	// RootID defines the root node's ID
	RootID = "root"
)

// Node represents a node in the tree and provides methods to get a Parent or Child instance
type Node struct {
	SpaceID   string
	ParentID  string
	ID        string
	Name      string
	Blobsize  int64
	BlobID    string
	owner     *userpb.UserId
	Exists    bool
	SpaceRoot *Node

	lu PathLookup
}

// PathLookup defines the interface for the lookup component
type PathLookup interface {
	InternalRoot() string
	InternalPath(spaceID, nodeID string) string
	Path(ctx context.Context, n *Node) (path string, err error)
	ShareFolder() string
}

// New returns a new instance of Node
func New(spaceID, id, parentID, name string, blobsize int64, blobID string, owner *userpb.UserId, lu PathLookup) *Node {
	if blobID == "" {
		blobID = uuid.New().String()
	}
	return &Node{
		SpaceID:  spaceID,
		ID:       id,
		ParentID: parentID,
		Name:     name,
		Blobsize: blobsize,
		owner:    owner,
		lu:       lu,
		BlobID:   blobID,
	}
}

// ChangeOwner sets the owner of n to newOwner
func (n *Node) ChangeOwner(new *userpb.UserId) (err error) {
	rootNodePath := n.SpaceRoot.InternalPath()
	n.SpaceRoot.owner = new

	var attribs = map[string]string{xattrs.OwnerIDAttr: new.OpaqueId,
		xattrs.OwnerIDPAttr:  new.Idp,
		xattrs.OwnerTypeAttr: utils.UserTypeToString(new.Type)}

	if err := xattrs.SetMultiple(rootNodePath, attribs); err != nil {
		return err
	}

	return
}

// SetMetadata populates a given key with its value.
// Note that consumers should be aware of the metadata options on xattrs.go.
func (n *Node) SetMetadata(key string, val string) (err error) {
	nodePath := n.InternalPath()
	if err := xattrs.Set(nodePath, key, val); err != nil {
		return errors.Wrap(err, "Decomposedfs: could not set extended attribute")
	}
	return nil
}

// RemoveMetadata removes a given key
func (n *Node) RemoveMetadata(key string) (err error) {
	if err = xattrs.Remove(n.InternalPath(), key); err == nil || xattrs.IsAttrUnset(err) {
		return nil
	}
	return err
}

// GetMetadata reads the metadata for the given key
func (n *Node) GetMetadata(key string) (val string, err error) {
	nodePath := n.InternalPath()
	if val, err = xattrs.Get(nodePath, key); err != nil {
		return "", errors.Wrap(err, "Decomposedfs: could not get extended attribute")
	}
	return val, nil
}

// WriteAllNodeMetadata writes the Node metadata to disk
func (n *Node) WriteAllNodeMetadata() (err error) {
	attribs := make(map[string]string)

	attribs[xattrs.ParentidAttr] = n.ParentID
	attribs[xattrs.NameAttr] = n.Name
	attribs[xattrs.BlobIDAttr] = n.BlobID
	attribs[xattrs.BlobsizeAttr] = strconv.FormatInt(n.Blobsize, 10)

	nodePath := n.InternalPath()
	return xattrs.SetMultiple(nodePath, attribs)
}

// WriteOwner writes the space owner
func (n *Node) WriteOwner(owner *userpb.UserId) error {
	n.SpaceRoot.owner = owner
	attribs := map[string]string{
		xattrs.OwnerIDAttr:   owner.OpaqueId,
		xattrs.OwnerIDPAttr:  owner.Idp,
		xattrs.OwnerTypeAttr: utils.UserTypeToString(owner.Type),
	}
	nodeRootPath := n.SpaceRoot.InternalPath()
	if err := xattrs.SetMultiple(nodeRootPath, attribs); err != nil {
		return err
	}
	n.SpaceRoot.owner = owner
	return nil
}

// ReadNode creates a new instance from an id and checks if it exists
func ReadNode(ctx context.Context, lu PathLookup, spaceID, nodeID string, canListDisabledSpace bool) (n *Node, err error) {

	// read space root
	r := &Node{
		SpaceID: spaceID,
		lu:      lu,
		ID:      spaceID,
	}
	r.SpaceRoot = r
	r.owner, err = r.readOwner()
	switch {
	case xattrs.IsNotExist(err):
		return r, nil // swallow not found, the node defaults to exists = false
	case err != nil:
		return nil, err
	}
	r.Exists = true

	if !canListDisabledSpace && r.IsDisabled() {
		// no permission = not found
		return nil, errtypes.NotFound(spaceID)
	}

	// check if this is a space root
	if spaceID == nodeID {
		return r, nil
	}

	// read node
	n = &Node{
		SpaceID:   spaceID,
		lu:        lu,
		ID:        nodeID,
		SpaceRoot: r,
	}

	nodePath := n.InternalPath()

	// lookup name in extended attributes
	n.Name, err = xattrs.Get(nodePath, xattrs.NameAttr)
	switch {
	case xattrs.IsNotExist(err):
		return n, nil // swallow not found, the node defaults to exists = false
	case err != nil:
		return nil, err
	}

	n.Exists = true

	// lookup blobID in extended attributes
	n.BlobID, err = xattrs.Get(nodePath, xattrs.BlobIDAttr)
	switch {
	case xattrs.IsNotExist(err):
		return n, nil // swallow not found, the node defaults to exists = false
	case err != nil:
		return nil, err
	}

	// Lookup blobsize
	n.Blobsize, err = ReadBlobSizeAttr(nodePath)
	switch {
	case xattrs.IsNotExist(err):
		return n, nil // swallow not found, the node defaults to exists = false
	case err != nil:
		return nil, err
	}

	// lookup parent id in extended attributes
	n.ParentID, err = xattrs.Get(nodePath, xattrs.ParentidAttr)
	switch {
	case xattrs.IsAttrUnset(err):
		return nil, errtypes.InternalError(err.Error())
	case xattrs.IsNotExist(err):
		return n, nil // swallow not found, the node defaults to exists = false
	case err != nil:
		return nil, errtypes.InternalError(err.Error())
	}

	// TODO why do we stat the parent? to determine if the current node is in the trash we would need to traverse all parents...
	// we need to traverse all parents for permissions anyway ...
	// - we can compare to space root owner with the current user
	// - we can compare the share permissions on the root for spaces, which would work for managers
	// - for non managers / owners we need to traverse all path segments because an intermediate node might have been shared
	// - if we want to support negative acls we need to traverse the path for all users (but the owner)
	// for trashed items we need to check all parents
	// - one of them might have the trash suffix ...
	// - options:
	//   - move deleted nodes in a trash folder that is still part of the tree (aka freedesktop org trash spec)
	//     - shares should still be removed, which requires traversing all trashed children ... and it should be undoable ...
	//     - what if a trashed file is restored? will child items be accessible by a share?
	//   - compare paths of trash root items and the trashed file?
	//     - to determine the relative path of a file we would need to traverse all intermediate nodes anyway
	//   - recursively mark all children as trashed ... async ... it is ok when that is not synchronous
	//     - how do we pick up if an error occurs? write a journal somewhere? activity log / delta?
	//     - stat requests will not pick up trashed items at all
	//   - recursively move all children into the trash folder?
	//     - no need to write an additional trash entry
	//     - can be made more robust with a journal
	//     - same recursion mechanism can be used to purge items? sth we still need to do
	//   - flag the two above options with dtime
	_, err = os.Stat(n.ParentInternalPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errtypes.NotFound(err.Error())
		}
		return nil, err
	}

	return
}

// The os error is buried inside the fs.PathError error
func isNotDir(err error) bool {
	if perr, ok := err.(*fs.PathError); ok {
		if serr, ok2 := perr.Err.(syscall.Errno); ok2 {
			return serr == syscall.ENOTDIR
		}
	}
	return false
}

func readChildNodeFromLink(path string) (string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return "", err
	}
	nodeID := strings.TrimLeft(link, "/.")
	nodeID = strings.ReplaceAll(nodeID, "/", "")
	return nodeID, nil
}

// Child returns the child node with the given name
func (n *Node) Child(ctx context.Context, name string) (*Node, error) {
	spaceID := n.SpaceID
	if spaceID == "" && n.ParentID == "root" {
		spaceID = n.ID
	} else if n.SpaceRoot != nil {
		spaceID = n.SpaceRoot.ID
	}
	nodeID, err := readChildNodeFromLink(filepath.Join(n.InternalPath(), name))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || isNotDir(err) {

			c := &Node{
				SpaceID:   spaceID,
				lu:        n.lu,
				ParentID:  n.ID,
				Name:      name,
				SpaceRoot: n.SpaceRoot,
			}
			return c, nil // if the file does not exist we return a node that has Exists = false
		}

		return nil, errors.Wrap(err, "decomposedfs: Wrap: readlink error")
	}

	var c *Node
	c, err = ReadNode(ctx, n.lu, spaceID, nodeID, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not read child node")
	}
	c.SpaceRoot = n.SpaceRoot

	return c, nil
}

// Parent returns the parent node
func (n *Node) Parent() (p *Node, err error) {
	if n.ParentID == "" {
		return nil, fmt.Errorf("decomposedfs: root has no parent")
	}
	p = &Node{
		SpaceID:   n.SpaceID,
		lu:        n.lu,
		ID:        n.ParentID,
		SpaceRoot: n.SpaceRoot,
	}

	// parentPath := n.lu.InternalPath(spaceID, n.ParentID)
	parentPath := p.InternalPath()

	// lookup parent id in extended attributes
	if p.ParentID, err = xattrs.Get(parentPath, xattrs.ParentidAttr); err != nil {
		p.ParentID = ""
		return
	}
	// lookup name in extended attributes
	if p.Name, err = xattrs.Get(parentPath, xattrs.NameAttr); err != nil {
		p.Name = ""
		p.ParentID = ""
		return
	}

	// check node exists
	if _, err := os.Stat(parentPath); err == nil {
		p.Exists = true
	}
	return
}

// Owner returns the space owner
func (n *Node) Owner() *userpb.UserId {
	return n.SpaceRoot.owner
}

// readOwner reads the owner from the extended attributes of the space root
// in case either owner id or owner idp are unset we return an error and an empty owner object
func (n *Node) readOwner() (*userpb.UserId, error) {

	owner := &userpb.UserId{}

	rootNodePath := n.SpaceRoot.InternalPath()
	// lookup parent id in extended attributes
	var attr string
	var err error
	// lookup ID in extended attributes
	attr, err = xattrs.Get(rootNodePath, xattrs.OwnerIDAttr)
	switch {
	case err == nil:
		owner.OpaqueId = attr
	case xattrs.IsAttrUnset(err):
		// ignore
	default:
		return nil, err
	}

	// lookup IDP in extended attributes
	attr, err = xattrs.Get(rootNodePath, xattrs.OwnerIDPAttr)
	switch {
	case err == nil:
		owner.Idp = attr
	case xattrs.IsAttrUnset(err):
		// ignore
	default:
		return nil, err
	}

	// lookup type in extended attributes
	attr, err = xattrs.Get(rootNodePath, xattrs.OwnerTypeAttr)
	switch {
	case err == nil:
		owner.Type = utils.UserTypeMap(attr)
	case xattrs.IsAttrUnset(err):
		// ignore
	default:
		return nil, err
	}

	// owner is an optional property
	if owner.Idp == "" && owner.OpaqueId == "" {
		return nil, nil
	}
	return owner, nil
}

// PermissionSet returns the permission set for the current user
// the parent nodes are not taken into account
func (n *Node) PermissionSet(ctx context.Context) provider.ResourcePermissions {
	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		appctx.GetLogger(ctx).Debug().Interface("node", n).Msg("no user in context, returning default permissions")
		return NoPermissions()
	}
	if utils.UserEqual(u.Id, n.SpaceRoot.Owner()) {
		return OwnerPermissions()
	}
	// read the permissions for the current user from the acls of the current node
	if np, err := n.ReadUserPermissions(ctx, u); err == nil {
		return np
	}
	return NoPermissions()
}

// InternalPath returns the internal path of the Node
func (n *Node) InternalPath() string {
	return n.lu.InternalPath(n.SpaceID, n.ID)
}

// ParentInternalPath returns the internal path of the parent of the current node
func (n *Node) ParentInternalPath() string {
	return n.lu.InternalPath(n.SpaceID, n.ParentID)
}

// LockFilePath returns the internal path of the lock file of the node
func (n *Node) LockFilePath() string {
	return n.InternalPath() + ".lock"
}

// CalculateEtag returns a hash of fileid + tmtime (or mtime)
func CalculateEtag(nodeID string, tmTime time.Time) (string, error) {
	return calculateEtag(nodeID, tmTime)
}

// calculateEtag returns a hash of fileid + tmtime (or mtime)
func calculateEtag(nodeID string, tmTime time.Time) (string, error) {
	h := md5.New()
	if _, err := io.WriteString(h, nodeID); err != nil {
		return "", err
	}
	if tb, err := tmTime.UTC().MarshalBinary(); err == nil {
		if _, err := h.Write(tb); err != nil {
			return "", err
		}
	} else {
		return "", err
	}
	return fmt.Sprintf(`"%x"`, h.Sum(nil)), nil
}

// SetMtime sets the mtime and atime of a node
func (n *Node) SetMtime(ctx context.Context, mtime string) error {
	sublog := appctx.GetLogger(ctx).With().Interface("node", n).Logger()
	if mt, err := parseMTime(mtime); err == nil {
		nodePath := n.InternalPath()
		// updating mtime also updates atime
		if err := os.Chtimes(nodePath, mt, mt); err != nil {
			sublog.Error().Err(err).
				Time("mtime", mt).
				Msg("could not set mtime")
			return errors.Wrap(err, "could not set mtime")
		}
	} else {
		sublog.Error().Err(err).
			Str("mtime", mtime).
			Msg("could not parse mtime")
		return errors.Wrap(err, "could not parse mtime")
	}
	return nil
}

// SetEtag sets the temporary etag of a node if it differs from the current etag
func (n *Node) SetEtag(ctx context.Context, val string) (err error) {
	sublog := appctx.GetLogger(ctx).With().Interface("node", n).Logger()
	nodePath := n.InternalPath()
	var tmTime time.Time
	if tmTime, err = n.GetTMTime(); err != nil {
		// no tmtime, use mtime
		var fi os.FileInfo
		if fi, err = os.Lstat(nodePath); err != nil {
			return
		}
		tmTime = fi.ModTime()
	}
	var etag string
	if etag, err = calculateEtag(n.ID, tmTime); err != nil {
		return
	}

	// sanitize etag
	val = fmt.Sprintf("\"%s\"", strings.Trim(val, "\""))
	if etag == val {
		sublog.Debug().
			Str("etag", val).
			Msg("ignoring request to update identical etag")
		return nil
	}
	// etag is only valid until the calculated etag changes, is part of propagation
	return xattrs.Set(nodePath, xattrs.TmpEtagAttr, val)
}

// SetFavorite sets the favorite for the current user
// TODO we should not mess with the user here ... the favorites is now a user specific property for a file
// that cannot be mapped to extended attributes without leaking who has marked a file as a favorite
// it is a specific case of a tag, which is user individual as well
// TODO there are different types of tags
// 1. public that are managed by everyone
// 2. private tags that are only visible to the user
// 3. system tags that are only visible to the system
// 4. group tags that are only visible to a group ...
// urgh ... well this can be solved using different namespaces
// 1. public = p:
// 2. private = u:<uid>: for user specific
// 3. system = s: for system
// 4. group = g:<gid>:
// 5. app? = a:<aid>: for apps?
// obviously this only is secure when the u/s/g/a namespaces are not accessible by users in the filesystem
// public tags can be mapped to extended attributes
func (n *Node) SetFavorite(uid *userpb.UserId, val string) error {
	nodePath := n.InternalPath()
	// the favorite flag is specific to the user, so we need to incorporate the userid
	fa := fmt.Sprintf("%s:%s:%s@%s", xattrs.FavPrefix, utils.UserTypeToString(uid.GetType()), uid.GetOpaqueId(), uid.GetIdp())
	return xattrs.Set(nodePath, fa, val)
}

// IsDir returns true if the note is a directory
func (n *Node) IsDir() bool {
	nodePath := n.InternalPath()
	if fi, err := os.Lstat(nodePath); err == nil {
		if fi.IsDir() {
			if _, err = xattrs.Get(nodePath, xattrs.ReferenceAttr); err != nil {
				return true
			}
		}
	}
	return false
}

// AsResourceInfo return the node as CS3 ResourceInfo
func (n *Node) AsResourceInfo(ctx context.Context, rp *provider.ResourcePermissions, mdKeys []string, returnBasename bool) (ri *provider.ResourceInfo, err error) {
	sublog := appctx.GetLogger(ctx).With().Interface("node", n.ID).Logger()

	var fn string
	nodePath := n.InternalPath()

	var fi os.FileInfo

	nodeType := provider.ResourceType_RESOURCE_TYPE_INVALID
	if fi, err = os.Lstat(nodePath); err != nil {
		return
	}

	var target string
	switch {
	case fi.IsDir():
		if target, err = xattrs.Get(nodePath, xattrs.ReferenceAttr); err == nil {
			nodeType = provider.ResourceType_RESOURCE_TYPE_REFERENCE
		} else {
			nodeType = provider.ResourceType_RESOURCE_TYPE_CONTAINER
		}
	case fi.Mode().IsRegular():
		nodeType = provider.ResourceType_RESOURCE_TYPE_FILE
	case fi.Mode()&os.ModeSymlink != 0:
		nodeType = provider.ResourceType_RESOURCE_TYPE_SYMLINK
		// TODO reference using ext attr on a symlink
		// nodeType = provider.ResourceType_RESOURCE_TYPE_REFERENCE
	}

	id := &provider.ResourceId{StorageId: n.SpaceID, OpaqueId: n.ID}

	if returnBasename {
		fn = n.Name
	} else {
		fn, err = n.lu.Path(ctx, n)
		if err != nil {
			return nil, err
		}
	}

	var parentID *provider.ResourceId
	if p, err := n.Parent(); err == nil {
		parentID = &provider.ResourceId{
			StorageId: p.SpaceID,
			OpaqueId:  p.ID,
		}
	}

	ri = &provider.ResourceInfo{
		Id:            id,
		Path:          fn,
		Type:          nodeType,
		MimeType:      mime.Detect(nodeType == provider.ResourceType_RESOURCE_TYPE_CONTAINER, fn),
		Size:          uint64(n.Blobsize),
		Target:        target,
		PermissionSet: rp,
		Owner:         n.Owner(),
		ParentId:      parentID,
	}

	if nodeType == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
		ts, err := n.GetTreeSize()
		if err == nil {
			ri.Size = ts
		} else {
			ri.Size = 0 // make dirs always return 0 if it is unknown
			sublog.Debug().Err(err).Msg("could not read treesize")
		}
	}

	// TODO make etag of files use fileid and checksum

	var tmTime time.Time
	if tmTime, err = n.GetTMTime(); err != nil {
		// no tmtime, use mtime
		tmTime = fi.ModTime()
	}

	// use temporary etag if it is set
	if b, err := xattrs.Get(nodePath, xattrs.TmpEtagAttr); err == nil {
		ri.Etag = fmt.Sprintf(`"%x"`, b) // TODO why do we convert string(b)? is the temporary etag stored as string? -> should we use bytes? use hex.EncodeToString?
	} else if ri.Etag, err = calculateEtag(n.ID, tmTime); err != nil {
		sublog.Debug().Err(err).Msg("could not calculate etag")
	}

	// mtime uses tmtime if present
	// TODO expose mtime and tmtime separately?
	un := tmTime.UnixNano()
	ri.Mtime = &types.Timestamp{
		Seconds: uint64(un / 1000000000),
		Nanos:   uint32(un % 1000000000),
	}

	mdKeysMap := make(map[string]struct{})
	for _, k := range mdKeys {
		mdKeysMap[k] = struct{}{}
	}

	var returnAllKeys bool
	if _, ok := mdKeysMap["*"]; len(mdKeys) == 0 || ok {
		returnAllKeys = true
	}

	metadata := map[string]string{}

	// read favorite flag for the current user
	if _, ok := mdKeysMap[FavoriteKey]; returnAllKeys || ok {
		favorite := ""
		if u, ok := ctxpkg.ContextGetUser(ctx); ok {
			// the favorite flag is specific to the user, so we need to incorporate the userid
			if uid := u.GetId(); uid != nil {
				fa := fmt.Sprintf("%s:%s:%s@%s", xattrs.FavPrefix, utils.UserTypeToString(uid.GetType()), uid.GetOpaqueId(), uid.GetIdp())
				if val, err := xattrs.Get(nodePath, fa); err == nil {
					sublog.Debug().
						Str("favorite", fa).
						Msg("found favorite flag")
					favorite = val
				}
			} else {
				sublog.Error().Err(errtypes.UserRequired("userrequired")).Msg("user has no id")
			}
		} else {
			sublog.Error().Err(errtypes.UserRequired("userrequired")).Msg("error getting user from ctx")
		}
		metadata[FavoriteKey] = favorite
	}
	// read locks
	if _, ok := mdKeysMap[LockdiscoveryKey]; returnAllKeys || ok {
		if n.hasLocks(ctx) {
			err = readLocksIntoOpaque(ctx, n, ri)
			if err != nil {
				sublog.Debug().Err(errtypes.InternalError("lockfail"))
			}
		}
	}

	// share indicator
	if _, ok := mdKeysMap[ShareTypesKey]; returnAllKeys || ok {
		if n.hasUserShares(ctx) {
			metadata[ShareTypesKey] = UserShareType
		}
	}

	// checksums
	if _, ok := mdKeysMap[ChecksumsKey]; (nodeType == provider.ResourceType_RESOURCE_TYPE_FILE) && (returnAllKeys || ok) {
		// TODO which checksum was requested? sha1 adler32 or md5? for now hardcode sha1?
		readChecksumIntoResourceChecksum(ctx, nodePath, storageprovider.XSSHA1, ri)
		readChecksumIntoOpaque(ctx, nodePath, storageprovider.XSMD5, ri)
		readChecksumIntoOpaque(ctx, nodePath, storageprovider.XSAdler32, ri)
	}
	// quota
	if _, ok := mdKeysMap[QuotaKey]; (nodeType == provider.ResourceType_RESOURCE_TYPE_CONTAINER) && returnAllKeys || ok {
		if n.SpaceRoot != nil && n.SpaceRoot.InternalPath() != "" {
			readQuotaIntoOpaque(ctx, n.SpaceRoot.InternalPath(), ri)
		}
	}

	// only read the requested metadata attributes
	attrs, err := xattr.List(nodePath)
	if err != nil {
		sublog.Error().Err(err).Msg("error getting list of extended attributes")
	} else {
		for i := range attrs {
			// filter out non-custom properties
			if !strings.HasPrefix(attrs[i], xattrs.MetadataPrefix) {
				continue
			}
			// only read when key was requested
			k := attrs[i][len(xattrs.MetadataPrefix):]
			if _, ok := mdKeysMap[k]; returnAllKeys || ok {
				if val, err := xattrs.Get(nodePath, attrs[i]); err == nil {
					metadata[k] = val
				} else {
					sublog.Error().Err(err).
						Str("entry", attrs[i]).
						Msg("error retrieving xattr metadata")
				}
			}

		}
	}
	ri.ArbitraryMetadata = &provider.ArbitraryMetadata{
		Metadata: metadata,
	}

	sublog.Debug().
		Interface("ri", ri).
		Msg("AsResourceInfo")

	return ri, nil
}

func readChecksumIntoResourceChecksum(ctx context.Context, nodePath, algo string, ri *provider.ResourceInfo) {
	v, err := xattrs.Get(nodePath, xattrs.ChecksumPrefix+algo)
	switch {
	case err == nil:
		ri.Checksum = &provider.ResourceChecksum{
			Type: storageprovider.PKG2GRPCXS(algo),
			Sum:  hex.EncodeToString([]byte(v)),
		}
	case xattrs.IsAttrUnset(err):
		appctx.GetLogger(ctx).Debug().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("checksum not set")
	case xattrs.IsNotExist(err):
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("file not fount")
	default:
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("could not read checksum")
	}
}

func readChecksumIntoOpaque(ctx context.Context, nodePath, algo string, ri *provider.ResourceInfo) {
	v, err := xattrs.Get(nodePath, xattrs.ChecksumPrefix+algo)
	switch {
	case err == nil:
		if ri.Opaque == nil {
			ri.Opaque = &types.Opaque{
				Map: map[string]*types.OpaqueEntry{},
			}
		}
		ri.Opaque.Map[algo] = &types.OpaqueEntry{
			Decoder: "plain",
			Value:   []byte(hex.EncodeToString([]byte(v))),
		}
	case xattrs.IsAttrUnset(err):
		appctx.GetLogger(ctx).Debug().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("checksum not set")
	case xattrs.IsNotExist(err):
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("file not fount")
	default:
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Str("algorithm", algo).Msg("could not read checksum")
	}
}

// quota is always stored on the root node
func readQuotaIntoOpaque(ctx context.Context, nodePath string, ri *provider.ResourceInfo) {
	v, err := xattrs.Get(nodePath, xattrs.QuotaAttr)
	switch {
	case err == nil:
		// make sure we have a proper signed int
		// we use the same magic numbers to indicate:
		// -1 = uncalculated
		// -2 = unknown
		// -3 = unlimited
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			if ri.Opaque == nil {
				ri.Opaque = &types.Opaque{
					Map: map[string]*types.OpaqueEntry{},
				}
			}
			ri.Opaque.Map[QuotaKey] = &types.OpaqueEntry{
				Decoder: "plain",
				Value:   []byte(v),
			}
		} else {
			appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Str("quota", v).Msg("malformed quota")
		}
	case xattrs.IsAttrUnset(err):
		appctx.GetLogger(ctx).Debug().Err(err).Str("nodepath", nodePath).Msg("quota not set")
	case xattrs.IsNotExist(err):
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Msg("file not found when reading quota")
	default:
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Msg("could not read quota")
	}
}

// HasPropagation checks if the propagation attribute exists and is set to "1"
func (n *Node) HasPropagation() (propagation bool) {
	if b, err := xattrs.Get(n.InternalPath(), xattrs.PropagationAttr); err == nil {
		return b == "1"
	}
	return false
}

// GetTMTime reads the tmtime from the extended attributes
func (n *Node) GetTMTime() (tmTime time.Time, err error) {
	var b string
	if b, err = xattrs.Get(n.InternalPath(), xattrs.TreeMTimeAttr); err != nil {
		return
	}
	return time.Parse(time.RFC3339Nano, b)
}

// SetTMTime writes the UTC tmtime to the extended attributes or removes the attribute if nil is passed
func (n *Node) SetTMTime(t *time.Time) (err error) {
	if t == nil {
		err = xattrs.Remove(n.InternalPath(), xattrs.TreeMTimeAttr)
		if xattrs.IsAttrUnset(err) {
			return nil
		}
		return err
	}
	return xattrs.Set(n.InternalPath(), xattrs.TreeMTimeAttr, t.UTC().Format(time.RFC3339Nano))
}

// GetDTime reads the dtime from the extended attributes
func (n *Node) GetDTime() (tmTime time.Time, err error) {
	var b string
	if b, err = xattrs.Get(n.InternalPath(), xattrs.DTimeAttr); err != nil {
		return
	}
	return time.Parse(time.RFC3339Nano, b)
}

// SetDTime writes the UTC dtime to the extended attributes or removes the attribute if nil is passed
func (n *Node) SetDTime(t *time.Time) (err error) {
	if t == nil {
		err = xattrs.Remove(n.InternalPath(), xattrs.DTimeAttr)
		if xattrs.IsAttrUnset(err) {
			return nil
		}
		return err
	}
	return xattrs.Set(n.InternalPath(), xattrs.DTimeAttr, t.UTC().Format(time.RFC3339Nano))
}

// IsDisabled returns true when the node has a dmtime attribute set
// only used to check if a space is disabled
// FIXME confusing with the trash logic
func (n *Node) IsDisabled() bool {
	if _, err := n.GetDTime(); err == nil {
		return true
	}
	return false
}

// GetTreeSize reads the treesize from the extended attributes
func (n *Node) GetTreeSize() (treesize uint64, err error) {
	var b string
	if b, err = xattrs.Get(n.InternalPath(), xattrs.TreesizeAttr); err != nil {
		return
	}
	return strconv.ParseUint(b, 10, 64)
}

// SetTreeSize writes the treesize to the extended attributes
func (n *Node) SetTreeSize(ts uint64) (err error) {
	return n.SetMetadata(xattrs.TreesizeAttr, strconv.FormatUint(ts, 10))
}

// SetChecksum writes the checksum with the given checksum type to the extended attributes
func (n *Node) SetChecksum(csType string, h hash.Hash) (err error) {
	return n.SetMetadata(xattrs.ChecksumPrefix+csType, string(h.Sum(nil)))
}

// UnsetTempEtag removes the temporary etag attribute
func (n *Node) UnsetTempEtag() (err error) {
	err = xattrs.Remove(n.InternalPath(), xattrs.TmpEtagAttr)
	if xattrs.IsAttrUnset(err) {
		return nil
	}
	return err
}

// ReadUserPermissions will assemble the permissions for the current user on the given node without parent nodes
func (n *Node) ReadUserPermissions(ctx context.Context, u *userpb.User) (ap provider.ResourcePermissions, err error) {
	// check if the current user is the owner
	if utils.UserEqual(u.Id, n.Owner()) {
		appctx.GetLogger(ctx).Debug().Str("node", n.ID).Msg("user is owner, returning owner permissions")
		return OwnerPermissions(), nil
	}

	ap = provider.ResourcePermissions{}

	// for an efficient group lookup convert the list of groups to a map
	// groups are just strings ... groupnames ... or group ids ??? AAARGH !!!
	groupsMap := make(map[string]bool, len(u.Groups))
	for i := range u.Groups {
		groupsMap[u.Groups[i]] = true
	}

	var g *provider.Grant

	// we read all grantees from the node
	var grantees []string
	if grantees, err = n.ListGrantees(ctx); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Interface("node", n).Msg("error listing grantees")
		return NoPermissions(), err
	}

	// instead of making n getxattr syscalls we are going to list the acls and filter them here
	// we have two options here:
	// 1. we can start iterating over the acls / grants on the node or
	// 2. we can iterate over the number of groups
	// The current implementation tries to be defensive for cases where users have hundreds or thousands of groups, so we iterate over the existing acls.
	userace := xattrs.GrantUserAcePrefix + u.Id.OpaqueId
	userFound := false
	for i := range grantees {
		switch {
		// we only need to find the user once
		case !userFound && grantees[i] == userace:
			g, err = n.ReadGrant(ctx, grantees[i])
		case strings.HasPrefix(grantees[i], xattrs.GrantGroupAcePrefix): // only check group grantees
			gr := strings.TrimPrefix(grantees[i], xattrs.GrantGroupAcePrefix)
			if groupsMap[gr] {
				g, err = n.ReadGrant(ctx, grantees[i])
			} else {
				// no need to check attribute
				continue
			}
		default:
			// no need to check attribute
			continue
		}

		switch {
		case err == nil:
			AddPermissions(&ap, g.GetPermissions())
		case xattrs.IsAttrUnset(err):
			err = nil
			appctx.GetLogger(ctx).Error().Interface("node", n).Str("grant", grantees[i]).Interface("grantees", grantees).Msg("grant vanished from node after listing")
			// continue with next segment
		default:
			appctx.GetLogger(ctx).Error().Err(err).Interface("node", n).Str("grant", grantees[i]).Msg("error reading permissions")
			// continue with next segment
		}
	}

	appctx.GetLogger(ctx).Debug().Interface("permissions", ap).Interface("node", n).Interface("user", u).Msg("returning aggregated permissions")
	return ap, nil
}

// ListGrantees lists the grantees of the current node
// We don't want to wast time and memory by creating grantee objects.
// The function will return a list of opaque strings that can be used to make a ReadGrant call
func (n *Node) ListGrantees(ctx context.Context) (grantees []string, err error) {
	var attrs []string

	if attrs, err = xattr.List(n.InternalPath()); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("node", n.ID).Msg("error listing attributes")
		return nil, err
	}
	for i := range attrs {
		if strings.HasPrefix(attrs[i], xattrs.GrantPrefix) {
			grantees = append(grantees, attrs[i])
		}
	}
	return
}

// ReadGrant reads a CS3 grant
func (n *Node) ReadGrant(ctx context.Context, grantee string) (g *provider.Grant, err error) {
	var b string
	if b, err = xattrs.Get(n.InternalPath(), grantee); err != nil {
		return nil, err
	}
	var e *ace.ACE
	if e, err = ace.Unmarshal(strings.TrimPrefix(grantee, xattrs.GrantPrefix), []byte(b)); err != nil {
		return nil, err
	}
	return e.Grant(), nil
}

// ListGrants lists all grants of the current node.
func (n *Node) ListGrants(ctx context.Context) ([]*provider.Grant, error) {
	grantees, err := n.ListGrantees(ctx)
	if err != nil {
		return nil, err
	}

	grants := make([]*provider.Grant, 0, len(grantees))
	for _, g := range grantees {
		grant, err := n.ReadGrant(ctx, g)
		if err != nil {
			appctx.GetLogger(ctx).
				Error().
				Err(err).
				Str("node", n.ID).
				Str("grantee", g).
				Msg("error reading grant")
			continue
		}
		grants = append(grants, grant)
	}
	return grants, nil
}

// ReadBlobSizeAttr reads the blobsize from the xattrs
func ReadBlobSizeAttr(path string) (int64, error) {
	attr, err := xattrs.Get(path, xattrs.BlobsizeAttr)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading blobsize xattr")
	}
	blobSize, err := strconv.ParseInt(attr, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid blobsize xattr format")
	}
	return blobSize, nil
}

// ReadBlobIDAttr reads the blobsize from the xattrs
func ReadBlobIDAttr(path string) (string, error) {
	attr, err := xattrs.Get(path, xattrs.BlobIDAttr)
	if err != nil {
		return "", errors.Wrapf(err, "error reading blobid xattr")
	}
	return attr, nil
}

func (n *Node) hasUserShares(ctx context.Context) bool {
	g, err := n.ListGrantees(ctx)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("hasUserShares: listGrantees")
		return false
	}

	for i := range g {
		if strings.HasPrefix(g[i], xattrs.GrantUserAcePrefix) {
			return true
		}
	}
	return false
}

func parseMTime(v string) (t time.Time, err error) {
	p := strings.SplitN(v, ".", 2)
	var sec, nsec int64
	if sec, err = strconv.ParseInt(p[0], 10, 64); err == nil {
		if len(p) > 1 {
			nsec, err = strconv.ParseInt(p[1], 10, 64)
		}
	}
	return time.Unix(sec, nsec), err
}

// FindStorageSpaceRoot calls n.Parent() and climbs the tree
// until it finds the space root node and adds it to the node
func (n *Node) FindStorageSpaceRoot() error {
	if n.SpaceRoot != nil {
		return nil
	}
	var err error
	// remember the node we ask for and use parent to climb the tree
	parent := n
	for {
		if IsSpaceRoot(parent) {
			n.SpaceRoot = parent
			break
		}
		if parent, err = parent.Parent(); err != nil {
			return err
		}
	}
	return nil
}

// IsSpaceRoot checks if the node is a space root
func IsSpaceRoot(r *Node) bool {
	path := r.InternalPath()
	if _, err := xattrs.Get(path, xattrs.SpaceNameAttr); err == nil {
		return true
	}
	return false
}

// CheckQuota checks if both disk space and available quota are sufficient
// Overwrite must be set to true if the new file replaces the old file e.g.
// when creating a new file version. In such a case the function will
// reduce the used bytes by the old file size and then add the new size.
// If overwrite is false oldSize will be ignored.
var CheckQuota = func(spaceRoot *Node, overwrite bool, oldSize, newSize uint64) (quotaSufficient bool, err error) {
	used, _ := spaceRoot.GetTreeSize()
	if !enoughDiskSpace(spaceRoot.InternalPath(), newSize) {
		return false, errtypes.InsufficientStorage("disk full")
	}
	quotaByteStr, _ := xattrs.Get(spaceRoot.InternalPath(), xattrs.QuotaAttr)
	if quotaByteStr == "" || quotaByteStr == QuotaUnlimited {
		// if quota is not set, it means unlimited
		return true, nil
	}
	quotaByte, _ := strconv.ParseUint(quotaByteStr, 10, 64)
	if overwrite {
		if quotaByte < used-oldSize+newSize {
			return false, errtypes.InsufficientStorage("quota exceeded")
		}
		// if total is smaller than used, total-used could overflow and be bigger than fileSize
	} else if newSize > quotaByte-used || quotaByte < used {
		return false, errtypes.InsufficientStorage("quota exceeded")
	}
	return true, nil
}

func enoughDiskSpace(path string, fileSize uint64) bool {
	avalB, err := GetAvailableSize(path)
	if err != nil {
		return false
	}
	return avalB > fileSize
}
