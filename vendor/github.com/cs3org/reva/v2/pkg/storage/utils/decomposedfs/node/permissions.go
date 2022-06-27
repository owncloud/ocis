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
	"strings"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

// NoPermissions represents an empty set of permissions
func NoPermissions() provider.ResourcePermissions {
	return provider.ResourcePermissions{}
}

// NoOwnerPermissions defines permissions for nodes that don't have an owner set, eg the root node
func NoOwnerPermissions() provider.ResourcePermissions {
	return provider.ResourcePermissions{
		Stat: true,
	}
}

// ShareFolderPermissions defines permissions for the shared jail
func ShareFolderPermissions() provider.ResourcePermissions {
	return provider.ResourcePermissions{
		// read permissions
		ListContainer:        true,
		Stat:                 true,
		InitiateFileDownload: true,
		GetPath:              true,
		GetQuota:             true,
		ListFileVersions:     true,
	}
}

// OwnerPermissions defines permissions for nodes owned by the user
func OwnerPermissions() provider.ResourcePermissions {
	return provider.ResourcePermissions{
		// all permissions
		AddGrant:             true,
		CreateContainer:      true,
		Delete:               true,
		GetPath:              true,
		GetQuota:             true,
		InitiateFileDownload: true,
		InitiateFileUpload:   true,
		ListContainer:        true,
		ListFileVersions:     true,
		ListGrants:           true,
		ListRecycle:          true,
		Move:                 true,
		PurgeRecycle:         true,
		RemoveGrant:          true,
		RestoreFileVersion:   true,
		RestoreRecycleItem:   true,
		Stat:                 true,
		UpdateGrant:          true,
	}
}

// Permissions implements permission checks
type Permissions struct {
	lu PathLookup
}

// NewPermissions returns a new Permissions instance
func NewPermissions(lu PathLookup) *Permissions {
	return &Permissions{
		lu: lu,
	}
}

// AssemblePermissions will assemble the permissions for the current user on the given node, taking into account all parent nodes
func (p *Permissions) AssemblePermissions(ctx context.Context, n *Node) (ap provider.ResourcePermissions, err error) {
	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		appctx.GetLogger(ctx).Debug().Interface("node", n.ID).Msg("no user in context, returning default permissions")
		return NoPermissions(), nil
	}

	// check if the current user is the owner
	if utils.UserEqual(u.Id, n.Owner()) {
		lp, err := n.lu.Path(ctx, n)
		if err == nil && lp == n.lu.ShareFolder() {
			return ShareFolderPermissions(), nil
		}
		appctx.GetLogger(ctx).Debug().Str("node", n.ID).Msg("user is owner, returning owner permissions")
		return OwnerPermissions(), nil
	}
	// determine root

	rn := n.SpaceRoot

	cn := n

	ap = provider.ResourcePermissions{}

	// for an efficient group lookup convert the list of groups to a map
	// groups are just strings ... groupnames ... or group ids ??? AAARGH !!!
	groupsMap := make(map[string]bool, len(u.Groups))
	for i := range u.Groups {
		groupsMap[u.Groups[i]] = true
	}

	// for all segments, starting at the leaf
	for cn.ID != rn.ID {
		if np, err := cn.ReadUserPermissions(ctx, u); err == nil {
			AddPermissions(&ap, &np)
		} else {
			appctx.GetLogger(ctx).Error().Err(err).Interface("node", cn.ID).Msg("error reading permissions")
			// continue with next segment
		}
		if cn, err = cn.Parent(); err != nil {
			return ap, errors.Wrap(err, "Decomposedfs: error getting parent "+cn.ParentID)
		}
	}

	// for the root node
	if np, err := cn.ReadUserPermissions(ctx, u); err == nil {
		AddPermissions(&ap, &np)
	} else {
		appctx.GetLogger(ctx).Error().Err(err).Interface("node", cn.ID).Msg("error reading root node permissions")
	}

	appctx.GetLogger(ctx).Debug().Interface("permissions", ap).Interface("node", n.ID).Interface("user", u).Msg("returning agregated permissions")
	return ap, nil
}

