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
	"encoding/json"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	tusd "github.com/tus/tusd/pkg/handler"
)

var defaultFilePerm = os.FileMode(0664)

// PermissionsChecker defines an interface for checking permissions on a Node
type PermissionsChecker interface {
	AssemblePermissions(ctx context.Context, n *node.Node) (ap provider.ResourcePermissions, err error)
}

// New returns a new processing instance
func New(ctx context.Context, info tusd.FileInfo, lu *lookup.Lookup, tp Tree, p PermissionsChecker, fsRoot string, pub events.Publisher, async bool, tknopts options.TokenOptions) (upload *Upload, err error) {

	log := appctx.GetLogger(ctx)
	log.Debug().Interface("info", info).Msg("Decomposedfs: NewUpload")

	if info.MetaData["filename"] == "" {
		return nil, errors.New("Decomposedfs: missing filename in metadata")
	}
	if info.MetaData["dir"] == "" {
		return nil, errors.New("Decomposedfs: missing dir in metadata")
	}

	n, err := lu.NodeFromSpaceID(ctx, &provider.ResourceId{
		SpaceId:  info.Storage["SpaceRoot"],
		OpaqueId: info.Storage["SpaceRoot"],
	})
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error getting space root node")
	}

	n, err = lookupNode(ctx, n, filepath.Join(info.MetaData["dir"], info.MetaData["filename"]), lu)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	log.Debug().Interface("info", info).Interface("node", n).Msg("Decomposedfs: resolved filename")

	// the parent owner will become the new owner
	parent, perr := n.Parent()
	if perr != nil {
		return nil, errors.Wrap(perr, "Decomposedfs: error getting parent "+n.ParentID)
	}

	// check permissions
	var (
		checkNode *node.Node
		path      string
	)
	if n.Exists {
		// check permissions of file to be overwritten
		checkNode = n
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}})
	} else {
		// check permissions of parent
		checkNode = parent
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}, Path: n.Name})
	}
	rp, err := p.AssemblePermissions(ctx, checkNode)
	switch {
	case err != nil:
		return nil, errtypes.InternalError(err.Error())
	case !rp.InitiateFileUpload:
		return nil, errtypes.PermissionDenied(path)
	}

	// are we trying to overwriting a folder with a file?
	if n.Exists && n.IsDir() {
		return nil, errtypes.PreconditionFailed("resource is not a file")
	}

	// check lock
	if info.MetaData["lockid"] != "" {
		ctx = ctxpkg.ContextSetLockID(ctx, info.MetaData["lockid"])
	}
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	info.ID = uuid.New().String()

	binPath := filepath.Join(fsRoot, "uploads", info.ID)
	usr := ctxpkg.ContextMustGetUser(ctx)

	var (
		spaceRoot string
		ok        bool
	)
	if info.Storage != nil {
		if spaceRoot, ok = info.Storage["SpaceRoot"]; !ok {
			spaceRoot = n.SpaceRoot.ID
		}
	} else {
		spaceRoot = n.SpaceRoot.ID
	}

	info.Storage = map[string]string{
		"Type":    "OCISStore",
		"BinPath": binPath,

		"NodeId":              n.ID,
		"NodeParentId":        n.ParentID,
		"NodeName":            n.Name,
		"SpaceRoot":           spaceRoot,
		"SpaceOwnerOrManager": info.Storage["SpaceOwnerOrManager"],

		"Idp":      usr.Id.Idp,
		"UserId":   usr.Id.OpaqueId,
		"UserType": utils.UserTypeToString(usr.Id.Type),
		"UserName": usr.Username,

		"LogLevel": log.GetLevel().String(),
	}
	// Create binary file in the upload folder with no content
	log.Debug().Interface("info", info).Msg("Decomposedfs: built storage info")
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	u := buildUpload(ctx, info, binPath, filepath.Join(fsRoot, "uploads", info.ID+".info"), lu, tp, pub, async, tknopts)

	// writeInfo creates the file by itself if necessary
	err = u.writeInfo()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Get returns the Upload for the given upload id
