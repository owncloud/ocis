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

package storage

import (
	"context"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

type skipTouchPropagationKey struct{}

// ContextSkipTouchPropagation tells TouchFile not to propagate.
func ContextSkipTouchPropagation(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTouchPropagationKey{}, true)
}

// SkipTouchPropagation reports whether TouchFile should skip propagation.
func SkipTouchPropagation(ctx context.Context) bool {
	skip, _ := ctx.Value(skipTouchPropagationKey{}).(bool)
	return skip
}

// OrphanChecker defines the interface for FS implementations that can resolve a resource's metadata.
type OrphanChecker interface {
	// IsOrphaned reports whether the referenced resource exists but its metadata is unreadable.
	IsOrphaned(ctx context.Context, ref *provider.Reference) bool
}

// NodeCreator defines the interface for FS implementations whose PrepareUpload
// creates a new file's node from UploadInfo.ParentID and Name, so an upload need
// not TouchFile it first.
type NodeCreator interface {
	// PrepareCreatesNode reports whether PrepareUpload creates a missing node.
	PrepareCreatesNode() bool
}

// UploadSession is the interface that storage drivers need to return whan listing upload sessions.
type UploadSession interface {
	// ID returns the upload id
	ID() string
	// Filename returns the filename of the file
	Filename() string
	// Size returns the size of the upload
	Size() int64
	// Offset returns the current offset
	Offset() int64
	// Reference returns a reference for the file being uploaded. May be absolute id based or relative to e.g. a space root
	Reference() provider.Reference
	// Executant returns the userid of the user that created the upload
	Executant() userpb.UserId
	// SpaceOwner returns the owner of a space if set. optional
	SpaceOwner() *userpb.UserId
	// Expires returns the time when the upload can no longer be used
	Expires() time.Time

	// IsProcessing returns true if postprocessing has not finished, yet
	// The actual postprocessing state is tracked in the postprocessing service.
	IsProcessing() bool

	// Purge allows completely removing an upload.
	Purge(ctx context.Context)

	// ScanData returns the scan data for the UploadSession
	ScanData() (string, time.Time)
}

// UploadSessionFilter can be used to filter upload sessions
type UploadSessionFilter struct {
	ID         *string
	Processing *bool
	Expired    *bool
	HasVirus   *bool
	// Orphaned filters sessions by whether their target node can still be
	// resolved. Evaluating it requires reading the node metadata of every
	// session, so it is only evaluated when set.
	Orphaned *bool
}
