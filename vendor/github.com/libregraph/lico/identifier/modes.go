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

package identifier

const (
	// ModeLogonUsernameEmptyPasswordCookie is the logon mode which requires a
	// username which matches the currently signed in user in the cookie and an
	// empty password.
	ModeLogonUsernameEmptyPasswordCookie = "0"
	// ModeLogonUsernamePassword is the logon mode which requires a username
	// and a password.
	ModeLogonUsernamePassword = "1"
)

const (
	// MustBeSignedIn is a authorize mode which tells the authorization code,
	// that it is expected to have a signed in user and everything else should
	// be treated as error.
	MustBeSignedIn = "must"
)

const (
	// StateModeEndSession is a state mode which selects end session specific
	// actions when processing state requests.
	StateModeEndSession = "0"
)
