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

package ocmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cs3org/reva/v2/pkg/appctx"
)

type configData struct {
	Enabled       bool            `json:"enabled" xml:"enabled"`
	APIVersion    string          `json:"apiVersion" xml:"apiVersion"`
	Host          string          `json:"host" xml:"host"`
	Endpoint      string          `json:"endPoint" xml:"endPoint"`
	Provider      string          `json:"provider" xml:"provider"`
	ResourceTypes []resourceTypes `json:"resourceTypes" xml:"resourceTypes"`
}

type resourceTypes struct {
	Name       string                 `json:"name"`
	ShareTypes []string               `json:"shareTypes"`
	Protocols  resourceTypesProtocols `json:"protocols"`
}

type resourceTypesProtocols struct {
	Webdav string `json:"webdav"`
}

type configHandler struct {
	c configData
}

func (h *configHandler) init(c *Config) {
	h.c = c.Config
	if h.c.APIVersion == "" {
		h.c.APIVersion = "1.0-proposal1"
	}
	if h.c.Host == "" {
		h.c.Host = "localhost"
	}
	if h.c.Provider == "" {
		h.c.Provider = "cernbox"
	}
	h.c.Enabled = true
	if len(c.Prefix) > 0 {
		h.c.Endpoint = fmt.Sprintf("https://%s/%s", h.c.Host, c.Prefix)
	} else {
		h.c.Endpoint = fmt.Sprintf("https://%s", h.c.Host)
	}
	h.c.ResourceTypes = []resourceTypes{{
		Name:       "file",
		ShareTypes: []string{"user"},
		Protocols: resourceTypesProtocols{
			Webdav: fmt.Sprintf("/%s/ocm_webdav", h.c.Provider),
		},
	}}
}

func (h *configHandler) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.GetLogger(r.Context())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		indentedConf, _ := json.MarshalIndent(h.c, "", "   ")
		if _, err := w.Write(indentedConf); err != nil {
			log.Err(err).Msg("Error writing to ResponseWriter")
		}

	})
}
