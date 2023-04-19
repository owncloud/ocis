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
	"strings"
	"time"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ocsconv "github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/filelocks"
	"github.com/cs3org/reva/v2/pkg/storage/utils/templates"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	_spaceTypePersonal = "personal"
	_spaceTypeProject  = "project"
	spaceTypeShare     = "share"
	spaceTypeAny       = "*"
	spaceIDAny         = "*"
	userIDAny          = "*"

	quotaUnrestricted = 0
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
	if req.Type == _spaceTypePersonal {
		spaceID = req.GetOwner().GetId().GetOpaqueId()
		alias = templates.WithSpacePropertiesAndUser(u, req.Type, req.Name, fs.o.PersonalSpaceAliasTemplate)
	}

	root, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true, nil, false) // will fall into `Exists` case below
	switch {
	case err != nil:
		return nil, err
	case !fs.p.CreateSpace(ctx, spaceID):
		return nil, errtypes.PermissionDenied(spaceID)
	case root.Exists:
		return nil, errtypes.AlreadyExists("decomposedfs: spaces: space already exists")
	}

	// create a directory node
	root.SetType(provider.ResourceType_RESOURCE_TYPE_CONTAINER)
	rootPath := root.InternalPath()

	if err := os.MkdirAll(rootPath, 0700); err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error creating node")
	}

	if err := root.WriteAllNodeMetadata(ctx); err != nil {
		return nil, err
	}
	var owner *userv1beta1.UserId
	if req.GetOwner() != nil && req.GetOwner().GetId() != nil {
		owner = req.GetOwner().GetId()
	} else {
		owner = &userv1beta1.UserId{OpaqueId: spaceID, Type: userv1beta1.UserType_USER_TYPE_SPACE_OWNER}
	}
	if err := root.WriteOwner(owner); err != nil {
		return nil, err
	}

	err = fs.updateIndexes(ctx, req.GetOwner().GetId().GetOpaqueId(), req.Type, root.ID)
	if err != nil {
		return nil, err
	}

	metadata := make(node.Attributes, 6)

	// always enable propagation on the storage space root
	// mark the space root node as the end of propagation
	metadata.SetString(prefixes.PropagationAttr, "1")
	metadata.SetString(prefixes.NameAttr, req.Name)
	metadata.SetString(prefixes.SpaceNameAttr, req.Name)

	if req.Type != "" {
		metadata.SetString(prefixes.SpaceTypeAttr, req.Type)
	}

	if q := req.GetQuota(); q != nil {
		// set default space quota
		if fs.o.MaxQuota != quotaUnrestricted && q.GetQuotaMaxBytes() > fs.o.MaxQuota {
			return nil, errtypes.BadRequest("decompsedFS: requested quota is higher than allowed")
		}
		metadata.SetInt64(prefixes.QuotaAttr, int64(q.QuotaMaxBytes))
	} else if fs.o.MaxQuota != quotaUnrestricted {
		// If no quota was requested but a max quota was set then the the storage space has a quota
		// of max quota.
		metadata.SetInt64(prefixes.QuotaAttr, int64(fs.o.MaxQuota))
	}

	if description != "" {
		metadata.SetString(prefixes.SpaceDescriptionAttr, description)
	}

	if alias != "" {
		metadata.SetString(prefixes.SpaceAliasAttr, alias)
	}

	if err := root.SetXattrs(metadata, true); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, utils.SpaceGrant, struct{ SpaceType string }{SpaceType: req.Type})

	if req.Type != _spaceTypePersonal {
		u := ctxpkg.ContextMustGetUser(ctx)
		if err := fs.AddGrant(ctx, &provider.Reference{
			ResourceId: &provider.ResourceId{
				SpaceId:  spaceID,
				OpaqueId: spaceID,
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

	space, err := fs.storageSpaceFromNode(ctx, root, true)
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

// ReadSpaceAndNodeFromIndexLink reads a symlink and parses space and node id if the link has the correct format, eg:
// ../../spaces/4c/510ada-c86b-4815-8820-42cdf82c3d51/nodes/4c/51/0a/da/-c86b-4815-8820-42cdf82c3d51
// ../../spaces/4c/510ada-c86b-4815-8820-42cdf82c3d51/nodes/4c/51/0a/da/-c86b-4815-8820-42cdf82c3d51.T.2022-02-24T12:35:18.196484592Z
func ReadSpaceAndNodeFromIndexLink(link string) (string, string, error) {
	// ../../../spaces/sp/ace-id/nodes/sh/or/tn/od/eid
	// 0  1  2  3      4  5      6     7  8  9  10  11
	parts := strings.Split(link, string(filepath.Separator))
	if len(parts) != 12 || parts[0] != ".." || parts[1] != ".." || parts[2] != ".." || parts[3] != "spaces" || parts[6] != "nodes" {
		return "", "", errtypes.InternalError("malformed link")
	}
	return strings.Join(parts[4:6], ""), strings.Join(parts[7:12], ""), nil
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
		spaceID         = spaceIDAny
		nodeID          = spaceIDAny
		requestedUserID = userIDAny
	)

	spaceTypes := map[string]struct{}{}

	for i := range filter {
		switch filter[i].Type {
		case provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE:
			switch filter[i].GetSpaceType() {
			case "+mountpoint":
				// TODO include mount poits
			case "+grant":
				// TODO include grants
			default:
				spaceTypes[filter[i].GetSpaceType()] = struct{}{}
			}
		case provider.ListStorageSpacesRequest_Filter_TYPE_ID:
			_, spaceID, nodeID, _ = storagespace.SplitID(filter[i].GetId().OpaqueId)
			if strings.Contains(nodeID, "/") {
				return []*provider.StorageSpace{}, nil
			}
		case provider.ListStorageSpacesRequest_Filter_TYPE_USER:
			// TODO: refactor this to GetUserId() in cs3
			requestedUserID = filter[i].GetUser().GetOpaqueId()
		case provider.ListStorageSpacesRequest_Filter_TYPE_OWNER:
			// TODO: improve further by not evaluating shares
			requestedUserID = filter[i].GetOwner().GetOpaqueId()
		}
	}
	if len(spaceTypes) == 0 {
		spaceTypes[spaceTypeAny] = struct{}{}
	}

	authenticatedUserID := ctxpkg.ContextMustGetUser(ctx).GetId().GetOpaqueId()

	if !fs.p.ListSpacesOfUser(ctx, requestedUserID) {
		return nil, errtypes.PermissionDenied(fmt.Sprintf("user %s is not allowed to list spaces of other users", authenticatedUserID))
	}

	checkNodePermissions := fs.MustCheckNodePermissions(ctx, unrestricted)

	spaces := []*provider.StorageSpace{}
	// build the glob path, eg.
	// /path/to/root/spaces/{spaceType}/{spaceId}
	// /path/to/root/spaces/personal/nodeid
	// /path/to/root/spaces/shared/nodeid

	if spaceID != spaceIDAny && nodeID != spaceIDAny {
		// try directly reading the node
		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true, nil, false) // permission to read disabled space is checked later
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("id", nodeID).Msg("could not read node")
			return nil, err
		}
		if !n.Exists {
			// return empty list
			return spaces, nil
		}
		space, err := fs.storageSpaceFromNode(ctx, n, checkNodePermissions)
		if err != nil {
			return nil, err
		}
		// filter space types
		_, ok1 := spaceTypes[spaceTypeAny]
		_, ok2 := spaceTypes[space.SpaceType]
		if ok1 || ok2 {
			spaces = append(spaces, space)
		}
		// TODO: filter user id
		return spaces, nil
	}

	matches := map[string]struct{}{}

	if requestedUserID != userIDAny {
		path := filepath.Join(fs.o.Root, "indexes", "by-user-id", requestedUserID, nodeID)
		m, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		for _, match := range m {
			link, err := os.Readlink(match)
			if err != nil {
				continue
			}
			matches[link] = struct{}{}
		}
	}

	if requestedUserID == userIDAny {
		for spaceType := range spaceTypes {
			path := filepath.Join(fs.o.Root, "indexes", "by-type", spaceType, nodeID)
			m, err := filepath.Glob(path)
			if err != nil {
				return nil, err
			}
			for _, match := range m {
				link, err := os.Readlink(match)
				if err != nil {
					continue
				}
				matches[link] = struct{}{}
			}
		}
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

	for match := range matches {
		var err error
		// TODO introduce metadata.IsLockFile(path)
		// do not investigate flock files any further. They indicate file locks but are not relevant here.
		if strings.HasSuffix(match, filelocks.LockFileSuffix) {
			continue
		}
		// skip metadata files
		if fs.lu.MetadataBackend().IsMetaFile(match) {
			continue
		}
		// always read link in case storage space id != node id
		spaceID, nodeID, err = ReadSpaceAndNodeFromIndexLink(match)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("match", match).Msg("could not read link, skipping")
			continue
		}

		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true, nil, true)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("id", nodeID).Msg("could not read node, skipping")
			continue
		}

		if !n.Exists {
			continue
		}

		space, err := fs.storageSpaceFromNode(ctx, n, checkNodePermissions)
		if err != nil {
			switch err.(type) {
			case errtypes.IsPermissionDenied:
				// ok
			case errtypes.NotFound:
				// ok
			default:
				appctx.GetLogger(ctx).Error().Err(err).Str("id", nodeID).Msg("could not convert to storage space")
			}
			continue
		}

		// FIXME type share evolved to grant on the edge branch ... make it configurable if the driver should support them or not for now ... ignore type share
		if space.SpaceType == spaceTypeShare {
			numShares++
			// do not list shares as spaces for the owner
			continue
		}

		// TODO apply more filters
		_, ok1 := spaceTypes[spaceTypeAny]
		_, ok2 := spaceTypes[space.SpaceType]
		if ok1 || ok2 {
			spaces = append(spaces, space)
		}
	}
	// if there are no matches (or they happened to be spaces for the owner) and the node is a child return a space
	if len(matches) <= numShares && nodeID != spaceID {
		// try node id
		n, err := node.ReadNode(ctx, fs.lu, spaceID, nodeID, true, nil, false) // permission to read disabled space is checked in storageSpaceFromNode
		if err != nil {
			return nil, err
		}
		if n.Exists {
			space, err := fs.storageSpaceFromNode(ctx, n, checkNodePermissions)
			if err != nil {
				return nil, err
			}
			spaces = append(spaces, space)
		}
	}

	return spaces, nil

}

