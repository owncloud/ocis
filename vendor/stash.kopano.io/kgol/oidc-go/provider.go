/*
 * Copyright 2019 Kopano
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/desertbit/timer"
	"gopkg.in/square/go-jose.v2"
)

// Provider represents an OpenID Connect server's configuration.
type Provider struct {
	mutex sync.RWMutex

	initialized bool
	ready       chan struct{}
	started     chan error
	cancel      context.CancelFunc

	issuer       string
	wellKnownURI *url.URL
	jwksURI      *url.URL

	logger     logger
	httpClient *http.Client
	httpHeader http.Header

	wellKnown *WellKnown
	jwks      *jose.JSONWebKeySet
}

// ProviderConfig bundles configuration for a Provider.
type ProviderConfig struct {
	HTTPClient   *http.Client
	HTTPHeader   http.Header
	WellKnownURI *url.URL
	Logger       logger
}

// DefaultProviderConfig is the Provider configuration uses when none was
// explicitly specified.
var DefaultProviderConfig = &ProviderConfig{}

// ProviderDefinition holds immutable provider information.
type ProviderDefinition struct {
	WellKnown *WellKnown
	JWKS      *jose.JSONWebKeySet
}

// A ProviderError is returned for OIDC Provider errors.
type ProviderError struct {
	Err error // The actual error
}

func wrapAsProviderError(err error) error {
	if err == nil {
		return nil
	}

	return &ProviderError{
		Err: err,
	}
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("oidc provider error: %v", e.Err)
}

// These are the errors that can be returned in ProviderError.Err.
var (
	ErrAllreadyInitialized = errors.New("already initialized")
	ErrNotInitialized      = errors.New("not initialized")
	ErrWrongInitialization = errors.New("wrong initialization")
	ErrIssuerMismatch      = errors.New("issuer mismatch")
)

// NewProvider uses OpenID Connect discovery to create a Provider.
func NewProvider(issuer *url.URL, config *ProviderConfig) (*Provider, error) {
	if config == nil {
		config = DefaultProviderConfig
	}

	p := &Provider{
		issuer: issuer.String(),

		httpClient: config.HTTPClient,
		httpHeader: config.HTTPHeader,
	}

	if config.WellKnownURI != nil {
		p.wellKnownURI = config.WellKnownURI
	} else {
		relativeWellKnownURI, err := url.Parse("/.well-known/openid-configuration")
		if err != nil {
			return nil, err
		}
		p.wellKnownURI = issuer.ResolveReference(relativeWellKnownURI)
	}
	if config.Logger != nil {
		p.logger = config.Logger
	} else {
		p.logger = DefaultLogger
	}

	return p, nil
}

// Initialize initializes the associated Provider with the provided Context. If
// updates and/or errors channels apre provided, those channels receive any
// update or update error from the tasks resulting from the initialization. Any
// of thes channels can be nil, disabling the corresponding events being sent.
func (p *Provider) Initialize(ctx context.Context, updates chan *ProviderDefinition, errors chan error) error {
	p.mutex.Lock()
	if p.initialized {
		p.mutex.Unlock()
		return wrapAsProviderError(ErrAllreadyInitialized)
	}

	c, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.initialized = true

	started := make(chan error, 1)
	p.started = started
	go p.start(c, started, updates, errors)
	p.mutex.Unlock()

	err := <-started
	return wrapAsProviderError(err)
}

// Shutdown stops the associated Provider and waits for it to do so.
func (p *Provider) Shutdown() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.initialized {
		return wrapAsProviderError(ErrNotInitialized)
	}

	p.cancel()
	err := <-p.started

	p.cancel = nil
	p.started = nil
	p.initialized = false
	p.ready = nil

	if err == context.Canceled {
		return nil
	}
	return wrapAsProviderError(err)
}

// Ready returns a channel that's closed when the associated Provider is ready.
func (p *Provider) Ready() <-chan struct{} {
	p.mutex.RLock()
	ready := p.ready
	p.mutex.RUnlock()

	return ready
}

func (p *Provider) start(ctx context.Context, started chan error, updates chan *ProviderDefinition, errors chan error) {
	p.mutex.Lock()
	if !p.initialized || started != p.started {
		p.mutex.Unlock()
		started <- ErrWrongInitialization
		return
	}

	readystate := false
	ready := make(chan struct{})
	p.ready = ready
	p.mutex.Unlock()
	started <- nil

	var wellKnown *WellKnown
	var jwks *jose.JSONWebKeySet

	var ignore error
	dLoad := true
	dUpdated := false
	dExpireTimer := timer.NewTimer(DefaultJSONFetchExpiry)
	kLoad := true
	kUpdated := false
	kExpireTimer := timer.NewTimer(DefaultJSONFetchExpiry)
	for {
		ignore = nil
		dUpdated = false
		kUpdated = false

		if dLoad {
			dst := WellKnown{}
			p.logger.Printf("fetching OIDC provider discover document: %v\n", p.wellKnownURI)
			expires, err := fetchJSON(ctx, p.wellKnownURI, &dst, p.httpClient, p.httpHeader)
			if err != nil {
				ignore = fmt.Errorf("failed to fetch discover document: %v", err)
				if errors == nil {
					p.logger.Printf("OIDC provider %v\n", ignore)
				}
			} else {
				wellKnown = &dst
				dUpdated = true
			}
			dLoad = false
			dExpireTimer.Reset(expires)
			p.logger.Printf("ODIC provider discover document loaded, expires: %v\n", expires)
		}
		if wellKnown != nil && kLoad {
			dst := jose.JSONWebKeySet{}
			if wellKnown.JwksURI != "" {
				jwksURI, err := url.Parse(wellKnown.JwksURI)
				if err != nil {
					ignore = fmt.Errorf("discover document invalid jwks_uri: %v", err)
					if errors == nil {
						p.logger.Printf("OIDC provider %v\n", ignore)
					}
				} else {
					p.logger.Printf("fetching OIDC provider jwks: %v", wellKnown.JwksURI)
					expires, err := fetchJSON(ctx, jwksURI, &dst, p.httpClient, p.httpHeader)
					if err != nil {
						ignore = fmt.Errorf("failed to fetch jwks: %v", err)
						if errors == nil {
							p.logger.Printf("OIDC provider %v\n", ignore)
						}
					} else {
						jwks = &dst
						kUpdated = true
					}
					kLoad = false
					kExpireTimer.Reset(expires)
					p.logger.Printf("OIDC provider jwks loaded, expires: %v\n", expires)
				}
			}
		}

		p.mutex.Lock()
		if dUpdated {
			if wellKnown.Issuer != p.issuer {
				if errors == nil {
					p.logger.Printf("OIDC provider issuer mismatch: %v != %v\n", wellKnown.Issuer, p.issuer)
				}
				ignore = ErrIssuerMismatch
			}

			if ignore == nil {
				p.logger.Printf("OIDC provider discover document updated\n")
				p.wellKnown = wellKnown
			}
		}
		if kUpdated {
			if ignore == nil {
				p.logger.Printf("ODIC provider jwks updated\n")
				p.jwks = jwks
			}
		}

		p.mutex.Unlock()

		if updates != nil && ignore == nil && (dUpdated || kUpdated) {
			p.logger.Printf("OIDC provider triggering update")
			updates <- &ProviderDefinition{
				WellKnown: wellKnown,
				JWKS:      jwks,
			}
		} else if errors != nil && ignore != nil {
			p.logger.Printf("OIDC provider triggering errors")
			errors <- wrapAsProviderError(ignore)
		}

		if !readystate {
			if p.wellKnown != nil && p.jwks != nil {
				readystate = true
				close(ready)
			}
		}

		select {
		case <-ctx.Done():
			started <- ctx.Err()
			return
		case <-dExpireTimer.C:
			dLoad = true
		case <-kExpireTimer.C:
			kLoad = true
		}
	}
}
