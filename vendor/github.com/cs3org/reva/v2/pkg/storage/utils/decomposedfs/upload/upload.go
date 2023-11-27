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

// Package upload handles the processing of uploads.
// In general this is the lifecycle of an upload from the perspective of a storageprovider:
// 1. To start an upload a client makes a call to InitializeUpload which will return protocols and urls that he can use to append bytes to the upload.
// 2. When the client has sent all bytes the tusd handler will call a PreFinishResponseCallback which marks the end of the transfer and the start of postprocessing.
// 3. When async uploads are enabled the storageprovider emits an BytesReceived event, otherwise a FileUploaded event and the upload lifcycle ends.
// 4. During async postprocessing the uploaded bytes might be read at the upload URL to determine the outcome of the postprocessing steps
// 5. To handle async postprocessing the storageporvider has to listen to multiple events:
//   - PostprocessingFinished determines what should happen with the upload:
//   - abort - the upload is cancelled but the bytes are kept in the upload folder, eg. when antivirus scanning encounters an error
//     then what? can the admin retrigger the upload?
//   - continue - the upload is moved to its final destination (eventually being marked with pp results)
//   - delete - the file and the upload should be deleted
//   - RestartPostprocessing
//   - PostprocessingStepFinished is used to set scan data on an upload
//
// 6. The storageprovider emits an UploadReady event that can be used by eg. the search or thumbnails services to do update their metadata.
//
// There are two interesting scenarios:
// 1. Two concurrent requests try to create the same file
// 2. Two concurrent requests try to overwrite the same file
// The first step to upload a file is making an InitiateUpload call to the storageprovider via CS3. It will return an upload id that can be used to append bytes to the upload.
// With an upload id clients can append bytes to the upload.
// When all bytes have been received tusd will call PreFinishResponseCallback on the storageprovider.
// The storageprovider cannot use the tus upload metadata to persist a postprocessing status we have to store the processing status on a revision node.
// On disk the layout for a node consists of the actual node metadata and revision nodes.
// The revision nodes are used to capture the different revsions ...
// * so every uploed always creates a revision node first?
// * and in PreFinishResponseCallback we update or create? the actual node? or do we create the node in the InitiateUpload call?
// * We need to skip unfinished revisions when listing versions?
// The size diff is always calculated when updating the node
//
// ## Client considerations
// When do we propagate the etag? Currently, already when an upload is in postprocessing ... why? because we update the node when all bytes are transferred?
// Does the client expect an etag change when it uploads a file? it should not ... sync and uploads are independent last someone explained it to me
// postprocessing könnte den content ändern und damit das etag
//
// When the client finishes transferring all bytes it gets the 'future' etag of the resource which it currently stores as the etag for the file in its local db.
// When the next propfind happens before postprocessing finishes the client would see the old etag and download the old version. Then, when postprocessing causes
// the next etag change, the client will download the file it previously uploaded.
//
// For the new file scenario, the desktop client would delete the uploaded file locally, when it is not listed in the next propfind.
//
// The graph api exposes pending uploads explicitly using the pendingOperations property, which carries a pendingContentUpdate resource with a
// queuedDateTime property: Date and time the pending binary operation was queued in UTC time. Read-only.
//
// So, until clients learn to keep track of their uploads we need to return 425 when an upload is in progress ಠ_ಠ
package upload

import (
	"context"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/grpc/services/storageprovider"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/tree"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	tusd "github.com/tus/tusd/pkg/handler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/upload")
}

