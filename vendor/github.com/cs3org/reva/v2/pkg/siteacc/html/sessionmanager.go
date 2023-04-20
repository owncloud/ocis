// Copyright 2018-2020 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package html

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// SessionManager manages HTML sessions.
type SessionManager struct {
	conf *config.Configuration
	log  *zerolog.Logger

	sessions map[string]*Session

	sessionName string

	mutex sync.Mutex
}

func (mngr *SessionManager) initialize(name string, conf *config.Configuration, log *zerolog.Logger) error {
	if name == "" {
		return errors.Errorf("no session name provided")
	}
	mngr.sessionName = name

	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	mngr.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	mngr.log = log

	mngr.sessions = make(map[string]*Session, 100)

	return nil
}

// HandleRequest performs all session-related tasks during an HTML request. Always returns a valid session object.
func (mngr *SessionManager) HandleRequest(w http.ResponseWriter, r *http.Request) (*Session, error) {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	var session *Session
	var sessionErr error

	// Try to get the session ID from the request; if none has been set yet, a new one will be assigned
	cookie, err := r.Cookie(mngr.sessionName)
	if err == nil {
		session = mngr.findSession(cookie.Value)
		if session != nil {
			mngr.logSessionInfo(session, r, "existing session found")

			// Verify the request against the session: If it is invalid, set an error; if the session has expired, create a new one; if it has already passed its halftime, migrate to a new one
			if err := session.VerifyRequest(r, mngr.conf.Webserver.VerifyRemoteAddress); err == nil {
				if session.HasExpired() {
					// The session has expired, so a new one needs to be created
					session = nil

					mngr.logSessionInfo(session, r, "session expired")
				} else if session.HalftimePassed() {
					// The session has passed its halftime, so migrate it to a new one (makes hijacking session IDs harder)
					session, err = mngr.migrateSession(session, r)
					if err != nil {
						session = nil
						sessionErr = errors.Wrap(err, "unable to migrate session")
					}

					mngr.logSessionInfo(session, r, "session migrated")
				}
			} else {
				session = nil
				sessionErr = errors.Wrap(err, "invalid session")

				mngr.logSessionInfo(session, r, "session invalid (verify failed)")
			}
		}
	} else if err != http.ErrNoCookie {
		// The session cookie exists but seems to be invalid, so set an error
		session = nil
		sessionErr = errors.Wrap(err, "unable to get the session ID from the client")

		mngr.logSessionInfo(session, r, fmt.Sprintf("session cookie error: %v", err))
	}

	if session == nil {
		// No session found for the client, so create a new one; this will always succeed
		session = mngr.createSession(r)

		mngr.logSessionInfo(session, r, "assigned new session")
	}

	// Store the session ID on the client side
	session.Save(mngr.conf.Webserver.URL, w)

	return session, sessionErr
}

// PurgeSessions removes any expired sessions.
func (mngr *SessionManager) PurgeSessions() {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	var expiredSessions []string
	for id, session := range mngr.sessions {
		if session.HasExpired() {
			expiredSessions = append(expiredSessions, id)
		}
	}

	for _, id := range expiredSessions {
		delete(mngr.sessions, id)
	}
}

func (mngr *SessionManager) createSession(r *http.Request) *Session {
	session := NewSession(mngr.sessionName, time.Duration(mngr.conf.Webserver.SessionTimeout)*time.Second, r)
	mngr.sessions[session.ID] = session
	return session
}

func (mngr *SessionManager) findSession(id string) *Session {
	if session, ok := mngr.sessions[id]; ok {
		return session
	}
	return nil
}

func (mngr *SessionManager) migrateSession(session *Session, r *http.Request) (*Session, error) {
	sessionNew := mngr.createSession(r)

	// Carry over the old session information, thus preserving the existing session
	sessionNew.MigrationID = session.ID
	sessionNew.Data = session.Data

	if user := session.LoggedInUser(); user != nil {
		sessionNew.LoginUser(user.Account, user.Site)
	} else {
		sessionNew.LogoutUser()
	}

	// Delete the old session
	delete(mngr.sessions, session.ID)

	return sessionNew, nil
}

func (mngr *SessionManager) logSessionInfo(session *Session, r *http.Request, info string) {
	if mngr.conf.Webserver.LogSessions {
		if session != nil {
			mngr.log.Debug().Str("id", session.ID).Str("address", r.RemoteAddr).Str("path", r.URL.Path).Msg(info)
		} else {
			mngr.log.Debug().Str("address", r.RemoteAddr).Str("path", r.URL.Path).Msg(info)
		}
	}
}

// NewSessionManager creates a new session manager.
func NewSessionManager(name string, conf *config.Configuration, log *zerolog.Logger) (*SessionManager, error) {
	mngr := &SessionManager{}
	if err := mngr.initialize(name, conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the session manager")
	}
	return mngr, nil
}
