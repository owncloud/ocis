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

package tus

import (
	"context"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registry.Register("tus", New)
}

type TusConfig struct {
	CorsEnabled          bool   `mapstructure:"cors_enabled"`
	CorsAllowOrigin      string `mapstructure:"cors_allow_origin"`
	CorsAllowCredentials bool   `mapstructure:"cors_allow_credentials"`
	CorsAllowMethods     string `mapstructure:"cors_allow_methods"`
	CorsAllowHeaders     string `mapstructure:"cors_allow_headers"`
	CorsMaxAge           string `mapstructure:"cors_max_age"`
	CorsExposeHeaders    string `mapstructure:"cors_expose_headers"`
}

type manager struct {
	conf      *TusConfig
	publisher events.Publisher
}

func parseConfig(m map[string]interface{}) (*TusConfig, error) {
	c := &TusConfig{}
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
	}, nil
}

func (m *manager) Handler(fs storage.FS) (http.Handler, error) {
	composable, ok := fs.(storage.ComposableFS)
	if !ok {
		return nil, errtypes.NotSupported("file system does not support the tus protocol")
	}

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()

	// let the composable storage tell tus which extensions it supports
	composable.UseIn(composer)

	config := tusd.Config{
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		Logger:                log.New(appctx.GetLogger(context.Background()), "", 0),
	}

	if m.conf.CorsEnabled {
		allowOrigin, err := regexp.Compile(m.conf.CorsAllowOrigin)
		if m.conf.CorsAllowOrigin != "" && err != nil {
			return nil, err
		}

		config.Cors = &tusd.CorsConfig{
			Disable:          false,
			AllowOrigin:      allowOrigin,
			AllowCredentials: m.conf.CorsAllowCredentials,
			AllowMethods:     m.conf.CorsAllowMethods,
			AllowHeaders:     m.conf.CorsAllowHeaders,
			MaxAge:           m.conf.CorsMaxAge,
			ExposeHeaders:    m.conf.CorsExposeHeaders,
		}
	}

	handler, err := tusd.NewUnroutedHandler(config)
	if err != nil {
		return nil, err
	}

	if usl, ok := fs.(storage.UploadSessionLister); ok {
		// We can currently only send updates if the fs is decomposedfs as we read very specific keys from the storage map of the tus info
		go func() {
			for {
				ev := <-handler.CompleteUploads
				// We should be able to get the upload progress with fs.GetUploadProgress, but currently tus will erase the info files
				// so we create a Progress instance here that is used to read the correct properties
				ups, err := usl.ListUploadSessions(context.Background(), storage.UploadSessionFilter{ID: &ev.Upload.ID})
				if err != nil {
					appctx.GetLogger(context.Background()).Error().Err(err).Str("session", ev.Upload.ID).Msg("failed to list upload session")
				} else {
					up := ups[0]
					executant := up.Executant()
					ref := up.Reference()
					if m.publisher != nil {
						if err := datatx.EmitFileUploadedEvent(up.SpaceOwner(), &executant, &ref, m.publisher); err != nil {
							appctx.GetLogger(context.Background()).Error().Err(err).Msg("failed to publish FileUploaded event")
						}
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
			setHeaders(fs, w, r)
			handler.PostFile(w, r)
		case "HEAD":
			handler.HeadFile(w, r)
		case "PATCH":
			metrics.UploadsActive.Add(1)
			defer func() {
				metrics.UploadsActive.Sub(1)
			}()
			// set etag, mtime and file id
			setHeaders(fs, w, r)
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

func setHeaders(fs storage.FS, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := path.Base(r.URL.Path)
	datastore, ok := fs.(tusd.DataStore)
	if !ok {
		appctx.GetLogger(ctx).Error().Interface("fs", fs).Msg("storage is not a tus datastore")
		return
	}
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
	expires := info.MetaData["expires"]
	if expires != "" {
		w.Header().Set(net.HeaderTusUploadExpires, expires)
	}
	resourceid := provider.ResourceId{
		StorageId: info.MetaData["providerID"],
		SpaceId:   info.Storage["SpaceRoot"],
		OpaqueId:  info.Storage["NodeId"],
	}
	w.Header().Set(net.HeaderOCFileID, storagespace.FormatResourceID(resourceid))
}
