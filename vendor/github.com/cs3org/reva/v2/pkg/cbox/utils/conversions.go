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

package utils

import (
	"strings"
	"time"

	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
)

// DBShare stores information about user and public shares.
type DBShare struct {
	ID           string
	UIDOwner     string
	UIDInitiator string
	Prefix       string
	ItemSource   string
	ItemType     string
	ShareWith    string
	Token        string
	Expiration   string
	Permissions  int
	ShareType    int
	ShareName    string
	STime        int
	FileTarget   string
	State        int
}

// FormatGrantee formats a CS3API grantee to a string
func FormatGrantee(g *provider.Grantee) (int, string) {
	var granteeType int
	var formattedID string
	switch g.Type {
	case provider.GranteeType_GRANTEE_TYPE_USER:
		granteeType = 0
		formattedID = FormatUserID(g.GetUserId())
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		granteeType = 1
		formattedID = FormatGroupID(g.GetGroupId())
	default:
		granteeType = -1
	}
	return granteeType, formattedID
}

// ExtractGrantee retrieves the CS3API grantee from a formatted string
func ExtractGrantee(t int, g string) *provider.Grantee {
	var grantee provider.Grantee
	switch t {
	case 0:
		grantee.Type = provider.GranteeType_GRANTEE_TYPE_USER
		grantee.Id = &provider.Grantee_UserId{UserId: ExtractUserID(g)}
	case 1:
		grantee.Type = provider.GranteeType_GRANTEE_TYPE_GROUP
		grantee.Id = &provider.Grantee_GroupId{GroupId: ExtractGroupID(g)}
	default:
		grantee.Type = provider.GranteeType_GRANTEE_TYPE_INVALID
	}
	return &grantee
}

// ResourceTypeToItem maps a resource type to a string
func ResourceTypeToItem(r provider.ResourceType) string {
	switch r {
	case provider.ResourceType_RESOURCE_TYPE_FILE:
		return "file"
	case provider.ResourceType_RESOURCE_TYPE_CONTAINER:
		return "folder"
	case provider.ResourceType_RESOURCE_TYPE_REFERENCE:
		return "reference"
	case provider.ResourceType_RESOURCE_TYPE_SYMLINK:
		return "symlink"
	default:
		return ""
	}
}

// ResourceTypeToItemInt maps a resource type to an integer
func ResourceTypeToItemInt(r provider.ResourceType) int {
	switch r {
	case provider.ResourceType_RESOURCE_TYPE_CONTAINER:
		return 0
	case provider.ResourceType_RESOURCE_TYPE_FILE:
		return 1
	default:
		return -1
	}
}

// SharePermToInt maps read/write permissions to an integer
func SharePermToInt(p *provider.ResourcePermissions) int {
	var perm int
	switch {
	case p.InitiateFileUpload && !p.InitiateFileDownload:
		perm = 4
	case p.InitiateFileUpload:
		perm = 15
	case p.InitiateFileDownload:
		perm = 1
	}
	// TODO map denials and resharing; currently, denials are mapped to 0
	return perm
}

// IntTosharePerm retrieves read/write permissions from an integer
func IntTosharePerm(p int, itemType string) *provider.ResourcePermissions {
	perms, _ := conversions.NewPermissions(p)
	return conversions.RoleFromOCSPermissions(perms).CS3ResourcePermissions()
}

// IntToShareState retrieves the received share state from an integer
func IntToShareState(g int) collaboration.ShareState {
	switch g {
	case 0:
		return collaboration.ShareState_SHARE_STATE_PENDING
	case 1:
		return collaboration.ShareState_SHARE_STATE_ACCEPTED
	case -1:
		return collaboration.ShareState_SHARE_STATE_REJECTED
	default:
		return collaboration.ShareState_SHARE_STATE_INVALID
	}
}

// FormatUserID formats a CS3API user ID to a string
func FormatUserID(u *userpb.UserId) string {
	return u.OpaqueId
}

// ExtractUserID retrieves a CS3API user ID from a string
func ExtractUserID(u string) *userpb.UserId {
	t := userpb.UserType_USER_TYPE_PRIMARY
	if strings.HasPrefix(u, "guest:") {
		t = userpb.UserType_USER_TYPE_LIGHTWEIGHT
	} else if strings.Contains(u, "@") {
		t = userpb.UserType_USER_TYPE_FEDERATED
	}
	return &userpb.UserId{OpaqueId: u, Type: t}
}

// FormatGroupID formats a CS3API group ID to a string
func FormatGroupID(u *grouppb.GroupId) string {
	return u.OpaqueId
}

// ExtractGroupID retrieves a CS3API group ID from a string
func ExtractGroupID(u string) *grouppb.GroupId {
	return &grouppb.GroupId{OpaqueId: u}
}

// ConvertToCS3Share converts a DBShare to a CS3API collaboration share
func ConvertToCS3Share(s DBShare) *collaboration.Share {
	ts := &typespb.Timestamp{
		Seconds: uint64(s.STime),
	}
	return &collaboration.Share{
		Id: &collaboration.ShareId{
			OpaqueId: s.ID,
		},
		//ResourceId:  &provider.Reference{StorageId: s.Prefix, NodeId: s.ItemSource},
		ResourceId: &provider.ResourceId{
			SpaceId:  s.Prefix,
			OpaqueId: s.ItemSource,
		},
		Permissions: &collaboration.SharePermissions{Permissions: IntTosharePerm(s.Permissions, s.ItemType)},
		Grantee:     ExtractGrantee(s.ShareType, s.ShareWith),
		Owner:       ExtractUserID(s.UIDOwner),
		Creator:     ExtractUserID(s.UIDInitiator),
		Ctime:       ts,
		Mtime:       ts,
	}
}

// ConvertToCS3ReceivedShare converts a DBShare to a CS3API collaboration received share
func ConvertToCS3ReceivedShare(s DBShare) *collaboration.ReceivedShare {
	return &collaboration.ReceivedShare{
		Share:      ConvertToCS3Share(s),
		State:      IntToShareState(s.State),
		MountPoint: &provider.Reference{Path: strings.TrimLeft(s.FileTarget, "/")},
	}
}

// ConvertToCS3PublicShare converts a DBShare to a CS3API public share
func ConvertToCS3PublicShare(s DBShare) *link.PublicShare {
	ts := &typespb.Timestamp{
		Seconds: uint64(s.STime),
	}
	pwd := false
	if s.ShareWith != "" {
		pwd = true
	}
	var expires *typespb.Timestamp
	if s.Expiration != "" {
		t, err := time.Parse("2006-01-02 15:04:05", s.Expiration)
		if err == nil {
			expires = &typespb.Timestamp{
				Seconds: uint64(t.Unix()),
			}
		}
	}
	return &link.PublicShare{
		Id: &link.PublicShareId{
			OpaqueId: s.ID,
		},
		ResourceId: &provider.ResourceId{
			SpaceId:  s.Prefix,
			OpaqueId: s.ItemSource,
		},
		Permissions:       &link.PublicSharePermissions{Permissions: IntTosharePerm(s.Permissions, s.ItemType)},
		Owner:             ExtractUserID(s.UIDOwner),
		Creator:           ExtractUserID(s.UIDInitiator),
		Token:             s.Token,
		DisplayName:       s.ShareName,
		PasswordProtected: pwd,
		Expiration:        expires,
		Ctime:             ts,
		Mtime:             ts,
	}
}