// MustCheckNodePermissions checks if permission checks are needed to be performed when user requests spaces
func (fs *Decomposedfs) MustCheckNodePermissions(ctx context.Context, unrestricted bool) bool {
	// canListAllSpaces indicates if the user has the permission from the global user role
	canListAllSpaces := fs.p.ListAllSpaces(ctx)
	// unrestricted is the param which indicates if the user wants to list all spaces or only the spaces he is part of
	// if a user lists all spaces unrestricted and doesn't have the permissions from the role, we need to check
	// the nodePermissions and this will return a spaces list where the user has access to
	// we can only skip the NodePermissions check if both values are true
	if canListAllSpaces && unrestricted {
		return false
	}
	return true
}

// UpdateStorageSpace updates a storage space
func (fs *Decomposedfs) UpdateStorageSpace(ctx context.Context, req *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	var restore bool
	if req.Opaque != nil {
		_, restore = req.Opaque.Map["restore"]
	}

	space := req.StorageSpace
	_, spaceID, _, _ := storagespace.SplitID(space.Id.OpaqueId)

	metadata := make(node.Attributes, 5)
	if space.Name != "" {
		metadata.SetString(prefixes.NameAttr, space.Name)
		metadata.SetString(prefixes.SpaceNameAttr, space.Name)
	}

	if space.Quota != nil {
		if fs.o.MaxQuota != quotaUnrestricted && fs.o.MaxQuota < space.Quota.QuotaMaxBytes {
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{Code: v1beta11.Code_CODE_INVALID_ARGUMENT, Message: "decompsedFS: requested quota is higher than allowed"},
			}, nil
		} else if fs.o.MaxQuota != quotaUnrestricted && space.Quota.QuotaMaxBytes == quotaUnrestricted {
			// If the caller wants to unrestrict the space we give it the maximum allowed quota.
			space.Quota.QuotaMaxBytes = fs.o.MaxQuota
		}
		metadata.SetInt64(prefixes.QuotaAttr, int64(space.Quota.QuotaMaxBytes))
	}

	// TODO also return values which are not in the request
	if space.Opaque != nil {
		if description, ok := space.Opaque.Map["description"]; ok {
			metadata[prefixes.SpaceDescriptionAttr] = description.Value
		}
		if alias := utils.ReadPlainFromOpaque(space.Opaque, "spaceAlias"); alias != "" {
			metadata.SetString(prefixes.SpaceAliasAttr, alias)
		}
		if image := utils.ReadPlainFromOpaque(space.Opaque, "image"); image != "" {
			imageID, err := storagespace.ParseID(image)
			if err != nil {
				return &provider.UpdateStorageSpaceResponse{
					Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND, Message: "decomposedFS: space image resource not found"},
				}, nil
			}
			metadata.SetString(prefixes.SpaceImageAttr, imageID.OpaqueId)
		}
		if readme := utils.ReadPlainFromOpaque(space.Opaque, "readme"); readme != "" {
			readmeID, err := storagespace.ParseID(readme)
			if err != nil {
				return &provider.UpdateStorageSpaceResponse{
					Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND, Message: "decomposedFS: space readme resource not found"},
				}, nil
			}
			metadata.SetString(prefixes.SpaceReadmeAttr, readmeID.OpaqueId)
		}
	}

	// check which permissions are needed
	spaceNode, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true, nil, false)
	if err != nil {
		return nil, err
	}

	if !spaceNode.Exists {
		return &provider.UpdateStorageSpaceResponse{
			Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND},
		}, nil
	}

	sp, err := fs.p.AssemblePermissions(ctx, spaceNode)
	if err != nil {
		return &provider.UpdateStorageSpaceResponse{
			Status: &v1beta11.Status{Code: v1beta11.Code_CODE_INVALID, Message: err.Error()},
		}, nil

	}

	if !restore && len(metadata) == 0 && !IsViewer(sp) {
		// you may land here when making an update request without changes
		// check if user has access to the drive before continuing
		return &provider.UpdateStorageSpaceResponse{
			Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND},
		}, nil
	}

	if !IsManager(sp) {
		// We are not a space manager. We need to check for additional permissions.
		k := []string{prefixes.NameAttr, prefixes.SpaceDescriptionAttr}
		if !IsEditor(sp) {
			k = append(k, prefixes.SpaceReadmeAttr, prefixes.SpaceAliasAttr, prefixes.SpaceImageAttr)
		}

		if mapHasKey(metadata, k...) && !fs.p.ManageSpaceProperties(ctx, spaceID) {
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{Code: v1beta11.Code_CODE_PERMISSION_DENIED},
			}, nil
		}

		if restore && !fs.p.SpaceAbility(ctx, spaceID) {
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{Code: v1beta11.Code_CODE_NOT_FOUND},
			}, nil
		}
	}

	if mapHasKey(metadata, prefixes.QuotaAttr) {
		typ, err := spaceNode.SpaceRoot.Xattr(prefixes.SpaceTypeAttr)
		if err != nil {
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{
					Code:    v1beta11.Code_CODE_INTERNAL,
					Message: "space has no type",
				},
			}, nil
		}

		if !fs.p.SetSpaceQuota(ctx, spaceID, string(typ)) {
			return &provider.UpdateStorageSpaceResponse{
				Status: &v1beta11.Status{Code: v1beta11.Code_CODE_PERMISSION_DENIED},
			}, nil
		}
	}

	err = spaceNode.SetXattrs(metadata, true)
	if err != nil {
		return nil, err
	}

	if restore {
		if err := spaceNode.SetDTime(nil); err != nil {
			return nil, err
		}
	}

	// send back the updated data from the storage
	updatedSpace, err := fs.storageSpaceFromNode(ctx, spaceNode, false)
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

	_, spaceID, _, err := storagespace.SplitID(req.Id.GetOpaqueId())
	if err != nil {
		return err
	}

	n, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, true, nil, false) // permission to read disabled space is checked later
	if err != nil {
		return err
	}

	st, err := n.SpaceRoot.XattrString(prefixes.SpaceTypeAttr)
	if err != nil {
		return errtypes.InternalError(fmt.Sprintf("space %s does not have a spacetype, possible corrupt decompsedfs", n.ID))
	}

	// - a User with the "delete-all-spaces" permission can delete any space
	// - spaces of type personal can also be deleted by users with the "delete-all-home-spaces" permission
	// - otherwise a space can be deleted by its manager (i.e. users have the "remove" grant)
	switch {
	case fs.p.DeleteAllSpaces(ctx):
		// We are allowed to delete any space, no further permission checks needed
	case st == "personal":
		if !fs.p.DeleteAllHomeSpaces(ctx) {
			return errtypes.PermissionDenied(fmt.Sprintf("user is not allowed to delete home space %s", n.ID))
		}
	default:
		// managers and users with 'space ability' permission are allowed to disable or purge a drive
		if !fs.p.SpaceAbility(ctx, spaceID) {
			rp, err := fs.p.AssemblePermissions(ctx, n)
			if err != nil || !IsManager(rp) {
				return errtypes.PermissionDenied(fmt.Sprintf("user is not allowed to delete spaces %s", n.ID))
			}
		}
	}

	if purge {
		if !n.IsDisabled() {
			return errtypes.NewErrtypeFromStatus(status.NewInvalid(ctx, "can't purge enabled space"))
		}

		spaceType, err := n.XattrString(prefixes.SpaceTypeAttr)
		if err != nil {
			return err
		}
		// remove type index
		spaceTypePath := filepath.Join(fs.o.Root, "indexes", "by-type", spaceType, spaceID)
		if err := os.Remove(spaceTypePath); err != nil {
			return err
		}

		// remove space metadata
		if err := os.RemoveAll(fs.getSpaceRoot(spaceID)); err != nil {
			return err
		}

		// FIXME remove space blobs

		return nil
	}

	// mark as disabled by writing a dtime attribute
	dtime := time.Now()
	return n.SetDTime(&dtime)
}

