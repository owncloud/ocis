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

package provider

import (
	"net/http"
	"time"

	"github.com/libregraph/lico/config"
)

// Config defines a Provider's configuration settings.
type Config struct {
	Config *config.Config

	IssuerIdentifier       string
	WellKnownPath          string
	JwksPath               string
	AuthorizationPath      string
	TokenPath              string
	UserInfoPath           string
	EndSessionPath         string
	CheckSessionIframePath string
	RegistrationPath       string

	BrowserStateCookiePath     string
	BrowserStateCookieName     string
	BrowserStateCookieSameSite http.SameSite

	SessionCookiePath     string
	SessionCookieName     string
	SessionCookieSameSite http.SameSite

	AccessTokenDuration  time.Duration
	IDTokenDuration      time.Duration
	RefreshTokenDuration time.Duration
}
