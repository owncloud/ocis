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

//go:build ceph
// +build ceph

package cephfs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"

	cephfs2 "github.com/ceph/go-ceph/cephfs"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/maxymania/go-system/posix_acl"
)

var perms = map[rune][]string{
	'r': {
		"Stat",
		"GetPath",
		"GetQuota",
		"InitiateFileDownload",
		"ListGrants",
	},
	'w': {
		"AddGrant",
		"CreateContainer",
		"Delete",
		"InitiateFileUpload",
		"Move",
		"RemoveGrant",
		"PurgeRecycle",
		"RestoreFileVersion",
		"RestoreRecycleItem",
		"UpdateGrant",
	},
	'x': {
		"ListRecycle",
		"ListContainer",
		"ListFileVersions",
	},
}

const (
	aclXattr = "system.posix_acl_access"
)

var op2int = map[rune]uint16{'r': 4, 'w': 2, 'x': 1}

func getPermissionSet(user *User, stat *cephfs2.CephStatx, mount Mount, path string) (perm *provider.ResourcePermissions) {
	perm = &provider.ResourcePermissions{}

	if int64(stat.Uid) == user.UidNumber || int64(stat.Gid) == user.GidNumber {
		updatePerms(perm, "rwx", false)
		return
	}

	acls := &posix_acl.Acl{}
	var xattr []byte
	var err error
	if xattr, err = mount.GetXattr(path, aclXattr); err != nil {
		return nil
	}
	acls.Decode(xattr)

	group, err := user.fs.getGroupByID(user.ctx, fmt.Sprint(stat.Gid))

	for _, acl := range acls.List {
		rwx := strings.Split(acl.String(), ":")[2]
		switch acl.GetType() {
		case posix_acl.ACL_USER:
			if int64(acl.GetID()) == user.UidNumber {
				updatePerms(perm, rwx, false)
			}
		case posix_acl.ACL_GROUP:
			if int64(acl.GetID()) == user.GidNumber || in(group.GroupName, user.Groups) {
				updatePerms(perm, rwx, false)
			}
		case posix_acl.ACL_MASK:
			updatePerms(perm, rwx, true)
		case posix_acl.ACL_OTHERS:
			updatePerms(perm, rwx, false)
		}
	}

	return
}

func (fs *cephfs) getFullPermissionSet(ctx context.Context, mount Mount, path string) (permList []*provider.Grant) {
	acls := &posix_acl.Acl{}
	var xattr []byte
	var err error
	if xattr, err = mount.GetXattr(path, aclXattr); err != nil {
		return nil
	}
	acls.Decode(xattr)

	for _, acl := range acls.List {
		rwx := strings.Split(acl.String(), ":")[2]
		switch acl.GetType() {
		case posix_acl.ACL_USER:
			user, err := fs.getUserByID(ctx, fmt.Sprint(acl.GetID()))
			if err != nil {
				return nil
			}
			userGrant := &provider.Grant{
				Grantee: &provider.Grantee{
					Type: provider.GranteeType_GRANTEE_TYPE_USER,
					Id:   &provider.Grantee_UserId{UserId: user.Id},
				},
				Permissions: &provider.ResourcePermissions{},
			}
			updatePerms(userGrant.Permissions, rwx, false)
			permList = append(permList, userGrant)
		case posix_acl.ACL_GROUP:
			group, err := fs.getGroupByID(ctx, fmt.Sprint(acl.GetID()))
			if err != nil {
				return nil
			}
			groupGrant := &provider.Grant{
				Grantee: &provider.Grantee{
					Type: provider.GranteeType_GRANTEE_TYPE_GROUP,
					Id:   &provider.Grantee_GroupId{GroupId: group.Id},
				},
				Permissions: &provider.ResourcePermissions{},
			}
			updatePerms(groupGrant.Permissions, rwx, false)
			permList = append(permList, groupGrant)
		}
	}

	return
}

/*
func permToIntRefl(p *provider.ResourcePermissions) (result uint16) {
	if p == nil { return 0b111 } //rwx

	item := reflect.ValueOf(p).Elem()
	for _, op := range "rwx" {
		for _, perm := range perms[op] {
			if item.FieldByName(perm).Bool() {
				result |= op2int[op]
				break //if value is 1 then bitwise OR can never change it again
			}
		}
	}

	return
}
*/