func (fs *Decomposedfs) updateIndexes(ctx context.Context, userID, spaceType, spaceID string) error {
	err := fs.linkStorageSpaceType(ctx, spaceType, spaceID)
	if err != nil {
		return err
	}
	return fs.linkSpaceByUser(ctx, userID, spaceID)
}

func (fs *Decomposedfs) linkSpaceByUser(ctx context.Context, userID, spaceID string) error {
	if userID == "" {
		return nil
	}
	// create user index dir
	// TODO: pathify userID
	if err := os.MkdirAll(filepath.Join(fs.o.Root, "indexes", "by-user-id", userID), 0700); err != nil {
		return err
	}

	err := os.Symlink("../../../spaces/"+lookup.Pathify(spaceID, 1, 2)+"/nodes/"+lookup.Pathify(spaceID, 4, 2), filepath.Join(fs.o.Root, "indexes/by-user-id", userID, spaceID))
	if err != nil {
		if isAlreadyExists(err) {
			appctx.GetLogger(ctx).Debug().Err(err).Str("space", spaceID).Str("user-id", userID).Msg("symlink already exists")
			// FIXME: is it ok to wipe this err if the symlink already exists?
			err = nil //nolint
		} else {
			// TODO how should we handle error cases here?
			appctx.GetLogger(ctx).Error().Err(err).Str("space", spaceID).Str("user-id", userID).Msg("could not create symlink")
		}
	}
	return nil
}

