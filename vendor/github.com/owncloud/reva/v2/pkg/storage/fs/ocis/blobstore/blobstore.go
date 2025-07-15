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

package blobstore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

// ErrBlobIDEmpty is returned when the BlobID is empty
var ErrBlobIDEmpty = fmt.Errorf("blobstore: BlobID is empty")

// Blobstore provides an interface to an filesystem based blobstore
type Blobstore struct {
	root string
}

// New returns a new Blobstore
func New(root string) (*Blobstore, error) {
	err := os.MkdirAll(root, 0700)
	if err != nil {
		return nil, err
	}

	return &Blobstore{
		root: root,
	}, nil
}

// Upload stores some data in the blobstore under the given key
func (bs *Blobstore) Upload(node *node.Node, source string) error {
	if node.BlobID == "" {
		return ErrBlobIDEmpty
	}

	dest := bs.Path(node)

	// ensure parent path exists
	if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: oCIS blobstore: error creating parent folders for blob")
	}

	if err := os.Rename(source, dest); err == nil {
		return nil
	}

	// Rename failed, file needs to be copied.
	file, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: oCIS blobstore: Can not open source file to upload")
	}
	defer file.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "could not open blob '%s' for writing", dest)
	}

	w := bufio.NewWriter(f)
	_, err = w.ReadFrom(file)
	if err != nil {
		return errors.Wrapf(err, "could not write blob '%s'", dest)
	}

	return w.Flush()
}

// Download retrieves a blob from the blobstore for reading
func (bs *Blobstore) Download(node *node.Node) (io.ReadCloser, error) {
	if node.BlobID == "" {
		return nil, ErrBlobIDEmpty
	}

	dest := bs.Path(node)
	file, err := os.Open(dest)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read blob '%s'", dest)
	}
	return file, nil
}

// Delete deletes a blob from the blobstore
func (bs *Blobstore) Delete(node *node.Node) error {
	if node.BlobID == "" {
		return ErrBlobIDEmpty
	}
	dest := bs.Path(node)
	if err := utils.RemoveItem(dest); err != nil {
		return errors.Wrapf(err, "could not delete blob '%s'", dest)
	}
	return nil
}

// List lists all blobs in the Blobstore
func (bs *Blobstore) List() ([]*node.Node, error) {
	dirs, err := filepath.Glob(filepath.Join(bs.root, "spaces", "*", "*", "blobs", "*", "*", "*", "*", "*"))
	if err != nil {
		return nil, err
	}
	blobids := make([]*node.Node, 0, len(dirs))
	for _, d := range dirs {
		_, s, _ := strings.Cut(d, "spaces")
		spaceraw, blobraw, _ := strings.Cut(s, "blobs")
		blobids = append(blobids, &node.Node{
			SpaceID: strings.ReplaceAll(spaceraw, "/", ""),
			BlobID:  strings.ReplaceAll(blobraw, "/", ""),
		})
	}
	return blobids, nil
}

func (bs *Blobstore) Path(node *node.Node) string {
	return filepath.Join(
		bs.root,
		filepath.Clean(filepath.Join(
			"/", "spaces", lookup.Pathify(node.SpaceID, 1, 2), "blobs", lookup.Pathify(node.BlobID, 4, 2)),
		),
	)
}
