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

package eventsmiddleware

import (
	"time"

	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
)

// ContainerCreated converts the response to an event
func ContainerCreated(r *provider.CreateContainerResponse, req *provider.CreateContainerRequest, spaceOwner, executant *user.UserId) events.ContainerCreated {
	return events.ContainerCreated{
		SpaceOwner: spaceOwner,
		Executant:  executant,
		Ref:        req.Ref,
		Timestamp:  utils.TSNow(),
	}
}

// ShareCreated converts the response to an event
func ShareCreated(r *collaboration.CreateShareResponse, executant *user.UserId) events.ShareCreated {
	return events.ShareCreated{
		ShareID:        r.Share.GetId(),
		Executant:      executant,
		Sharer:         r.Share.Creator,
		GranteeUserID:  r.Share.GetGrantee().GetUserId(),
		GranteeGroupID: r.Share.GetGrantee().GetGroupId(),
		ItemID:         r.Share.ResourceId,
		CTime:          r.Share.Ctime,
		Permissions:    r.Share.Permissions,
	}
}

// ShareRemoved converts the response to an event
func ShareRemoved(r *collaboration.RemoveShareResponse, req *collaboration.RemoveShareRequest, executant *user.UserId) events.ShareRemoved {
	var (
		userid  *user.UserId
		groupid *group.GroupId
		rid     *provider.ResourceId
	)
	_ = utils.ReadJSONFromOpaque(r.Opaque, "granteeuserid", &userid)
	_ = utils.ReadJSONFromOpaque(r.Opaque, "granteegroupid", &userid)
	_ = utils.ReadJSONFromOpaque(r.Opaque, "resourceid", &rid)
	return events.ShareRemoved{
		Executant:      executant,
		ShareID:        req.Ref.GetId(),
		ShareKey:       req.Ref.GetKey(),
		GranteeUserID:  userid,
		GranteeGroupID: groupid,
		ItemID:         rid,
		Timestamp:      time.Now(),
	}
}

// ShareUpdated converts the response to an event
func ShareUpdated(r *collaboration.UpdateShareResponse, req *collaboration.UpdateShareRequest, executant *user.UserId) events.ShareUpdated {
	updated := ""
	if req.Field.GetPermissions() != nil {
		updated = "permissions"
	} else if req.Field.GetDisplayName() != "" {
		updated = "displayname"
	}
	return events.ShareUpdated{
		Executant:      executant,
		ShareID:        r.Share.Id,
		ItemID:         r.Share.ResourceId,
		Permissions:    r.Share.Permissions,
		GranteeUserID:  r.Share.GetGrantee().GetUserId(),
		GranteeGroupID: r.Share.GetGrantee().GetGroupId(),
		Sharer:         r.Share.Creator,
		MTime:          r.Share.Mtime,
		Updated:        updated,
	}
}

// ReceivedShareUpdated converts the response to an event
func ReceivedShareUpdated(r *collaboration.UpdateReceivedShareResponse, executant *user.UserId) events.ReceivedShareUpdated {
	return events.ReceivedShareUpdated{
		Executant:      executant,
		ShareID:        r.Share.Share.Id,
		ItemID:         r.Share.Share.ResourceId,
		Permissions:    r.Share.Share.Permissions,
		GranteeUserID:  r.Share.Share.GetGrantee().GetUserId(),
		GranteeGroupID: r.Share.Share.GetGrantee().GetGroupId(),
		Sharer:         r.Share.Share.Creator,
		MTime:          r.Share.Share.Mtime,
		State:          collaboration.ShareState_name[int32(r.Share.State)],
	}
}

// LinkCreated converts the response to an event
func LinkCreated(r *link.CreatePublicShareResponse, executant *user.UserId) events.LinkCreated {
	return events.LinkCreated{
		Executant:         executant,
		ShareID:           r.Share.Id,
		Sharer:            r.Share.Creator,
		ItemID:            r.Share.ResourceId,
		Permissions:       r.Share.Permissions,
		DisplayName:       r.Share.DisplayName,
		Expiration:        r.Share.Expiration,
		PasswordProtected: r.Share.PasswordProtected,
		CTime:             r.Share.Ctime,
		Token:             r.Share.Token,
	}
}

// LinkUpdated converts the response to an event
func LinkUpdated(r *link.UpdatePublicShareResponse, req *link.UpdatePublicShareRequest, executant *user.UserId) events.LinkUpdated {
	return events.LinkUpdated{
		Executant:         executant,
		ShareID:           r.Share.Id,
		Sharer:            r.Share.Creator,
		ItemID:            r.Share.ResourceId,
		Permissions:       r.Share.Permissions,
		DisplayName:       r.Share.DisplayName,
		Expiration:        r.Share.Expiration,
		PasswordProtected: r.Share.PasswordProtected,
		CTime:             r.Share.Ctime,
		Token:             r.Share.Token,
		FieldUpdated:      link.UpdatePublicShareRequest_Update_Type_name[int32(req.Update.GetType())],
	}
}

