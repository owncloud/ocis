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
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ocsconv "github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/storage/utils/templates"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	spaceTypePersonal = "personal"
	// spaceTypeProject  = "project"
	spaceTypeShare = "share"
	spaceTypeAny   = "*"
	spaceIDAny     = "*"
)

// CreateStorageSpace creates a storage space
func (fs *Decomposedfs) CreateStorageSpace(ctx context.Context, req *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {

	// "everything is a resource" this is the unique ID for the Space resource.
	spaceID := uuid.New().String()
	// allow sending a space id
	if reqSpaceID := utils.ReadPlainFromOpaque(req.Opaque, "spaceid"); reqSpaceID != "" {
		spaceID = reqSpaceID
	}
	// allow sending a space description
	description := utils.ReadPlainFromOpaque(req.Opaque, "description")
	// allow sending a spaceAlias
	alias := utils.ReadPlainFromOpaque(req.Opaque, "spaceAlias")
	u := ctxpkg.ContextMustGetUser(ctx)
	if alias == "" {
		alias = templates.WithSpacePropertiesAndUser(u, req.Type, req.Name, fs.o.GeneralSpaceAliasTemplate)
	}
	// TODO enforce a uuid?
	// TODO clarify if we want to enforce a single personal storage space or if we want to allow sending the spaceid
	if req.Type == spaceTypePersonal {
		spaceID = req.GetOwner().GetId().GetOpaqueId()
		alias = templates.WithSpacePropertiesAndUser(u, req.Type, req.Name, fs.o.PersonalSpaceAliasTemplate)
	}

	root, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true) // will fall into `Exists` case below
	if err == nil && root.Exists {
		return nil, errtypes.AlreadyExists("decomposedfs: spaces: space already exists")
	}

	if !fs.canCreateSpace(ctx, spaceID) {
		return nil, errtypes.PermissionDenied(spaceID)
	}

	// create a directory node
	rootPath := root.InternalPath()
	if err = os.MkdirAll(rootPath, 0700); err != nil {
		return nil, errors.Wrap(err, "decomposedfs: error creating node")
	}

	if err := root.WriteAllNodeMetadata(); err != nil {
		return nil, err
	}
	if req.GetOwner() != nil && req.GetOwner().GetId() != nil {
		if err := root.WriteOwner(req.GetOwner().GetId()); err != nil {
			return nil, err
		}
	}

	err = fs.linkStorageSpaceType(ctx, req.Type, root.ID)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]string, 6)

	// always enable propagation on the storage space root
	// mark the space root node as the end of propagation
	metadata[xattrs.PropagationAttr] = "1"
	metadata[xattrs.SpaceNameAttr] = req.Name

	if req.Type != "" {
		metadata[xattrs.SpaceTypeAttr] = req.Type
	}

	if q := req.GetQuota(); q != nil {
		// set default space quota
		metadata[xattrs.QuotaAttr] = strconv.FormatUint(q.QuotaMaxBytes, 10)
	}

	if description != "" {
		metadata[xattrs.SpaceDescriptionAttr] = description
	}

	if alias != "" {
		metadata[xattrs.SpaceAliasAttr] = alias
	}

	if err := xattrs.SetMultiple(root.InternalPath(), metadata); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, utils.SpaceGrant, struct{}{})

	if req.Type != spaceTypePersonal {
		u := ctxpkg.ContextMustGetUser(ctx)
		if err := fs.AddGrant(ctx, &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: spaceID,
				OpaqueId:  spaceID,
			},
		}, &provider.Grant{
			Grantee: &provider.Grantee{
				Type: provider.GranteeType_GRANTEE_TYPE_USER,
				Id: &provider.Grantee_UserId{
					UserId: u.Id,
				},
			},
			Permissions: ocsconv.NewManagerRole().CS3ResourcePermissions(),
		}); err != nil {
			return nil, err
		}
	}

	space, err := fs.storageSpaceFromNode(ctx, root, root.InternalPath(), false, false)
	if err != nil {
		return nil, err
	}

	resp := &provider.CreateStorageSpaceResponse{
		Status: &v1beta11.Status{
			Code: v1beta11.Code_CODE_OK,
		},
		StorageSpace: space,
	}
	return resp, nil
}

