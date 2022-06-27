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

package decomposedfs

import (
	"context"
	"path/filepath"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/ace"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/xattr"
)

// DenyGrant denies access to a resource.
func (fs *Decomposedfs) DenyGrant(ctx context.Context, ref *provider.Reference, g *provider.Grantee) error {
	return errtypes.NotSupported("decomposedfs: not supported")
}

// AddGrant adds a grant to a resource
func (fs *Decomposedfs) AddGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) (err error) {
	log := appctx.GetLogger(ctx)
	log.Debug().Interface("ref", ref).Interface("grant", g).Msg("AddGrant()")
	node, grant, err := fs.loadGrant(ctx, ref, g)
	if err != nil {
		return err
	}

	if grant != nil {
		// grant exists -> go to UpdateGrant
		// TODO: should we hard error in this case?
		return fs.UpdateGrant(ctx, ref, g)
	}

	// add acting user
	u := ctxpkg.ContextMustGetUser(ctx)
	g.Creator = u.GetId()

	owner := node.Owner()
	grants, err := node.ListGrants(ctx)
	if err != nil {
		return err
	}

	// If the owner is empty and there are no grantees then we are dealing with a just created project space.
	// In this case we don't need to check for permissions and just add the grant since this will be the project
	// manager.
	// When the owner is empty but grants are set then we do want to check the grants.
	// However, if we are trying to edit an existing grant we do not have to check for permission if the user owns the grant
	// TODO: find a better to check this
	if !(len(grants) == 0 && (owner == nil || owner.OpaqueId == "")) {
		ok, err := fs.p.HasPermission(ctx, node, func(rp *provider.ResourcePermissions) bool {
			return rp.AddGrant
		})
		switch {
		case err != nil:
			return errtypes.InternalError(err.Error())
		case !ok:
			return errtypes.PermissionDenied(filepath.Join(node.ParentID, node.Name))
		}
	}

	return fs.storeGrant(ctx, node, g)
}

// ListGrants lists the grants on the specified resource
func (fs *Decomposedfs) ListGrants(ctx context.Context, ref *provider.Reference) (grants []*provider.Grant, err error) {
	var node *node.Node
	if node, err = fs.lu.NodeFromResource(ctx, ref); err != nil {
		return
	}
	if !node.Exists {
		err = errtypes.NotFound(filepath.Join(node.ParentID, node.Name))
		return
	}

	listGrants, err := fs.p.HasPermission(ctx, node, func(rp *provider.ResourcePermissions) bool {
		return rp.ListGrants
	})
	if err != nil {
		return nil, errtypes.InternalError(err.Error())
	}

	log := appctx.GetLogger(ctx)
	np := node.InternalPath()
	var attrs []string
	if attrs, err = xattr.List(np); err != nil {
		log.Error().Err(err).Msg("error listing attributes")
		return nil, err
	}

	log.Debug().Interface("attrs", attrs).Msg("read attributes")

	aces := extractACEsFromAttrs(ctx, np, attrs)

	uid := ctxpkg.ContextMustGetUser(ctx).GetId()
	grants = make([]*provider.Grant, 0, len(aces))
	for i := range aces {
		g := aces[i].Grant()

		// you may list your own grants even without listgrants permission
		if !listGrants && !utils.UserEqual(g.Creator, uid) && !utils.UserEqual(g.Grantee.GetUserId(), uid) {
			continue
		}

		grants = append(grants, g)
	}

	return grants, nil
}

// RemoveGrant removes a grant from resource
func (fs *Decomposedfs) RemoveGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) (err error) {
	node, grant, err := fs.loadGrant(ctx, ref, g)
	if err != nil {
		return err
	}

	// you are allowed to remove grants if you created them yourself or have the proper permission
	if !utils.UserEqual(grant.Creator, ctxpkg.ContextMustGetUser(ctx).GetId()) {
		ok, err := fs.p.HasPermission(ctx, node, func(rp *provider.ResourcePermissions) bool {
			return rp.RemoveGrant
		})
		switch {
		case err != nil:
			return errtypes.InternalError(err.Error())
		case !ok:
			return errtypes.PermissionDenied(filepath.Join(node.ParentID, node.Name))
		}
	}

	// check lock
	if err := node.CheckLock(ctx); err != nil {
		return err
	}

	var attr string
	if g.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_GROUP {
		attr = xattrs.GrantGroupAcePrefix + g.Grantee.GetGroupId().OpaqueId
	} else {
		attr = xattrs.GrantUserAcePrefix + g.Grantee.GetUserId().OpaqueId
	}

	if err = xattrs.Remove(node.InternalPath(), attr); err != nil {
		return
	}

	return fs.tp.Propagate(ctx, node)
}

