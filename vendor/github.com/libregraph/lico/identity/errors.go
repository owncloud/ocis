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

package identity

import (
	"net/url"
)

// IsHandledError is an error which tells that the backend has handled
// the request and all further handling should stop
type IsHandledError struct {
}

// Error implements the error interface.
func (err *IsHandledError) Error() string {
	return "is_handled"
}

// RedirectError is an error which backends can return if a
// redirection is required.
type RedirectError struct {
	id          string
	redirectURI *url.URL
}

// NewRedirectError creates a new corresponding error with the
// provided id and redirect URL.
func NewRedirectError(id string, redirectURI *url.URL) *RedirectError {
	return &RedirectError{
		id:          id,
		redirectURI: redirectURI,
	}
}

// Error implements the error interface.
func (err *RedirectError) Error() string {
	return err.id
}

// RedirectURI returns the redirection URL of the accociated error.
func (err *RedirectError) RedirectURI() *url.URL {
	return err.redirectURI
}

// LoginRequiredError which backends can return to indicate that sign-in is
// required.
type LoginRequiredError struct {
	id        string
	signInURI *url.URL
}

// NewLoginRequiredError creates a new corresponding error with the provided id.
func NewLoginRequiredError(id string, signInURI *url.URL) *LoginRequiredError {
	return &LoginRequiredError{
		id:        id,
		signInURI: signInURI,
	}
}

// Error implements the error interface.
func (err *LoginRequiredError) Error() string {
	return err.id
}

// SignInURI returns the sign-in URL of the accociated error.
func (err *LoginRequiredError) SignInURI() *url.URL {
	return err.signInURI
}
