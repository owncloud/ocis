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
	"encoding/json"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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

// Upload processes the upload
// it implements tus tusd.Upload interface https://tus.io/protocols/resumable-upload.html#core-protocol
// it also implements its termination extension as specified in https://tus.io/protocols/resumable-upload.html#termination
// it also implements its creation-defer-length extension as specified in https://tus.io/protocols/resumable-upload.html#creation
// it also implements its concatenation extension as specified in https://tus.io/protocols/resumable-upload.html#concatenation
type Upload struct {
	// we use a struct field on the upload as tus pkg will give us an empty context.Background
	Ctx context.Context
	// info stores the current information about the upload
	Info tusd.FileInfo
	// node for easy access
	Node *node.Node
	// infoPath is the path to the .info file
	infoPath string
	// binPath is the path to the binary file (which has no extension)
	binPath string
	// lu and tp needed for file operations
	lu *lookup.Lookup
	tp Tree
	// versionsPath will be empty if there was no file before
	versionsPath string
	// sizeDiff size difference between new and old file version
	sizeDiff int64
	// and a logger as well
	log zerolog.Logger
	// publisher used to publish events
	pub events.Publisher
	// async determines if uploads shoud be done asynchronously
	async bool
	// tknopts hold token signing information
	tknopts options.TokenOptions
}

func buildUpload(ctx context.Context, info tusd.FileInfo, binPath string, infoPath string, lu *lookup.Lookup, tp Tree, pub events.Publisher, async bool, tknopts options.TokenOptions) *Upload {
	return &Upload{
		Info:     info,
		binPath:  binPath,
		infoPath: infoPath,
		lu:       lu,
		tp:       tp,
		Ctx:      ctx,
		pub:      pub,
		async:    async,
		tknopts:  tknopts,
		log: appctx.GetLogger(ctx).
			With().
			Interface("info", info).
			Str("binPath", binPath).
			Logger(),
	}
}

// Cleanup cleans the upload
func Cleanup(upload *Upload, failure bool, keepUpload bool) {
	ctx, span := tracer.Start(upload.Ctx, "Cleanup")
	defer span.End()
	upload.cleanup(failure, !keepUpload, !keepUpload)

	// unset processing status
	if upload.Node != nil { // node can be nil when there was an error before it was created (eg. checksum-mismatch)
		if err := upload.Node.UnmarkProcessing(ctx, upload.Info.ID); err != nil {
			upload.log.Info().Str("path", upload.Node.InternalPath()).Err(err).Msg("unmarking processing failed")
		}
	}
}

// WriteChunk writes the stream from the reader to the given offset of the upload
func (upload *Upload) WriteChunk(_ context.Context, offset int64, src io.Reader) (int64, error) {
	ctx, span := tracer.Start(upload.Ctx, "WriteChunk")
	defer span.End()
	_, subspan := tracer.Start(ctx, "os.OpenFile")
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
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

	upload.Info.Offset += n
	return n, upload.writeInfo()
}

// GetInfo returns the FileInfo
func (upload *Upload) GetInfo(_ context.Context) (tusd.FileInfo, error) {
	return upload.Info, nil
}

// GetReader returns an io.Reader for the upload
func (upload *Upload) GetReader(_ context.Context) (io.Reader, error) {
	_, span := tracer.Start(upload.Ctx, "GetReader")
	defer span.End()
	return os.Open(upload.binPath)
}

