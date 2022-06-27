/*
 * Copyright 2019 Kopano and its licensors
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
)

// SetX509KeyPairToTLSConfig reads and parses a public/private key pair from a
// pair of files and adds the resulting certificate to the provided TLs config.
// If the provided TLS config is nil, a new empty one will be created and
// returned.
func SetX509KeyPairToTLSConfig(certFile, keyFile string, config *tls.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return config, err
	}

	if config == nil {
		config = &tls.Config{}
	}
	config.Certificates = []tls.Certificate{cert}

	return config, nil
}
