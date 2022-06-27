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

package eosfs

import (
	"context"
	"io"
	"os"
	"path"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/pkg/errors"
)

func (fs *eosfs) Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, _ storage.UploadFinishedFunc) error {
	p, err := fs.resolve(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "eos: error resolving reference")
	}

	if fs.isShareFolder(ctx, p) {
		return errtypes.PermissionDenied("eos: cannot upload under the virtual share folder")
	}

	if chunking.IsChunked(p) {
		var assembledFile string
		p, assembledFile, err = fs.chunkHandler.WriteChunk(p, r)
		if err != nil {
			return err
		}
		if p == "" {
			return errtypes.PartialContent(ref.String())
		}
		fd, err := os.Open(assembledFile)
		if err != nil {
			return errors.Wrap(err, "eos: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)
		r = fd
	}

	fn := fs.wrap(ctx, p)

	u, err := getUser(ctx)
	if err != nil {
		return errors.Wrap(err, "eos: no user in ctx")
	}

	// We need the auth corresponding to the parent directory
	// as the file might not exist at the moment
	auth, err := fs.getUserAuth(ctx, u, path.Dir(fn))
	if err != nil {
		return err
	}
	return fs.c.Write(ctx, auth, fn, r)
}

func (fs *eosfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	fn, _, err := fs.resolveRefAndGetAuth(ctx, ref)
	if err != nil {
		return nil, err
	}
	fn, err = fs.unwrap(ctx, fn)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"simple": fn,
	}, nil
}
