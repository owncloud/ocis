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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pquerna/cachecontrol"
)

// Basic HTTP related global settings.
var (
	DefaultHTTPClient       *http.Client
	DefaultHTTPHeader       http.Header
	DefaultMaxJSONFetchSize int64 = 5 * 1024 * 1024 // 5 MiB
	DefaultJSONFetchExpiry        = time.Minute * 1
	DefaultJSONFetchRetry         = time.Second * 3
)

func fetchJSON(ctx context.Context, u *url.URL, dst interface{}, client *http.Client, header http.Header) (time.Duration, error) {
	if client == nil {
		client = DefaultHTTPClient
		if client == nil {
			client = http.DefaultClient
		}
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return DefaultJSONFetchRetry, fmt.Errorf("failed create fetch JSON request: %v", err)
	}
	if header == nil {
		header = DefaultHTTPHeader
	}
	if header != nil {
		for h, values := range header {
			for _, v := range values {
				req.Header.Add(h, v)
			}
		}
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return DefaultJSONFetchRetry, fmt.Errorf("failed to fetch JSON: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return DefaultJSONFetchRetry, fmt.Errorf("failed to fetch JSON (status: %d)", res.StatusCode)
	}

	_, expires, _ := cachecontrol.CachableResponse(req, res, cachecontrol.Options{})
	err = json.NewDecoder(io.LimitReader(res.Body, DefaultMaxJSONFetchSize)).Decode(dst)
	if err != nil {
		return DefaultJSONFetchRetry, fmt.Errorf("failed to fetch JSON: %v", err)
	}

	expirationDuration := expires.Sub(time.Now())
	if expirationDuration < DefaultJSONFetchRetry {
		if err == nil {
			expirationDuration = DefaultJSONFetchExpiry
		} else {
			expirationDuration = DefaultJSONFetchRetry
		}
	}

	return expirationDuration, err
}