func Get(ctx context.Context, id string, lu *lookup.Lookup, tp Tree, fsRoot string, pub events.Publisher, async bool, tknopts options.TokenOptions) (*Upload, error) {
	infoPath := filepath.Join(fsRoot, "uploads", id+".info")

	info := tusd.FileInfo{}
	data, err := os.ReadFile(infoPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = tusd.ErrNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	stat, err := os.Stat(info.Storage["BinPath"])
	if err != nil {
		return nil, err
	}

	info.Offset = stat.Size()

	u := &userpb.User{
		Id: &userpb.UserId{
			Idp:      info.Storage["Idp"],
			OpaqueId: info.Storage["UserId"],
			Type:     utils.UserTypeMap(info.Storage["UserType"]),
		},
		Username: info.Storage["UserName"],
	}

	ctx = ctxpkg.ContextSetUser(ctx, u)
	// TODO configure the logger the same way ... store and add traceid in file info

	var opts []logger.Option
	opts = append(opts, logger.WithLevel(info.Storage["LogLevel"]))
	opts = append(opts, logger.WithWriter(os.Stderr, logger.ConsoleMode))
	l := logger.New(opts...)

	sub := l.With().Int("pid", os.Getpid()).Logger()

	ctx = appctx.WithLogger(ctx, &sub)

	up := buildUpload(ctx, info, info.Storage["BinPath"], infoPath, lu, tp, pub, async, tknopts)
	up.versionsPath = info.MetaData["versionsPath"]
	up.sizeDiff, _ = strconv.ParseInt(info.MetaData["sizeDiff"], 10, 64)
	return up, nil
}

// CreateNodeForUpload will create the target node for the Upload
func CreateNodeForUpload(upload *Upload, initAttrs node.Attributes) (*node.Node, error) {
	fi, err := os.Stat(upload.binPath)
	if err != nil {
		return nil, err
	}

	fsize := fi.Size()
	spaceID := upload.Info.Storage["SpaceRoot"]
	n := node.New(
		spaceID,
		upload.Info.Storage["NodeId"],
		upload.Info.Storage["NodeParentId"],
		upload.Info.Storage["NodeName"],
		fsize,
		upload.Info.ID,
		provider.ResourceType_RESOURCE_TYPE_FILE,
		nil,
		upload.lu,
	)
	n.SpaceRoot, err = node.ReadNode(upload.Ctx, upload.lu, spaceID, spaceID, false, nil, false)
	if err != nil {
		return nil, err
	}

	// check lock
	if err := n.CheckLock(upload.Ctx); err != nil {
		return nil, err
	}

	var f *lockedfile.File
	switch n.ID {
	case "":
		f, err = initNewNode(upload, n, uint64(fsize))
	default:
		f, err = updateExistingNode(upload, n, spaceID, uint64(fsize))
	}
	defer func() {
		if err := f.Close(); err != nil {
			appctx.GetLogger(upload.Ctx).Error().Err(err).Str("nodeid", n.ID).Str("parentid", n.ParentID).Msg("could not close lock")
		}
	}()
	if err != nil {
		return nil, err
	}

	// overwrite technical information
	initAttrs.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	initAttrs.SetString(prefixes.ParentidAttr, n.ParentID)
	initAttrs.SetString(prefixes.NameAttr, n.Name)
	initAttrs.SetString(prefixes.BlobIDAttr, n.BlobID)
	initAttrs.SetInt64(prefixes.BlobsizeAttr, n.Blobsize)
	initAttrs.SetString(prefixes.StatusPrefix, node.ProcessingStatus+upload.Info.ID)

	// update node metadata with new blobid etc
	err = n.SetXattrs(initAttrs, false)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: could not write metadata")
	}

	// overwrite mtime if requested
	if upload.Info.MetaData["mtime"] != "" {
		if err := n.SetMtimeString(upload.Info.MetaData["mtime"]); err != nil {
			return nil, err
		}
	}

	// add etag to metadata
	tmtime, err := n.GetMTime()
	if err != nil {
		return nil, err
	}
	upload.Info.MetaData["etag"], _ = node.CalculateEtag(n.ID, tmtime)

	// update nodeid for later
	upload.Info.Storage["NodeId"] = n.ID
	if err := upload.writeInfo(); err != nil {
		return nil, err
	}

	return n, nil
}