func validateRequest(ctx context.Context, size int64, uploadMetadata Metadata, n *node.Node) error {
	if err := n.CheckLock(ctx); err != nil {
		return err
	}

	if _, err := node.CheckQuota(ctx, n.SpaceRoot, true, uint64(n.Blobsize), uint64(size)); err != nil {
		return err
	}

	mtime, err := n.GetMTime(ctx)
	if err != nil {
		return err
	}
	currentEtag, err := node.CalculateEtag(n, mtime)
	if err != nil {
		return err
	}

	// When the if-match header was set we need to check if the
	// etag still matches before finishing the upload.
	if uploadMetadata.HeaderIfMatch != "" {
		if uploadMetadata.HeaderIfMatch != currentEtag {
			return errtypes.Aborted("etag mismatch")
		}
	}

	// When the if-none-match header was set we need to check if any of the
	// etags matches before finishing the upload.
	if uploadMetadata.HeaderIfNoneMatch != "" {
		if uploadMetadata.HeaderIfNoneMatch == "*" {
			return errtypes.Aborted("etag mismatch, resource exists")
		}
		for _, ifNoneMatchTag := range strings.Split(uploadMetadata.HeaderIfNoneMatch, ",") {
			if ifNoneMatchTag == currentEtag {
				return errtypes.Aborted("etag mismatch")
			}
		}
	}

	// When the if-unmodified-since header was set we need to check if the
	// etag still matches before finishing the upload.
	if uploadMetadata.HeaderIfUnmodifiedSince != "" {
		if err != nil {
			return errtypes.InternalError(fmt.Sprintf("failed to read mtime of node: %s", err))
		}
		ifUnmodifiedSince, err := time.Parse(time.RFC3339Nano, uploadMetadata.HeaderIfUnmodifiedSince)
		if err != nil {
			return errtypes.InternalError(fmt.Sprintf("failed to parse if-unmodified-since time: %s", err))
		}

		if mtime.After(ifUnmodifiedSince) {
			return errtypes.Aborted("if-unmodified-since mismatch")
		}
	}
	return nil
}

func openExistingNode(ctx context.Context, lu *lookup.Lookup, n *node.Node) (*lockedfile.File, error) {
	nodePath := n.InternalPath()

	// create and read lock existing node metadata
	log := appctx.GetLogger(ctx)
	log.Info().Str("nodepath", nodePath).Msg("grabbing lock for node")
	f, err := lockedfile.OpenFile(lu.MetadataBackend().LockfilePath(nodePath), os.O_RDONLY, 0600)
	log.Info().Str("nodepath", nodePath).Msg("got lock")

	return f, err
}

func initNewNode(ctx context.Context, lu *lookup.Lookup, uploadID, mtime string, n *node.Node) (*lockedfile.File, error) {
	nodePath := n.InternalPath()
	// create folder structure (if needed)
	if err := os.MkdirAll(filepath.Dir(nodePath), 0700); err != nil {
		return nil, err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentPath(), n.Name)
	relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))

	log := appctx.GetLogger(ctx).With().Str("childNameLink", childNameLink).Str("relativeNodePath", relativeNodePath).Logger()
	log.Info().Msg("initNewNode: creating symlink")

	// create and write lock new node metadata
	log.Info().Str("nodepath", nodePath).Msg("grabbing lock for node")
	f, err := lockedfile.OpenFile(lu.MetadataBackend().LockfilePath(nodePath), os.O_RDWR|os.O_CREATE, 0600)
	log.Info().Str("nodepath", nodePath).Msg("got lock")
	if err != nil {
		return nil, err
	}

	// FIXME if this is removed links to files will be dangling, causing subsequest stats to files to fail
	// we also need to touch the actual node file here it stores the mtime of the resource
	h, err := os.OpenFile(nodePath, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return f, err
	}
	h.Close()

	if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
		log.Info().Err(err).Msg("initNewNode: symlink failed")
		if errors.Is(err, iofs.ErrExist) {
			log.Info().Err(err).Msg("initNewNode: symlink already exists")
			return f, errtypes.AlreadyExists(n.Name)
		}
		return f, errors.Wrap(err, "Decomposedfs: could not symlink child entry")
	}
	log.Info().Msg("initNewNode: symlink created")

	attrs := node.Attributes{}
	attrs.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	attrs.SetString(prefixes.ParentidAttr, n.ParentID)
	attrs.SetString(prefixes.NameAttr, n.Name)
	attrs.SetString(prefixes.MTimeAttr, mtime) // TODO use mtime

	// here we set the status the first time.
	attrs.SetString(prefixes.StatusPrefix, node.ProcessingStatus+uploadID)

	// update node metadata with basic metadata
	err = n.SetXattrsWithContext(ctx, attrs, false)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: could not write metadata")
	}
	return f, nil
}