// FinishUpload finishes an upload and moves the file to the internal destination
func (upload *Upload) FinishUpload(_ context.Context) error {
	ctx, span := tracer.Start(upload.Ctx, "FinishUpload")
	defer span.End()
	// set lockID to context
	if upload.Info.MetaData["lockid"] != "" {
		upload.Ctx = ctxpkg.ContextSetLockID(upload.Ctx, upload.Info.MetaData["lockid"])
	}

	log := appctx.GetLogger(upload.Ctx)

	// calculate the checksum of the written bytes
	// they will all be written to the metadata later, so we cannot omit any of them
	// TODO only calculate the checksum in sync that was requested to match, the rest could be async ... but the tests currently expect all to be present
	// TODO the hashes all implement BinaryMarshaler so we could try to persist the state for resumable upload. we would neet do keep track of the copied bytes ...
	sha1h := sha1.New()
	md5h := md5.New()
	adler32h := adler32.New()
	{
		_, subspan := tracer.Start(ctx, "os.Open")
		f, err := os.Open(upload.binPath)
		subspan.End()
		if err != nil {
			// we can continue if no oc checksum header is set
			log.Info().Err(err).Str("binPath", upload.binPath).Msg("error opening binPath")
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
	if upload.Info.MetaData["checksum"] != "" {
		var err error
		parts := strings.SplitN(upload.Info.MetaData["checksum"], " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		switch parts[0] {
		case "sha1":
			err = upload.checkHash(parts[1], sha1h)
		case "md5":
			err = upload.checkHash(parts[1], md5h)
		case "adler32":
			err = upload.checkHash(parts[1], adler32h)
		default:
			err = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if err != nil {
			Cleanup(upload, true, false)
			return err
		}
	}

	// update checksums
	attrs := node.Attributes{
		prefixes.ChecksumPrefix + "sha1":    sha1h.Sum(nil),
		prefixes.ChecksumPrefix + "md5":     md5h.Sum(nil),
		prefixes.ChecksumPrefix + "adler32": adler32h.Sum(nil),
	}

	n, err := CreateNodeForUpload(upload, attrs)
	if err != nil {
		Cleanup(upload, true, false)
		return err
	}

	upload.Node = n

	if upload.pub != nil {
		u, _ := ctxpkg.ContextGetUser(upload.Ctx)
		s, err := upload.URL(upload.Ctx)
		if err != nil {
			return err
		}

		if err := events.Publish(upload.pub, events.BytesReceived{
			UploadID:      upload.Info.ID,
			URL:           s,
			SpaceOwner:    n.SpaceOwnerOrManager(upload.Ctx),
			ExecutingUser: u,
			ResourceID:    &provider.ResourceId{SpaceId: n.SpaceID, OpaqueId: n.ID},
			Filename:      upload.Info.Storage["NodeName"],
			Filesize:      uint64(upload.Info.Size),
		}); err != nil {
			return err
		}
	}

	if !upload.async {
		// handle postprocessing synchronously
		err = upload.Finalize()
		Cleanup(upload, err != nil, false)
		if err != nil {
			return err
		}
	}

	return upload.tp.Propagate(upload.Ctx, n, upload.sizeDiff)
}

// Terminate terminates the upload
func (upload *Upload) Terminate(_ context.Context) error {
	upload.cleanup(true, true, true)
	return nil
}

// DeclareLength updates the upload length information
func (upload *Upload) DeclareLength(_ context.Context, length int64) error {
	upload.Info.Size = length
	upload.Info.SizeIsDeferred = false
	return upload.writeInfo()
}

// ConcatUploads concatenates multiple uploads
func (upload *Upload) ConcatUploads(_ context.Context, uploads []tusd.Upload) (err error) {
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, partialUpload := range uploads {
		fileUpload := partialUpload.(*Upload)

		src, err := os.Open(fileUpload.binPath)
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

// writeInfo updates the entire information. Everything will be overwritten.
func (upload *Upload) writeInfo() error {
	_, span := tracer.Start(upload.Ctx, "writeInfo")
	defer span.End()
	data, err := json.Marshal(upload.Info)
	if err != nil {
		return err
	}
	return os.WriteFile(upload.infoPath, data, defaultFilePerm)
}

// Finalize finalizes the upload (eg moves the file to the internal destination)
func (upload *Upload) Finalize() (err error) {
	ctx, span := tracer.Start(upload.Ctx, "Finalize")
	defer span.End()
	n := upload.Node
	if n == nil {
		var err error
		n, err = node.ReadNode(ctx, upload.lu, upload.Info.Storage["SpaceRoot"], upload.Info.Storage["NodeId"], false, nil, false)
		if err != nil {
			return err
		}
		upload.Node = n
	}

	// upload the data to the blobstore
	_, subspan := tracer.Start(ctx, "WriteBlob")
	err = upload.tp.WriteBlob(n, upload.binPath)
	subspan.End()
	if err != nil {
		return errors.Wrap(err, "failed to upload file to blostore")
	}

	return nil
}

func (upload *Upload) checkHash(expected string, h hash.Hash) error {
	if expected != hex.EncodeToString(h.Sum(nil)) {
		return errtypes.ChecksumMismatch(fmt.Sprintf("invalid checksum: expected %s got %x", upload.Info.MetaData["checksum"], h.Sum(nil)))
	}
	return nil
}

// cleanup cleans up after the upload is finished
func (upload *Upload) cleanup(cleanNode, cleanBin, cleanInfo bool) {
	if cleanNode && upload.Node != nil {
		switch p := upload.versionsPath; p {
		case "":
			// remove node
			if err := utils.RemoveItem(upload.Node.InternalPath()); err != nil {
				upload.log.Info().Str("path", upload.Node.InternalPath()).Err(err).Msg("removing node failed")
			}

			// no old version was present - remove child entry
			src := filepath.Join(upload.Node.ParentPath(), upload.Node.Name)
			if err := os.Remove(src); err != nil {
				upload.log.Info().Str("path", upload.Node.ParentPath()).Err(err).Msg("removing node from parent failed")
			}

			// remove node from upload as it no longer exists
			upload.Node = nil
		default:

			if err := upload.lu.CopyMetadata(upload.Ctx, p, upload.Node.InternalPath(), func(attributeName string) bool {
				return strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
					attributeName == prefixes.TypeAttr ||
					attributeName == prefixes.BlobIDAttr ||
					attributeName == prefixes.BlobsizeAttr
			}); err != nil {
				upload.log.Info().Str("versionpath", p).Str("nodepath", upload.Node.InternalPath()).Err(err).Msg("renaming version node failed")
			}

			if err := os.RemoveAll(p); err != nil {
				upload.log.Info().Str("versionpath", p).Str("nodepath", upload.Node.InternalPath()).Err(err).Msg("error removing version")
			}

		}
	}

	if cleanBin {
		if err := os.Remove(upload.binPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
			upload.log.Error().Str("path", upload.binPath).Err(err).Msg("removing upload failed")
		}
	}

	if cleanInfo {
		if err := os.Remove(upload.infoPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
			upload.log.Error().Str("path", upload.infoPath).Err(err).Msg("removing upload info failed")
		}
	}
}

// URL returns a url to download an upload
func (upload *Upload) URL(_ context.Context) (string, error) {
	type transferClaims struct {
		jwt.StandardClaims
		Target string `json:"target"`
	}

	u := joinurl(upload.tknopts.DownloadEndpoint, "tus/", upload.Info.ID)
	ttl := time.Duration(upload.tknopts.TransferExpires) * time.Second
	claims := transferClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Audience:  "reva",
			IssuedAt:  time.Now().Unix(),
		},
		Target: u,
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tkn, err := t.SignedString([]byte(upload.tknopts.TransferSharedSecret))
	if err != nil {
		return "", errors.Wrapf(err, "error signing token with claims %+v", claims)
	}

	return joinurl(upload.tknopts.DataGatewayEndpoint, tkn), nil
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
