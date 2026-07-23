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

package upload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/owncloud/reva/v2/pkg/errtypes"
)

// coordinatedUpload is the TUS interface adapter between the tusd library and the coordinator.
// It implements tusd.Upload, TerminatableUpload, LengthDeclarableUpload, and ConcatableUpload.
// All real upload logic lives in coordinator; this type only translates tusd callbacks.
type coordinatedUpload struct {
	session Session
	coord   *coordinator
}

func (u *coordinatedUpload) GetInfo(ctx context.Context) (tusd.FileInfo, error) {
	return u.session.GetInfo(ctx)
}

func (u *coordinatedUpload) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return u.session.GetReader(ctx)
}

func (u *coordinatedUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	return u.session.WriteChunk(ctx, offset, src)
}

func (u *coordinatedUpload) Terminate(ctx context.Context) error {
	u.session.Cleanup(true, true)
	// Node may not exist if Terminate is called before FinishUpload.
	ref := u.session.Reference()
	if ref.ResourceId.GetOpaqueId() != "" {
		_ = u.coord.fs.MarkProcessing(ctx, &ref, false, u.session.ID())
		if !u.session.NodeExists() {
			_, _ = u.coord.fs.Delete(ctx, &ref)
		}
	}
	return nil
}

func (u *coordinatedUpload) DeclareLength(ctx context.Context, length int64) error {
	u.session.SetSize(length)
	u.session.SetSizeIsDeferred(false)
	return u.session.Persist(ctx)
}

func (u *coordinatedUpload) ConcatUploads(ctx context.Context, partials []tusd.Upload) error {
	file, err := os.OpenFile(u.session.BinPath(), os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, partial := range partials {
		cu, ok := partial.(*coordinatedUpload)
		if !ok {
			return fmt.Errorf("coordinator: unexpected partial type %T", partial)
		}
		src, err := cu.session.GetReader(ctx)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(file, src)
		src.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func (u *coordinatedUpload) FinishUpload(ctx context.Context) error {
	err := u.coord.finishUpload(u.session.Context(ctx), u.session)
	switch err.(type) {
	case nil:
		return nil
	case errtypes.ResourceProcessing, errtypes.TooEarly:
		return tusd.NewError("ERR_TOO_EARLY", err.Error(), http.StatusTooEarly)
	case errtypes.Aborted:
		return tusd.NewError("ERR_PRECONDITION_FAILED", err.Error(), http.StatusPreconditionFailed)
	case errtypes.PreconditionFailed:
		return tusd.NewError("ERR_PRECONDITION_FAILED", err.Error(), http.StatusMethodNotAllowed)
	case errtypes.Locked:
		return tusd.NewError("ERR_LOCKED", err.Error(), http.StatusLocked)
	case errtypes.BadRequest:
		return tusd.NewError("ERR_BAD_REQUEST", err.Error(), http.StatusBadRequest)
	case errtypes.ChecksumMismatch:
		return tusd.NewError("ERR_CHECKSUM_MISMATCH", err.Error(), errtypes.StatusChecksumMismatch)
	default:
		return err
	}
}

// UseIn registers the coordinator as the TUS data store in the composer.
func (c *coordinator) UseIn(composer *tusd.StoreComposer) {
	composer.UseCore(c)
	composer.UseTerminater(c)
	composer.UseConcater(c)
	composer.UseLengthDeferrer(c)
}

// NewUpload is not supported; uploads are initiated via the CS3 API.
func (c *coordinator) NewUpload(_ context.Context, _ tusd.FileInfo) (tusd.Upload, error) {
	return nil, errNotImplemented
}

// GetUpload returns the upload session wrapped in a coordinatedUpload so the
// TUS FinishUpload hook runs the coordinator path rather than the legacy one.
func (c *coordinator) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	session, err := c.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &coordinatedUpload{session: session, coord: c}, nil
}

func (c *coordinator) AsTerminatableUpload(up tusd.Upload) tusd.TerminatableUpload {
	return up.(*coordinatedUpload)
}

func (c *coordinator) AsLengthDeclarableUpload(up tusd.Upload) tusd.LengthDeclarableUpload {
	return up.(*coordinatedUpload)
}

func (c *coordinator) AsConcatableUpload(up tusd.Upload) tusd.ConcatableUpload {
	return up.(*coordinatedUpload)
}
