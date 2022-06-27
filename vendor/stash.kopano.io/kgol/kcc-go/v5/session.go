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
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// SessionAutorefreshInterval defines the interval when sessions are auto
	// refreshed automatically.
	SessionAutorefreshInterval = 4 * time.Minute
	// SessionExpirationGrace defines the duration after SessionAutorefreshInterval
	// when a session was not refreshed and can be considered non active.
	SessionExpirationGrace = 2 * time.Minute
)

// KCSessionID is the type for Kopano Core session IDs.
type KCSessionID uint64

func (sid KCSessionID) String() string {
	return strconv.FormatUint(uint64(sid), 10)
}

// KCNoSessionID define the value to use for KCSessionID when there is no session.
const KCNoSessionID KCSessionID = 0

// Session holds the data structures to keep a session open on the accociated
// Kopano server.
type Session struct {
	id         KCSessionID
	serverGUID string
	active     bool
	when       time.Time

	mutex     sync.RWMutex
	ctx       context.Context
	ctxCancel context.CancelFunc
	c         *KCC

	autoRefresh (chan bool)
}

// NewSession connects to the provided server with the provided parameters,
// creates a new Session which will be automatically refreshed until detroyed.
func NewSession(ctx context.Context, c *KCC, username, password string) (*Session, error) {
	if c == nil {
		c = NewKCC(nil)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := c.Logon(ctx, username, password, 0)
	if err != nil {
		return nil, fmt.Errorf("create session logon failed: %v", err)
	}
	if resp.Er != KCSuccess {
		return nil, fmt.Errorf("create session logon mapi error: %v", resp.Er)
	}
	if resp.SessionID == KCNoSessionID {
		return nil, fmt.Errorf("create session logon returned invalid session ID")
	}
	if resp.ServerGUID == "" {
		return nil, fmt.Errorf("create session logon return invalid server GUID")
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	s := &Session{
		id:         resp.SessionID,
		serverGUID: resp.ServerGUID,

		active: true,
		when:   time.Now(),

		ctx:       sessionCtx,
		ctxCancel: cancel,
		c:         c,
	}

	err = s.StartAutoRefresh()
	return s, err
}

// NewSSOSession connects to the provided server with the provided parameters,
// creates a new Session which will be automatically refreshed until detroyed.
func NewSSOSession(ctx context.Context, c *KCC, prefix SSOType, username string, input []byte, sessionID KCSessionID) (*Session, error) {
	if c == nil {
		c = NewKCC(nil)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := c.SSOLogon(ctx, prefix, username, input, sessionID, 0)
	if err != nil {
		return nil, fmt.Errorf("create session sso logon failed: %v", err)
	}
	if resp.Er != KCSuccess {
		return nil, fmt.Errorf("create session sso logon mapi error: %v", resp.Er)
	}
	if resp.SessionID == KCNoSessionID {
		return nil, fmt.Errorf("create session sso logon returned invalid session ID")
	}
	if resp.ServerGUID == "" {
		return nil, fmt.Errorf("create session sso logon return invalid server GUID")
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	s := &Session{
		id:         resp.SessionID,
		serverGUID: resp.ServerGUID,

		active: true,
		when:   time.Now(),

		ctx:       sessionCtx,
		ctxCancel: cancel,
		c:         c,
	}

	err = s.StartAutoRefresh()
	return s, err
}

// CreateSession creates a new Session without the server using the provided
// data.
func CreateSession(ctx context.Context, c *KCC, id KCSessionID, serverGUID string, active bool) (*Session, error) {
	if id == KCNoSessionID {
		return nil, fmt.Errorf("create session with invalid session ID")
	}
	if serverGUID == "" {
		return nil, fmt.Errorf("create session with invalid server GUID")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	s := &Session{
		id:         id,
		serverGUID: serverGUID,

		active: active,

		ctx:       sessionCtx,
		ctxCancel: cancel,
		c:         c,
	}

	if active {
		s.when = time.Now()
	}

	return s, nil
}

// Context returns the accociated Session's context.
func (s *Session) Context() context.Context {
	return s.ctx
}

// IsActive retruns true when the accociated Session is not destroyed, if the
// last refresh was successful and if the last activity is recent enough.
func (s *Session) IsActive() bool {
	s.mutex.RLock()
	active := s.active
	when := s.when
	s.mutex.RUnlock()

	return active && !when.Before(time.Now().Add(-(SessionAutorefreshInterval + SessionExpirationGrace)))
}

// ID returns the accociated Session's ID.
func (s *Session) ID() KCSessionID {
	return s.id
}

// Destroy logs off the accociated Session at the accociated Server and stops
// auto refreshing by canceling the accociated Session's Context. An error is
// retruned if the logoff request fails.
func (s *Session) Destroy(ctx context.Context, logoff bool) error {
	s.mutex.Lock()
	if !s.active {
		s.mutex.Unlock()
		return nil
	}
	s.active = false
	s.mutex.Unlock()
	s.ctxCancel()

	if logoff {
		resp, err := s.c.Logoff(ctx, s.id)
		if err != nil {
			return fmt.Errorf("logoff session logoff failed: %v", err)
		}

		if resp.Er != KCSuccess {
			return fmt.Errorf("logoff session logoff error: %v", resp.Er)
		}
	}

	return nil
}

func (s *Session) String() string {
	return fmt.Sprintf("Session(%s@%s)", s.id, s.serverGUID)
}

// Refresh triggers a server call to let the server know that the accociated
// session is still active.
func (s *Session) Refresh() error {
	s.mutex.RLock()
	active := s.active
	s.mutex.RUnlock()
	if !active {
		return nil
	}

	resp, err := s.c.ResolveUsername(s.ctx, "SYSTEM", s.id)
	if err != nil {
		return fmt.Errorf("refresh session resolveUsername failed: %v", err)
	}
	if resp.Er != KCSuccess {
		return fmt.Errorf("refresh session resolveUsername mapi error: %v", resp.Er)
	}
	s.mutex.Lock()
	s.when = time.Now()
	s.mutex.Unlock()

	return nil
}

// StartAutoRefresh enables auto refresh of the accociated session.
func (s *Session) StartAutoRefresh() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.autoRefresh != nil {
		close(s.autoRefresh)
	}
	s.autoRefresh = make(chan bool, 1)

	return s.runAutoRefresh(s.autoRefresh)
}

// StopAutoRefresh stops a running auto refresh of the accociated session.
func (s *Session) StopAutoRefresh() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.autoRefresh != nil {
		close(s.autoRefresh)
		s.autoRefresh = nil
	}

	return nil
}

func (s *Session) runAutoRefresh(stop chan bool) error {
	ctx := s.Context()
	ticker := time.NewTicker(SessionAutorefreshInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.StopAutoRefresh()
				return
			case <-ticker.C:
				err := s.Refresh()
				if err != nil {
					s.Destroy(ctx, err != KCERR_END_OF_SESSION)
					s.StopAutoRefresh()
				}
			case <-stop:
				return
			}
		}
	}()

	return nil
}