// UpdateGrant updates a grant on a resource
// TODO remove AddGrant or UpdateGrant grant from CS3 api, redundant? tracked in https://github.com/cs3org/cs3apis/issues/92
func (fs *Decomposedfs) UpdateGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	log := appctx.GetLogger(ctx)
	log.Debug().Interface("ref", ref).Interface("grant", g).Msg("UpdateGrant()")

	node, grant, err := fs.loadGrant(ctx, ref, g)
	if err != nil {
		return err
	}

	if grant == nil {
		// grant not found
		// TODO: fallback to AddGrant?
		return errtypes.NotFound(g.Grantee.GetUserId().GetOpaqueId())
	}

	// You may update a grant when you have the UpdateGrant permission or created the grant (regardless what your permissions are now)
	if !utils.UserEqual(grant.Creator, ctxpkg.ContextMustGetUser(ctx).GetId()) {
		ok, err := fs.p.HasPermission(ctx, node, func(rp *provider.ResourcePermissions) bool {
			return rp.UpdateGrant
		})
		switch {
		case err != nil:
			return errtypes.InternalError(err.Error())
		case !ok:
			return errtypes.PermissionDenied(filepath.Join(node.ParentID, node.Name))
		}
	}

	return fs.storeGrant(ctx, node, g)
}

// checks if the given grant exists and returns it. Nil grant means it doesn't exist
func (fs *Decomposedfs) loadGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) (*node.Node, *provider.Grant, error) {
	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		return nil, nil, err
	}
	if !n.Exists {
		return nil, nil, errtypes.NotFound(filepath.Join(n.ParentID, n.Name))
	}

	grants, err := n.ListGrants(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, grant := range grants {
		switch grant.Grantee.GetType() {
		case provider.GranteeType_GRANTEE_TYPE_USER:
			if g.Grantee.GetUserId().GetOpaqueId() == grant.Grantee.GetUserId().GetOpaqueId() {
				return n, grant, nil
			}
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			if g.Grantee.GetGroupId().GetOpaqueId() == grant.Grantee.GetGroupId().GetOpaqueId() {
				return n, grant, nil
			}
		}
	}

	return n, nil, nil
}

func (fs *Decomposedfs) storeGrant(ctx context.Context, n *node.Node, g *provider.Grant) error {
	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return err
	}

	e := ace.FromGrant(g)
	principal, value := e.Marshal()
	if err := n.SetMetadata(xattrs.GrantPrefix+principal, string(value)); err != nil {
		return err
	}

	// when a grant is added to a space, do not add a new space under "shares"
	if spaceGrant := ctx.Value(utils.SpaceGrant); spaceGrant == nil {
		err := fs.linkStorageSpaceType(ctx, spaceTypeShare, n.ID)
		if err != nil {
			return err
		}
	}

	return fs.tp.Propagate(ctx, n)
}

// extractACEsFromAttrs reads ACEs in the list of attrs from the node
func extractACEsFromAttrs(ctx context.Context, fsfn string, attrs []string) (entries []*ace.ACE) {
	log := appctx.GetLogger(ctx)
	entries = []*ace.ACE{}
	for i := range attrs {
		if strings.HasPrefix(attrs[i], xattrs.GrantPrefix) {
			var value string
			var err error
			if value, err = xattrs.Get(fsfn, attrs[i]); err != nil {
				log.Error().Err(err).Str("attr", attrs[i]).Msg("could not read attribute")
				continue
			}
			var e *ace.ACE
			principal := attrs[i][len(xattrs.GrantPrefix):]
			if e, err = ace.Unmarshal(principal, []byte(value)); err != nil {
				log.Error().Err(err).Str("principal", principal).Str("attr", attrs[i]).Msg("could not unmarshal ace")
				continue
			}
			entries = append(entries, e)
		}
	}
	return
}
