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

package lookup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

const (
	_spaceTypePersonal = "personal"
)

func init() {
	tracer = otel.Tracer("github.com/owncloud/reva/pkg/storage/utils/decomposedfs/lookup")
}

// Lookup implements transformations from filepath to node and back
type Lookup struct {
	Options *options.Options

	metadataBackend metadata.Backend
	tm              node.TimeManager
}

// New returns a new Lookup instance
func New(b metadata.Backend, o *options.Options, tm node.TimeManager) *Lookup {
	return &Lookup{
		Options:         o,
		metadataBackend: b,
		tm:              tm,
	}
}

// MetadataBackend returns the metadata backend
func (lu *Lookup) MetadataBackend() metadata.Backend {
	return lu.metadataBackend
}

func (lu *Lookup) ReadBlobIDAndSizeAttr(ctx context.Context, path string, attrs node.Attributes) (string, int64, error) {
	blobID := ""
	blobSize := int64(0)
	var err error

	if attrs != nil {
		blobID = attrs.String(prefixes.BlobIDAttr)
		if blobID != "" {
			blobSize, err = attrs.Int64(prefixes.BlobsizeAttr)
			if err != nil {
				return "", 0, err
			}
		}
	} else {
		attrs, err := lu.metadataBackend.All(ctx, path)
		if err != nil {
			return "", 0, errors.Wrapf(err, "error reading blobid xattr")
		}
		nodeAttrs := node.Attributes(attrs)
		blobID = nodeAttrs.String(prefixes.BlobIDAttr)
		blobSize, err = nodeAttrs.Int64(prefixes.BlobsizeAttr)
		if err != nil {
			return "", 0, errors.Wrapf(err, "error reading blobsize xattr")
		}
	}
	return blobID, blobSize, nil
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

func (lu *Lookup) NodeIDFromParentAndName(ctx context.Context, parent *node.Node, name string) (string, error) {
	nodeID, err := readChildNodeFromLink(filepath.Join(parent.InternalPath(), name))
	if err != nil {
		return "", errors.Wrap(err, "decomposedfs: Wrap: readlink error")
	}
	return nodeID, nil
}

// TypeFromPath returns the type of the node at the given path
func (lu *Lookup) TypeFromPath(ctx context.Context, path string) provider.ResourceType {
	// Try to read from xattrs
	typeAttr, err := lu.metadataBackend.GetInt64(ctx, path, prefixes.TypeAttr)
	if err == nil {
		return provider.ResourceType(int32(typeAttr))
	}

	t := provider.ResourceType_RESOURCE_TYPE_INVALID
	// Fall back to checking on disk
	fi, err := os.Lstat(path)
	if err != nil {
		return t
	}

	switch {
	case fi.IsDir():
		if _, err = lu.metadataBackend.Get(ctx, path, prefixes.ReferenceAttr); err == nil {
			t = provider.ResourceType_RESOURCE_TYPE_REFERENCE
		} else {
			t = provider.ResourceType_RESOURCE_TYPE_CONTAINER
		}
	case fi.Mode().IsRegular():
		t = provider.ResourceType_RESOURCE_TYPE_FILE
	case fi.Mode()&os.ModeSymlink != 0:
		t = provider.ResourceType_RESOURCE_TYPE_SYMLINK
		// TODO reference using ext attr on a symlink
		// nodeType = provider.ResourceType_RESOURCE_TYPE_REFERENCE
	}
	return t
}

// NodeFromResource takes in a request path or request id and converts it to a Node
func (lu *Lookup) NodeFromResource(ctx context.Context, ref *provider.Reference) (*node.Node, error) {
	ctx, span := tracer.Start(ctx, "NodeFromResource")
	defer span.End()

	if ref.ResourceId != nil {
		// check if a storage space reference is used
		// currently, the decomposed fs uses the root node id as the space id
		n, err := lu.NodeFromID(ctx, ref.ResourceId)
		if err != nil {
			return nil, err
		}
		// is this a relative reference?
		if ref.Path != "" {
			p := filepath.Clean(ref.Path)
			if p != "." && p != "/" {
				// walk the relative path
				n, err = lu.WalkPath(ctx, n, p, false, func(ctx context.Context, n *node.Node) error { return nil })
				if err != nil {
					return nil, err
				}
				n.SpaceID = ref.ResourceId.SpaceId
			}
		}
		return n, nil
	}

	// reference is invalid
	return nil, fmt.Errorf("invalid reference %+v. resource_id must be set", ref)
}

// NodeFromID returns the internal path for the id
func (lu *Lookup) NodeFromID(ctx context.Context, id *provider.ResourceId) (n *node.Node, err error) {
	ctx, span := tracer.Start(ctx, "NodeFromID")
	defer span.End()
	if id == nil {
		return nil, fmt.Errorf("invalid resource id %+v", id)
	}
	if id.OpaqueId == "" {
		// The Resource references the root of a space
		return lu.NodeFromSpaceID(ctx, id.SpaceId)
	}
	return node.ReadNode(ctx, lu, id.SpaceId, id.OpaqueId, false, nil, false)
}

// Pathify segments the beginning of a string into depth segments of width length
// Pathify("aabbccdd", 3, 1) will return "a/a/b/bccdd"
func Pathify(id string, depth, width int) string {
	b := strings.Builder{}
	i := 0
	for ; i < depth; i++ {
		if len(id) <= i*width+width {
			break
		}
		b.WriteString(id[i*width : i*width+width])
		b.WriteRune(filepath.Separator)
	}
	b.WriteString(id[i*width:])
	return b.String()
}

// NodeFromSpaceID converts a resource id into a Node
func (lu *Lookup) NodeFromSpaceID(ctx context.Context, spaceID string) (n *node.Node, err error) {
	node, err := node.ReadNode(ctx, lu, spaceID, spaceID, false, nil, false)
	if err != nil {
		return nil, err
	}

	node.SpaceRoot = node
	return node, nil
}

// GenerateSpaceID generates a new space id and alias
func (lu *Lookup) GenerateSpaceID(spaceType string, owner *user.User) (string, error) {
	switch spaceType {
	case _spaceTypePersonal:
		return owner.Id.OpaqueId, nil
	default:
		return uuid.New().String(), nil
	}
}

// Path returns the path for node
func (lu *Lookup) Path(ctx context.Context, n *node.Node, hasPermission node.PermissionFunc) (p string, err error) {
	root := n.SpaceRoot
	var child *node.Node
	for n.ID != root.ID {
		p = filepath.Join(n.Name, p)
		child = n
		if n, err = n.Parent(ctx); err != nil {
			appctx.GetLogger(ctx).
				Error().Err(err).
				Str("path", p).
				Str("spaceid", child.SpaceID).
				Str("nodeid", child.ID).
				Str("parentid", child.ParentID).
				Msg("Path()")
			return
		}

		if !hasPermission(n) {
			break
		}
	}
	p = filepath.Join("/", p)
	return
}

// WalkPath calls n.Child(segment) on every path segment in p starting at the node r.
// If a function f is given it will be executed for every segment node, but not the root node r.
// If followReferences is given the current visited reference node is replaced by the referenced node.
func (lu *Lookup) WalkPath(ctx context.Context, r *node.Node, p string, followReferences bool, f func(ctx context.Context, n *node.Node) error) (*node.Node, error) {
	segments := strings.Split(strings.Trim(p, "/"), "/")
	var err error
	for i := range segments {
		if r, err = r.Child(ctx, segments[i]); err != nil {
			return r, err
		}

		if followReferences {
			if attrBytes, err := r.Xattr(ctx, prefixes.ReferenceAttr); err == nil {
				realNodeID := attrBytes
				ref, err := refFromCS3(realNodeID)
				if err != nil {
					return nil, err
				}

				r, err = lu.NodeFromID(ctx, ref.ResourceId)
				if err != nil {
					return nil, err
				}
			}
		}
		if r.IsSpaceRoot(ctx) {
			r.SpaceRoot = r
		}

		if !r.Exists && i < len(segments)-1 {
			return r, errtypes.NotFound(segments[i])
		}
		if f != nil {
			if err = f(ctx, r); err != nil {
				return r, err
			}
		}
	}
	return r, nil
}

// InternalRoot returns the internal storage root directory
func (lu *Lookup) InternalRoot() string {
	return lu.Options.Root
}

// InternalPath returns the internal path for a given ID
func (lu *Lookup) InternalPath(spaceID, nodeID string) string {
	return filepath.Join(lu.Options.Root, "spaces", Pathify(spaceID, 1, 2), "nodes", Pathify(nodeID, 4, 2))
}

// // ReferenceFromAttr returns a CS3 reference from xattr of a node.
// // Supported formats are: "cs3:storageid/nodeid"
// func ReferenceFromAttr(b []byte) (*provider.Reference, error) {
// 	return refFromCS3(b)
// }

// refFromCS3 creates a CS3 reference from a set of bytes. This method should remain private
// and only be called after validation because it can potentially panic.
func refFromCS3(b []byte) (*provider.Reference, error) {
	parts := string(b[4:])
	return &provider.Reference{
		ResourceId: &provider.ResourceId{
			StorageId: strings.Split(parts, "/")[0],
			OpaqueId:  strings.Split(parts, "/")[1],
		},
	}, nil
}

// CopyMetadata copies all extended attributes from source to target.
// The optional filter function can be used to filter by attribute name, e.g. by checking a prefix
// For the source file, a shared lock is acquired.
// NOTE: target resource will be write locked!
func (lu *Lookup) CopyMetadata(ctx context.Context, src, target string, filter func(attributeName string, value []byte) (newValue []byte, copy bool), acquireTargetLock bool) (err error) {
	// Acquire a read log on the source node
	// write lock existing node before reading treesize or tree time
	lock, err := lockedfile.OpenFile(lu.MetadataBackend().LockfilePath(src), os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if err != nil {
		return errors.Wrap(err, "xattrs: Unable to lock source to read")
	}
	defer func() {
		rerr := lock.Close()

		// if err is non nil we do not overwrite that
		if err == nil {
			err = rerr
		}
	}()

	return lu.CopyMetadataWithSourceLock(ctx, src, target, filter, lock, acquireTargetLock)
}

// CopyMetadataWithSourceLock copies all extended attributes from source to target.
// The optional filter function can be used to filter by attribute name, e.g. by checking a prefix
// For the source file, a matching lockedfile is required.
// NOTE: target resource will be write locked!
func (lu *Lookup) CopyMetadataWithSourceLock(ctx context.Context, sourcePath, targetPath string, filter func(attributeName string, value []byte) (newValue []byte, copy bool), lockedSource *lockedfile.File, acquireTargetLock bool) (err error) {
	switch {
	case lockedSource == nil:
		return errors.New("no lock provided")
	case lockedSource.File.Name() != lu.MetadataBackend().LockfilePath(sourcePath):
		return errors.New("lockpath does not match filepath")
	}

	attrs, err := lu.metadataBackend.All(ctx, sourcePath)
	if err != nil {
		return err
	}

	newAttrs := make(map[string][]byte, 0)
	for attrName, val := range attrs {
		if filter != nil {
			var ok bool
			if val, ok = filter(attrName, val); !ok {
				continue
			}
		}
		newAttrs[attrName] = val
	}

	return lu.MetadataBackend().SetMultiple(ctx, targetPath, newAttrs, acquireTargetLock)
}

// TimeManager returns the time manager
func (lu *Lookup) TimeManager() node.TimeManager {
	return lu.tm
}

// DetectBackendOnDisk returns the name of the metadata backend being used on disk
func DetectBackendOnDisk(root string) string {
	matches, _ := filepath.Glob(filepath.Join(root, "spaces", "*", "*"))
	if len(matches) > 0 {
		base := matches[len(matches)-1]
		spaceid := strings.ReplaceAll(
			strings.TrimPrefix(base, filepath.Join(root, "spaces")),
			"/", "")
		spaceRoot := Pathify(spaceid, 4, 2)
		_, err := os.Stat(filepath.Join(base, "nodes", spaceRoot+".mpk"))
		if err == nil {
			return "mpk"
		}
		_, err = os.Stat(filepath.Join(base, "nodes", spaceRoot+".ini"))
		if err == nil {
			return "ini"
		}
	}
	return "xattrs"
}
