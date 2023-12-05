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

package tree

import (
	"errors"
	"io"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
)

//go:generate make --no-print-directory -C ../../../../.. mockery NAME=Blobstore

// Blobstore defines an interface for storing blobs in a blobstore
type Blobstore interface {
	Upload(node *node.Node, source string) error
	Download(node *node.Node) (io.ReadCloser, error)
	Delete(node *node.Node) error
}

// BlobstoreMover is used to move a file from the upload to the final destination
type BlobstoreMover interface {
	MoveBlob(n *node.Node, source, bucket, key string) error
}

var ErrBlobstoreCannotMove = errors.New("cannot move")
