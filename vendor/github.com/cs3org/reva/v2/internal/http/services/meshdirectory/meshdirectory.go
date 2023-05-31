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

package meshdirectory

import (
	"encoding/json"
	"fmt"
	"net/http"

	meshdirectoryweb "github.com/sciencemesh/meshdirectory-web"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"

	"github.com/cs3org/reva/v2/internal/http/services/ocmd"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/rhttp/router"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/cs3org/reva/v2/pkg/rhttp/global"
	"github.com/mitchellh/mapstructure"
)

func init() {
	global.Register("meshdirectory", New)
}

type config struct {
	Prefix     string `mapstructure:"prefix"`
	GatewaySvc string `mapstructure:"gatewaysvc"`
}

func (c *config) init() {
	c.GatewaySvc = sharedconf.GetGatewaySVC(c.GatewaySvc)

	if c.Prefix == "" {
		c.Prefix = "meshdir"
	}
}

type svc struct {
	conf *config
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New returns a new Mesh Directory HTTP service
func New(m map[string]interface{}, log *zerolog.Logger) (global.Service, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	c.init()

	service := &svc{
		conf: c,
	}
	return service, nil
}

// Service prefix
func (s *svc) Prefix() string {
	return s.conf.Prefix
}

// Unprotected endpoints
func (s *svc) Unprotected() []string {
	return []string{"/"}
}

// Close performs cleanup.
func (s *svc) Close() error {
	return nil
}

func (s *svc) getClient() (gateway.GatewayAPIClient, error) {
	selector, err := pool.GatewaySelector(s.conf.GatewaySvc)
	if err != nil {
		return nil, err
	}

	return selector.Next()
}

func (s *svc) serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	gatewayClient, err := s.getClient()
	if err != nil {
		ocmd.WriteError(w, r, ocmd.APIErrorServerError,
			fmt.Sprintf("error getting grpc client on addr: %v", s.conf.GatewaySvc), err)
		return
	}

	providers, err := gatewayClient.ListAllProviders(ctx, &providerv1beta1.ListAllProvidersRequest{})
	if err != nil {
		ocmd.WriteError(w, r, ocmd.APIErrorServerError, "error listing all providers", err)
		return
	}

	jsonResponse, err := json.Marshal(providers.Providers)
	if err != nil {
		ocmd.WriteError(w, r, ocmd.APIErrorServerError, "error marshalling providers data", err)
		return
	}

	// Write response
	_, err = w.Write(jsonResponse)
	if err != nil {
		ocmd.WriteError(w, r, ocmd.APIErrorServerError, "error writing providers data", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP service handler
func (s *svc) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = router.ShiftPath(r.URL.Path)
		switch head {
		case "providers":
			s.serveJSON(w, r)
			return
		default:
			r.URL.Path = head + r.URL.Path
			meshdirectoryweb.ServeMeshDirectorySPA(w, r)
			return
		}
	})
}
