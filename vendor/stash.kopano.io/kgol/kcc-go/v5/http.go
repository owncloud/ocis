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
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/publicsuffix"
)

// DefaultHTTPClient is the default Client as used by KCC for HTTP SOAP requests.
var DefaultHTTPClient *http.Client

var (
	defaultHTTPTimeoutSeconds         int64 = 10
	defaultHTTPMaxIdleConns                 = 100
	defaultHTTPMaxIdleConnsPerHost          = 100
	defaultHTTPIdleConnTimeoutSeconds int64 = 90
	defaultHTTPDialTimeoutSeconds     int64 = 30
	defaultHTTPKeepAliveSeconds       int64 = 120
	defaultHTTPDualStack                    = true
)

var defaultHTTPTransport *http.Transport

func init() {
	debug = os.Getenv("KCC_GO_DEBUG") != ""

	if s := os.Getenv("KCC_GO_HTTP_TIMEOUT"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			defaultHTTPTimeoutSeconds = n
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_MAX_IDLE_CONNS"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 0); err == nil {
			defaultHTTPMaxIdleConns = int(n)
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_MAX_IDLE_CONNS_PER_HOST"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 0); err == nil {
			defaultHTTPMaxIdleConnsPerHost = int(n)
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_IDLE_CONN_TIMEOUT"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			defaultHTTPIdleConnTimeoutSeconds = n
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_DIAL_TIMEOUT"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			defaultHTTPDialTimeoutSeconds = n
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_KEEPALIVE"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			defaultHTTPKeepAliveSeconds = n
		}
	}
	if s := os.Getenv("KCC_GO_HTTP_DUALSTACK"); s != "" {
		switch s {
		case "off", "false", "no":
			defaultHTTPDualStack = false
		case "on", "true", "yes":
			defaultHTTPDualStack = true
		}
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic(err)
	}

	defaultHTTPTransport = http.DefaultTransport.(*http.Transport).Clone()
	defaultHTTPTransport.DialContext = (&net.Dialer{
		Timeout:   time.Duration(defaultHTTPDialTimeoutSeconds) * time.Second,
		KeepAlive: time.Duration(defaultHTTPKeepAliveSeconds) * time.Second,
		DualStack: defaultHTTPDualStack,
	}).DialContext
	defaultHTTPTransport.MaxIdleConns = defaultHTTPMaxIdleConns
	defaultHTTPTransport.MaxIdleConnsPerHost = defaultHTTPMaxIdleConnsPerHost
	defaultHTTPTransport.IdleConnTimeout = time.Duration(defaultHTTPIdleConnTimeoutSeconds) * time.Second

	DefaultHTTPClient = &http.Client{
		Jar:       jar,
		Timeout:   time.Duration(defaultHTTPTimeoutSeconds) * time.Second,
		Transport: defaultHTTPTransport,
	}

	if debug {
		fmt.Printf("HTTP client: %+v\n", DefaultHTTPClient)
		fmt.Printf("HTTP client transport: %+v\n", defaultHTTPTransport)
	}
}
