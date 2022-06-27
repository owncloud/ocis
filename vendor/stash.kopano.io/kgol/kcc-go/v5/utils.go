/*
 * Copyright 2017-2019 Kopano and its licensors
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

package kcc

import (
	"crypto/tls"
	"net/http"
)

// useX509KeyPair enables TLS client authentication for the provided
// http.transport using the provided certificate and private key. The files must
// contain PEM encoded data.
func useX509KeyPair(transport *http.Transport, certFile, keyFile string) error {
	config := transport.TLSClientConfig
	if config == nil {
		config = &tls.Config{}
	} else {
		config = config.Clone()
	}
	if _, err := SetX509KeyPairToTLSConfig(certFile, keyFile, config); err != nil {
		return err
	}

	transport.TLSClientConfig = config
	return nil
}
