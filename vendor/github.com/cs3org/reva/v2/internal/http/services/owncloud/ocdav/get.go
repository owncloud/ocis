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

package ocdav

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/grpc/services/storageprovider"
	"github.com/cs3org/reva/v2/internal/http/services/datagateway"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/errors"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/spacelookup"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func (s *svc) handlePathGet(w http.ResponseWriter, r *http.Request, ns string) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "get")
	defer span.End()

	fn := path.Join(ns, r.URL.Path)

	sublog := appctx.GetLogger(ctx).With().Str("path", fn).Str("svc", "ocdav").Str("handler", "get").Logger()

	client, err := s.getClient()
	if err != nil {
		sublog.Error().Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	space, status, err := spacelookup.LookUpStorageSpaceForPath(ctx, client, fn)
	if err != nil {
		sublog.Error().Err(err).Str("path", fn).Msg("failed to look up storage space")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if status.Code != rpc.Code_CODE_OK {
		errors.HandleErrorStatus(&sublog, w, status)
		return
	}

	s.handleGet(ctx, w, r, spacelookup.MakeRelativeReference(space, fn, false), "spaces", sublog)
}

func (s *svc) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request, ref *provider.Reference, dlProtocol string, log zerolog.Logger) {
	client, err := s.getClient()
	if err != nil {
		log.Error().Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sReq := &provider.StatRequest{Ref: ref}
	sRes, err := client.Stat(ctx, sReq)
	switch {
	case err != nil:
		log.Error().Err(err).Msg("error sending grpc stat request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	case sRes.Status.Code != rpc.Code_CODE_OK:
		errors.HandleErrorStatus(&log, w, sRes.Status)
		return
	case sRes.Info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER:
		log.Warn().Msg("resource is a folder and cannot be downloaded")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	dReq := &provider.InitiateFileDownloadRequest{Ref: ref}
	dRes, err := client.InitiateFileDownload(ctx, dReq)
	if err != nil {
		log.Error().Err(err).Msg("error initiating file download")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if dRes.Status.Code != rpc.Code_CODE_OK {
		errors.HandleErrorStatus(&log, w, dRes.Status)
		return
	}

	var ep, token string
	for _, p := range dRes.Protocols {
		if p.Protocol == dlProtocol {
			ep, token = p.DownloadEndpoint, p.Token
		}
	}

	httpReq, err := rhttp.NewRequest(ctx, http.MethodGet, ep, nil)
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set(datagateway.TokenTransportHeader, token)

	if r.Header.Get(net.HeaderRange) != "" {
		httpReq.Header.Set(net.HeaderRange, r.Header.Get(net.HeaderRange))
	}

	httpClient := s.client

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("error performing http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusPartialContent {
		w.WriteHeader(httpRes.StatusCode)
		return
	}

	info := sRes.Info

	w.Header().Set(net.HeaderContentType, info.MimeType)
	w.Header().Set(net.HeaderContentDisposistion, "attachment; filename*=UTF-8''"+
		path.Base(info.Path)+"; filename=\""+path.Base(info.Path)+"\"")
	w.Header().Set(net.HeaderETag, info.Etag)
	w.Header().Set(net.HeaderOCFileID, storagespace.FormatResourceID(*info.Id))
	w.Header().Set(net.HeaderOCETag, info.Etag)
	t := utils.TSToTime(info.Mtime).UTC()
	lastModifiedString := t.Format(time.RFC1123Z)
	w.Header().Set(net.HeaderLastModified, lastModifiedString)

	if httpRes.StatusCode == http.StatusPartialContent {
		w.Header().Set(net.HeaderContentRange, httpRes.Header.Get(net.HeaderContentRange))
		w.Header().Set(net.HeaderContentLength, httpRes.Header.Get(net.HeaderContentLength))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.Header().Set(net.HeaderContentLength, strconv.FormatUint(info.Size, 10))
	}
	if info.Checksum != nil {
		w.Header().Set(net.HeaderOCChecksum, fmt.Sprintf("%s:%s", strings.ToUpper(string(storageprovider.GRPC2PKGXS(info.Checksum.Type))), info.Checksum.Sum))
	}
	var c int64
	if c, err = io.Copy(w, httpRes.Body); err != nil {
		log.Error().Err(err).Msg("error finishing copying data to response")
	}
	if httpRes.Header.Get(net.HeaderContentLength) != "" {
		i, err := strconv.ParseInt(httpRes.Header.Get(net.HeaderContentLength), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("content-length", httpRes.Header.Get(net.HeaderContentLength)).Msg("invalid content length in datagateway response")
		}
		if i != c {
			log.Error().Int64("content-length", i).Int64("transferred-bytes", c).Msg("content length vs transferred bytes mismatch")
		}
	}
	// TODO we need to send the If-Match etag in the GET to the datagateway to prevent race conditions between stating and reading the file
}

func (s *svc) handleSpacesGet(w http.ResponseWriter, r *http.Request, spaceID string) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "spaces_get")
	defer span.End()

	sublog := appctx.GetLogger(ctx).With().Str("path", r.URL.Path).Str("spaceid", spaceID).Str("handler", "get").Logger()

	ref, err := spacelookup.MakeStorageSpaceReference(spaceID, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.handleGet(ctx, w, r, &ref, "spaces", sublog)
}