func (fs *Decomposedfs) canListAllSpaces(ctx context.Context) bool {
	user := ctxpkg.ContextMustGetUser(ctx)
	checkRes, err := fs.permissionsClient.CheckPermission(ctx, &cs3permissions.CheckPermissionRequest{
		Permission: "list-all-spaces",
		SubjectRef: &cs3permissions.SubjectReference{
			Spec: &cs3permissions.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
	})
	if err != nil {
		return false
	}

	return checkRes.Status.Code == v1beta11.Code_CODE_OK
}

// returns true when the user in the context can create a space / resource with storageID and nodeID set to his user opaqueID
func (fs *Decomposedfs) canCreateSpace(ctx context.Context, spaceID string) bool {
	user := ctxpkg.ContextMustGetUser(ctx)
	checkRes, err := fs.permissionsClient.CheckPermission(ctx, &cs3permissions.CheckPermissionRequest{
		Permission: "create-space",
		SubjectRef: &cs3permissions.SubjectReference{
			Spec: &cs3permissions.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
		Ref: &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: spaceID,
				// OpaqueId is the same, no need to transfer it
			},
		},
	})
	if err != nil {
		return false
	}

	return checkRes.Status.Code == v1beta11.Code_CODE_OK
}

// ReadSpaceAndNodeFromSpaceTypeLink reads a symlink and parses space and node id if the link has the correct format, eg:
// ../../spaces/4c/510ada-c86b-4815-8820-42cdf82c3d51/nodes/4c/51/0a/da/-c86b-4815-8820-42cdf82c3d51
// ../../spaces/4c/510ada-c86b-4815-8820-42cdf82c3d51/nodes/4c/51/0a/da/-c86b-4815-8820-42cdf82c3d51.T.2022-02-24T12:35:18.196484592Z
func ReadSpaceAndNodeFromSpaceTypeLink(path string) (string, string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return "", "", err
	}
	// ../../spaces/sp/ace-id/nodes/sh/or/tn/od/eid
	// 0  1  2      3  4      5     6  7  8  9  10
	parts := strings.Split(link, string(filepath.Separator))
	if len(parts) != 11 || parts[0] != ".." || parts[1] != ".." || parts[2] != "spaces" || parts[5] != "nodes" {
		return "", "", errtypes.InternalError("malformed link")
	}
	return strings.Join(parts[3:5], ""), strings.Join(parts[6:11], ""), nil
}

