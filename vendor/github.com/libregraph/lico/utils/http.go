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

package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"github.com/libregraph/lico/version"
)

const (
	defaultHTTPTimeout               = 30 * time.Second
	defaultHTTPKeepAlive             = 30 * time.Second
	defaultHTTPMaxIdleConns          = 100
	defaultHTTPIdleConnTimeout       = 90 * time.Second
	defaultHTTPTLSHandshakeTimeout   = 10 * time.Second
	defaultHTTPExpectContinueTimeout = 1 * time.Second
)

// DefaultHTTPUserAgent is the User-Agent Header which should be used when
// making HTTP requests.
var DefaultHTTPUserAgent = "LibreGraph-Connect/" + version.Version

// HTTPTransportWithTLSClientConfig creates a new http.Transport with sane
// default settings using the provided tls.Config.
func HTTPTransportWithTLSClientConfig(tlsClientConfig *tls.Config) *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultHTTPTimeout,
			KeepAlive: defaultHTTPKeepAlive,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          defaultHTTPMaxIdleConns,
		IdleConnTimeout:       defaultHTTPIdleConnTimeout,
		TLSHandshakeTimeout:   defaultHTTPTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultHTTPExpectContinueTimeout,
	}
	if tlsClientConfig != nil {
		transport.TLSClientConfig = tlsClientConfig
		err := http2.ConfigureTransport(transport)
		if err != nil {
			panic(err)
		}
	}

	return transport
}

// DefaultTLSConfig returns a new tls.Config.
func DefaultTLSConfig() *tls.Config {
	return &tls.Config{
		ClientSessionCache: tls.NewLRUClientSessionCache(0),
	}
}

// InsecureSkipVerifyTLSConfig returns a new tls.Config which does skip TLS verification.
func InsecureSkipVerifyTLSConfig() *tls.Config {
	config := DefaultTLSConfig()
	config.InsecureSkipVerify = true

	return config
}

// DefaultHTTPClient is a http.Client with a timeout set.
var DefaultHTTPClient = &http.Client{
	Timeout:   defaultHTTPTimeout,
	Transport: HTTPTransportWithTLSClientConfig(DefaultTLSConfig()),
}

// InsecureHTTPClient is a http.Client with a timeout set and with TLS
// verification disabled.
var InsecureHTTPClient = &http.Client{
	Timeout:   defaultHTTPTimeout,
	Transport: HTTPTransportWithTLSClientConfig(InsecureSkipVerifyTLSConfig()),
}