// LinkAccessed converts the response to an event
func LinkAccessed(r *link.GetPublicShareByTokenResponse, executant *user.UserId) events.LinkAccessed {
	return events.LinkAccessed{
		Executant:         executant,
		ShareID:           r.Share.Id,
		Sharer:            r.Share.Creator,
		ItemID:            r.Share.ResourceId,
		Permissions:       r.Share.Permissions,
		DisplayName:       r.Share.DisplayName,
		Expiration:        r.Share.Expiration,
		PasswordProtected: r.Share.PasswordProtected,
		CTime:             r.Share.Ctime,
		Token:             r.Share.Token,
	}
}

// LinkAccessFailed converts the response to an event
func LinkAccessFailed(r *link.GetPublicShareByTokenResponse, req *link.GetPublicShareByTokenRequest, executant *user.UserId) events.LinkAccessFailed {
	e := events.LinkAccessFailed{
		Executant: executant,
		Status:    r.Status.Code,
		Message:   r.Status.Message,
		Timestamp: utils.TSNow(),
	}
	if r.Share != nil {
		e.ShareID = r.Share.Id
		e.Token = r.Share.Token
	}
	return e
}

// LinkRemoved converts the response to an event
func LinkRemoved(r *link.RemovePublicShareResponse, req *link.RemovePublicShareRequest, executant *user.UserId) events.LinkRemoved {
	return events.LinkRemoved{
		Executant:  executant,
		ShareID:    req.Ref.GetId(),
		ShareToken: req.Ref.GetToken(),
		Timestamp:  utils.TSNow(),
	}
}

// FileTouched converts the response to an event
func FileTouched(r *provider.TouchFileResponse, req *provider.TouchFileRequest, spaceOwner, executant *user.UserId) events.FileTouched {
	return events.FileTouched{
		SpaceOwner: spaceOwner,
		Executant:  executant,
		Ref:        req.Ref,
		Timestamp:  utils.TSNow(),
	}
}

// FileUploaded converts the response to an event
func FileUploaded(r *provider.InitiateFileUploadResponse, req *provider.InitiateFileUploadRequest, spaceOwner, executant *user.UserId) events.FileUploaded {
	return events.FileUploaded{
		SpaceOwner: spaceOwner,
		Executant:  executant,
		Ref:        req.Ref,
		Timestamp:  utils.TSNow(),
	}
}

// FileDownloaded converts the response to an event
func FileDownloaded(r *provider.InitiateFileDownloadResponse, req *provider.InitiateFileDownloadRequest, executant *user.UserId) events.FileDownloaded {
	return events.FileDownloaded{
		Executant: executant,
		Ref:       req.Ref,
		Timestamp: utils.TSNow(),
	}
}

// ItemTrashed converts the response to an event
func ItemTrashed(r *provider.DeleteResponse, req *provider.DeleteRequest, spaceOwner, executant *user.UserId) events.ItemTrashed {
	opaqueID := utils.ReadPlainFromOpaque(r.Opaque, "opaque_id")
	return events.ItemTrashed{
		SpaceOwner: spaceOwner,
		Executant:  executant,
		Ref:        req.Ref,
		ID: &provider.ResourceId{
			StorageId: req.Ref.GetResourceId().GetStorageId(),
			SpaceId:   req.Ref.GetResourceId().GetSpaceId(),
			OpaqueId:  opaqueID,
		},
		Timestamp: utils.TSNow(),
	}
}

// ItemMoved converts the response to an event
func ItemMoved(r *provider.MoveResponse, req *provider.MoveRequest, spaceOwner, executant *user.UserId) events.ItemMoved {
	return events.ItemMoved{
		SpaceOwner:   spaceOwner,
		Executant:    executant,
		Ref:          req.Destination,
		OldReference: req.Source,
		Timestamp:    utils.TSNow(),
	}
}

// ItemPurged converts the response to an event
func ItemPurged(r *provider.PurgeRecycleResponse, req *provider.PurgeRecycleRequest, executant *user.UserId) events.ItemPurged {
	return events.ItemPurged{
		Executant: executant,
		Ref:       req.Ref,
		Timestamp: utils.TSNow(),
	}
}

// ItemRestored converts the response to an event
func ItemRestored(r *provider.RestoreRecycleItemResponse, req *provider.RestoreRecycleItemRequest, spaceOwner, executant *user.UserId) events.ItemRestored {
	ref := req.Ref
	if req.RestoreRef != nil {
		ref = req.RestoreRef
	}
	return events.ItemRestored{
		SpaceOwner:   spaceOwner,
		Executant:    executant,
		Ref:          ref,
		OldReference: req.Ref,
		Key:          req.Key,
		Timestamp:    utils.TSNow(),
	}
}