// ListStorageSpaces returns a list of StorageSpaces.
// The list can be filtered by space type or space id.
// Spaces are persisted with symlinks in /spaces/<type>/<spaceid> pointing to ../../nodes/<nodeid>, the root node of the space
// The spaceid is a concatenation of storageid + "!" + nodeid
func (fs *Decomposedfs) ListStorageSpaces(ctx context.Context, filter []*provider.ListStorageSpacesRequest_Filter, unrestricted bool) ([]*provider.StorageSpace, error) {
	// TODO check filters

	// TODO when a space symlink is broken delete the space for cleanup
	// read permissions are deduced from the node?

	// TODO for absolute references this actually requires us to move all user homes into a subfolder of /nodes/root,
	// e.g. /nodes/root/<space type> otherwise storage space names might collide even though they are of different types
	// /nodes/root/personal/foo and /nodes/root/shares/foo might be two very different spaces, a /nodes/root/foo is not expressive enough
	// we would not need /nodes/root if access always happened via spaceid+relative path

	var (
		spaceID = spaceIDAny
		nodeID  = spaceIDAny
	)

	spaceTypes := []string{}

	for i := range filter {
		switch filter[i].Type {
		case provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE:
			switch filter[i].GetSpaceType() {
			case "+mountpoint":
				// TODO include mount poits
			case "+grant":
				// TODO include grants
			default:
				spaceTypes = append(spaceTypes, filter[i].GetSpaceType())
			}
		case provider.ListStorageSpacesRequest_Filter_TYPE_ID:
			spaceID, nodeID, _ = storagespace.SplitID(filter[i].GetId().OpaqueId)
			if strings.Contains(nodeID, "/") {
				return []*provider.StorageSpace{}, nil
			}
		}
	}
	if len(spaceTypes) == 0 {
		spaceTypes = []string{spaceTypeAny}
	}

	canListAllSpaces := fs.canListAllSpaces(ctx)

	spaces := []*provider.StorageSpace{}
	// build the glob path, eg.
	// /path/to/root/spaces/{spaceType}/{spaceId}
	// /path/to/root/spaces/personal/nodeid
	// /path/to/root/spaces/shared/nodeid

	if spaceID != spaceIDAny && nodeID != spaceIDAny {
		// try directly reading the node
		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true) // permission to read disabled space is checked later
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("id", nodeID).Msg("could not read node")
			return nil, err
		}
		if !n.Exists {
			// return empty list
			return spaces, nil
		}
		space, err := fs.storageSpaceFromNode(ctx, n, n.InternalPath(), canListAllSpaces, unrestricted)
		if err != nil {
			return nil, err
		}
		// filter space types
		for _, spaceType := range spaceTypes {
			if spaceType == spaceTypeAny || spaceType == space.SpaceType {
				spaces = append(spaces, space)
			}
		}

		return spaces, nil
	}

	matches := []string{}
	for _, spaceType := range spaceTypes {
		path := filepath.Join(fs.o.Root, "spacetypes", spaceType, nodeID)
		m, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m...)
	}

	// FIXME if the space does not exist try a node as the space root.

	// But then the whole /spaces/{spaceType}/{spaceid} becomes obsolete
	// we can alway just look up by nodeid
	// -> no. The /spaces folder is used for efficient lookup by type, otherwise we would have
	//    to iterate over all nodes and read the type from extended attributes
	// -> but for lookup by id we can use the node directly.
	// But what about sharding nodes by space?
	// an efficient lookup would be possible if we received a spaceid&opaqueid in the request
	// the personal spaces must also use the nodeid and not the name

	numShares := 0

	for i := range matches {
		var err error
		// do not investigate flock files any further. They indicate file locks but are not relevant here.
		if strings.HasSuffix(matches[i], ".flock") {
			continue
		}
		// always read link in case storage space id != node id
		spaceID, nodeID, err = ReadSpaceAndNodeFromSpaceTypeLink(matches[i])
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("match", matches[i]).Msg("could not read link, skipping")
			continue
		}

		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("id", nodeID).Msg("could not read node, skipping")
			continue
		}

		if !n.Exists {
			continue
		}

		spaceType := filepath.Base(filepath.Dir(matches[i]))

		// FIXME type share evolved to grant on the edge branch ... make it configurable if the driver should support them or not for now ... ignore type share
		if spaceType == spaceTypeShare {
			numShares++
			// do not list shares as spaces for the owner
			continue
		}

		// TODO apply more filters
		space, err := fs.storageSpaceFromNode(ctx, n, matches[i], canListAllSpaces, unrestricted)
		if err != nil {
			if _, ok := err.(errtypes.IsPermissionDenied); !ok {
				appctx.GetLogger(ctx).Error().Err(err).Interface("node", n).Msg("could not convert to storage space")
			}
			continue
		}
		spaces = append(spaces, space)

	}
	// if there are no matches (or they happened to be spaces for the owner) and the node is a child return a space
	if len(matches) <= numShares && nodeID != spaceID {
		// try node id
		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true) // permission to read disabled space is checked in storageSpaceFromNode
		if err != nil {
			return nil, err
		}
		if n.Exists {
			space, err := fs.storageSpaceFromNode(ctx, n, n.InternalPath(), canListAllSpaces, unrestricted)
			if err != nil {
				return nil, err
			}
			spaces = append(spaces, space)
		}
	}

	return spaces, nil

}

