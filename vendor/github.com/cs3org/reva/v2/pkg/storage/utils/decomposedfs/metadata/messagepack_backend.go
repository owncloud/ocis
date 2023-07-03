// Copyright 2018-2023 CERN
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

package metadata

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cs3org/reva/v2/pkg/storage/cache"
	"github.com/pkg/xattr"
	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/shamaton/msgpack/v2"
	"go.opentelemetry.io/otel/codes"
)

// MessagePackBackend persists the attributes in messagepack format inside the file
type MessagePackBackend struct {
	rootPath  string
	metaCache cache.FileMetadataCache
}

type readWriteCloseSeekTruncater interface {
	io.ReadWriteCloser
	io.Seeker
	Truncate(int64) error
}

// NewMessagePackBackend returns a new MessagePackBackend instance
func NewMessagePackBackend(rootPath string, o cache.Config) MessagePackBackend {
	return MessagePackBackend{
		rootPath:  filepath.Clean(rootPath),
		metaCache: cache.GetFileMetadataCache(o.Store, o.Nodes, o.Database, "filemetadata:", time.Duration(o.TTL)*time.Second, o.Size),
	}
}

// Name returns the name of the backend
func (MessagePackBackend) Name() string { return "messagepack" }

// All reads all extended attributes for a node
func (b MessagePackBackend) All(ctx context.Context, path string) (map[string][]byte, error) {
	return b.loadAttributes(ctx, path, nil)
}

// Get an extended attribute value for the given key
func (b MessagePackBackend) Get(ctx context.Context, path, key string) ([]byte, error) {
	attribs, err := b.loadAttributes(ctx, path, nil)
	if err != nil {
		return []byte{}, err
	}
	val, ok := attribs[key]
	if !ok {
		return []byte{}, &xattr.Error{Op: "mpk.get", Path: path, Name: key, Err: xattr.ENOATTR}
	}
	return val, nil
}

