// Copyright 2018-2020 CERN
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

package data

import (
	"encoding/json"
	"sort"

	"github.com/cs3org/reva/v2/pkg/mentix/utils/network"
	"github.com/pkg/errors"
)

// SiteInformation holds the most basic information about a site.
type SiteInformation struct {
	ID       string
	Name     string
	FullName string
}

// QueryAvailableSites uses Mentix to query a list of all available (registered) sites.
func QueryAvailableSites(mentixHost, dataEndpoint string) ([]SiteInformation, error) {
	mentixURL, err := network.GenerateURL(mentixHost, dataEndpoint, network.URLParams{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate Mentix URL")
	}

	data, err := network.ReadEndpoint(mentixURL, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the Mentix endpoint")
	}

	// Decode the data into a simplified, reduced data type
	type siteData struct {
		Sites []SiteInformation
	}
	sites := siteData{}
	if err := json.Unmarshal(data, &sites); err != nil {
		return nil, errors.Wrap(err, "error while decoding the JSON data")
	}

	// Sort the sites alphabetically by their names
	sort.Slice(sites.Sites, func(i, j int) bool {
		return sites.Sites[i].Name < sites.Sites[j].Name
	})

	return sites.Sites, nil
}

// QuerySiteName uses Mentix to query the name of a site given by its ID.
func QuerySiteName(siteID string, fullName bool, mentixHost, dataEndpoint string) (string, error) {
	sites, err := QueryAvailableSites(mentixHost, dataEndpoint)
	if err != nil {
		return "", err
	}

	index := len(sites)
	for i, site := range sites {
		if site.ID == siteID {
			index = i
			break
		}
	}

	if index != len(sites) {
		if fullName {
			return sites[index].FullName, nil
		}

		return sites[index].Name, nil
	}

	return "", errors.Errorf("no site with ID %v found", siteID)
}