func permToInt(rp *provider.ResourcePermissions) (result uint16) {
	if rp == nil {
		return 0b111 // rwx
	}
	if rp.Stat || rp.GetPath || rp.GetQuota || rp.ListGrants || rp.InitiateFileDownload {
		result |= 4
	}
	if rp.CreateContainer || rp.Move || rp.Delete || rp.InitiateFileUpload || rp.AddGrant || rp.UpdateGrant ||
		rp.RemoveGrant || rp.DenyGrant || rp.RestoreFileVersion || rp.PurgeRecycle || rp.RestoreRecycleItem {
		result |= 2
	}
	if rp.ListRecycle || rp.ListContainer || rp.ListFileVersions {
		result |= 1
	}

	return
}

const (
	updateGrant = iota
	removeGrant = iota
)

func (fs *cephfs) changePerms(ctx context.Context, mt Mount, grant *provider.Grant, path string, method int) (err error) {
	buf, err := mt.GetXattr(path, aclXattr)
	if err != nil {
		return
	}
	acls := &posix_acl.Acl{}
	acls.Decode(buf)
	var sid posix_acl.AclSID

	switch grant.Grantee.Type {
	case provider.GranteeType_GRANTEE_TYPE_USER:
		var user *userpb.User
		if user, err = fs.getUserByOpaqueID(ctx, grant.Grantee.GetUserId().OpaqueId); err != nil {
			return
		}
		sid.SetUid(uint32(user.UidNumber))
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		var group *grouppb.Group
		if group, err = fs.getGroupByOpaqueID(ctx, grant.Grantee.GetGroupId().OpaqueId); err != nil {
			return
		}
		sid.SetGid(uint32(group.GidNumber))
	default:
		return errors.New("cephfs: invalid grantee type")
	}

	var found = false
	var i int
	for i = range acls.List {
		if acls.List[i].AclSID == sid {
			found = true
		}
	}

	if method == updateGrant {
		if found {
			acls.List[i].Perm |= permToInt(grant.Permissions)
			if acls.List[i].Perm == 0 { // remove empty grant
				acls.List = append(acls.List[:i], acls.List[i+1:]...)
			}
		} else {
			acls.List = append(acls.List, posix_acl.AclElement{
				AclSID: sid,
				Perm:   permToInt(grant.Permissions),
			})
		}
	} else { //removeGrant
		if found {
			acls.List[i].Perm &^= permToInt(grant.Permissions) //bitwise and-not, to clear bits on Perm
			if acls.List[i].Perm == 0 {                        // remove empty grant
				acls.List = append(acls.List[:i], acls.List[i+1:]...)
			}
		}
	}

	err = mt.SetXattr(path, aclXattr, acls.Encode(), 0)

	return
}

/*
func updatePermsRefl(rp *provider.ResourcePermissions, acl string, unset bool) {
	if rp == nil { return }
	for _, t := range "rwx" {
		if strings.ContainsRune(acl, t) {
			for _, i := range perms[t] {
				reflect.ValueOf(rp).Elem().FieldByName(i).SetBool(true)
			}
		} else if unset {
			for _, i := range perms[t] {
				reflect.ValueOf(rp).Elem().FieldByName(i).SetBool(false)
			}
		}
	}
}
*/

func updatePerms(rp *provider.ResourcePermissions, acl string, unset bool) {
	if rp == nil {
		return
	}
	if strings.ContainsRune(acl, 'r') {
		rp.Stat = true
		rp.GetPath = true
		rp.GetQuota = true
		rp.InitiateFileDownload = true
		rp.ListGrants = true
	} else if unset {
		rp.Stat = false
		rp.GetPath = false
		rp.GetQuota = false
		rp.InitiateFileDownload = false
		rp.ListGrants = false
	}
	if strings.ContainsRune(acl, 'w') {
		rp.AddGrant = true
		rp.DenyGrant = true
		rp.CreateContainer = true
		rp.Delete = true
		rp.InitiateFileUpload = true
		rp.Move = true
		rp.RemoveGrant = true
		rp.PurgeRecycle = true
		rp.RestoreFileVersion = true
		rp.RestoreRecycleItem = true
		rp.UpdateGrant = true
	} else if unset {
		rp.AddGrant = false
		rp.DenyGrant = false
		rp.CreateContainer = false
		rp.Delete = false
		rp.InitiateFileUpload = false
		rp.Move = false
		rp.RemoveGrant = false
		rp.PurgeRecycle = false
		rp.RestoreFileVersion = false
		rp.RestoreRecycleItem = false
		rp.UpdateGrant = false
	}
	if strings.ContainsRune(acl, 'x') {
		rp.ListRecycle = true
		rp.ListContainer = true
		rp.ListFileVersions = true
	} else if unset {
		rp.ListRecycle = false
		rp.ListContainer = false
		rp.ListFileVersions = false
	}
}
