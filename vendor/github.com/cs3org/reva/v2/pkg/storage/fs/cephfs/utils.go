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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cephfs2 "github.com/ceph/go-ceph/cephfs"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// Mount type
type Mount = *cephfs2.MountInfo

// Statx type
type Statx = *cephfs2.CephStatx

var dirPermFull = uint32(0777)
var dirPermDefault = uint32(0775)
var filePermDefault = uint32(0660)

func closeDir(directory *cephfs2.Directory) {
	if directory != nil {
		_ = directory.Close()
	}
}

func closeFile(file *cephfs2.File) {
	if file != nil {
		_ = file.Close()
	}
}

func destroyCephConn(mt Mount, perm *cephfs2.UserPerm) *cacheVal {
	if perm != nil {
		perm.Destroy()
	}
	if mt != nil {
		_ = mt.Release()
	}
	return nil
}

func deleteFile(mount *cephfs2.MountInfo, path string) {
	_ = mount.Unlink(path)
}

func isDir(t provider.ResourceType) bool {
	return t == provider.ResourceType_RESOURCE_TYPE_CONTAINER
}

func (fs *cephfs) makeFIDPath(fid string) string {
	return "" //filepath.Join(fs.conf.EIDFolder, fid)
}

func (fs *cephfs) makeFID(absolutePath string, inode string) (rid *provider.ResourceId, err error) {
	sum := md5.New()
	sum.Write([]byte(absolutePath))
	fid := fmt.Sprintf("%s-%s", hex.EncodeToString(sum.Sum(nil)), inode)
	rid = &provider.ResourceId{OpaqueId: fid}

	_ = fs.adminConn.adminMount.Link(absolutePath, fs.makeFIDPath(fid))
	_ = fs.adminConn.adminMount.SetXattr(absolutePath, xattrEID, []byte(fid), 0)

	return
}

func (fs *cephfs) getFIDPath(cv *cacheVal, path string) (fid string, err error) {
	var buffer []byte
	if buffer, err = cv.mount.GetXattr(path, xattrEID); err != nil {
		return
	}

	return fs.makeFIDPath(string(buffer)), err
}

func calcChecksum(filepath string, mt Mount, stat Statx) (checksum string, err error) {
	file, err := mt.Open(filepath, os.O_RDONLY, 0)
	defer closeFile(file)
	if err != nil {
		return
	}
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	checksum = hex.EncodeToString(hash.Sum(nil))
	// we don't care if they fail, the checksum will just be recalculated if an error happens
	_ = mt.SetXattr(filepath, xattrMd5ts, []byte(strconv.FormatInt(stat.Mtime.Sec, 10)), 0)
	_ = mt.SetXattr(filepath, xattrMd5, []byte(checksum), 0)

	return
}

func resolveRevRef(mt Mount, ref *provider.Reference, revKey string) (str string, err error) {
	var buf []byte
	if ref.GetResourceId() != nil {
		str, err = mt.Readlink(filepath.Join(snap, revKey, ref.ResourceId.OpaqueId))
		if err != nil {
			return "", fmt.Errorf("cephfs: invalid reference %+v", ref)
		}
	} else if str = ref.GetPath(); str != "" {
		buf, err = mt.GetXattr(str, xattrEID)
		if err != nil {
			return
		}
		str, err = mt.Readlink(filepath.Join(snap, revKey, string(buf)))
		if err != nil {
			return
		}
	} else {
		return "", fmt.Errorf("cephfs: empty reference %+v", ref)
	}

	return filepath.Join(snap, revKey, str), err
}

func removeLeadingSlash(path string) string {
	return filepath.Join(".", path)
}

func addLeadingSlash(path string) string {
	return filepath.Join("/", path)
}

func in(lookup string, list []string) bool {
	for _, item := range list {
		if item == lookup {
			return true
		}
	}
	return false
}

func pathGenerator(path string, reverse bool, str chan string) {
	if reverse {
		str <- path
		for i := range path {
			if path[len(path)-i-1] == filepath.Separator {
				str <- path[:len(path)-i-1]
			}
		}
	} else {
		for i := range path {
			if path[i] == filepath.Separator {
				str <- path[:i]
			}
		}
		str <- path
	}

	close(str)
}

func walkPath(path string, f func(string) error, reverse bool) (err error) {
	paths := make(chan string)
	go pathGenerator(path, reverse, paths)
	for path := range paths {
		if path == "" {
			continue
		}
		if err = f(path); err != nil && err.Error() != errFileExists && err.Error() != errNotFound {
			break
		} else {
			err = nil
		}
	}

	return
}

func (fs *cephfs) writeIndex(oid string, value string) (err error) {
	return fs.adminConn.radosIO.WriteFull(oid, []byte(value))
}

func (fs *cephfs) removeIndex(oid string) error {
	return fs.adminConn.radosIO.Delete(oid)
}

func (fs *cephfs) resolveIndex(oid string) (fullPath string, err error) {
	var i int
	var currPath strings.Builder
	root := string(filepath.Separator)
	offset := uint64(0)
	io := fs.adminConn.radosIO
	bsize := 4096
	buffer := make([]byte, bsize)
	for {
		for { //read object
			i, err = io.Read(oid, buffer, offset)
			offset += uint64(bsize)
			currPath.Write(buffer)
			if err == nil && i >= bsize {
				buffer = buffer[:0]
				continue
			} else {
				offset = 0
				break
			}
		}
		if err != nil {
			return
		}

		ss := strings.SplitN(currPath.String(), string(filepath.Separator), 2)
		if len(ss) != 2 {
			if currPath.String() == root {
				return
			}

			return "", fmt.Errorf("cephfs: entry id is not in the form of \"parentID/entryname\"")
		}
		parentOID := ss[0]
		entryName := ss[1]
		fullPath = filepath.Join(entryName, fullPath)
		oid = parentOID
		currPath.Reset()
	}
}
