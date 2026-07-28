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

// This file holds the whole tusd adapter surface: the per-upload
// coordinatedUpload type plus the coordinator methods that exist only to satisfy
// tusd's DataStore interfaces. Upload logic itself lives in coordinator.go.

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

var errNotImplemented = tusd.NewError("ERR_NOT_IMPLEMENTED", "use InitiateUpload on the CS3 API to start a new upload", http.StatusNotImplemented)

// coordinatedUpload adapts a single upload session to the tusd.Upload interface
// family.
//
// It exists because tusd's per-upload methods (WriteChunk, FinishUpload) carry no
// upload id, so the receiver itself must identify the upload. The coordinator is
// shared process-wide and cannot hold that state without racing between
// concurrent uploads.
//
// This type only translates; all upload logic lives in coordinator.
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

// FinishUpload runs the coordinator's finish path once all bytes have arrived.
//
// The context tusd provides carries no user, so the session rebuilds the one
// recorded at initiate time. Errors are mapped to tusd errors because tusd turns
// error types it does not recognise into a bare 500.
func (u *coordinatedUpload) FinishUpload(ctx context.Context) error {
	err := u.coord.finishUpload(u.session.Context(ctx), u.session)

	switch err.(type) {
	case nil:
		return nil
	case errtypes.AlreadyExists:
		return tusd.NewError("ERR_ALREADY_EXISTS", err.Error(), http.StatusConflict)
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

// Terminate discards an upload: it drops the staged files and, when this upload
// created the node, removes it again so a cancelled upload leaves nothing behind.
func (u *coordinatedUpload) Terminate(ctx context.Context) error {
	u.session.Cleanup(true, true)

	// Terminate can run before the node was created, leaving nothing to undo.
	ref := u.session.Reference()
	if ref.GetResourceId().GetOpaqueId() == "" {
		return nil
	}

	_ = u.coord.fs.MarkProcessing(ctx, &ref, false, u.session.ID())
	if !u.session.NodeExists() {
		_, _ = u.coord.fs.Delete(ctx, &ref)
	}
	return nil
}

// DeclareLength records the total size for uploads initiated without one
// (creation-defer-length), so the finish path knows when all bytes have arrived.
func (u *coordinatedUpload) DeclareLength(ctx context.Context, length int64) error {
	u.session.SetSize(length)
	u.session.SetSizeIsDeferred(false)
	return u.session.Persist(ctx)
}

// ConcatUploads appends the staged bytes of the partial uploads to this upload,
// in the order given, implementing the TUS concatenation extension.
func (u *coordinatedUpload) ConcatUploads(ctx context.Context, partials []tusd.Upload) error {
	file, err := os.OpenFile(u.session.BinPath(), os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, partial := range partials {
		cu, ok := partial.(*coordinatedUpload)
		if !ok {
			return fmt.Errorf("coordinator: unexpected partial upload type %T", partial)
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

// UseIn registers the coordinator as tusd's data store, in place of the driver's
// own store. This is what makes the TUS data path run through the coordinator for
// every driver, so finishing an upload ends in driver.CommitUpload.
//
// The opt-in extensions mirror what the drivers registered before, so no TUS
// capability is lost: without them tusd answers those requests with 501.
func (c *coordinator) UseIn(composer *tusd.StoreComposer) {
	composer.UseCore(c)
	composer.UseTerminater(c)
	composer.UseConcater(c)
	composer.UseLengthDeferrer(c)
}

// NewUpload is unsupported: tusd may not create sessions on its own. Uploads are
// always started through the CS3 InitiateUpload call, which performs the
// permission, quota and lock checks that tusd knows nothing about.
func (c *coordinator) NewUpload(_ context.Context, _ tusd.FileInfo) (tusd.Upload, error) {
	return nil, errNotImplemented
}

// GetUpload loads the session with the given id and wraps it so the TUS data
// path runs through the coordinator. This is the only place a coordinatedUpload
// is constructed.
func (c *coordinator) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	session, err := c.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &coordinatedUpload{session: session, coord: c}, nil
}

// The As* methods let tusd reach the extension interfaces on an upload it holds as
// a plain tusd.Upload. Every upload we hand out is a *coordinatedUpload, which
// implements all three.

func (c *coordinator) AsTerminatableUpload(up tusd.Upload) tusd.TerminatableUpload {
	return up.(*coordinatedUpload)
}

func (c *coordinator) AsLengthDeclarableUpload(up tusd.Upload) tusd.LengthDeclarableUpload {
	return up.(*coordinatedUpload)
}

func (c *coordinator) AsConcatableUpload(up tusd.Upload) tusd.ConcatableUpload {
	return up.(*coordinatedUpload)
}
