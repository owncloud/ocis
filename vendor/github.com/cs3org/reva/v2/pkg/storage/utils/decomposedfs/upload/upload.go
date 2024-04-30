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

package upload

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/upload")
}

// Tree is used to manage a tree hierarchy
type Tree interface {
	Setup() error

	GetMD(ctx context.Context, node *node.Node) (os.FileInfo, error)
	ListFolder(ctx context.Context, node *node.Node) ([]*node.Node, error)
	// CreateHome(owner *userpb.UserId) (n *node.Node, err error)
	CreateDir(ctx context.Context, node *node.Node) (err error)
	// CreateReference(ctx context.Context, node *node.Node, targetURI *url.URL) error
	Move(ctx context.Context, oldNode *node.Node, newNode *node.Node) (err error)
	Delete(ctx context.Context, node *node.Node) (err error)
	RestoreRecycleItemFunc(ctx context.Context, spaceid, key, trashPath string, target *node.Node) (*node.Node, *node.Node, func() error, error)
	PurgeRecycleItemFunc(ctx context.Context, spaceid, key, purgePath string) (*node.Node, func() error, error)

	WriteBlob(node *node.Node, binPath string) error
	ReadBlob(node *node.Node) (io.ReadCloser, error)
	DeleteBlob(node *node.Node) error

	Propagate(ctx context.Context, node *node.Node, sizeDiff int64) (err error)
}

var defaultFilePerm = os.FileMode(0664)

// WriteChunk writes the stream from the reader to the given offset of the upload
func (session *OcisSession) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	ctx, span := tracer.Start(session.Context(ctx), "WriteChunk")
	defer span.End()
	_, subspan := tracer.Start(ctx, "os.OpenFile")
	file, err := os.OpenFile(session.binPath(), os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	subspan.End()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// calculate cheksum here? needed for the TUS checksum extension. https://tus.io/protocols/resumable-upload.html#checksum
	// TODO but how do we get the `Upload-Checksum`? WriteChunk() only has a context, offset and the reader ...
	// It is sent with the PATCH request, well or in the POST when the creation-with-upload extension is used
	// but the tus handler uses a context.Background() so we cannot really check the header and put it in the context ...
	_, subspan = tracer.Start(ctx, "io.Copy")
	n, err := io.Copy(file, src)
	subspan.End()

	// If the HTTP PATCH request gets interrupted in the middle (e.g. because
	// the user wants to pause the upload), Go's net/http returns an io.ErrUnexpectedEOF.
	// However, for the ocis driver it's not important whether the stream has ended
	// on purpose or accidentally.
	if err != nil && err != io.ErrUnexpectedEOF {
		return n, err
	}

	// update upload.Session.Offset so subsequent code flow can use it.
	// No need to persist the session as the offset is determined by stating the blob in the GetUpload / ReadSession codepath.
	// The session offset is written to disk in FinishUpload
	session.info.Offset += n
	return n, nil
}

// GetInfo returns the FileInfo
func (session *OcisSession) GetInfo(_ context.Context) (tusd.FileInfo, error) {
	return session.ToFileInfo(), nil
}

// GetReader returns an io.Reader for the upload
func (session *OcisSession) GetReader(ctx context.Context) (io.Reader, error) {
	_, span := tracer.Start(session.Context(ctx), "GetReader")
	defer span.End()
	return os.Open(session.binPath())
}

