/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package authorities

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Registry implements the registry for registered authorities.
type Registry struct {
	mutex sync.RWMutex

	baseURI *url.URL

	defaultID   string
	authorities map[string]AuthorityRegistration

	logger logrus.FieldLogger
}

// NewRegistry creates a new authorizations Registry with the provided parameters.
func NewRegistry(ctx context.Context, baseURI *url.URL, registrationConfFilepath string, logger logrus.FieldLogger) (*Registry, error) {
	registryData := &authorityRegistryData{}

	if registrationConfFilepath != "" {
		logger.Debugf("parsing authorities registration conf from %v", registrationConfFilepath)
		registryFile, err := ioutil.ReadFile(registrationConfFilepath)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(registryFile, registryData)
		if err != nil {
			return nil, err
		}
	}

	r := &Registry{
		baseURI: baseURI,

		authorities: make(map[string]AuthorityRegistration),

		logger: logger,
	}

	var defaultAuthorityRegistrationData *authorityRegistrationData
	var defaultAuthority AuthorityRegistration
	for _, registrationData := range registryData.Authorities {
		var authority AuthorityRegistration
		var validateErr error

		if registrationData.ID == "" {
			registrationData.ID = registrationData.Name
			r.logger.WithField("id", registrationData.ID).Warnln("authority has no id, using name")
		}

		switch registrationData.AuthorityType {
		case AuthorityTypeOIDC:
			authority, validateErr = newOIDCAuthorityRegistration(r, registrationData)
		case AuthorityTypeSAML2:
			authority, validateErr = newSAML2AuthorityRegistration(r, registrationData)
		}

		fields := logrus.Fields{
			"id":             registrationData.ID,
			"authority_type": registrationData.AuthorityType,
			"insecure":       registrationData.Insecure,
			"trusted":        registrationData.Trusted,
			"default":        registrationData.Default,
			"alias_required": registrationData.IdentityAliasRequired,
		}

		if validateErr != nil {
			logger.WithError(validateErr).WithFields(fields).Warnln("skipped registration of invalid authority entry")
			continue
		}

		if authority == nil {
			logger.WithFields(fields).Warnln("skipped registration of authority of unknown type")
			continue
		}

		if registerErr := r.Register(authority); registerErr != nil {
			logger.WithError(registerErr).WithFields(fields).Warnln("skipped registration of invalid authority")
			continue
		}

		if registrationData.Default || defaultAuthorityRegistrationData == nil {
			if defaultAuthorityRegistrationData == nil || !defaultAuthorityRegistrationData.Default {
				defaultAuthorityRegistrationData = registrationData
				defaultAuthority = authority
			} else {
				logger.Warnln("ignored default authority flag since already have a default")
			}
		} else {
			// TODO(longsleep): Implement authority selection.
			logger.Warnln("non-default additional authorities are not supported yet")
		}

		go func() {
			if initializeErr := authority.Initialize(ctx, r); initializeErr != nil {
				logger.WithError(initializeErr).WithFields(fields).Warnln("failed to initialize authority")
			}
		}()

		logger.WithFields(fields).Debugln("registered authority")
	}

	if defaultAuthority != nil {
		if defaultAuthorityRegistrationData.Default {
			r.defaultID = defaultAuthorityRegistrationData.ID
			logger.WithField("id", defaultAuthorityRegistrationData.ID).Infoln("using external default authority")
		} else {
			logger.Warnln("non-default authorities are not supported yet")
		}
	}

	return r, nil
}

// Register validates the provided authority registration and adds the authority
// to the accociated registry if valid. Returns error otherwise.
func (r *Registry) Register(authority AuthorityRegistration) error {
	id := authority.ID()
	if id == "" {
		return errors.New("no authority id")
	}

	if err := authority.Validate(); err != nil {
		return fmt.Errorf("authority data validation error: %w", err)
	}

	switch authority.AuthorityType() {
	case AuthorityTypeOIDC:
		// breaks
	case AuthorityTypeSAML2:
		// breaks
	default:
		return fmt.Errorf("unknown authority type: %v", authority.AuthorityType())
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.authorities[id] = authority

	return nil
}

// Lookup returns and validates the authority Detail information for the provided
// parameters from the accociated authority registry.
func (r *Registry) Lookup(ctx context.Context, authorityID string) (*Details, error) {
	registration, ok := r.Get(ctx, authorityID)
	if !ok {
		return nil, fmt.Errorf("unknown authority id: %v", authorityID)
	}

	details := registration.Authority()
	return details, nil
}

// Get returns the registered authorities registration for the provided client ID.
func (r *Registry) Get(ctx context.Context, authorityID string) (AuthorityRegistration, bool) {
	if authorityID == "" {
		return nil, false
	}

	// Lookup authority registration.
	r.mutex.RLock()
	registration, ok := r.authorities[authorityID]
	r.mutex.RUnlock()

	return registration, ok
}

// Find returns the first registered authority that satisfies the provided
// selector function.
func (r *Registry) Find(ctx context.Context, selector func(authority AuthorityRegistration) bool) (AuthorityRegistration, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, authority := range r.authorities {
		if selector(authority) {
			return authority, true
		}
	}

	return nil, false
}

// Default returns the default authority from the associated registry if any.
func (r *Registry) Default(ctx context.Context) *Details {
	authority, _ := r.Lookup(ctx, r.defaultID)
	return authority
}
