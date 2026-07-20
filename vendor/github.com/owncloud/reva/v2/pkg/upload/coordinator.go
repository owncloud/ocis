// Copyright 2018-2024 CERN
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

// Package upload provides the driver-agnostic upload coordinator.
package upload

import (
	"context"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	"github.com/owncloud/reva/v2/pkg/storage"
)

// Coordinator owns the upload lifecycle.
type Coordinator interface {
	// InitiateUpload returns a list of protocols with urls that can be used to append bytes to a new upload session.
	InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error)
}

// coordinator is the concrete implementation of Coordinator.
type coordinator struct {
	fs storage.FS
}

// NewCoordinator constructs a coordinator backed by the given storage driver.
func NewCoordinator(fs storage.FS) *coordinator {
	return &coordinator{fs: fs}
}

// InitiateUpload returns a list of protocols with urls that can be used to append bytes to a new upload session.
//
// For now this delegates straight to the underlying storage driver, preserving
// existing behaviour. It lets us wire the coordinator into the storageprovider
// and exercise the seam end-to-end before porting the driver-agnostic upload
// logic into the coordinator itself.
func (c *coordinator) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	return c.fs.InitiateUpload(ctx, ref, uploadLength, metadata)
}
