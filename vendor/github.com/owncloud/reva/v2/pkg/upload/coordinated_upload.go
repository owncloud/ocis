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
	"os"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
)

// coordinatedUpload is the object the TUS handler talks to during an upload.
// TUS is a resumable upload protocol: clients upload files in chunks via HTTP PATCH,
// can resume after interruption, and optionally delete (Terminate), declare size
// upfront (DeclareLength), or assemble partial uploads (ConcatUploads).
//
// To plug into the TUS handler, a type must implement:
//   - tusd.Upload: GetInfo, GetReader, WriteChunk, FinishUpload (core read/write)
//   - tusd.TerminatableUpload: Terminate (DELETE support)
//   - tusd.LengthDeclarableUpload: DeclareLength (deferred-length uploads)
//   - tusd.ConcatableUpload: ConcatUploads (parallel chunk concatenation)
//
// coordinatedUpload owns all upload lifecycle logic (checksums, MarkProcessing,
// event publishing). It delegates raw data access to the Session.
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
	ref := u.session.Reference()
	u.session.Cleanup(true, true)
	_ = u.coord.fs.MarkProcessing(ctx, &ref, false, u.session.ID())
	if !u.session.NodeExists() {
		_, _ = u.coord.fs.Delete(ctx, &ref)
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
	if err := checksumAndFinish(ctx, u.session); err != nil {
		u.coord.rollback(ctx, u.session)
		return err
	}
	// Persist checksums so the postprocessing handler can read them after BytesReceived.
	if err := u.session.Persist(ctx); err != nil {
		u.coord.rollback(ctx, u.session)
		return err
	}

	metrics.UploadProcessing.Inc()
	metrics.UploadSessionsBytesReceived.Inc()

	if !u.coord.async {
		return u.coord.commitSync(ctx, u.session)
	}

	if u.session.Size() > 0 {
		s, err := u.session.URL(ctx)
		if err != nil {
			u.coord.rollback(ctx, u.session)
			return err
		}
		if err := events.Publish(ctx, u.coord.pub, events.BytesReceived{
			UploadID:   u.session.ID(),
			URL:        s,
			SpaceOwner: u.session.SpaceOwner(),
			ExecutingUser: &user.User{
				Id: &user.UserId{
					Type:     u.session.Executant().Type,
					Idp:      u.session.Executant().Idp,
					OpaqueId: u.session.Executant().OpaqueId,
				},
			},
			ResourceID: &provider.ResourceId{
				StorageId: u.session.ProviderID(),
				SpaceId:   u.session.SpaceID(),
				OpaqueId:  u.session.NodeID(),
			},
			Filename:          u.session.Filename(),
			Filesize:          uint64(u.session.Size()),
			ImpersonatingUser: impersonatingUser(ctx),
		}); err != nil {
			u.coord.rollback(ctx, u.session)
			return err
		}
	}
	return nil
}

// checksumAndFinish computes and validates checksums, then stores them on the session.
// Used by both the TUS and simple PUT paths.
func checksumAndFinish(ctx context.Context, session Session) error {
	sha1h, md5h, adler32h, err := calculateChecksums(ctx, session.BinPath())
	if err != nil {
		return err
	}
	info, err := session.GetInfo(ctx)
	if err != nil {
		return err
	}
	if checksum := info.MetaData["checksum"]; checksum != "" {
		parts := strings.SplitN(checksum, " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		var checkErr error
		switch parts[0] {
		case "sha1":
			checkErr = checkHash(parts[1], sha1h)
		case "md5":
			checkErr = checkHash(parts[1], md5h)
		case "adler32":
			checkErr = checkHash(parts[1], adler32h)
		default:
			checkErr = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if checkErr != nil {
			session.Cleanup(true, true)
			return checkErr
		}
	}
	session.SetChecksums(sha1h.Sum(nil), md5h.Sum(nil), adler32h.Sum(nil))
	return nil
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
