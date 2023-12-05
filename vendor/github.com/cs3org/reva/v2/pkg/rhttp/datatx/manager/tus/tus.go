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

// Package tus implements a data tx manager that handles uploads using the TUS protocol.
// reva storage drivers should implement the hasTusDatastore interface by using composition
// of an upstream tusd.DataStore. If necessary they can also implement a tusd.DataStore directly.
package tus

import (
	"context"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/cache"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registry.Register("tus", New)
}

type manager struct {
	conf      *cache.Config
	publisher events.Publisher
	statCache cache.StatCache
}

func parseConfig(m map[string]interface{}) (*cache.Config, error) {
	c := &cache.Config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New returns a datatx manager implementation that relies on HTTP PUT/GET.
func New(m map[string]interface{}, publisher events.Publisher) (datatx.DataTX, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}
	return &manager{
		conf:      c,
		publisher: publisher,
		statCache: cache.GetStatCache(c.Store, c.Nodes, c.Database, c.Table, time.Duration(c.TTL)*time.Second, c.Size),
	}, nil
}

func (m *manager) Handler(fs storage.FS) (http.Handler, error) {
	zlog, err := logger.FromConfig(&logger.LogConf{
		Output: "stderr",
		Mode:   "json",
		Level:  "error", // FIXME introduce shared config for logging
	})
	if err != nil {
		return nil, errtypes.NotSupported("could not initialize log")
	}

	composer := tusd.NewStoreComposer()

	config := tusd.Config{
		StoreComposer: composer,
		PreUploadCreateCallback: func(hook tusd.HookEvent) error {
			return errors.New("uploads must be created with a cs3 InitiateUpload call")
		},
		NotifyCompleteUploads: true,
		Logger:                log.New(zlog, "", 0),
	}

	var dataStore tusd.DataStore

	cb, ok := fs.(hasTusDatastore)
	if ok {
		dataStore = cb.GetDataStore()
		composable, ok := dataStore.(composable)
		if !ok {
			return nil, errtypes.NotSupported("tus datastore is not composable")
		}
		composable.UseIn(composer)
		config.PreFinishResponseCallback = cb.PreFinishResponseCallback
	} else {
		composable, ok := fs.(composable)
		if !ok {
			return nil, errtypes.NotSupported("storage driver does not support the tus protocol")
		}

		// let the composable storage tell tus which extensions it supports
		composable.UseIn(composer)
		dataStore, ok = fs.(tusd.DataStore)
		if !ok {
			return nil, errtypes.NotSupported("storage driver does not support the tus datastore")
		}
	}

	handler, err := tusd.NewUnroutedHandler(config)
	if err != nil {
		return nil, err
	}

	usl, ok := fs.(storage.UploadSessionLister)
	if ok {
		// We can currently only send updates if the fs is decomposedfs as we read very specific keys from the storage map of the tus info
		go func() {
			for {
				ev := <-handler.CompleteUploads
				// We should be able to get the upload progress with fs.GetUploadProgress, but currently tus will erase the info files
				// so we create a Progress instance here that is used to read the correct properties
				sessions, err := usl.ListUploadSessions(context.Background(), storage.UploadSessionFilter{ID: &ev.Upload.ID})
				if err != nil {
					appctx.GetLogger(context.Background()).Error().Err(err).Str("id", ev.Upload.ID).Msg("failed to list upload session for upload")
					continue
				}
				if len(sessions) != 1 {
					appctx.GetLogger(context.Background()).Error().Err(err).Str("id", ev.Upload.ID).Msg("no upload session found")
					continue
				}
				us := sessions[0]

				executant := us.Executant()
				ref := us.Reference()
				datatx.InvalidateCache(&executant, &ref, m.statCache)
				if m.publisher != nil {
					if err := datatx.EmitFileUploadedEvent(us.SpaceOwner(), &executant, &ref, m.publisher); err != nil {
						appctx.GetLogger(context.Background()).Error().Err(err).Msg("failed to publish FileUploaded event")
					}
				}
			}
		}()
	}

	h := handler.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		// https://github.com/tus/tus-resumable-upload-protocol/blob/master/protocol.md#x-http-method-override
		if r.Header.Get("X-HTTP-Method-Override") != "" {
			method = r.Header.Get("X-HTTP-Method-Override")
		}

		switch method {
		case "POST":
			metrics.UploadsActive.Add(1)
			defer func() {
				metrics.UploadsActive.Sub(1)
			}()
			// set etag, mtime and file id
			setHeaders(dataStore, usl, w, r)
			handler.PostFile(w, r)
		case "HEAD":
			handler.HeadFile(w, r)
		case "PATCH":
			metrics.UploadsActive.Add(1)
			defer func() {
				metrics.UploadsActive.Sub(1)
			}()
			// set etag, mtime and file id
			setHeaders(dataStore, usl, w, r)
			handler.PatchFile(w, r)
		case "DELETE":
			handler.DelFile(w, r)
		case "GET":
			metrics.DownloadsActive.Add(1)
			defer func() {
				metrics.DownloadsActive.Sub(1)
			}()
			// NOTE: this is breaking change - allthought it does not seem to be used
			// We can make a switch here depending on some header value if that is needed
			// download.GetOrHeadFile(w, r, fs, "")
			handler.GetFile(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))

	return h, nil
}

// Composable is the interface that a struct needs to implement
// to be composable, so that it can support the TUS methods
type composable interface {
	UseIn(composer *tusd.StoreComposer)
}

type hasTusDatastore interface {
	PreFinishResponseCallback(hook tusd.HookEvent) error
	GetDataStore() tusd.DataStore
}

func setHeaders(datastore tusd.DataStore, usl storage.UploadSessionLister, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := path.Base(r.URL.Path)
	upload, err := datastore.GetUpload(ctx, id)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("could not get upload from storage")
		return
	}
	info, err := upload.GetInfo(ctx)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("could not get upload info for upload")
		return
	}
	var expires string

	var resourceid provider.ResourceId
	var uploadSession storage.UploadSession
	if usl != nil {
		sessions, err := usl.ListUploadSessions(ctx, storage.UploadSessionFilter{ID: &id})
		if err != nil {
			appctx.GetLogger(context.Background()).Error().Err(err).Str("id", id).Msg("failed to list upload session for upload")
			return
		}
		if len(sessions) != 1 {
			appctx.GetLogger(context.Background()).Error().Err(err).Str("id", id).Msg("no upload session found")
			return
		}
		uploadSession = sessions[0]

		t := time.Time{}
		if uploadSession.Expires() != t {
			expires = uploadSession.Expires().Format(net.RFC1123)
		}

		reference := uploadSession.Reference()
		resourceid = *reference.GetResourceId()
	}

	// FIXME expires should be part of the tus handler
	// fallback for outdated storageproviders that implement a tus datastore
	if expires == "" {
		expires = info.MetaData["expires"]
	}

	if expires != "" {
		w.Header().Set(net.HeaderTusUploadExpires, expires)
	}

	// fallback for outdated storageproviders that implement a tus datastore
	if resourceid.GetStorageId() == "" {
		resourceid.StorageId = info.MetaData["providerID"]
	}
	if resourceid.GetSpaceId() == "" {
		resourceid.SpaceId = info.MetaData["SpaceRoot"]
	}
	if resourceid.GetOpaqueId() == "" {
		resourceid.OpaqueId = info.MetaData["NodeId"]
	}

	w.Header().Set(net.HeaderOCFileID, storagespace.FormatResourceID(resourceid))
}
