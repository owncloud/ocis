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
	"net/http"
	"path"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/propfind"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp/router"
)

// PublicFileHandler handles requests on a shared file. it needs to be wrapped in a collection
type PublicFileHandler struct {
	namespace string
}

func (h *PublicFileHandler) init(ns string) error {
	h.namespace = path.Join("/", ns)
	return nil
}

// Handler handles requests
func (h *PublicFileHandler) Handler(s *svc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := appctx.GetLogger(ctx)
		token, relativePath := router.ShiftPath(r.URL.Path)

		base := path.Join(ctx.Value(net.CtxKeyBaseURI).(string), token)
		ctx = context.WithValue(ctx, net.CtxKeyBaseURI, base)
		r = r.WithContext(ctx)

		log.Debug().Str("relativePath", relativePath).Msg("PublicFileHandler func")

		if relativePath != "" && relativePath != "/" {
			// accessing the file

			switch r.Method {
			case MethodPropfind:
				s.handlePropfindOnToken(w, r, h.namespace, false)
			case http.MethodGet:
				s.handlePathGet(w, r, h.namespace)
			case http.MethodOptions:
				s.handleOptions(w, r)
			case http.MethodHead:
				s.handlePathHead(w, r, h.namespace)
			case http.MethodPut:
				s.handlePathPut(w, r, h.namespace)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			// accessing the virtual parent folder
			switch r.Method {
			case MethodPropfind:
				s.handlePropfindOnToken(w, r, h.namespace, true)
			case http.MethodOptions:
				s.handleOptions(w, r)
			case http.MethodHead:
				s.handlePathHead(w, r, h.namespace)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	})
}

// ns is the namespace that is prefixed to the path in the cs3 namespace
func (s *svc) handlePropfindOnToken(w http.ResponseWriter, r *http.Request, ns string, onContainer bool) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "token_propfind")
	defer span.End()

	tokenStatInfo := ctx.Value(tokenStatInfoKey{}).(*provider.ResourceInfo)
	sublog := appctx.GetLogger(ctx).With().Interface("tokenStatInfo", tokenStatInfo).Logger()
	sublog.Debug().Msg("handlePropfindOnToken")

	dh := r.Header.Get(net.HeaderDepth)
	depth, err := net.ParseDepth(dh)
	if err != nil {
		sublog.Debug().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pf, status, err := propfind.ReadPropfind(r.Body)
	if err != nil {
		sublog.Debug().Err(err).Msg("error reading propfind request")
		w.WriteHeader(status)
		return
	}

	infos := s.getPublicFileInfos(onContainer, depth == net.DepthZero, tokenStatInfo)

	prefer := net.ParsePrefer(r.Header.Get("prefer"))
	returnMinimal := prefer[net.HeaderPreferReturn] == "minimal"

	propRes, err := propfind.MultistatusResponse(ctx, &pf, infos, s.c.PublicURL, ns, nil, returnMinimal)
	if err != nil {
		sublog.Error().Err(err).Msg("error formatting propfind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(net.HeaderDav, "1, 3, extended-mkcol")
	w.Header().Set(net.HeaderContentType, "application/xml; charset=utf-8")
	w.Header().Set(net.HeaderVary, net.HeaderPrefer)
	if returnMinimal {
		w.Header().Set(net.HeaderPreferenceApplied, "return=minimal")
	}
	w.WriteHeader(http.StatusMultiStatus)
	if _, err := w.Write(propRes); err != nil {
		sublog.Err(err).Msg("error writing response")
	}
}

// there are only two possible entries
// 1. the non existing collection
// 2. the shared file
func (s *svc) getPublicFileInfos(onContainer, onlyRoot bool, i *provider.ResourceInfo) []*provider.ResourceInfo {
	infos := []*provider.ResourceInfo{}
	if onContainer {
		// copy link-share data if present
		// we don't copy everything because the checksum should not be present
		var o *typesv1beta1.Opaque
		if i.Opaque != nil && i.Opaque.Map != nil && i.Opaque.Map["link-share"] != nil {
			o = &typesv1beta1.Opaque{
				Map: map[string]*typesv1beta1.OpaqueEntry{
					"link-share": i.Opaque.Map["link-share"],
				},
			}
		}
		// always add collection
		infos = append(infos, &provider.ResourceInfo{
			// Opaque carries the link-share data we need when rendering the collection root href
			Opaque: o,
			Path:   path.Dir(i.Path),
			Type:   provider.ResourceType_RESOURCE_TYPE_CONTAINER,
		})
		if onlyRoot {
			return infos
		}
	}

	// add the file info
	infos = append(infos, i)

	return infos
}