// UpdateStorageSpace updates a storage space
func (fs *Decomposedfs) UpdateStorageSpace(ctx context.Context, req *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	var restore bool
	if req.Opaque != nil {
		_, restore = req.Opaque.Map["restore"]
	}

	space := req.StorageSpace
	spaceID, _, _ := storagespace.SplitID(space.Id.OpaqueId)

	node, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true) // permission to read disabled space will be checked later
	if err != nil {
		return nil, err
	}

	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		return nil, fmt.Errorf("decomposedfs: spaces: contextual user not found")
	}
	space.Owner = u

	metadata := make(map[string]string, 5)
	if space.Name != "" {
		metadata[xattrs.SpaceNameAttr] = space.Name
	}

	if space.Quota != nil {
		metadata[xattrs.QuotaAttr] = strconv.FormatUint(space.Quota.QuotaMaxBytes, 10)
	}

	// TODO also return values which are not in the request
	hasDescription := false
	if space.Opaque != nil {
		if description, ok := space.Opaque.Map["description"]; ok {
			metadata[xattrs.SpaceDescriptionAttr] = string(description.Value)
			hasDescription = true
		}
		if alias := utils.ReadPlainFromOpaque(space.Opaque, "spaceAlias"); alias != "" {
			metadata[xattrs.SpaceAliasAttr] = alias
		}
		if image := utils.ReadPlainFromOpaque(space.Opaque, "image"); image != "" {
			imageID, err := storagespace.ParseID(image)
			if err != nil {
				return &provider.UpdateStorageSpaceResponse{
					Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND, Message: "decomposedFS: space image resource not found"},
				}, nil
			}
			metadata[xattrs.SpaceImageAttr] = imageID.OpaqueId
		}
		if readme := utils.ReadPlainFromOpaque(space.Opaque, "readme"); readme != "" {
			readmeID, err := storagespace.ParseID(readme)
			if err != nil {
				return &provider.UpdateStorageSpaceResponse{
					Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND, Message: "decomposedFS: space readme resource not found"},
				}, nil
			}
			metadata[xattrs.SpaceReadmeAttr] = readmeID.OpaqueId
		}
	}

	// TODO change the permission handling
	// these three attributes need manager permissions
	if space.Name != "" || hasDescription || restore {
		err = fs.checkManagerPermission(ctx, node)
	}
	if err != nil {
		if restore {
			// a disabled space is invisible to non admins
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND, Message: err.Error()},
			}, nil
		}
		return &provider.UpdateStorageSpaceResponse{
			Status: &v1beta11.Status{Code: v1beta11.Code_CODE_PERMISSION_DENIED, Message: err.Error()},
		}, nil
	}
	// all other attributes need editor permissions
	err = fs.checkEditorPermission(ctx, node)
	if err != nil {
		return &provider.UpdateStorageSpaceResponse{
			Status: &v1beta11.Status{Code: v1beta11.Code_CODE_PERMISSION_DENIED, Message: err.Error()},
		}, nil
	}

	err = xattrs.SetMultiple(node.InternalPath(), metadata)
	if err != nil {
		return nil, err
	}

	if restore {
		if err := node.SetDTime(nil); err != nil {
			return nil, err
		}
	}

	// send back the updated data from the storage
	updatedSpace, err := fs.storageSpaceFromNode(ctx, node, node.InternalPath(), false, false)
	if err != nil {
		return nil, err
	}

	return &provider.UpdateStorageSpaceResponse{
		Status:       &v1beta11.Status{Code: v1beta11.Code_CODE_OK},
		StorageSpace: updatedSpace,
	}, nil
}

// DeleteStorageSpace deletes a storage space
func (fs *Decomposedfs) DeleteStorageSpace(ctx context.Context, req *provider.DeleteStorageSpaceRequest) error {
	opaque := req.Opaque
	var purge bool
	if opaque != nil {
		_, purge = opaque.Map["purge"]
	}

	spaceID := req.Id.OpaqueId

	n, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true) // permission to read disabled space is checked later
	if err != nil {
		return err
	}

	// only managers are allowed to disable or purge a drive
	if err := fs.checkManagerPermission(ctx, n); err != nil {
		return errtypes.PermissionDenied(fmt.Sprintf("user is not allowed to delete spaces %s", n.ID))
	}

	if purge {
		if !n.IsDisabled() {
			return errtypes.NewErrtypeFromStatus(status.NewInvalidArg(ctx, "can't purge enabled space"))
		}

		spaceType, err := n.GetMetadata(xattrs.SpaceTypeAttr)
		if err != nil {
			return err
		}
		// remove type index
		spaceTypePath := filepath.Join(fs.o.Root, "spacetypes", spaceType, spaceID)
		if err := os.Remove(spaceTypePath); err != nil {
			return err
		}

		// remove space metadata
		if err := os.RemoveAll(filepath.Join(fs.o.Root, "spaces", lookup.Pathify(spaceID, 1, 2))); err != nil {
			return err
		}

		// FIXME remove space blobs

		return nil
	}

	// mark as disabled by writing a dtime attribute
	dtime := time.Now()
	return n.SetDTime(&dtime)
}

