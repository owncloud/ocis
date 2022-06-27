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
)

func (p *Provider) setBrowserStateCookie(rw http.ResponseWriter, value string) error {
	cookie := http.Cookie{
		Name:  p.browserStateCookieName,
		Value: value,

		Path:     p.browserStateCookiePath,
		Secure:   true,
		HttpOnly: false, // This Cookie is intended to be read by Javascript.
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(rw, &cookie)

	return nil
}

func (p *Provider) removeBrowserStateCookie(rw http.ResponseWriter) error {
	cookie := http.Cookie{
		Name: p.browserStateCookieName,

		Path:     p.browserStateCookiePath,
		Secure:   true,
		HttpOnly: false, // This Cookie is intended to be read by Javascript.
		SameSite: http.SameSiteNoneMode,

		Expires: farPastExpiryTime,
	}
	http.SetCookie(rw, &cookie)

	return nil
}

func (p *Provider) setSessionCookie(rw http.ResponseWriter, value string) error {
	cookie := http.Cookie{
		Name:  p.sessionCookieName,
		Value: value,

		Path:     p.sessionCookiePath,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(rw, &cookie)

	return nil
}

func (p *Provider) getSessionCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.sessionCookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (p *Provider) removeSessionCookie(rw http.ResponseWriter) error {
	cookie := http.Cookie{
		Name: p.sessionCookieName,

		Path:     p.sessionCookiePath,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,

		Expires: farPastExpiryTime,
	}
	http.SetCookie(rw, &cookie)

	return nil
}