func CreateRevisionNode(ctx context.Context, lu *lookup.Lookup, revisionNode *node.Node) (*lockedfile.File, error) {
	revisionPath := revisionNode.InternalPath()
	log := appctx.GetLogger(ctx)

	// write lock existing node before reading any metadata
	log.Info().Str("revisionPath", revisionPath).Msg("grabbing lock for node")
	f, err := lockedfile.OpenFile(lu.MetadataBackend().LockfilePath(revisionPath), os.O_RDWR|os.O_CREATE, 0600)
	log.Info().Str("revisionPath", revisionPath).Msg("got lock")
	if err != nil {
		return nil, err
	}

	// FIXME if this is removed listing revisions breaks because it globs the dir but then filters all metadata files
	// we also need to touch the versions node here to list revisions
	h, err := os.OpenFile(revisionPath, os.O_CREATE /*|os.O_EXCL*/, 0600) // we have to allow overwriting revisions to be oc10 compatible
	if err != nil {
		return f, err
	}
	h.Close()
	return f, nil
}

func SetNodeToUpload(ctx context.Context, lu *lookup.Lookup, n *node.Node, uploadMetadata Metadata) (int64, error) {

	nodePath := n.InternalPath()
	// lock existing node metadata
	nh, err := lockedfile.OpenFile(lu.MetadataBackend().LockfilePath(nodePath), os.O_RDWR, 0600)
	if err != nil {
		return 0, err
	}
	defer nh.Close()
	// read nodes

	n, err = node.ReadNode(ctx, lu, n.SpaceID, n.ID, false, n.SpaceRoot, true)
	if err != nil {
		return 0, err
	}

	sizeDiff := uploadMetadata.BlobSize - n.Blobsize

	// TODO set blobid ind size ... do we need to do this? the node is passed by reference so subsequent calls might rely on this
	n.BlobID = uploadMetadata.BlobID
	n.Blobsize = uploadMetadata.BlobSize

	rm := RevisionMetadata{
		MTime:           uploadMetadata.MTime,
		BlobID:          uploadMetadata.BlobID,
		BlobSize:        uploadMetadata.BlobSize,
		ChecksumSHA1:    uploadMetadata.ChecksumSHA1,
		ChecksumMD5:     uploadMetadata.ChecksumMD5,
		ChecksumADLER32: uploadMetadata.ChecksumADLER32,
	}

	if rm.MTime == "" {
		rm.MTime = time.Now().UTC().Format(time.RFC3339Nano)
	}

	// update node
	err = WriteRevisionMetadataToNode(ctx, n, rm)
	if err != nil {
		return 0, errors.Wrap(err, "Decomposedfs: could not write metadata")
	}

	return sizeDiff, nil
}

type RevisionMetadata struct {
	MTime           string
	BlobID          string
	BlobSize        int64
	ChecksumSHA1    []byte
	ChecksumMD5     []byte
	ChecksumADLER32 []byte
}

func WriteRevisionMetadataToNode(ctx context.Context, n *node.Node, revisionMetadata RevisionMetadata) error {
	attrs := node.Attributes{}
	attrs.SetString(prefixes.BlobIDAttr, revisionMetadata.BlobID)
	attrs.SetInt64(prefixes.BlobsizeAttr, revisionMetadata.BlobSize)
	attrs.SetString(prefixes.MTimeAttr, revisionMetadata.MTime)
	attrs[prefixes.ChecksumPrefix+storageprovider.XSSHA1] = revisionMetadata.ChecksumSHA1
	attrs[prefixes.ChecksumPrefix+storageprovider.XSMD5] = revisionMetadata.ChecksumMD5
	attrs[prefixes.ChecksumPrefix+storageprovider.XSAdler32] = revisionMetadata.ChecksumADLER32

	return n.SetXattrsWithContext(ctx, attrs, false)
}

