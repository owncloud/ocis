// Copyright 2018-2023 CERN
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

package metadata

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/filelocks"
	"github.com/pkg/errors"
	"github.com/pkg/xattr"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// XattrsBackend stores the file attributes in extended attributes
type XattrsBackend struct{}

// Name returns the name of the backend
func (XattrsBackend) Name() string { return "xattrs" }

// Get an extended attribute value for the given key
// No file locking is involved here as reading a single xattr is
// considered to be atomic.
func (b XattrsBackend) Get(filePath, key string) ([]byte, error) {
	return xattr.Get(filePath, key)
}

// GetInt64 reads a string as int64 from the xattrs
func (b XattrsBackend) GetInt64(filePath, key string) (int64, error) {
	attr, err := b.Get(filePath, key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(string(attr), 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// List retrieves a list of names of extended attributes associated with the
// given path in the file system.
func (XattrsBackend) List(filePath string) (attribs []string, err error) {
	attrs, err := xattr.List(filePath)
	if err == nil {
		return attrs, nil
	}

	f, err := lockedfile.OpenFile(filePath+filelocks.LockFileSuffix, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer cleanupLockfile(f)

	return xattr.List(filePath)
}

// All reads all extended attributes for a node, protected by a
// shared file lock
func (b XattrsBackend) All(filePath string) (attribs map[string][]byte, err error) {
	attrNames, err := b.List(filePath)

	if err != nil {
		return nil, err
	}

	var (
		xerrs = 0
		xerr  error
	)
	// error handling: Count if there are errors while reading all attribs.
	// if there were any, return an error.
	attribs = make(map[string][]byte, len(attrNames))
	for _, name := range attrNames {
		var val []byte
		if val, xerr = xattr.Get(filePath, name); xerr != nil {
			xerrs++
		} else {
			attribs[name] = val
		}
	}

	if xerrs > 0 {
		err = errors.Wrap(xerr, "Failed to read all xattrs")
	}

	return attribs, err
}

// Set sets one attribute for the given path
func (b XattrsBackend) Set(path string, key string, val []byte) (err error) {
	return b.SetMultiple(path, map[string][]byte{key: val}, true)
}

// SetMultiple sets a set of attribute for the given path
func (XattrsBackend) SetMultiple(path string, attribs map[string][]byte, acquireLock bool) (err error) {
	if acquireLock {
		err := os.MkdirAll(filepath.Dir(path), 0600)
		if err != nil {
			return err
		}
		lockedFile, err := lockedfile.OpenFile(path+filelocks.LockFileSuffix, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer cleanupLockfile(lockedFile)
	}

	// error handling: Count if there are errors while setting the attribs.
	// if there were any, return an error.
	var (
		xerrs = 0
		xerr  error
	)
	for key, val := range attribs {
		if xerr = xattr.Set(path, key, val); xerr != nil {
			// log
			xerrs++
		}
	}
	if xerrs > 0 {
		return errors.Wrap(xerr, "Failed to set all xattrs")
	}

	return nil
}

// Remove an extended attribute key
func (XattrsBackend) Remove(filePath string, key string) (err error) {
	lockedFile, err := lockedfile.OpenFile(filePath+filelocks.LockFileSuffix, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer cleanupLockfile(lockedFile)

	return xattr.Remove(filePath, key)
}

// IsMetaFile returns whether the given path represents a meta file
func (XattrsBackend) IsMetaFile(path string) bool { return strings.HasSuffix(path, ".meta.lock") }

// Purge purges the data of a given path
func (XattrsBackend) Purge(path string) error { return nil }

// Rename moves the data for a given path to a new path
func (XattrsBackend) Rename(oldPath, newPath string) error { return nil }

// MetadataPath returns the path of the file holding the metadata for the given path
func (XattrsBackend) MetadataPath(path string) string { return path }

// LockfilePath returns the path of the lock file
func (XattrsBackend) LockfilePath(path string) string { return path + ".mlock" }

func cleanupLockfile(f *lockedfile.File) {
	_ = f.Close()
	_ = os.Remove(f.Name())
}

// AllWithLockedSource reads all extended attributes from the given reader.
// The path argument is used for storing the data in the cache
func (b XattrsBackend) AllWithLockedSource(path string, _ io.Reader) (map[string][]byte, error) {
	return b.All(path)
}