func initNewNode(upload *Upload, n *node.Node, fsize uint64) (*lockedfile.File, error) {
	n.ID = uuid.New().String()

	// create folder structure (if needed)
	if err := os.MkdirAll(filepath.Dir(n.InternalPath()), 0700); err != nil {
		return nil, err
	}

	// create and write lock new node metadata
	f, err := lockedfile.OpenFile(upload.lu.MetadataBackend().MetadataPath(n.InternalPath()), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	switch upload.lu.MetadataBackend().(type) {
	case metadata.MessagePackBackend:
		// for the ini and metadata backend we also need to touch the actual node file here.
		// it stores the mtime of the resource, which must not change when we update the ini file
		h, err := os.OpenFile(n.InternalPath(), os.O_CREATE, 0600)
		if err != nil {
			return f, err
		}
		h.Close()
	case metadata.XattrsBackend:
		// nothing to do
	}

	if _, err := node.CheckQuota(n.SpaceRoot, false, 0, fsize); err != nil {
		return f, err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentPath(), n.Name)
	link, err := os.Readlink(childNameLink)
	if err == nil && link != "../"+n.ID {
		if err := os.Remove(childNameLink); err != nil {
			return f, errors.Wrap(err, "Decomposedfs: could not remove symlink child entry")
		}
	}
	if errors.Is(err, iofs.ErrNotExist) || link != "../"+n.ID {
		relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
		if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
			return f, errors.Wrap(err, "Decomposedfs: could not symlink child entry")
		}
	}

	// on a new file the sizeDiff is the fileSize
	upload.sizeDiff = int64(fsize)
	upload.Info.MetaData["sizeDiff"] = strconv.Itoa(int(upload.sizeDiff))
	return f, nil
}

func updateExistingNode(upload *Upload, n *node.Node, spaceID string, fsize uint64) (*lockedfile.File, error) {
	old, _ := node.ReadNode(upload.Ctx, upload.lu, spaceID, n.ID, false, nil, false)
	if _, err := node.CheckQuota(n.SpaceRoot, true, uint64(old.Blobsize), fsize); err != nil {
		return nil, err
	}

	tmtime, err := old.GetTMTime()
	if err != nil {
		return nil, err
	}

	// When the if-match header was set we need to check if the
	// etag still matches before finishing the upload.
	if ifMatch, ok := upload.Info.MetaData["if-match"]; ok {
		targetEtag, err := node.CalculateEtag(n.ID, tmtime)
		switch {
		case err != nil:
			return nil, errtypes.InternalError(err.Error())
		case ifMatch != targetEtag:
			return nil, errtypes.Aborted("etag mismatch")
		}
	}

	upload.versionsPath = upload.lu.InternalPath(spaceID, n.ID+node.RevisionIDDelimiter+tmtime.UTC().Format(time.RFC3339Nano))
	upload.sizeDiff = int64(fsize) - old.Blobsize
	upload.Info.MetaData["versionsPath"] = upload.versionsPath
	upload.Info.MetaData["sizeDiff"] = strconv.Itoa(int(upload.sizeDiff))

	targetPath := n.InternalPath()

	// write lock existing node before reading treesize or tree time
	f, err := lockedfile.OpenFile(upload.lu.MetadataBackend().MetadataPath(targetPath), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// create version node
	if _, err := os.Create(upload.versionsPath); err != nil {
		return f, err
	}

	// copy blob metadata to version node
	if err := upload.lu.CopyMetadataWithSourceLock(targetPath, upload.versionsPath, func(attributeName string) bool {
		return strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
			attributeName == prefixes.TypeAttr ||
			attributeName == prefixes.BlobIDAttr ||
			attributeName == prefixes.BlobsizeAttr
	}, f); err != nil {
		return f, err
	}

	// keep mtime from previous version
	if err := os.Chtimes(upload.versionsPath, tmtime, tmtime); err != nil {
		return f, errtypes.InternalError(fmt.Sprintf("failed to change mtime of version node: %s", err))
	}

	// update mtime of current version
	mtime := time.Now()
	if err := os.Chtimes(n.InternalPath(), mtime, mtime); err != nil {
		return nil, err
	}

	return f, nil
}

// lookupNode looks up nodes by path.
// This method can also handle lookups for paths which contain chunking information.
func lookupNode(ctx context.Context, spaceRoot *node.Node, path string, lu *lookup.Lookup) (*node.Node, error) {
	p := path
	isChunked := chunking.IsChunked(path)
	if isChunked {
		chunkInfo, err := chunking.GetChunkBLOBInfo(path)
		if err != nil {
			return nil, err
		}
		p = chunkInfo.Path
	}

	n, err := lu.WalkPath(ctx, spaceRoot, p, true, func(ctx context.Context, n *node.Node) error { return nil })
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	if isChunked {
		n.Name = filepath.Base(path)
	}
	return n, nil
}