// FileVersionRestored converts the response to an event
func FileVersionRestored(r *provider.RestoreFileVersionResponse, req *provider.RestoreFileVersionRequest, spaceOwner, executant *user.UserId) events.FileVersionRestored {
	return events.FileVersionRestored{
		SpaceOwner: spaceOwner,
		Executant:  executant,
		Ref:        req.Ref,
		Key:        req.Key,
		Timestamp:  utils.TSNow(),
	}
}

// SpaceCreated converts the response to an event
func SpaceCreated(r *provider.CreateStorageSpaceResponse, executant *user.UserId) events.SpaceCreated {
	return events.SpaceCreated{
		Executant: executant,
		ID:        r.StorageSpace.Id,
		Owner:     extractOwner(r.StorageSpace.Owner),
		Root:      r.StorageSpace.Root,
		Name:      r.StorageSpace.Name,
		Type:      r.StorageSpace.SpaceType,
		Quota:     r.StorageSpace.Quota,
		MTime:     r.StorageSpace.Mtime,
	}
}

// SpaceRenamed converts the response to an event
func SpaceRenamed(r *provider.UpdateStorageSpaceResponse, req *provider.UpdateStorageSpaceRequest, executant *user.UserId) events.SpaceRenamed {
	return events.SpaceRenamed{
		Executant: executant,
		ID:        r.StorageSpace.Id,
		Owner:     extractOwner(r.StorageSpace.Owner),
		Name:      r.StorageSpace.Name,
		Timestamp: utils.TSNow(),
	}
}

// SpaceUpdated converts the response to an event
func SpaceUpdated(r *provider.UpdateStorageSpaceResponse, req *provider.UpdateStorageSpaceRequest, executant *user.UserId) events.SpaceUpdated {
	return events.SpaceUpdated{
		Executant: executant,
		ID:        r.StorageSpace.Id,
		Space:     r.StorageSpace,
		Timestamp: utils.TSNow(),
	}
}

// SpaceEnabled converts the response to an event
func SpaceEnabled(r *provider.UpdateStorageSpaceResponse, req *provider.UpdateStorageSpaceRequest, executant *user.UserId) events.SpaceEnabled {
	return events.SpaceEnabled{
		Executant: executant,
		ID:        r.StorageSpace.Id,
		Owner:     extractOwner(r.StorageSpace.Owner),
		Timestamp: utils.TSNow(),
	}
}

// SpaceShared converts the response to an event
// func SpaceShared(req *provider.AddGrantRequest, executant, sharer *user.UserId, grantee *provider.Grantee) events.SpaceShared {
func SpaceShared(r *provider.AddGrantResponse, req *provider.AddGrantRequest, executant *user.UserId) events.SpaceShared {
	id := storagespace.FormatStorageID(req.Ref.ResourceId.StorageId, req.Ref.ResourceId.SpaceId)
	return events.SpaceShared{
		Executant:      executant,
		Creator:        req.Grant.Creator,
		GranteeUserID:  req.Grant.GetGrantee().GetUserId(),
		GranteeGroupID: req.Grant.GetGrantee().GetGroupId(),
		ID:             &provider.StorageSpaceId{OpaqueId: id},
		Timestamp:      time.Now(),
	}
}

// SpaceUnshared  converts the response to an event
func SpaceUnshared(r *provider.RemoveGrantResponse, req *provider.RemoveGrantRequest, executant *user.UserId) events.SpaceUnshared {
	id := storagespace.FormatStorageID(req.Ref.ResourceId.StorageId, req.Ref.ResourceId.SpaceId)
	return events.SpaceUnshared{
		Executant:      executant,
		GranteeUserID:  req.Grant.GetGrantee().GetUserId(),
		GranteeGroupID: req.Grant.GetGrantee().GetGroupId(),
		ID:             &provider.StorageSpaceId{OpaqueId: id},
		Timestamp:      time.Now(),
	}
}

// SpaceDisabled converts the response to an event
func SpaceDisabled(r *provider.DeleteStorageSpaceResponse, req *provider.DeleteStorageSpaceRequest, executant *user.UserId) events.SpaceDisabled {
	return events.SpaceDisabled{
		Executant: executant,
		ID:        req.Id,
		Timestamp: time.Now(),
	}
}

// SpaceDeleted converts the response to an event
func SpaceDeleted(r *provider.DeleteStorageSpaceResponse, req *provider.DeleteStorageSpaceRequest, executant *user.UserId) events.SpaceDeleted {
	var final map[string]provider.ResourcePermissions
	_ = utils.ReadJSONFromOpaque(r.GetOpaque(), "grants", &final)
	return events.SpaceDeleted{
		Executant:    executant,
		ID:           req.Id,
		SpaceName:    utils.ReadPlainFromOpaque(r.GetOpaque(), "spacename"),
		FinalMembers: final,
		Timestamp:    time.Now(),
	}
}

func extractOwner(u *user.User) *user.UserId {
	if u != nil {
		return u.Id
	}
	return nil
}