// GetInt64 reads a string as int64 from the xattrs
func (b MessagePackBackend) GetInt64(ctx context.Context, path, key string) (int64, error) {
	attribs, err := b.loadAttributes(ctx, path, nil)
	if err != nil {
		return 0, err
	}
	val, ok := attribs[key]
	if !ok {
		return 0, &xattr.Error{Op: "mpk.get", Path: path, Name: key, Err: xattr.ENOATTR}
	}
	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// List retrieves a list of names of extended attributes associated with the
// given path in the file system.
func (b MessagePackBackend) List(ctx context.Context, path string) ([]string, error) {
	attribs, err := b.loadAttributes(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	keys := []string{}
	for k := range attribs {
		keys = append(keys, k)
	}
	return keys, nil
}

// Set sets one attribute for the given path
func (b MessagePackBackend) Set(ctx context.Context, path, key string, val []byte) error {
	return b.SetMultiple(ctx, path, map[string][]byte{key: val}, true)
}

// SetMultiple sets a set of attribute for the given path
func (b MessagePackBackend) SetMultiple(ctx context.Context, path string, attribs map[string][]byte, acquireLock bool) error {
	return b.saveAttributes(ctx, path, attribs, nil, acquireLock)
}

// Remove an extended attribute key
func (b MessagePackBackend) Remove(ctx context.Context, path, key string) error {
	return b.saveAttributes(ctx, path, nil, []string{key}, true)
}

// AllWithLockedSource reads all extended attributes from the given reader (if possible).
// The path argument is used for storing the data in the cache
func (b MessagePackBackend) AllWithLockedSource(ctx context.Context, path string, source io.Reader) (map[string][]byte, error) {
	return b.loadAttributes(ctx, path, source)
}

func (b MessagePackBackend) saveAttributes(ctx context.Context, path string, setAttribs map[string][]byte, deleteAttribs []string, acquireLock bool) error {
	var (
		err error
		f   readWriteCloseSeekTruncater
	)
	ctx, span := tracer.Start(ctx, "saveAttributes")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}()

	metaPath := b.MetadataPath(path)
	if acquireLock {
		_, subspan := tracer.Start(ctx, "lockedfile.OpenFile")
		f, err = lockedfile.OpenFile(metaPath, os.O_RDWR|os.O_CREATE, 0600)
		subspan.End()
	} else {
		_, subspan := tracer.Start(ctx, "os.OpenFile")
		f, err = os.OpenFile(metaPath, os.O_RDWR|os.O_CREATE, 0600)
		subspan.End()
	}
	if err != nil {
		return err
	}
	defer f.Close()

	// Invalidate cache early
	_, subspan := tracer.Start(ctx, "metaCache.RemoveMetadata")
	_ = b.metaCache.RemoveMetadata(b.cacheKey(path))
	subspan.End()

	// Read current state
	_, subspan = tracer.Start(ctx, "io.ReadAll")
	var msgBytes []byte
	msgBytes, err = io.ReadAll(f)
	subspan.End()
	if err != nil {
		return err
	}
	attribs := map[string][]byte{}
	if len(msgBytes) > 0 {
		err = msgpack.Unmarshal(msgBytes, &attribs)
		if err != nil {
			return err
		}
	}

	// set new metadata
	for key, val := range setAttribs {
		attribs[key] = val
	}
	for _, key := range deleteAttribs {
		delete(attribs, key)
	}

	// Truncate file
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, subspan = tracer.Start(ctx, "f.Truncate")
	err = f.Truncate(0)
	subspan.End()
	if err != nil {
		return err
	}

	// Write new metadata to file
	var d []byte
	d, err = msgpack.Marshal(attribs)
	if err != nil {
		return err
	}
	_, subspan = tracer.Start(ctx, "f.Write")
	_, err = f.Write(d)
	subspan.End()
	if err != nil {
		return err
	}

	_, subspan = tracer.Start(ctx, "metaCache.PushToCache")
	err = b.metaCache.PushToCache(b.cacheKey(path), attribs)
	subspan.End()
	return err
}

func (b MessagePackBackend) loadAttributes(ctx context.Context, path string, source io.Reader) (map[string][]byte, error) {
	ctx, span := tracer.Start(ctx, "loadAttributes")
	defer span.End()
	attribs := map[string][]byte{}
	err := b.metaCache.PullFromCache(b.cacheKey(path), &attribs)
	if err == nil {
		return attribs, err
	}

	metaPath := b.MetadataPath(path)
	if source == nil {
		_, subspan := tracer.Start(ctx, "lockedfile.Open")
		source, err = lockedfile.Open(metaPath)
		subspan.End()
		// // No cached entry found. Read from storage and store in cache
		if err != nil {
			if os.IsNotExist(err) {
				// some of the caller rely on ENOTEXISTS to be returned when the
				// actual file (not the metafile) does not exist in order to
				// determine whether a node exists or not -> stat the actual node
				_, subspan := tracer.Start(ctx, "os.Stat")
				_, err := os.Stat(path)
				subspan.End()
				if err != nil {
					return nil, err
				}
				return attribs, nil // no attributes set yet
			}
		}
		defer source.(*lockedfile.File).Close()
	}

	_, subspan := tracer.Start(ctx, "io.ReadAll")
	msgBytes, err := io.ReadAll(source)
	subspan.End()
	if err != nil {
		return nil, err
	}
	if len(msgBytes) > 0 {
		err = msgpack.Unmarshal(msgBytes, &attribs)
		if err != nil {
			return nil, err
		}
	}

	_, subspan = tracer.Start(ctx, "metaCache.PushToCache")
	err = b.metaCache.PushToCache(b.cacheKey(path), attribs)
	subspan.End()
	if err != nil {
		return nil, err
	}

	return attribs, nil
}

// IsMetaFile returns whether the given path represents a meta file
func (MessagePackBackend) IsMetaFile(path string) bool { return strings.HasSuffix(path, ".mpk") }

// Purge purges the data of a given path
func (b MessagePackBackend) Purge(path string) error {
	if err := b.metaCache.RemoveMetadata(b.cacheKey(path)); err != nil {
		return err
	}
	return os.Remove(b.MetadataPath(path))
}

// Rename moves the data for a given path to a new path
func (b MessagePackBackend) Rename(oldPath, newPath string) error {
	data := map[string][]byte{}
	err := b.metaCache.PullFromCache(b.cacheKey(oldPath), &data)
	if err == nil {
		err = b.metaCache.PushToCache(b.cacheKey(newPath), data)
		if err != nil {
			return err
		}
	}
	err = b.metaCache.RemoveMetadata(b.cacheKey(oldPath))
	if err != nil {
		return err
	}

	return os.Rename(b.MetadataPath(oldPath), b.MetadataPath(newPath))
}

// MetadataPath returns the path of the file holding the metadata for the given path
func (MessagePackBackend) MetadataPath(path string) string { return path + ".mpk" }

func (b MessagePackBackend) cacheKey(path string) string {
	// rootPath is guaranteed to have no trailing slash
	// the cache key shouldn't begin with a slash as some stores drop it which can cause
	// confusion
	return strings.TrimPrefix(path, b.rootPath+"/")
}