func ReadNode(ctx context.Context, lu *lookup.Lookup, uploadMetadata Metadata) (*node.Node, error) {
	var n *node.Node
	var err error
	if uploadMetadata.NodeID == "" {
		p, err := node.ReadNode(ctx, lu, uploadMetadata.SpaceRoot, uploadMetadata.NodeParentID, false, nil, true)
		if err != nil {
			return nil, err
		}
		n, err = p.Child(ctx, uploadMetadata.Filename)
		if err != nil {
			return nil, err
		}
	} else {
		n, err = node.ReadNode(ctx, lu, uploadMetadata.SpaceRoot, uploadMetadata.NodeID, false, nil, true)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

// Cleanup cleans the upload
func Cleanup(ctx context.Context, lu *lookup.Lookup, n *node.Node, uploadID, revision string, failure bool) {
	ctx, span := tracer.Start(ctx, "Cleanup")
	defer span.End()

	if n != nil { // node can be nil when there was an error before it was created (eg. checksum-mismatch)
		if failure {
			removeRevision(ctx, lu, n, revision)
		}
		// unset processing status
		if err := n.UnmarkProcessing(ctx, uploadID); err != nil {
			log := appctx.GetLogger(ctx)
			log.Info().Str("path", n.InternalPath()).Err(err).Msg("unmarking processing failed")
		}
	}
}

// removeRevision cleans up after the upload is finished
func removeRevision(ctx context.Context, lu *lookup.Lookup, n *node.Node, revision string) {
	log := appctx.GetLogger(ctx)
	nodePath := n.InternalPath()
	revisionPath := node.JoinRevisionKey(nodePath, revision)
	// remove revision
	if err := utils.RemoveItem(revisionPath); err != nil {
		log.Info().Str("path", revisionPath).Err(err).Msg("removing revision failed")
	}
	// purge revision metadata to clean up cache
	if err := lu.MetadataBackend().Purge(revisionPath); err != nil {
		log.Info().Str("path", revisionPath).Err(err).Msg("purging revision metadata failed")
	}

	if n.BlobID == "" { // FIXME ... this is brittle
		// no old version was present - remove child entry symlink from directory
		src := filepath.Join(n.ParentPath(), n.Name)
		if err := os.Remove(src); err != nil {
			log.Info().Str("path", n.ParentPath()).Err(err).Msg("removing node from parent failed")
		}

		// delete node
		if err := utils.RemoveItem(nodePath); err != nil {
			log.Info().Str("path", nodePath).Err(err).Msg("removing node failed")
		}

		// purge node metadata to clean up cache
		if err := lu.MetadataBackend().Purge(nodePath); err != nil {
			log.Info().Str("path", nodePath).Err(err).Msg("purging node metadata failed")
		}
	}
}

// Finalize finalizes the upload (eg moves the file to the internal destination)
func Finalize(ctx context.Context, blobstore tree.Blobstore, revision string, info tusd.FileInfo, n *node.Node, blobID string) error {
	_, span := tracer.Start(ctx, "Finalize")
	defer span.End()

	rn := n.RevisionNode(ctx, revision)
	rn.BlobID = blobID
	var err error
	if mover, ok := blobstore.(tree.BlobstoreMover); ok {
		err = mover.MoveBlob(rn, "", info.Storage["Bucket"], info.Storage["Key"])
		switch err {
		case nil:
			return nil
		case tree.ErrBlobstoreCannotMove:
			// fallback below
		default:
			return err
		}
	}

	// upload the data to the blobstore
	_, subspan := tracer.Start(ctx, "WriteBlob")
	err = blobstore.Upload(rn, info.Storage["Path"]) // FIXME where do we read from
	subspan.End()
	if err != nil {
		return errors.Wrap(err, "failed to upload file to blobstore")
	}

	// FIXME use a reader
	return nil
}