// AddPermissions merges a set of permissions into another
// TODO we should use a bitfield for this ...
func AddPermissions(l *provider.ResourcePermissions, r *provider.ResourcePermissions) {
	l.AddGrant = l.AddGrant || r.AddGrant
	l.CreateContainer = l.CreateContainer || r.CreateContainer
	l.Delete = l.Delete || r.Delete
	l.GetPath = l.GetPath || r.GetPath
	l.GetQuota = l.GetQuota || r.GetQuota
	l.InitiateFileDownload = l.InitiateFileDownload || r.InitiateFileDownload
	l.InitiateFileUpload = l.InitiateFileUpload || r.InitiateFileUpload
	l.ListContainer = l.ListContainer || r.ListContainer
	l.ListFileVersions = l.ListFileVersions || r.ListFileVersions
	l.ListGrants = l.ListGrants || r.ListGrants
	l.ListRecycle = l.ListRecycle || r.ListRecycle
	l.Move = l.Move || r.Move
	l.PurgeRecycle = l.PurgeRecycle || r.PurgeRecycle
	l.RemoveGrant = l.RemoveGrant || r.RemoveGrant
	l.RestoreFileVersion = l.RestoreFileVersion || r.RestoreFileVersion
	l.RestoreRecycleItem = l.RestoreRecycleItem || r.RestoreRecycleItem
	l.Stat = l.Stat || r.Stat
	l.UpdateGrant = l.UpdateGrant || r.UpdateGrant
}

// HasPermission call check() for every node up to the root until check returns true
func (p *Permissions) HasPermission(ctx context.Context, n *Node, check func(*provider.ResourcePermissions) bool) (can bool, err error) {

	var u *userv1beta1.User
	var perms *provider.ResourcePermissions
	if u, perms = p.getUserAndPermissions(ctx, n); perms != nil {
		return check(perms), nil
	}

	// determine root
	/*
		if err = n.FindStorageSpaceRoot(); err != nil {
			return false, err
		}
	*/

	// for an efficient group lookup convert the list of groups to a map
	// groups are just strings ... groupnames ... or group ids ??? AAARGH !!!
	groupsMap := make(map[string]bool, len(u.Groups))
	for i := range u.Groups {
		groupsMap[u.Groups[i]] = true
	}

	// for all segments, starting at the leaf
	cn := n
	for cn.ID != n.SpaceRoot.ID {
		if ok := nodeHasPermission(ctx, cn, groupsMap, u.Id.OpaqueId, check); ok {
			return true, nil
		}

		if cn, err = cn.Parent(); err != nil {
			return false, errors.Wrap(err, "Decomposedfs: error getting parent "+cn.ParentID)
		}
	}

	// also check permissions on root, eg. for for project spaces
	return nodeHasPermission(ctx, cn, groupsMap, u.Id.OpaqueId, check), nil
}

func nodeHasPermission(ctx context.Context, cn *Node, groupsMap map[string]bool, userid string, check func(*provider.ResourcePermissions) bool) (ok bool) {

	var grantees []string
	var err error
	if grantees, err = cn.ListGrantees(ctx); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Interface("node", cn.ID).Msg("error listing grantees")
		return false
	}

	userace := xattrs.GrantUserAcePrefix + userid
	userFound := false
	for i := range grantees {
		// we only need the find the user once per node
		var g *provider.Grant
		switch {
		case !userFound && grantees[i] == userace:
			g, err = cn.ReadGrant(ctx, grantees[i])
		case strings.HasPrefix(grantees[i], xattrs.GrantGroupAcePrefix):
			gr := strings.TrimPrefix(grantees[i], xattrs.GrantGroupAcePrefix)
			if groupsMap[gr] {
				g, err = cn.ReadGrant(ctx, grantees[i])
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
			appctx.GetLogger(ctx).Debug().Interface("node", cn.ID).Str("grant", grantees[i]).Interface("permissions", g.GetPermissions()).Msg("checking permissions")
			if check(g.GetPermissions()) {
				return true
			}
		case xattrs.IsAttrUnset(err):
			appctx.GetLogger(ctx).Error().Interface("node", cn.ID).Str("grant", grantees[i]).Interface("grantees", grantees).Msg("grant vanished from node after listing")
		default:
			appctx.GetLogger(ctx).Error().Err(err).Interface("node", cn.ID).Str("grant", grantees[i]).Msg("error reading permissions")
			return false
		}
	}

	return false
}

func (p *Permissions) getUserAndPermissions(ctx context.Context, n *Node) (*userv1beta1.User, *provider.ResourcePermissions) {
	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		appctx.GetLogger(ctx).Debug().Interface("node", n.ID).Msg("no user in context, returning default permissions")
		perms := NoPermissions()
		return nil, &perms
	}
	// check if the current user is the owner
	if utils.UserEqual(u.Id, n.Owner()) {
		appctx.GetLogger(ctx).Debug().Str("node", n.ID).Msg("user is owner, returning owner permissions")
		perms := OwnerPermissions()
		return u, &perms
	}
	return u, nil
}