// TODO: implement linkSpaceByGroup

func (fs *Decomposedfs) linkStorageSpaceType(ctx context.Context, spaceType string, spaceID string) error {
	if spaceType == "" {
		return nil
	}
	// create space type dir
	if err := os.MkdirAll(filepath.Join(fs.o.Root, "indexes", "by-type", spaceType), 0700); err != nil {
		return err
	}

	// link space in spacetypes
	err := os.Symlink("../../../spaces/"+lookup.Pathify(spaceID, 1, 2)+"/nodes/"+lookup.Pathify(spaceID, 4, 2), filepath.Join(fs.o.Root, "indexes", "by-type", spaceType, spaceID))
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

func (fs *Decomposedfs) storageSpaceFromNode(ctx context.Context, n *node.Node, checkPermissions bool) (*provider.StorageSpace, error) {
	user := ctxpkg.ContextMustGetUser(ctx)
	if checkPermissions {
		rp, err := fs.p.AssemblePermissions(ctx, n)
		switch {
		case err != nil:
			return nil, errtypes.InternalError(err.Error())
		case !rp.Stat:
			return nil, errtypes.NotFound(fmt.Sprintf("space %s not found", n.ID))
		}

		if n.SpaceRoot.IsDisabled() {
			rp, err := fs.p.AssemblePermissions(ctx, n)
			if err != nil || !IsManager(rp) {
				return nil, errtypes.PermissionDenied(fmt.Sprintf("user %s is not allowed to list deleted spaces %s", user.Username, n.ID))
			}
		}
	}

	var err error
	// TODO apply more filters
	var sname string
	if sname, err = n.SpaceRoot.XattrString(prefixes.SpaceNameAttr); err != nil {
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

	grantMap := make(map[string]*provider.ResourcePermissions, len(grants))
	grantExpiration := make(map[string]*types.Timestamp)
	groupMap := make(map[string]struct{})
	for _, g := range grants {
		var id string
		switch g.Grantee.Type {
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			id = g.Grantee.GetGroupId().OpaqueId
			groupMap[id] = struct{}{}
		case provider.GranteeType_GRANTEE_TYPE_USER:
			id = g.Grantee.GetUserId().OpaqueId
		default:
			continue
		}

		if g.Expiration != nil {
			// We are doing this check here because we want to remove expired grants "on access".
			// This way we don't have to have a cron job checking the grants in regular intervals.
			// The tradeof obviously is that this code is here.
			if isGrantExpired(g) {
				err := fs.RemoveGrant(ctx, &provider.Reference{
					ResourceId: &provider.ResourceId{
						SpaceId:  n.SpaceRoot.SpaceID,
						OpaqueId: n.ID},
				}, g)
				appctx.GetLogger(ctx).Error().Err(err).
					Str("space", n.SpaceRoot.ID).
					Str("grantee", id).
					Msg("failed to remove expired space grant")
				continue
			}
			grantExpiration[id] = g.Expiration
		}
		grantMap[id] = g.Permissions
	}

	grantMapJSON, err := json.Marshal(grantMap)
	if err != nil {
		return nil, err
	}

	grantExpirationMapJSON, err := json.Marshal(grantExpiration)
	if err != nil {
		return nil, err
	}

	groupMapJSON, err := json.Marshal(groupMap)
	if err != nil {
		return nil, err
	}

	ssID, err := storagespace.FormatReference(
		&provider.Reference{
			ResourceId: &provider.ResourceId{
				SpaceId:  n.SpaceRoot.SpaceID,
				OpaqueId: n.ID},
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
					Value:   grantMapJSON,
				},
				"grants_expirations": {
					Decoder: "json",
					Value:   grantExpirationMapJSON,
				},
				"groups": {
					Decoder: "json",
					Value:   groupMapJSON,
				},
			},
		},
		Id: &provider.StorageSpaceId{OpaqueId: ssID},
		Root: &provider.ResourceId{
			SpaceId:  n.SpaceRoot.SpaceID,
			OpaqueId: n.ID,
		},
		Name: sname,
		// SpaceType is read from xattr below
		// Mtime is set either as node.tmtime or as fi.mtime below
	}

	space.SpaceType, err = n.SpaceRoot.XattrString(prefixes.SpaceTypeAttr)
	if err != nil {
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
	} else if fi, err := os.Stat(n.InternalPath()); err == nil {
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

	spaceAttributes, err := n.Xattrs()
	if err != nil {
		return nil, err
	}

	// if quota is set try parsing it as int64, otherwise don't bother
	if q, err := spaceAttributes.Int64(prefixes.QuotaAttr); err == nil && q >= 0 {
		// make sure we have a proper signed int
		// we use the same magic numbers to indicate:
		// -1 = uncalculated
		// -2 = unknown
		// -3 = unlimited
		space.Quota = &provider.Quota{
			QuotaMaxBytes: uint64(q),
			QuotaMaxFiles: math.MaxUint64, // TODO MaxUInt64? = unlimited? why even max files? 0 = unlimited?
		}
	}
	if si := spaceAttributes.String(prefixes.SpaceImageAttr); si != "" {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "image", storagespace.FormatResourceID(
			provider.ResourceId{StorageId: space.Root.StorageId, SpaceId: space.Root.SpaceId, OpaqueId: si},
		))
	}
	if sd := spaceAttributes.String(prefixes.SpaceDescriptionAttr); sd != "" {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "description", sd)
	}
	if sr := spaceAttributes.String(prefixes.SpaceReadmeAttr); sr != "" {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "readme", storagespace.FormatResourceID(
			provider.ResourceId{StorageId: space.Root.StorageId, SpaceId: space.Root.SpaceId, OpaqueId: sr},
		))
	}
	if sa := spaceAttributes.String(prefixes.SpaceAliasAttr); sa != "" {
		space.Opaque = utils.AppendPlainToOpaque(space.Opaque, "spaceAlias", sa)
	}

	// add rootinfo
	ps, _ := n.SpaceRoot.PermissionSet(ctx)
	space.RootInfo, _ = n.SpaceRoot.AsResourceInfo(ctx, &ps, nil, nil, false)
	return space, nil
}

func mapHasKey(checkMap map[string][]byte, keys ...string) bool {
	for _, key := range keys {
		if _, hasKey := checkMap[key]; hasKey {
			return true
		}
	}
	return false
}

func isGrantExpired(g *provider.Grant) bool {
	if g.Expiration == nil {
		return false
	}
	return time.Now().After(time.Unix(int64(g.Expiration.Seconds), int64(g.Expiration.Nanos)))
}

func (fs *Decomposedfs) getSpaceRoot(spaceID string) string {
	return filepath.Join(fs.o.Root, "spaces", lookup.Pathify(spaceID, 1, 2))
}
