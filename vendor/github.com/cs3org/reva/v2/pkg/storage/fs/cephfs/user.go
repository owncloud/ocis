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
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/cs3org/reva/v2/pkg/errtypes"

	cephfs2 "github.com/ceph/go-ceph/cephfs"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctx2 "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/storage/utils/templates"
	"github.com/pkg/errors"
)

type callBack func(cb *cacheVal)

// User custom type to add functionality to current struct
type User struct {
	*userv1beta1.User
	fs   *cephfs
	ctx  context.Context
	home string
}

func (fs *cephfs) makeUser(ctx context.Context) *User {
	u := ctx2.ContextMustGetUser(ctx)
	home := filepath.Join(fs.conf.Root, templates.WithUser(u, fs.conf.UserLayout))
	return &User{u, fs, ctx, home}
}

func (user *User) absPath(path string) string {
	//shares will always be absolute to avoid prepending the user path to the path of the file's owner
	if !filepath.IsAbs(path) {
		path = filepath.Join(user.home, path)
	}

	return path
}

func (user *User) op(cb callBack) {
	conn := user.fs.conn
	if err := conn.lock.Acquire(conn.ctx, 1); err != nil {
		return
	}
	defer conn.lock.Release(1)

	val, found := conn.cache.Get(user.Id.OpaqueId)
	if !found {
		cvalue := newConn(user)
		if cvalue != nil {
			conn.cache.Set(user.Id.OpaqueId, cvalue, 1)
		} else {
			return
		}
		cb(cvalue)
		return
	}

	cb(val.(*cacheVal))
}

func (user *User) fileAsResourceInfo(cv *cacheVal, path string, stat *cephfs2.CephStatx, mdKeys []string) (ri *provider.ResourceInfo, err error) {
	var (
		_type  provider.ResourceType
		target string
		size   uint64
		buf    []byte
	)

	switch int(stat.Mode) & syscall.S_IFMT {
	case syscall.S_IFDIR:
		_type = provider.ResourceType_RESOURCE_TYPE_CONTAINER
		if buf, err = cv.mount.GetXattr(path, "ceph.dir.rbytes"); err == nil {
			size, err = strconv.ParseUint(string(buf), 10, 64)
		}
	case syscall.S_IFLNK:
		_type = provider.ResourceType_RESOURCE_TYPE_SYMLINK
		target, err = cv.mount.Readlink(path)
	case syscall.S_IFREG:
		_type = provider.ResourceType_RESOURCE_TYPE_FILE
		size = stat.Size
	default:
		return nil, errors.New("cephfs: unknown entry type")
	}

	if err != nil {
		return
	}

	var xattrs []string
	keys := make(map[string]bool, len(mdKeys))
	for _, key := range mdKeys {
		keys[key] = true
	}
	if keys["*"] || len(keys) == 0 {
		mdKeys = []string{}
		keys = map[string]bool{}
	}
	mx := make(map[string]string)
	if xattrs, err = cv.mount.ListXattr(path); err == nil {
		for _, xattr := range xattrs {
			if len(mdKeys) == 0 || keys[xattr] {
				if buf, err := cv.mount.GetXattr(path, xattr); err == nil {
					mx[xattr] = string(buf)
				}
			}
		}
	}

	//TODO(tmourati): Add entry id logic here

	var etag string
	if isDir(_type) {
		rctime, _ := cv.mount.GetXattr(path, "ceph.dir.rctime")
		etag = fmt.Sprint(stat.Inode) + ":" + string(rctime)
	} else {
		etag = fmt.Sprint(stat.Inode) + ":" + strconv.FormatInt(stat.Ctime.Sec, 10)
	}

	mtime := &typesv1beta1.Timestamp{
		Seconds: uint64(stat.Mtime.Sec),
		Nanos:   uint32(stat.Mtime.Nsec),
	}

	perms := getPermissionSet(user, stat, cv.mount, path)

	for key := range mx {
		if !strings.HasPrefix(key, xattrUserNs) {
			delete(mx, key)
		}
	}

	var checksum provider.ResourceChecksum
	var md5 string
	if _type == provider.ResourceType_RESOURCE_TYPE_FILE {
		md5tsBA, err := cv.mount.GetXattr(path, xattrMd5ts) //local error inside if scope
		if err == nil {
			md5ts, _ := strconv.ParseInt(string(md5tsBA), 10, 64)
			if stat.Mtime.Sec == md5ts {
				md5BA, err := cv.mount.GetXattr(path, xattrMd5)
				if err != nil {
					md5, err = calcChecksum(path, cv.mount, stat)
				} else {
					md5 = string(md5BA)
				}
			} else {
				md5, err = calcChecksum(path, cv.mount, stat)
			}
		} else {
			md5, err = calcChecksum(path, cv.mount, stat)
		}

		if err != nil && err.Error() == errPermissionDenied {
			checksum.Type = provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_UNSET
		} else if err != nil {
			return nil, errors.New("cephfs: error calculating checksum of file")
		} else {
			checksum.Type = provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_MD5
			checksum.Sum = md5
		}
	} else {
		checksum.Type = provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_UNSET
	}

	var ownerID *userv1beta1.UserId
	if stat.Uid != 0 {
		var owner *userv1beta1.User
		if int64(stat.Uid) != user.UidNumber {
			owner, err = user.fs.getUserByID(user.ctx, fmt.Sprint(stat.Uid))
		} else {
			owner = user.User
		}

		if owner == nil {
			return nil, errors.New("cephfs: error getting owner of entry: " + path)
		}

		ownerID = owner.Id
	} else {
		ownerID = &userv1beta1.UserId{OpaqueId: "root"}
	}

	ri = &provider.ResourceInfo{
		Type:              _type,
		Id:                &provider.ResourceId{OpaqueId: fmt.Sprint(stat.Inode)},
		Checksum:          &checksum,
		Etag:              etag,
		MimeType:          mime.Detect(isDir(_type), path),
		Mtime:             mtime,
		Path:              path,
		PermissionSet:     perms,
		Size:              size,
		Owner:             ownerID,
		Target:            target,
		ArbitraryMetadata: &provider.ArbitraryMetadata{Metadata: mx},
	}

	return
}

func (user *User) resolveRef(ref *provider.Reference) (str string, err error) {
	if ref == nil {
		return "", fmt.Errorf("cephfs: nil reference")
	}

	if str = ref.GetPath(); str == "" {
		return "", errtypes.NotSupported("cephfs: entry IDs not currently supported")
	}

	str = removeLeadingSlash(str) //path must be relative

	return
}
