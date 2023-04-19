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
	"net"
	"net/http"
)

// IsRequestFromTrustedSource checks if the provided requests remote address is
// one either one of the provided ips or in one of the provided networks.
func IsRequestFromTrustedSource(req *http.Request, ips []*net.IP, nets []*net.IPNet) (bool, error) {
	ipString, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return false, err
	}

	ip := net.ParseIP(ipString)

	for _, checkIP := range ips {
		if checkIP.Equal(ip) {
			return true, nil
		}
	}

	for _, checkNet := range nets {
		if checkNet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}
