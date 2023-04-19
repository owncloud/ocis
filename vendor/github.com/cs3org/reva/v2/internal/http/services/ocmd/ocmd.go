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
	"net/http"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp/global"
	"github.com/cs3org/reva/v2/pkg/rhttp/router"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/cs3org/reva/v2/pkg/smtpclient"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

func init() {
	global.Register("ocmd", New)
}

// Config holds the config options that need to be passed down to all ocdav handlers
type Config struct {
	SMTPCredentials  *smtpclient.SMTPCredentials `mapstructure:"smtp_credentials"`
	Prefix           string                      `mapstructure:"prefix"`
	Host             string                      `mapstructure:"host"`
	GatewaySvc       string                      `mapstructure:"gatewaysvc"`
	MeshDirectoryURL string                      `mapstructure:"mesh_directory_url"`
	Config           configData                  `mapstructure:"config"`
}

func (c *Config) init() {
	c.GatewaySvc = sharedconf.GetGatewaySVC(c.GatewaySvc)

	// if c.Prefix == "" {
	// 	c.Prefix = "ocm"
	// }
}

type svc struct {
	Conf                 *Config
	SharesHandler        *sharesHandler
	NotificationsHandler *notificationsHandler
	ConfigHandler        *configHandler
	InvitesHandler       *invitesHandler
	SendHandler          *sendHandler
}

// New returns a new ocmd object
func New(m map[string]interface{}, log *zerolog.Logger) (global.Service, error) {

	conf := &Config{}
	if err := mapstructure.Decode(m, conf); err != nil {
		return nil, err
	}
	conf.init()

	s := &svc{
		Conf: conf,
	}
	s.SharesHandler = new(sharesHandler)
	s.NotificationsHandler = new(notificationsHandler)
	s.ConfigHandler = new(configHandler)
	s.InvitesHandler = new(invitesHandler)
	s.SendHandler = new(sendHandler)
	s.SharesHandler.init(s.Conf)
	s.NotificationsHandler.init(s.Conf)
	log.Debug().Str("initializing ConfigHandler Host", s.Conf.Host)

	s.ConfigHandler.init(s.Conf)
	s.InvitesHandler.init(s.Conf)
	s.SendHandler.init(s.Conf)

	return s, nil
}

// Close performs cleanup.
func (s *svc) Close() error {
	return nil
}

func (s *svc) Prefix() string {
	return s.Conf.Prefix
}

func (s *svc) Unprotected() []string {
	return []string{"/invites/accept", "/shares", "/ocm-provider", "/notifications"}
}

func (s *svc) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		log := appctx.GetLogger(ctx)

		var head string
		head, r.URL.Path = router.ShiftPath(r.URL.Path)
		log.Debug().Str("head", head).Str("tail", r.URL.Path).Msg("http routing")

		switch head {
		case "ocm-provider":
			s.ConfigHandler.Handler().ServeHTTP(w, r)
			return
		case "shares":
			s.SharesHandler.Handler().ServeHTTP(w, r)
			return
		case "notifications":
			s.NotificationsHandler.Handler().ServeHTTP(w, r)
			return
		case "invites":
			s.InvitesHandler.Handler().ServeHTTP(w, r)
			return
		case "send":
			s.SendHandler.Handler().ServeHTTP(w, r)
		}

		log.Warn().Msg("request not handled")
		w.WriteHeader(http.StatusNotFound)
	})
}
