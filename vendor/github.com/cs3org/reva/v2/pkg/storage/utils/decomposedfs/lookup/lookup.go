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
	"path/filepath"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
)

// Lookup implements transformations from filepath to node and back
type Lookup struct {
	Options *options.Options
}

// NodeFromResource takes in a request path or request id and converts it to a Node
func (lu *Lookup) NodeFromResource(ctx context.Context, ref *provider.Reference) (*node.Node, error) {
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
				n.SpaceID = ref.ResourceId.StorageId
			}
		}
		return n, nil
	}

	// reference is invalid
	return nil, fmt.Errorf("invalid reference %+v. resource_id must be set", ref)
}

// NodeFromID returns the internal path for the id
func (lu *Lookup) NodeFromID(ctx context.Context, id *provider.ResourceId) (n *node.Node, err error) {
	if id == nil {
		return nil, fmt.Errorf("invalid resource id %+v", id)
	}
	if id.OpaqueId == "" {
		// The Resource references the root of a space
		return lu.NodeFromSpaceID(ctx, id)
	}
	return node.ReadNode(ctx, lu, id.StorageId, id.OpaqueId, false)
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

// NodeFromSpaceID converts a resource id without an opaque id into a Node
func (lu *Lookup) NodeFromSpaceID(ctx context.Context, id *provider.ResourceId) (n *node.Node, err error) {
	node, err := node.ReadNode(ctx, lu, id.StorageId, id.StorageId, false)
	if err != nil {
		return nil, err
	}

	node.SpaceRoot = node
	return node, nil
}

// Path returns the path for node
func (lu *Lookup) Path(ctx context.Context, n *node.Node) (p string, err error) {
	root := n.SpaceRoot
	for n.ID != root.ID {
		p = filepath.Join(n.Name, p)
		if n, err = n.Parent(); err != nil {
			appctx.GetLogger(ctx).
				Error().Err(err).
				Str("path", p).
				Interface("node", n).
				Msg("Path()")
			return
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
			if attrBytes, err := r.GetMetadata(xattrs.ReferenceAttr); err == nil {
				realNodeID := attrBytes
				ref, err := xattrs.ReferenceFromAttr([]byte(realNodeID))
				if err != nil {
					return nil, err
				}

				r, err = lu.NodeFromID(ctx, ref.ResourceId)
				if err != nil {
					return nil, err
				}
			}
		}
		if node.IsSpaceRoot(r) {
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

// ShareFolder returns the internal storage root directory
func (lu *Lookup) ShareFolder() string {
	return lu.Options.ShareFolder
}
