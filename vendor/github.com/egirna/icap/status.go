// Copyright 2011 Andy Balholm. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ICAP status codes.

package icap

import (
	"net/http"
)

var statusText = map[int]string{
	100: "Continue",
	204: "No Modifications",
	400: "Bad Request",
	404: "ICAP Service Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	500: "Server Error",
	501: "Method Not Implemented",
	502: "Bad Gateway",
	503: "Service Overloaded",
	505: "ICAP Version Not Supported",
}

// StatusText returns a text for the ICAP status code. It returns the empty string if the code is unknown.
func StatusText(code int) string {
	text, ok := statusText[code]
	if ok {
		return text
	}
	return http.StatusText(code)
}
