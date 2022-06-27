// Copyright 2018-2022 CERN
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
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

// Disk represents a disk metadata storage
type Disk struct {
	dataDir string
}

// NewDiskStorage returns a new disk storage instance
func NewDiskStorage(dataDir string) (s Storage, err error) {
	return &Disk{
		dataDir: dataDir,
	}, nil
}

// Init creates the metadata space
func (disk *Disk) Init(_ context.Context, _ string) (err error) {
	return os.MkdirAll(disk.dataDir, 0777)
}

// Backend returns the backend name of the storage
func (disk *Disk) Backend() string {
	return "disk"
}

// SimpleUpload stores a file on disk
func (disk *Disk) SimpleUpload(_ context.Context, uploadpath string, content []byte) error {
	return os.WriteFile(disk.targetPath(uploadpath), content, 0644)
}

// SimpleDownload reads a file from disk
func (disk *Disk) SimpleDownload(_ context.Context, downloadpath string) ([]byte, error) {
	return os.ReadFile(disk.targetPath(downloadpath))
}

// Delete deletes a path
func (disk *Disk) Delete(_ context.Context, path string) error {
	return os.Remove(disk.targetPath(path))
}

// ReadDir returns the resource infos in a given directory
func (disk *Disk) ReadDir(_ context.Context, p string) ([]string, error) {
	infos, err := ioutil.ReadDir(disk.targetPath(p))
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			return []string{}, nil
		}
		return nil, err
	}

	entries := make([]string, 0, len(infos))
	for _, entry := range infos {
		entries = append(entries, path.Join(p, entry.Name()))
	}
	return entries, nil
}

// MakeDirIfNotExist will create a root node in the metadata storage. Requires an authenticated context.
func (disk *Disk) MakeDirIfNotExist(_ context.Context, path string) error {
	return os.MkdirAll(disk.targetPath(path), 0777)
}

// CreateSymlink creates a symlink
func (disk *Disk) CreateSymlink(_ context.Context, oldname, newname string) error {
	return os.Symlink(oldname, disk.targetPath(newname))
}

// ResolveSymlink resolves a symlink
func (disk *Disk) ResolveSymlink(_ context.Context, path string) (string, error) {
	return os.Readlink(disk.targetPath(path))
}

func (disk *Disk) targetPath(p string) string {
	return path.Join(disk.dataDir, p)
}
