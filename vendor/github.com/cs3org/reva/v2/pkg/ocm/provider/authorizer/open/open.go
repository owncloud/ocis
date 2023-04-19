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

package open

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/provider"
	"github.com/cs3org/reva/v2/pkg/ocm/provider/authorizer/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	registry.Register("open", New)
}

// New returns a new authorizer object.
func New(m map[string]interface{}) (provider.Authorizer, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	c.init()

	f, err := os.ReadFile(c.Providers)
	if err != nil {
		return nil, err
	}
	providers := []*ocmprovider.ProviderInfo{}
	err = json.Unmarshal(f, &providers)
	if err != nil {
		return nil, err
	}

	a := &authorizer{}
	a.providers = a.getOCMProviders(providers)

	return a, nil
}

type config struct {
	// Users holds a path to a file containing json conforming the Users struct
	Providers string `mapstructure:"providers"`
}

func (c *config) init() {
	if c.Providers == "" {
		c.Providers = "/etc/revad/ocm-providers.json"
	}
}

type authorizer struct {
	providers []*ocmprovider.ProviderInfo
}

func (a *authorizer) GetInfoByDomain(ctx context.Context, domain string) (*ocmprovider.ProviderInfo, error) {
	for _, p := range a.providers {
		if strings.Contains(p.Domain, domain) {
			return p, nil
		}
	}
	return nil, errtypes.NotFound(domain)
}

func (a *authorizer) IsProviderAllowed(ctx context.Context, provider *ocmprovider.ProviderInfo) error {
	return nil
}

func (a *authorizer) ListAllProviders(ctx context.Context) ([]*ocmprovider.ProviderInfo, error) {
	return a.providers, nil
}

func (a *authorizer) getOCMProviders(providers []*ocmprovider.ProviderInfo) (po []*ocmprovider.ProviderInfo) {
	for _, p := range providers {
		_, err := a.getOCMHost(p)
		if err == nil {
			po = append(po, p)
		}
	}
	return
}

func (a *authorizer) getOCMHost(provider *ocmprovider.ProviderInfo) (string, error) {
	for _, s := range provider.Services {
		if s.Endpoint.Type.Name == "OCM" {
			ocmHost, err := url.Parse(s.Host)
			if err != nil {
				return "", errors.Wrap(err, "json: error parsing OCM host URL")
			}
			return ocmHost.Host, nil
		}
	}
	return "", errtypes.NotFound("OCM Host")
}
