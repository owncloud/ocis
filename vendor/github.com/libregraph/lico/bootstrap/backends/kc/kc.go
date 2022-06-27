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

package bskc

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	kcc "stash.kopano.io/kgol/kcc-go/v5"

	"github.com/libregraph/lico/bootstrap"
	"github.com/libregraph/lico/identifier"
	"github.com/libregraph/lico/identifier/backends/kc"
	"github.com/libregraph/lico/identity"
	"github.com/libregraph/lico/identity/managers"
	"github.com/libregraph/lico/utils"
	"github.com/libregraph/lico/version"
)

// Identity managers.
const (
	identityManagerName = "kc"
)

func Register() error {
	return bootstrap.RegisterIdentityManager(identityManagerName, NewIdentityManager)
}

func MustRegister() {
	if err := Register(); err != nil {
		panic(err)
	}
}

func NewIdentityManager(bs bootstrap.Bootstrap) (identity.Manager, error) {
	config := bs.Config()

	logger := config.Config.Logger

	if config.AuthorizationEndpointURI.String() != "" {
		return nil, fmt.Errorf("kc backend is incompatible with authorization-endpoint-uri parameter")
	}
	config.AuthorizationEndpointURI.Path = bs.MakeURIPath(bootstrap.APITypeSignin, "/identifier/_/authorize")

	if config.EndSessionEndpointURI.String() != "" {
		return nil, fmt.Errorf("kc backend is incompatible with endsession-endpoint-uri parameter")
	}
	config.EndSessionEndpointURI.Path = bs.MakeURIPath(bootstrap.APITypeSignin, "/identifier/_/endsession")

	if config.SignInFormURI.EscapedPath() == "" {
		config.SignInFormURI.Path = bs.MakeURIPath(bootstrap.APITypeSignin, "/identifier")
	}

	if config.SignedOutURI.EscapedPath() == "" {
		config.SignedOutURI.Path = bs.MakeURIPath(bootstrap.APITypeSignin, "/goodbye")
	}

	useGlobalSession := false
	globalSessionUsername := os.Getenv("KOPANO_SERVER_USERNAME")
	globalSessionPassword := os.Getenv("KOPANO_SERVER_PASSWORD")
	globalSessionClientCertificate := os.Getenv("KOPANO_CLIENT_CERTIFICATE")
	globalSessionClientPrivateKey := os.Getenv("KOPANO_CLIENT_PRIVATE_KEY")
	if globalSessionUsername == "" && (globalSessionClientCertificate != "" && globalSessionClientPrivateKey != "") {
		globalSessionUsername = "SYSTEM"
	}
	if globalSessionUsername != "" {
		useGlobalSession = true
	}

	var sessionTimeoutSeconds uint64 = 300 // 5 Minutes is the default.
	if sessionTimeoutSecondsString := os.Getenv("KOPANO_SERVER_SESSION_TIMEOUT"); sessionTimeoutSecondsString != "" {
		var sessionTimeoutSecondsErr error
		sessionTimeoutSeconds, sessionTimeoutSecondsErr = strconv.ParseUint(sessionTimeoutSecondsString, 10, 64)
		if sessionTimeoutSecondsErr != nil {
			return nil, fmt.Errorf("invalid KOPANO_SERVER_SESSION_TIMEOUT value: %v", sessionTimeoutSecondsErr)
		}
	}
	if !useGlobalSession && config.AccessTokenDurationSeconds+60 > sessionTimeoutSeconds {
		config.AccessTokenDurationSeconds = sessionTimeoutSeconds - 60
		config.Config.Logger.Warnf("limiting access token duration to %d seconds because of lower KOPANO_SERVER_SESSION_TIMEOUT", config.AccessTokenDurationSeconds)
	}
	// Update kcc defaults to our values.
	kcc.SessionAutorefreshInterval = time.Duration(sessionTimeoutSeconds-60) * time.Second
	kcc.SessionExpirationGrace = 2 * time.Minute // 2 Minutes grace until cleanup.

	// Setup kcc default HTTP client with our values.
	tlsClientConfig := config.TLSClientConfig
	kcc.DefaultHTTPClient = &http.Client{
		Timeout:   utils.DefaultHTTPClient.Timeout,
		Transport: utils.HTTPTransportWithTLSClientConfig(tlsClientConfig),
	}

	kopanoStorageServerClient := kcc.NewKCC(nil)
	if err := kopanoStorageServerClient.SetClientApp("konnect", version.Version); err != nil {
		return nil, fmt.Errorf("failed to initialize kc client: %v", err)
	}
	if useGlobalSession && (globalSessionClientCertificate != "" || globalSessionClientPrivateKey != "") {
		if globalSessionClientCertificate == "" {
			return nil, fmt.Errorf("invalid or empty KOPANO_CLIENT_CERTIFICATE value")
		}
		if globalSessionClientPrivateKey == "" {
			return nil, fmt.Errorf("invalid or empty KOPANO_CLIENT_PRIVATE_KEY value")
		}
		if tlsClientConfig == nil {
			return nil, fmt.Errorf("no TLS client config - this should not happen")
		}
		if _, err := kcc.SetX509KeyPairToTLSConfig(globalSessionClientCertificate, globalSessionClientPrivateKey, tlsClientConfig); err != nil {
			return nil, fmt.Errorf("failed to load/set kc client x509 certificate: %v", err)
		}
		logger.Infoln("kc server identifier backend initialized client for TLS authentication")
	}

	identifierBackend, identifierErr := kc.NewKCIdentifierBackend(
		config.Config,
		kopanoStorageServerClient,
		useGlobalSession,
		globalSessionUsername,
		globalSessionPassword,
	)
	if identifierErr != nil {
		return nil, fmt.Errorf("failed to create identifier backend: %v", identifierErr)
	}

	fullAuthorizationEndpointURL := bootstrap.WithSchemeAndHost(config.AuthorizationEndpointURI, config.IssuerIdentifierURI)
	fullSignInFormURL := bootstrap.WithSchemeAndHost(config.SignInFormURI, config.IssuerIdentifierURI)
	fullSignedOutEndpointURL := bootstrap.WithSchemeAndHost(config.SignedOutURI, config.IssuerIdentifierURI)

	activeIdentifier, err := identifier.NewIdentifier(&identifier.Config{
		Config: config.Config,

		BaseURI:         config.IssuerIdentifierURI,
		PathPrefix:      bs.MakeURIPath(bootstrap.APITypeSignin, ""),
		StaticFolder:    config.IdentifierClientPath,
		LogonCookieName: "__Secure-KKT", // Kopano-Konnect-Token
		ScopesConf:      config.IdentifierScopesConf,
		WebAppDisabled:  config.IdentifierClientDisabled,

		AuthorizationEndpointURI: fullAuthorizationEndpointURL,
		SignedOutEndpointURI:     fullSignedOutEndpointURL,

		DefaultBannerLogo:       config.IdentifierDefaultBannerLogo,
		DefaultSignInPageText:   config.IdentifierDefaultSignInPageText,
		DefaultUsernameHintText: config.IdentifierDefaultUsernameHintText,
		UILocales:               config.IdentifierUILocales,

		Backend: identifierBackend,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create identifier: %v", err)
	}
	err = activeIdentifier.SetKey(config.EncryptionSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid --encryption-secret parameter value for identifier: %v", err)
	}

	identityManagerConfig := &identity.Config{
		SignInFormURI: fullSignInFormURL,
		SignedOutURI:  fullSignedOutEndpointURL,

		Logger: logger,

		ScopesSupported: config.Config.AllowedScopes,
	}

	identifierIdentityManager := managers.NewIdentifierIdentityManager(identityManagerConfig, activeIdentifier)
	logger.Infoln("using identifier backed identity manager")

	return identifierIdentityManager, nil
}