func (fs *Decomposedfs) linkStorageSpaceType(ctx context.Context, spaceType string, spaceID string) error {
	// create space type dir
	if err := os.MkdirAll(filepath.Join(fs.o.Root, "spacetypes", spaceType), 0700); err != nil {
		return err
	}

	// link space in spacetypes
	err := os.Symlink("../../spaces/"+lookup.Pathify(spaceID, 1, 2)+"/nodes/"+lookup.Pathify(spaceID, 4, 2), filepath.Join(fs.o.Root, "spacetypes", spaceType, spaceID))
	if err != nil {
		if isAlreadyExists(err) {
			appctx.GetLogger(ctx).Debug().Err(err).Str("space", spaceID).Str("spacetype", spaceType).Msg("symlink already exists")
			// FIXME: is it ok to wipe this err if the symlink already exists?
			err = nil
		} else {
			// TODO how should we handle error cases here?
			appctx.GetLogger(ctx).Error().Err(err).Str("space", spaceID).Str("spacetype", spaceType).Msg("could not create symlink")
		}
	}

	return err
}

func (fs *Decomposedfs) storageSpaceFromNode(ctx context.Context, n *node.Node, nodePath string, canListAllSpaces bool, unrestricted bool) (*provider.StorageSpace, error) {
	user := ctxpkg.ContextMustGetUser(ctx)
	if !canListAllSpaces || !unrestricted {
		ok, err := node.NewPermissions(fs.lu).HasPermission(ctx, n, func(p *provider.ResourcePermissions) bool {
			return p.Stat
		})
		if err != nil || !ok {
			return nil, errtypes.PermissionDenied(fmt.Sprintf("user %s is not allowed to Stat the space %s", user.Username, n.ID))
		}

		if n.SpaceRoot.IsDisabled() {
			if err := fs.checkManagerPermission(ctx, n); err != nil {
				return nil, errtypes.PermissionDenied(fmt.Sprintf("user %s is not allowed to list deleted spaces %s", user.Username, n.ID))
			}
		}
	}

	var err error
	// TODO apply more filters
	var sname string
	if sname, err = n.SpaceRoot.GetMetadata(xattrs.SpaceNameAttr); err != nil {
		// FIXME: Is that a severe problem?
		appctx.GetLogger(ctx).Debug().Err(err).Msg("space does not have a name attribute")
	}

	/*
		if err := n.FindStorageSpaceRoot(); err != nil {
			return nil, err
		}
	*/

	// read the grants from the current node, not the root
	grants, err := n.ListGrants(ctx)
	if err != nil {
		return nil, err
	}

	m := make(map[string]*provider.ResourcePermissions, len(grants))
	for _, g := range grants {
		var id string
		switch g.Grantee.Type {
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			id = g.Grantee.GetGroupId().OpaqueId
		case provider.GranteeType_GRANTEE_TYPE_USER:
			id = g.Grantee.GetUserId().OpaqueId
		default:
			continue
		}

		m[id] = g.Permissions
	}
	marshalled, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	ssID, err := storagespace.FormatReference(
		&provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: n.SpaceRoot.SpaceID,
				OpaqueId:  n.ID},
		},
	)
	if err != nil {
		return nil, err
	}
	space := &provider.StorageSpace{
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"grants": {
					Decoder: "json",
					Value:   marshalled,
				},
			},
		},
		Id: &provider.StorageSpaceId{OpaqueId: ssID},
		Root: &provider.ResourceId{
			StorageId: n.SpaceRoot.SpaceID,
			OpaqueId:  n.ID,
		},
		Name: sname,
		// SpaceType is read from xattr below
		// Mtime is set either as node.tmtime or as fi.mtime below
	}

	if space.SpaceType, err = n.SpaceRoot.GetMetadata(xattrs.SpaceTypeAttr); err != nil {
		appctx.GetLogger(ctx).Debug().Err(err).Msg("space does not have a type attribute")
	}

	if n.SpaceRoot.IsDisabled() {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "trashed", "trashed")
	}

	if n.Owner() != nil && n.Owner().OpaqueId != "" {
		space.Owner = &userv1beta1.User{ // FIXME only return a UserID, not a full blown user object
			Id: n.Owner(),
		}
	}

	// we set the space mtime to the root item mtime
	// override the stat mtime with a tmtime if it is present
	var tmtime time.Time
	if tmt, err := n.GetTMTime(); err == nil {
		tmtime = tmt
		un := tmt.UnixNano()
		space.Mtime = &types.Timestamp{
			Seconds: uint64(un / 1000000000),
			Nanos:   uint32(un % 1000000000),
		}
	} else if fi, err := os.Stat(nodePath); err == nil {
		// fall back to stat mtime
		tmtime = fi.ModTime()
		un := fi.ModTime().UnixNano()
		space.Mtime = &types.Timestamp{
			Seconds: uint64(un / 1000000000),
			Nanos:   uint32(un % 1000000000),
		}
	}

	etag, err := node.CalculateEtag(n.ID, tmtime)
	if err != nil {
		return nil, err
	}
	space.Opaque.Map["etag"] = &types.OpaqueEntry{
		Decoder: "plain",
		Value:   []byte(etag),
	}

	spaceAttributes, err := xattrs.All(nodePath)
	if err != nil {
		return nil, err
	}

	// quota
	quotaAttr, ok := spaceAttributes[xattrs.QuotaAttr]
	if ok {
		// make sure we have a proper signed int
		// we use the same magic numbers to indicate:
		// -1 = uncalculated
		// -2 = unknown
		// -3 = unlimited
		if quota, err := strconv.ParseUint(quotaAttr, 10, 64); err == nil {
			space.Quota = &provider.Quota{
				QuotaMaxBytes: quota,
				QuotaMaxFiles: math.MaxUint64, // TODO MaxUInt64? = unlimited? why even max files? 0 = unlimited?
			}
		} else {
			return nil, err
		}
	}
	spaceImage, ok := spaceAttributes[xattrs.SpaceImageAttr]
	if ok {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "image", storagespace.FormatResourceID(
			provider.ResourceId{StorageId: space.Root.StorageId, OpaqueId: spaceImage},
		))
	}
	spaceDescription, ok := spaceAttributes[xattrs.SpaceDescriptionAttr]
	if ok {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "description", spaceDescription)
	}
	spaceReadme, ok := spaceAttributes[xattrs.SpaceReadmeAttr]
	if ok {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "readme", storagespace.FormatResourceID(
			provider.ResourceId{StorageId: space.Root.StorageId, OpaqueId: spaceReadme},
		))
	}
	spaceAlias, ok := spaceAttributes[xattrs.SpaceAliasAttr]
	if ok {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "spaceAlias", spaceAlias)
	}
	return space, nil
}

func (fs *Decomposedfs) checkManagerPermission(ctx context.Context, n *node.Node) error {
	// to update the space name or short description we need the manager role
	// current workaround: check if RemoveGrant Permission exists
	managerPerm, err := fs.p.HasPermission(ctx, n, func(rp *provider.ResourcePermissions) bool {
		return rp.RemoveGrant
	})
	switch {
	case err != nil:
		return errtypes.InternalError(err.Error())
	case !managerPerm:
		msg := fmt.Sprintf("not enough permissions to change attributes on %s", filepath.Join(n.ParentID, n.Name))
		return errtypes.PermissionDenied(msg)
	}
	return nil
}

func (fs *Decomposedfs) checkEditorPermission(ctx context.Context, n *node.Node) error {
	// current workaround: check if InitiateFileUpload Permission exists
	editorPerm, err := fs.p.HasPermission(ctx, n, func(rp *provider.ResourcePermissions) bool {
		return rp.InitiateFileUpload
	})
	switch {
	case err != nil:
		return errtypes.InternalError(err.Error())
	case !editorPerm:
		msg := fmt.Sprintf("not enough permissions to change attributes on %s", filepath.Join(n.ParentID, n.Name))
		return errtypes.PermissionDenied(msg)
	}
	return nil
}
