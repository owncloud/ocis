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
		Output: "stdout",
		Mode:   "console",
		Level:  "debug",
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

	umg, ok := fs.(storage.HasUploadMetadata)
	if ok {
		go func() {
			for {
				ev := <-handler.CompleteUploads
				info := ev.Upload
				um, err := umg.GetUploadMetadata(context.TODO(), info.ID) // TODO we need to pass in a context, maybe with tusd 2.0. IIRC the relvease notes mention using context in more places.
				if err != nil {
					appctx.GetLogger(context.Background()).Error().Err(err).Msg("failed to get upload metadata on publish FileUploaded event")
				}
				spaceOwner := um.GetSpaceOwner()
				executant := um.GetExecutantID()
				ref := um.GetReference()
				datatx.InvalidateCache(&executant, &ref, m.statCache)
				if m.publisher != nil {
					if err := datatx.EmitFileUploadedEvent(&spaceOwner, &executant, &ref, m.publisher); err != nil {
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
			setHeaders(dataStore, umg, w, r)
			handler.PostFile(w, r)
		case "HEAD":
			handler.HeadFile(w, r)
		case "PATCH":
			metrics.UploadsActive.Add(1)
			defer func() {
				metrics.UploadsActive.Sub(1)
			}()
			// set etag, mtime and file id
			setHeaders(dataStore, umg, w, r)
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

func setHeaders(datastore tusd.DataStore, umg storage.HasUploadMetadata, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := path.Base(r.URL.Path)
	u, err := datastore.GetUpload(ctx, id)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("could not get upload from storage")
		return
	}
	info, err := u.GetInfo(ctx)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("could not get upload info for upload")
		return
	}
	expires := ""
	resourceid := provider.ResourceId{}
	if umg != nil {
		um, err := umg.GetUploadMetadata(ctx, info.ID)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Msg("could not get upload info for upload")
			return
		}
		expires = um.GetExpires()
		resourceid = um.GetResourceID()
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
	if resourceid.StorageId == "" {
		resourceid.StorageId = info.MetaData["providerID"]
	}
	if resourceid.SpaceId == "" {
		resourceid.SpaceId = info.MetaData["SpaceRoot"]
	}
	if resourceid.OpaqueId == "" {
		resourceid.OpaqueId = info.MetaData["NodeId"]
	}

	w.Header().Set(net.HeaderOCFileID, storagespace.FormatResourceID(resourceid))
}
