/*
 * Copyright 2017-2020 Kopano and its licensors
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
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/crewjam/httperr"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/sirupsen/logrus"

	"github.com/libregraph/lico/identity/authorities/samlext"
	"github.com/libregraph/lico/utils"
)

var cleanWhitespaceRegexp = regexp.MustCompile(`\s+`)

type saml2AuthorityRegistration struct {
	registry *Registry
	data     *authorityRegistrationData

	discover         bool
	metadataEndpoint *url.URL

	mutex sync.RWMutex
	ready bool

	serviceProvider             *saml.ServiceProvider
	serviceProviderSigningCerts []*x509.Certificate
}

func newSAML2AuthorityRegistration(registry *Registry, registrationData *authorityRegistrationData) (*saml2AuthorityRegistration, error) {
	ar := &saml2AuthorityRegistration{
		registry: registry,
		data:     registrationData,
	}

	if ar.data.RawMetadataEndpoint != "" {
		if u, err := url.Parse(ar.data.RawMetadataEndpoint); err == nil {
			ar.metadataEndpoint = u
		} else {
			return nil, fmt.Errorf("invalid metadata_endpoint value: %w", err)
		}
	}

	if ar.data.EntityID == "" {
		baseURIString := registry.baseURI.String()
		metadataURI, _ := url.Parse(baseURIString + "/identifier/saml2/metadata") // Use our own meta data
		ar.data.EntityID = metadataURI.String()
	}

	if ar.data.Discover != nil {
		ar.discover = *ar.data.Discover
	}

	if !ar.discover {
		return nil, errors.New("saml2 must use discover")
	}

	return ar, nil
}

func (ar *saml2AuthorityRegistration) ID() string {
	return ar.data.ID
}

func (ar *saml2AuthorityRegistration) Name() string {
	return ar.data.Name
}

func (ar *saml2AuthorityRegistration) AuthorityType() string {
	return ar.data.AuthorityType
}

func (ar *saml2AuthorityRegistration) Authority() *Details {
	details := &Details{
		ID:            ar.data.ID,
		Name:          ar.data.Name,
		AuthorityType: ar.data.AuthorityType,

		Trusted:  ar.data.Trusted,
		Insecure: ar.data.Insecure,

		EndSessionEnabled: ar.data.EndSessionEnabled,

		registration: ar,
	}

	ar.mutex.RLock()
	details.ready = ar.ready
	ar.mutex.RUnlock()

	return details
}

func (ar *saml2AuthorityRegistration) Issuer() string {
	issuer := ar.serviceProvider.IDPMetadata.EntityID
	if issuer == "" {
		issuer = ar.metadataEndpoint.String()
	}
	return issuer
}

func (ar *saml2AuthorityRegistration) Validate() error {
	return nil
}

func (ar *saml2AuthorityRegistration) Initialize(ctx context.Context, registry *Registry) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	if ar.metadataEndpoint == nil {
		return fmt.Errorf("no metadata_endpoint set")
	}
	if ar.data.EntityID == "" {
		return fmt.Errorf("no entity_id set")
	}

	logger := registry.logger.WithFields(logrus.Fields{
		"id":   ar.data.ID,
		"type": AuthorityTypeSAML2,
	})

	var client *http.Client
	if ar.data.Insecure {
		client = utils.InsecureHTTPClient
	} else {
		client = utils.DefaultHTTPClient
	}

	baseURIString := registry.baseURI.String()
	acsURL, _ := url.Parse(baseURIString + "/identifier/saml2/acs")   // Assertion Consumer Service
	sloURL, _ := url.Parse(baseURIString + "/identifier/_/saml2/slo") // Single Logout Service

	go func() {
		var md *saml.EntityDescriptor
		var err error
		for {
			logger.Debugf("fetching SAML2 provider meta data: %s", ar.metadataEndpoint.String())
			md, err = func() (*saml.EntityDescriptor, error) {
				req, fetchErr := http.NewRequest(http.MethodGet, ar.metadataEndpoint.String(), nil)
				if fetchErr != nil {
					return nil, fetchErr
				}
				req = req.WithContext(ctx)
				req.Header.Set("User-Agent", utils.DefaultHTTPUserAgent)

				resp, fetchErr := client.Do(req)
				if fetchErr != nil {
					return nil, fetchErr
				}
				defer resp.Body.Close()
				if resp.StatusCode >= 400 {
					return nil, httperr.Response(*resp)
				}

				data, fetchErr := ioutil.ReadAll(resp.Body)
				if fetchErr != nil {
					return nil, fetchErr
				}
				return samlsp.ParseMetadata(data)
			}()
			if err != nil {
				logger.WithError(err).Errorln("error while saml2 provider meta data update")
			}

			select {
			case <-ctx.Done():
				return
			default:
			}

			if md != nil {
				for {
					var serviceProviderSigningCerts []*x509.Certificate
					serviceProviderSigningCerts, err = getCertsFromMetadata(md, "signing")
					if err != nil {
						break
					}
					if len(serviceProviderSigningCerts) == 0 {
						err = errors.New("no signing certificate in meta data")
						break
					}

					ar.mutex.Lock()

					ar.serviceProviderSigningCerts = serviceProviderSigningCerts
					ar.serviceProvider = &saml.ServiceProvider{
						EntityID:          ar.data.EntityID,
						AcsURL:            *acsURL,
						SloURL:            *sloURL,
						IDPMetadata:       md,
						AllowIDPInitiated: false,
						AuthnNameIDFormat: saml.TransientNameIDFormat,
					}

					ready := ar.ready
					if ar.serviceProvider != nil {
						ar.ready = true
					} else {
						ar.ready = false
					}
					if ready != ar.ready {
						if ar.ready {
							logger.Infoln("authority is now ready")
						} else {
							logger.Warnln("authority is no longer ready")
						}
					} else if !ar.ready {
						logger.Warnln("authority not ready")
					}

					ready = ar.ready

					ar.mutex.Unlock()

					if ready {
						logger.WithFields(logrus.Fields{
							"signing_certs": len(serviceProviderSigningCerts),
							"issuer":        ar.Issuer(),
						}).Debugln("SAML2 provider meta data loaded and initialized")
						return
					}

					break
				}
				if err != nil {
					logger.WithError(err).Errorln("error while initializing saml2 provider from meta data")
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(30 * time.Second):
				// breaks
			}
		}
	}()

	return nil
}

func (ar *saml2AuthorityRegistration) IdentityClaimValue(rawAssertion interface{}) (string, map[string]interface{}, error) {
	assertion, _ := rawAssertion.(*saml.Assertion)
	if assertion == nil {
		return "", nil, errors.New("invalid assertion data")
	}

	icn := ar.data.IdentityClaimName
	if icn == "" {
		icn = "uid" // TODO(longsleep): Use constant.
	}

	var cvs string
	var ok bool
	for _, attributeStatement := range assertion.AttributeStatements {
		for _, attr := range attributeStatement.Attributes {
			values := []string{}
			for _, value := range attr.Values {
				values = append(values, value.Value)
			}
			ar.registry.logger.WithFields(logrus.Fields{
				"FriendlyName": attr.FriendlyName,
				"Name":         attr.Name,
				"NameFormat":   attr.NameFormat,
				"Values":       values,
			}).Debugln("saml2 attributeStatement")

			if !ok {
				claimName := attr.FriendlyName
				if claimName == "" {
					claimName = attr.Name
				}
				if claimName == icn && len(values) == 1 {
					if attr.NameFormat != "urn:oasis:names:tc:SAML:2.0:attrname-format:basic" {
						ar.registry.logger.WithField("NameFormat", attr.NameFormat).Warnln("saml2 ignoring unsupported name format for identity claim name")
						continue
					}
					cvs = values[0]
					ok = true
				}
			}
		}
	}

	if !ok {
		return "", nil, errors.New("identity claim not found")
	}

	// Add extra external authority claims, for example SessionIndex.
	claims := make(map[string]interface{})
	for _, authnStatement := range assertion.AuthnStatements {
		ar.registry.logger.WithFields(logrus.Fields{
			"SessionNotOnOrAfter": authnStatement.SessionNotOnOrAfter,
			"SessionIndex":        authnStatement.SessionIndex,
		}).Debugln("saml2 authnStatement")
		if authnStatement.SessionIndex != "" {
			claims["SessionIndex"] = authnStatement.SessionIndex
			if authnStatement.SessionNotOnOrAfter != nil {
				if saml.TimeNow().After(*authnStatement.SessionNotOnOrAfter) {
					return "", nil, errors.New("session is expired")
				}
				claims["SessionNotOnOrAfter"] = authnStatement.SessionNotOnOrAfter
			}
		}
	}

	if assertion.Subject != nil {
		switch assertion.Subject.NameID.Format {
		case string(saml.TransientNameIDFormat):
			claims["TransientNameID"] = assertion.Subject.NameID.Value
		case string(saml.PersistentNameIDFormat):
			claims["PersistentNameID"] = assertion.Subject.NameID.Value
		case string(saml.UnspecifiedNameIDFormat):
			claims["UnspecifiedNameID"] = assertion.Subject.NameID.Value
		default:
			return "", nil, errors.New("nameid format must be transient")
		}
	} else {
		return "", nil, errors.New("subject not found")
	}

	// Convert claim value.
	whitelisted := false
	if ar.data.IdentityAliases != nil {
		if alias, ok := ar.data.IdentityAliases[cvs]; ok && alias != "" {
			cvs = alias
			whitelisted = true
		}
	}

	// Check whitelist.
	if ar.data.IdentityAliasRequired && !whitelisted {
		return "", nil, errors.New("identity claim has no alias")
	}

	return cvs, claims, nil
}

func (ar *saml2AuthorityRegistration) MakeRedirectAuthenticationRequestURL(state string) (*url.URL, map[string]interface{}, error) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	if !ar.ready {
		return nil, nil, errors.New("not ready")
	}

	authReq, err := ar.serviceProvider.MakeAuthenticationRequest(ar.serviceProvider.GetSSOBindingLocation(saml.HTTPRedirectBinding), saml.HTTPRedirectBinding, saml.HTTPPostBinding)
	if err != nil {
		return nil, nil, err
	}

	uri, err := authReq.Redirect(state, ar.serviceProvider)
	if err != nil {
		return nil, nil, err
	}

	return uri, map[string]interface{}{
		"rid": authReq.ID,
	}, nil
}

func (ar *saml2AuthorityRegistration) MakeRedirectEndSessionRequestURL(ref interface{}, state string) (*url.URL, map[string]interface{}, error) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	if !ar.ready {
		return nil, nil, errors.New("not ready")
	}

	logonRef := ref.(*string)
	var nameID string
	var nameIDFormat saml.NameIDFormat
	if logonRef != nil {
		logonRefParts := strings.SplitN(*logonRef, ":", 2)
		switch logonRefParts[0] {
		case "transient":
			nameIDFormat = saml.TransientNameIDFormat
		case "persistent":
			nameIDFormat = saml.PersistentNameIDFormat
		case "unspecified":
			nameIDFormat = saml.UnspecifiedNameIDFormat
		default:
			return nil, nil, fmt.Errorf("unsupported name id format prefix: %v", logonRefParts[0])
		}
		nameID = logonRefParts[1]
	}
	req, err := ar.serviceProvider.MakeLogoutRequest(ar.serviceProvider.GetSLOBindingLocation(saml.HTTPRedirectBinding), nameID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make redirect logout request: %w", err)
	}
	req.NameID.Format = string(nameIDFormat)

	lor := &samlext.LogoutRequest{
		LogoutRequest: req,
	}
	return lor.Redirect(state), nil, nil
}

func (ar *saml2AuthorityRegistration) MakeRedirectEndSessionResponseURL(rawReq interface{}, state string) (*url.URL, map[string]interface{}, error) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	if !ar.ready {
		return nil, nil, errors.New("not ready")
	}

	req, _ := rawReq.(*saml.LogoutRequest)
	if req == nil {
		return nil, nil, errors.New("invalid request data")
	}

	// NOTE(longsleep): This resonse currently always reports success.
	status := &saml.Status{
		StatusCode: saml.StatusCode{
			Value: saml.StatusSuccess,
		},
	}
	res, err := samlext.MakeLogoutResponse(ar.serviceProvider, req, status, saml.HTTPRedirectBinding)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make logout response: %w", err)
	}

	uri := res.Redirect(state)
	return uri, nil, nil
}

func (ar *saml2AuthorityRegistration) ParseStateResponse(req *http.Request, state string, extra map[string]interface{}) (interface{}, error) {
	requestID := extra["rid"].(string)

	return ar.serviceProvider.ParseResponse(req, []string{requestID})
}

func (ar *saml2AuthorityRegistration) ValidateIdpEndSessionRequest(req interface{}, state string) (bool, error) {
	slo := req.(*samlext.IdpLogoutRequest)

	if slo.Request == nil {
		return false, fmt.Errorf("request not set")
	}

	// NOTE(longsleep): We currently only support redirect binding (which uses a detached signature).
	if slo.Binding != saml.HTTPRedirectBinding {
		return false, fmt.Errorf("binding not supported")
	}

	// Only validate signature if signed.
	if slo.SigAlg == nil {
		return false, nil
	}

	ar.mutex.RLock()
	serviceProviderSigningCerts := ar.serviceProviderSigningCerts
	ready := ar.ready
	ar.mutex.RUnlock()

	if !ready {
		return false, errors.New("not ready")
	}

	if len(serviceProviderSigningCerts) == 0 {
		// No signing certs, cannot do anything.
		return false, nil
	}

	// Check if we are good.
	switch *slo.SigAlg {
	case dsig.RSASHA1SignatureMethod:
		ar.registry.logger.WithField("sig_alg", *slo.SigAlg).Warnln("saml2 insecure signature alg in idp logout request")
		if !ar.Authority().Insecure {
			return false, nil
		}

	default:
		// Let the rest pass, and decide later.
	}

	if len(slo.Signature) == 0 {
		return true, fmt.Errorf("signature data is empty")
	}

	// Get first certificate, and verify.
	if len(serviceProviderSigningCerts) > 1 {
		ar.registry.logger.Warnln("saml2 authority has multiple signing keys, using first")
	}
	pubKey := serviceProviderSigningCerts[0].PublicKey
	if verifyErr := slo.VerifySignature(pubKey); verifyErr != nil {
		return true, fmt.Errorf("signature verification failed: %w", verifyErr)
	}

	return true, nil
}

func (ar *saml2AuthorityRegistration) ValidateIdpEndSessionResponse(res interface{}, state string) (bool, error) {
	lor := res.(*samlext.IdpLogoutResponse)

	if lor.Response == nil {
		return false, fmt.Errorf("response not set")
	}

	// NOTE(longsleep): We currently only support redirect binding (which uses a detached signature).
	if lor.Binding != saml.HTTPRedirectBinding {
		return false, fmt.Errorf("binding not supported")
	}

	// Only validate signature if signed.
	if lor.SigAlg == nil {
		return false, nil
	}

	ar.mutex.RLock()
	serviceProviderSigningCerts := ar.serviceProviderSigningCerts
	ready := ar.ready
	ar.mutex.RUnlock()

	if !ready {
		return false, errors.New("not ready")
	}

	if len(serviceProviderSigningCerts) == 0 {
		// No signing certs, cannot do anything.
		return false, nil
	}

	// Check if we are good.
	switch *lor.SigAlg {
	case dsig.RSASHA1SignatureMethod:
		ar.registry.logger.WithField("sig_alg", *lor.SigAlg).Warnln("saml2 insecure signature alg in idp logout response")
		if !ar.Authority().Insecure {
			return false, nil
		}

	default:
		// Let the rest pass, and decide later.
	}

	if len(lor.Signature) == 0 {
		return true, fmt.Errorf("signature data is empty")
	}

	// Get first certificate, and verify.
	if len(serviceProviderSigningCerts) > 1 {
		ar.registry.logger.Warnln("saml2 authority has multiple signing keys, using first")
	}
	pubKey := serviceProviderSigningCerts[0].PublicKey
	if verifyErr := lor.VerifySignature(pubKey); verifyErr != nil {
		return true, fmt.Errorf("signature verification failed: %w", verifyErr)
	}

	return true, nil
}

func (ar *saml2AuthorityRegistration) Metadata() AuthorityMetadata {
	ar.mutex.RLock()
	sp := ar.serviceProvider
	ar.mutex.RUnlock()

	if sp == nil {
		return nil
	}

	metadata := sp.Metadata()

	// Set SLO to use redirect binding.
	metadata.SPSSODescriptors[0].SSODescriptor.SingleLogoutServices = []saml.Endpoint{
		{
			Binding:  saml.HTTPRedirectBinding,
			Location: sp.SloURL.String(),
		},
	}

	return metadata
}

func getCertsFromMetadata(md *saml.EntityDescriptor, use string) ([]*x509.Certificate, error) {
	var certStrs []string
	for _, idpSSODescriptor := range md.IDPSSODescriptors {
		for _, keyDescriptor := range idpSSODescriptor.KeyDescriptors {
			if keyDescriptor.Use == use {
				for _, cert := range keyDescriptor.KeyInfo.X509Data.X509Certificates {
					if cert.Data != "" {
						certStrs = append(certStrs, cert.Data)
					}
				}
			}
		}
	}

	// If there are no explicitly signing certs, just return the first non-empty cert we find.
	if len(certStrs) == 0 {
		for _, idpSSODescriptor := range md.IDPSSODescriptors {
			for _, keyDescriptor := range idpSSODescriptor.KeyDescriptors {
				if keyDescriptor.Use == "" {
					for _, cert := range keyDescriptor.KeyInfo.X509Data.X509Certificates {
						if cert.Data != "" {
							certStrs = append(certStrs, cert.Data)
						}
					}
					break
				}
			}
		}
	}

	var certs []*x509.Certificate

	for _, certStr := range certStrs {
		certStr = cleanWhitespaceRegexp.ReplaceAllString(certStr, "")
		certBytes, err := base64.StdEncoding.DecodeString(certStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		parsedCert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, parsedCert)
	}

	return certs, nil
}