// FinishUpload finishes an upload and moves the file to the internal destination
func (session *OcisSession) FinishUpload(ctx context.Context) error {
	ctx, span := tracer.Start(session.Context(ctx), "FinishUpload")
	defer span.End()
	log := appctx.GetLogger(ctx)

	// calculate the checksum of the written bytes
	// they will all be written to the metadata later, so we cannot omit any of them
	// TODO only calculate the checksum in sync that was requested to match, the rest could be async ... but the tests currently expect all to be present
	// TODO the hashes all implement BinaryMarshaler so we could try to persist the state for resumable upload. we would neet do keep track of the copied bytes ...
	sha1h := sha1.New()
	md5h := md5.New()
	adler32h := adler32.New()
	{
		_, subspan := tracer.Start(ctx, "os.Open")
		f, err := os.Open(session.binPath())
		subspan.End()
		if err != nil {
			// we can continue if no oc checksum header is set
			log.Info().Err(err).Str("binPath", session.binPath()).Msg("error opening binPath")
		}
		defer f.Close()

		r1 := io.TeeReader(f, sha1h)
		r2 := io.TeeReader(r1, md5h)

		_, subspan = tracer.Start(ctx, "io.Copy")
		_, err = io.Copy(adler32h, r2)
		subspan.End()
		if err != nil {
			log.Info().Err(err).Msg("error copying checksums")
		}
	}

	// compare if they match the sent checksum
	// TODO the tus checksum extension would do this on every chunk, but I currently don't see an easy way to pass in the requested checksum. for now we do it in FinishUpload which is also called for chunked uploads
	var err error
	if session.info.MetaData["checksum"] != "" {
		var err error
		parts := strings.SplitN(session.info.MetaData["checksum"], " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		switch parts[0] {
		case "sha1":
			err = checkHash(parts[1], sha1h)
		case "md5":
			err = checkHash(parts[1], md5h)
		case "adler32":
			err = checkHash(parts[1], adler32h)
		default:
			err = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if err != nil {
			session.store.Cleanup(ctx, session, true, false, false)
			return err
		}
	}

	// update checksums
	attrs := node.Attributes{
		prefixes.ChecksumPrefix + "sha1":    sha1h.Sum(nil),
		prefixes.ChecksumPrefix + "md5":     md5h.Sum(nil),
		prefixes.ChecksumPrefix + "adler32": adler32h.Sum(nil),
	}

	n, err := session.store.CreateNodeForUpload(session, attrs)
	if err != nil {
		session.store.Cleanup(ctx, session, true, false, false)
		return err
	}

	// increase the processing counter for every started processing
	// will be decreased in Cleanup()
	metrics.UploadProcessing.Inc()
	metrics.UploadSessionsBytesReceived.Inc()

	if session.store.pub != nil {
		u, _ := ctxpkg.ContextGetUser(ctx)
		s, err := session.URL(ctx)
		if err != nil {
			return err
		}

		if err := events.Publish(ctx, session.store.pub, events.BytesReceived{
			UploadID:      session.ID(),
			URL:           s,
			SpaceOwner:    n.SpaceOwnerOrManager(session.Context(ctx)),
			ExecutingUser: u,
			ResourceID:    &provider.ResourceId{SpaceId: n.SpaceID, OpaqueId: n.ID},
			Filename:      session.Filename(),
			Filesize:      uint64(session.Size()),
		}); err != nil {
			return err
		}
	}

	if !session.store.async {
		// handle postprocessing synchronously
		err = session.Finalize()
		session.store.Cleanup(ctx, session, err != nil, false, err == nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to upload")
			return err
		}
		metrics.UploadSessionsFinalized.Inc()
	}

	return session.store.tp.Propagate(ctx, n, session.SizeDiff())
}

// Terminate terminates the upload
func (session *OcisSession) Terminate(_ context.Context) error {
	session.Cleanup(true, true, true)
	return nil
}

// DeclareLength updates the upload length information
func (session *OcisSession) DeclareLength(ctx context.Context, length int64) error {
	session.info.Size = length
	session.info.SizeIsDeferred = false
	return session.Persist(session.Context(ctx))
}

// ConcatUploads concatenates multiple uploads
func (session *OcisSession) ConcatUploads(_ context.Context, uploads []tusd.Upload) (err error) {
	file, err := os.OpenFile(session.binPath(), os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, partialUpload := range uploads {
		fileUpload := partialUpload.(*OcisSession)

		src, err := os.Open(fileUpload.binPath())
		if err != nil {
			return err
		}
		defer src.Close()

		if _, err := io.Copy(file, src); err != nil {
			return err
		}
	}

	return
}

// Finalize finalizes the upload (eg moves the file to the internal destination)
func (session *OcisSession) Finalize() (err error) {
	ctx, span := tracer.Start(session.Context(context.Background()), "Finalize")
	defer span.End()

	revisionNode := &node.Node{SpaceID: session.SpaceID(), BlobID: session.ID(), Blobsize: session.Size()}

	// upload the data to the blobstore
	_, subspan := tracer.Start(ctx, "WriteBlob")
	err = session.store.tp.WriteBlob(revisionNode, session.binPath())
	subspan.End()
	if err != nil {
		return errors.Wrap(err, "failed to upload file to blobstore")
	}

	return nil
}

func checkHash(expected string, h hash.Hash) error {
	hash := hex.EncodeToString(h.Sum(nil))
	if expected != hash {
		return errtypes.ChecksumMismatch(fmt.Sprintf("invalid checksum: expected %s got %x", expected, hash))
	}
	return nil
}

func (session *OcisSession) removeNode(ctx context.Context) {
	n, err := session.Node(ctx)
	if err != nil {
		appctx.GetLogger(ctx).Error().Str("session", session.ID()).Err(err).Msg("getting node from session failed")
		return
	}
	if err := n.Purge(ctx); err != nil {
		appctx.GetLogger(ctx).Error().Str("nodepath", n.InternalPath()).Err(err).Msg("purging node failed")
	}
}

// cleanup cleans up after the upload is finished
func (session *OcisSession) Cleanup(revertNodeMetadata, cleanBin, cleanInfo bool) {
	ctx := session.Context(context.Background())

	if revertNodeMetadata {
		n, err := session.Node(ctx)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("node", n.ID).Str("sessionid", session.ID()).Msg("reading node for session failed")
		}
		if session.NodeExists() {
			p := session.info.MetaData["versionsPath"]
			if err := session.store.lu.CopyMetadata(ctx, p, n.InternalPath(), func(attributeName string, value []byte) (newValue []byte, copy bool) {
				return value, strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
					attributeName == prefixes.TypeAttr ||
					attributeName == prefixes.BlobIDAttr ||
					attributeName == prefixes.BlobsizeAttr ||
					attributeName == prefixes.MTimeAttr
			}, true); err != nil {
				appctx.GetLogger(ctx).Info().Str("versionpath", p).Str("nodepath", n.InternalPath()).Err(err).Msg("renaming version node failed")
			}

			if err := os.RemoveAll(p); err != nil {
				appctx.GetLogger(ctx).Info().Str("versionpath", p).Str("nodepath", n.InternalPath()).Err(err).Msg("error removing version")
			}

		} else {
			// if no other upload session is in progress (processing id != session id) or has finished (processing id == "")
			latestSession, err := n.ProcessingID(ctx)
			if err != nil {
				appctx.GetLogger(ctx).Error().Err(err).Str("node", n.ID).Str("sessionid", session.ID()).Msg("reading processingid for session failed")
			}
			if latestSession == session.ID() {
				// actually delete the node
				session.removeNode(ctx)
			}
			// FIXME else if the upload has become a revision, delete the revision, or if it is the last one, delete the node
		}
	}

	if cleanBin {
		if err := os.Remove(session.binPath()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			appctx.GetLogger(ctx).Error().Str("path", session.binPath()).Err(err).Msg("removing upload failed")
		}
	}

	if cleanInfo {
		if err := session.Purge(ctx); err != nil && !errors.Is(err, fs.ErrNotExist) {
			appctx.GetLogger(ctx).Error().Err(err).Str("session", session.ID()).Msg("removing upload info failed")
		}
	}
}

// URL returns a url to download an upload
func (session *OcisSession) URL(_ context.Context) (string, error) {
	type transferClaims struct {
		jwt.StandardClaims
		Target string `json:"target"`
	}

	u := joinurl(session.store.tknopts.DownloadEndpoint, "tus/", session.ID())
	ttl := time.Duration(session.store.tknopts.TransferExpires) * time.Second
	claims := transferClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Audience:  "reva",
			IssuedAt:  time.Now().Unix(),
		},
		Target: u,
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tkn, err := t.SignedString([]byte(session.store.tknopts.TransferSharedSecret))
	if err != nil {
		return "", errors.Wrapf(err, "error signing token with claims %+v", claims)
	}

	return joinurl(session.store.tknopts.DataGatewayEndpoint, tkn), nil
}

// replace with url.JoinPath after switching to go1.19
func joinurl(paths ...string) string {
	var s strings.Builder
	l := len(paths)
	for i, p := range paths {
		s.WriteString(p)
		if !strings.HasSuffix(p, "/") && i != l-1 {
			s.WriteString("/")
		}
	}

	return s.String()
}
