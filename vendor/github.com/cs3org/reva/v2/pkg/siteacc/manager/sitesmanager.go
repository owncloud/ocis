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

package manager

import (
	"strings"
	"sync"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// SitesManager is responsible for all sites related tasks.
type SitesManager struct {
	conf *config.Configuration
	log  *zerolog.Logger

	storage data.Storage

	sites data.Sites

	mutex sync.RWMutex
}

func (mngr *SitesManager) initialize(storage data.Storage, conf *config.Configuration, log *zerolog.Logger) error {
	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	mngr.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	mngr.log = log

	if storage == nil {
		return errors.Errorf("no storage provided")
	}
	mngr.storage = storage

	mngr.sites = make(data.Sites, 0, 32) // Reserve some space for sites
	mngr.readAllSites()

	return nil
}

func (mngr *SitesManager) readAllSites() {
	if sites, err := mngr.storage.ReadSites(); err == nil {
		mngr.sites = *sites
	} else {
		// Just warn when not being able to read sites
		mngr.log.Warn().Err(err).Msg("error while reading sites")
	}
}

func (mngr *SitesManager) writeAllSites() {
	if err := mngr.storage.WriteSites(&mngr.sites); err != nil {
		// Just warn when not being able to write sites
		mngr.log.Warn().Err(err).Msg("error while writing sites")
	}
}

// GetSite retrieves the site with the given ID, creating it first if necessary.
func (mngr *SitesManager) GetSite(id string, cloneSite bool) (*data.Site, error) {
	mngr.mutex.RLock()
	defer mngr.mutex.RUnlock()

	site, err := mngr.getSite(id)
	if err != nil {
		return nil, err
	}

	if cloneSite {
		site = site.Clone(false)
	}

	return site, nil
}

// FindSite returns the site specified by the ID if one exists.
func (mngr *SitesManager) FindSite(id string) *data.Site {
	site, _ := mngr.findSite(id)
	return site
}

// UpdateSite updates the site identified by the site ID; if no such site exists, one will be created first.
func (mngr *SitesManager) UpdateSite(siteData *data.Site) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	site, err := mngr.getSite(siteData.ID)
	if err != nil {
		return errors.Wrap(err, "site to update not found")
	}

	if err := site.Update(siteData, mngr.conf.Security.CredentialsPassphrase); err == nil {
		mngr.storage.SiteUpdated(site)
		mngr.writeAllSites()
	} else {
		return errors.Wrap(err, "error while updating site")
	}

	return nil
}

// CloneSites retrieves all sites currently stored by cloning the data, thus avoiding race conflicts and making outside modifications impossible.
func (mngr *SitesManager) CloneSites(eraseCredentials bool) data.Sites {
	mngr.mutex.RLock()
	defer mngr.mutex.RUnlock()

	clones := make(data.Sites, 0, len(mngr.sites))
	for _, site := range mngr.sites {
		clones = append(clones, site.Clone(eraseCredentials))
	}

	return clones
}

func (mngr *SitesManager) getSite(id string) (*data.Site, error) {
	site, err := mngr.findSite(id)
	if site == nil {
		site, err = mngr.createSite(id)
	}
	return site, err
}

func (mngr *SitesManager) createSite(id string) (*data.Site, error) {
	site, err := data.NewSite(id)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating site")
	}
	mngr.sites = append(mngr.sites, site)
	mngr.storage.SiteAdded(site)
	mngr.writeAllSites()
	return site, nil
}

func (mngr *SitesManager) findSite(id string) (*data.Site, error) {
	if len(id) == 0 {
		return nil, errors.Errorf("no search ID specified")
	}

	site := mngr.findSiteByPredicate(func(site *data.Site) bool { return strings.EqualFold(site.ID, id) })
	if site != nil {
		return site, nil
	}

	return nil, errors.Errorf("no site found matching the specified ID")
}

func (mngr *SitesManager) findSiteByPredicate(predicate func(*data.Site) bool) *data.Site {
	for _, site := range mngr.sites {
		if predicate(site) {
			return site
		}
	}
	return nil
}

// NewSitesManager creates a new sites manager instance.
func NewSitesManager(storage data.Storage, conf *config.Configuration, log *zerolog.Logger) (*SitesManager, error) {
	mngr := &SitesManager{}
	if err := mngr.initialize(storage, conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the sites manager")
	}
	return mngr, nil
}
