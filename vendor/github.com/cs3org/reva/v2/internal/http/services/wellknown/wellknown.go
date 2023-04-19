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

package wellknown

import (
	"net/http"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp/global"
	"github.com/cs3org/reva/v2/pkg/rhttp/router"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

func init() {
	global.Register("wellknown", New)
}

type config struct {
	Prefix                string `mapstructure:"prefix"`
	Issuer                string `mapstructure:"issuer"`
	AuthorizationEndpoint string `mapstructure:"authorization_endpoint"`
	JwksURI               string `mapstructure:"jwks_uri"`
	TokenEndpoint         string `mapstructure:"token_endpoint"`
	RevocationEndpoint    string `mapstructure:"revocation_endpoint"`
	IntrospectionEndpoint string `mapstructure:"introspection_endpoint"`
	UserinfoEndpoint      string `mapstructure:"userinfo_endpoint"`
	EndSessionEndpoint    string `mapstructure:"end_session_endpoint"`
}

func (c *config) init() {
	if c.Prefix == "" {
		c.Prefix = ".well-known"
	}
}

type svc struct {
	conf    *config
	handler http.Handler
}

// New returns a new webuisvc
func New(m map[string]interface{}, log *zerolog.Logger) (global.Service, error) {
	conf := &config{}
	if err := mapstructure.Decode(m, conf); err != nil {
		return nil, err
	}

	conf.init()

	s := &svc{
		conf: conf,
	}
	s.setHandler()
	return s, nil
}

func (s *svc) Close() error {
	return nil
}

func (s *svc) Prefix() string {
	return s.conf.Prefix
}

func (s *svc) Handler() http.Handler {
	return s.handler
}

func (s *svc) Unprotected() []string {
	return []string{
		"/openid-configuration",
	}
}

func (s *svc) setHandler() {
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.GetLogger(r.Context())
		var head string
		head, r.URL.Path = router.ShiftPath(r.URL.Path)
		log.Info().Msgf("wellknown routing: head=%s tail=%s", head, r.URL.Path)
		switch head {
		case "webfinger":
			s.doWebfinger(w, r)
		case "openid-configuration":
			s.doOpenidConfiguration(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
