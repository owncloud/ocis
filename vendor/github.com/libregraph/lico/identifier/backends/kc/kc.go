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

package kc

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"stash.kopano.io/kgol/kcc-go/v5"
	"stash.kopano.io/kgol/oidc-go"

	konnect "github.com/libregraph/lico"
	"github.com/libregraph/lico/config"
	"github.com/libregraph/lico/identifier/backends"
	"github.com/libregraph/lico/identifier/meta/scopes"
	"github.com/libregraph/lico/identity"
	"github.com/libregraph/lico/managers"
)

const (
	kcSessionMaxRetries = 3
	kcSessionRetryDelay = 50 * time.Millisecond

	kcIdentifierBackendName = "identifier-kc"
)

var kcSupportedScopes = []string{
	oidc.ScopeProfile,
	oidc.ScopeEmail,
	konnect.ScopeNumericID,
	konnect.ScopeUniqueUserID,
	konnect.ScopeRawSubject,
	ScopeKopanoGC,
}

var scopesMeta = &scopes.Scopes{
	Definitions: map[string]*scopes.Definition{
		ScopeKopanoGC: &scopes.Definition{
			Description: "Read and write your Kopano Groupware data",
		},
	},
}

// Property mappings for Kopano Server user meta data.
var (
	KCServerDefaultFamilyNameProperty = kcc.PR_SURNAME_A
	KCServerDefaultGivenNameProperty  = kcc.PR_GIVEN_NAME_A
)

// KCIdentifierBackend is a backend for the Identifier which connects to
// Kopano Core via kcc-go.
type KCIdentifierBackend struct {
	ctx context.Context

	c        *kcc.KCC
	username string
	password string

	globalSession      *kcc.Session
	globalSessionMutex sync.RWMutex
	useGlobalSession   bool // Wether konnect maintains a single global session
	sessions           cmap.ConcurrentMap

	identityManager identity.Manager
	oidcProvider    konnect.AccessTokenProvider

	logger logrus.FieldLogger
}

type kcUser struct {
	sub string

	user *kcc.User
}

func newKCUser(user *kcc.User, abeid kcc.ABEID) (*kcUser, error) {
	// NOTE(longsleep): KCC uses EntryIDs  which needs parsing to get the
	// unique user identifier. We take the users entryID and extract the
	// exID field and use its base64URL encoded value as subject.

	// Parse user EntryID.
	if abeid == nil {
		var err error
		if abeid, err = kcc.NewABEIDFromBase64([]byte(user.UserEntryID)); err != nil {
			return nil, fmt.Errorf("failed to parse user entry id: %v", err)
		}
	}

	// Fetch external ID and as subject. Since its binary, encode it base64URL.
	sub := exIDToString(abeid.ExID())

	return &kcUser{
		sub: sub,

		user: user,
	}, nil
}

func (u *kcUser) Subject() string {
	return u.sub
}

func (u *kcUser) Email() string {
	return u.user.MailAddress
}

func (u *kcUser) EmailVerified() bool {
	return true
}

func (u *kcUser) Name() string {
	return u.user.FullName
}

func (u *kcUser) FamilyName() string {
	var n string
	if u.user.Props != nil {
		n, _ = u.user.Props.Get(KCServerDefaultFamilyNameProperty)
	} else {
		n = u.splitFullName()[1]
	}
	return n
}

func (u *kcUser) GivenName() string {
	var n string
	if u.user.Props != nil {
		n, _ = u.user.Props.Get(KCServerDefaultGivenNameProperty)
	} else {
		n = u.splitFullName()[0]
	}
	return n
}

func (u *kcUser) ID() int64 {
	return int64(u.user.ID)
}

func (u *kcUser) UniqueID() string {
	// Provide our sub as unique ID since it is the exID part of the
	// users ABEID which is unique in whatever user backend the server uses.
	return u.sub
}

func (u *kcUser) Username() string {
	return u.user.Username
}

// BackendClaims return backend specific additional user claims.
func (u *kcUser) BackendClaims() map[string]interface{} {
	claims := make(map[string]interface{})
	claims[konnect.IdentifiedUserIDClaim] = u.user.UserEntryID
	claims[konnect.IdentifiedUserClaim] = u.sub

	return claims
}

func (u *kcUser) BackendScopes() []string {
	return nil
}

func (u *kcUser) RequiredScopes() []string {
	return nil
}

func (u *kcUser) splitFullName() [2]string {
	// TODO(longsleep): Cache this, instead of doing every time.
	parts := strings.SplitN(u.user.FullName, " ", 2)
	if len(parts) == 2 {
		return [2]string{parts[0], parts[1]}
	}
	return [2]string{"", ""}
}

// NewKCIdentifierBackend creates a new KCIdentifierBackend with the provided
// parameters.
func NewKCIdentifierBackend(c *config.Config, client *kcc.KCC, useGlobalSession bool, username string, password string) (*KCIdentifierBackend, error) {
	b := &KCIdentifierBackend{
		c: client,

		logger: c.Logger,

		sessions: cmap.New(),
	}

	// Store credentials if required.
	if useGlobalSession {
		b.username = username
		b.password = password
		b.useGlobalSession = useGlobalSession
	}

	b.logger.WithField("client", b.c.String()).Infoln("kc server identifier backend connection set up")

	return b, nil
}

// RegisterManagers registers the provided managers,
func (b *KCIdentifierBackend) RegisterManagers(mgrs *managers.Managers) error {
	b.identityManager = mgrs.Must("identity").(identity.Manager)
	if oidcManager, ok := mgrs.Get("oidc"); ok {
		b.oidcProvider = oidcManager.(konnect.AccessTokenProvider)
	}

	return nil
}

// RunWithContext implements the Backend interface. KCIdentifierBackends keep
// a session to the accociated Kopano Core client. This session is auto renewed
// and auto rerestablished and is bound to the provided Context.
func (b *KCIdentifierBackend) RunWithContext(ctx context.Context) error {
	b.ctx = ctx

	// Helper to keep dedicated session running.
	if b.useGlobalSession {
		b.logger.WithField("username", b.username).Infoln("kc server identifier global session enabled")

		go func() {
			retry := time.NewTimer(5 * time.Second)
			retry.Stop()
			refreshCh := make(chan bool, 1)
			for {
				b.setGlobalSession(nil)
				session, sessionErr := kcc.NewSession(ctx, b.c, b.username, b.password)
				if sessionErr != nil {
					b.logger.WithError(sessionErr).Errorln("failed to create kc server global session")
					retry.Reset(5 * time.Second)
				} else {
					b.logger.Debugf("kc server identifier global session established: %v", session)
					b.setGlobalSession(session)
					go func() {
						<-session.Context().Done()
						b.logger.Debugf("kc server identifier global session has ended: %v", session)
						refreshCh <- true
					}()
				}

				select {
				case <-refreshCh:
					// will retry instantly.
				case <-retry.C:
					// will retry instantly.
				case <-ctx.Done():
					// exit.
					return
				}
			}
		}()
	}

	// Helper to clean out old session data from memory.
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				expired := make([]string, 0)
				for entry := range b.sessions.IterBuffered() {
					session := entry.Val.(*kcc.Session)
					if session == nil || !session.IsActive() {
						expired = append(expired, entry.Key)
					}
				}
				for _, ref := range expired {
					b.sessions.Remove(ref)
				}
			case <-ctx.Done():
				// exit.
				return
			}
		}
	}()

	return nil
}

// Logon implements the Backend interface, enabling Logon with user name and
// password as provided. Requests are bound to the provided context.
func (b *KCIdentifierBackend) Logon(ctx context.Context, audience, username, password string) (bool, *string, *string, backends.UserFromBackend, error) {
	var logonFlags kcc.KCFlag
	logonFlags |= kcc.KOPANO_LOGON_NO_UID_AUTH
	if b.useGlobalSession {
		// Do not register session when loggon when that session is not used.
		logonFlags |= kcc.KOPANO_LOGON_NO_REGISTER_SESSION
	}

	response, err := b.c.Logon(ctx, username, password, logonFlags)
	if err != nil {
		return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon error: %v", err)
	}

	switch response.Er {
	case kcc.KCSuccess:
		// Session
		var session *kcc.Session
		if !b.useGlobalSession {
			// Register the session, we just created for later use. It will
			// eventually expire if not refreshed by other means.
			if response.SessionID != kcc.KCNoSessionID {
				session, err = kcc.CreateSession(b.ctx, b.c, response.SessionID, response.ServerGUID, true)
				if err != nil {
					return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon session error: %v", err)
				}
			} else {
				return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon missing session")
			}
		}

		// Resolve user details.
		// TODO(longsleep): Avoid extra resolve when logon response already
		// includes the required data (TODO in core).
		resolve, err := b.resolveUsername(ctx, username, session)
		if err != nil {
			return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon resolve error: %v", err)
		}

		user, err := newKCUser(&kcc.User{
			ID:          resolve.ID,
			Username:    username,
			UserEntryID: resolve.UserEntryID,
		}, nil)
		if err != nil {
			return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon resolve format error: %v", err)
		}

		// Use the users subject as user id.
		userID := user.Subject()

		sessionRef := identity.GetSessionRef(b.Name(), audience, userID)
		b.sessions.Set(*sessionRef, session)
		b.logger.WithFields(logrus.Fields{
			"session":  session,
			"ref":      *sessionRef,
			"username": user.Username(),
			"id":       userID,
			"abeid":    user.user.UserEntryID,
		}).Debugln("kc identifier backend logon")

		return true, &userID, sessionRef, user, nil

	case kcc.KCERR_LOGON_FAILED:
		return false, nil, nil, nil, nil
	}

	return false, nil, nil, nil, fmt.Errorf("kc identifier backend logon failed: %v", response.Er)
}

// ResolveUserByUsername implements the Beckend interface, providing lookup for user by
// providing the username. Requests are bound to the provided context.
func (b *KCIdentifierBackend) ResolveUserByUsername(ctx context.Context, username string) (backends.UserFromBackend, error) {
	// NOTE(longsleep): No session support here. This means resolving of users
	// by their user name always needs a global session.
	response, err := b.resolveUsername(ctx, username, nil)
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend resolve user error: %v", err)
	}

	switch response.Er {
	case kcc.KCSuccess:
		// success.

		return newKCUser(&kcc.User{
			ID:          response.ID,
			Username:    username,
			UserEntryID: response.UserEntryID,
		}, nil)

	case kcc.KCERR_NOT_FOUND:
		return nil, nil
	}

	return nil, fmt.Errorf("kc identifier backend get user failed: %v", response.Er)
}

// GetUser implements the Backend interface, providing user meta data retrieval
// for the user specified by the userID. Requests are bound to the provided
// context.
func (b *KCIdentifierBackend) GetUser(ctx context.Context, userEntryID string, sessionRef *string, requestedScopes map[string]bool) (backends.UserFromBackend, error) {
	abeid, err := kcc.NewABEIDFromBase64([]byte(userEntryID))
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend resolve session with invalid entry id: %v", err)
	}

	session, err := b.getSessionForUser(ctx, userEntryID, sessionRef, true, true, false)
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend resolve session error: %v", err)
	}

	response, err := b.getUser(ctx, userEntryID, session)
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend get user error: %v", err)
	}

	switch response.Er {
	case kcc.KCSuccess:
		// success.
		responseAbeid, _ := kcc.NewABEIDFromBase64([]byte(response.User.UserEntryID))
		if !kcc.ABEIDEqual(abeid, responseAbeid) {
			return nil, fmt.Errorf("kc identifier backend get user returned wrong user")
		}

		user, err := newKCUser(response.User, responseAbeid)
		if err != nil {
			return nil, fmt.Errorf("kc identifier backend get user failed to process user: %v", err)
		}

		return user, nil

	case kcc.KCERR_NOT_FOUND:
		return nil, nil
	}

	return nil, fmt.Errorf("kc identifier backend get user failed: %v", response.Er)
}

// RefreshSession implements the Backend interface providing refresh to KC session.
func (b *KCIdentifierBackend) RefreshSession(ctx context.Context, userID string, sessionRef *string, claims map[string]interface{}) error {
	userEntryID, _ := claims[konnect.IdentifiedUserIDClaim].(string)
	if userEntryID == "" {
		return fmt.Errorf("kc identifier backend refresh session missing claim")
	}

	_, err := kcc.NewABEIDFromBase64([]byte(userEntryID))
	if err != nil {
		return fmt.Errorf("kc identifier backend refresh session with invalid entry id: %v", err)
	}

	_, err = b.getSessionForUser(ctx, userEntryID, sessionRef, true, true, false)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"id":    userID,
			"abeid": userEntryID,
		}).Debugf("kc identifier backend refresh session failed: %v", err)
	}
	return err
}

// DestroySession implements the Backend interface providing destroy to KC session.
func (b *KCIdentifierBackend) DestroySession(ctx context.Context, sessionRef *string) error {
	session, err := b.getSessionForUser(ctx, "", sessionRef, false, false, true)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}

	return session.Destroy(ctx, true)
}

// UserClaims implements the Backend interface, providing user specific claims
// for the user specified by userID.
func (b *KCIdentifierBackend) UserClaims(userID string, authorizedScopes map[string]bool) map[string]interface{} {
	var claims map[string]interface{}

	if authorizedScope, _ := authorizedScopes[ScopeKopanoGC]; authorizedScope {
		claims = make(map[string]interface{})
		// Inject userID as ID claim.
		claims[KopanoGCIDClaim] = userID
	}

	return claims
}

// ScopesSupported implements the Backend interface, providing supported scopes
// when running this backend.
func (b *KCIdentifierBackend) ScopesSupported() []string {
	return kcSupportedScopes
}

// ScopesMeta implements the Backend interface, providing meta data for
// supported scopes.
func (b *KCIdentifierBackend) ScopesMeta() *scopes.Scopes {
	return scopesMeta
}

// Name implements the Backend interface.
func (b *KCIdentifierBackend) Name() string {
	return kcIdentifierBackendName
}

func (b *KCIdentifierBackend) resolveUsername(ctx context.Context, username string, session *kcc.Session) (*kcc.ResolveUserResponse, error) {
	result, err := b.withSessionAndRetry(ctx, session, func(ctx context.Context, session *kcc.Session) (interface{}, error, bool) {
		user, err := b.c.ResolveUsername(ctx, username, session.ID())
		if err != nil {
			return nil, err, true
		}

		if user.Er == kcc.KCERR_NOT_FOUND {
			return nil, user.Er, false
		}

		return user, nil, true
	})
	if err != nil {
		return nil, err
	}

	user := result.(*kcc.ResolveUserResponse)
	return user, err
}

func (b *KCIdentifierBackend) getUser(ctx context.Context, userEntryID string, session *kcc.Session) (*kcc.GetUserResponse, error) {
	result, err := b.withSessionAndRetry(ctx, session, func(ctx context.Context, session *kcc.Session) (interface{}, error, bool) {
		user, err := b.c.GetUser(ctx, userEntryID, session.ID())
		if err != nil {
			return nil, err, true
		}

		if user.Er == kcc.KCERR_NOT_FOUND {
			return nil, user.Er, false
		}

		return user, nil, true
	})
	if err != nil {
		return nil, err
	}

	user := result.(*kcc.GetUserResponse)
	return user, err
}

func (b *KCIdentifierBackend) getSessionForUser(ctx context.Context, userEntryID string, sessionRef *string, register bool, refresh bool, removeIfRegistered bool) (*kcc.Session, error) {
	if b.useGlobalSession {
		return nil, nil
	}
	if sessionRef == nil {
		return nil, nil
	}

	var session *kcc.Session
	if s, ok := b.sessions.Get(*sessionRef); ok {
		// Existing session.
		session = s.(*kcc.Session)
		var removed bool
		if removeIfRegistered {
			b.sessions.Remove(*sessionRef)
			removed = true
		}
		if refresh {
			// Refresh when requested to ensure it is still valid.
			// TODO(longsleep): Debounce session refresh, to avoid doing the
			// same session when it currently refreshes or just has been
			// refreshed.
			err := session.Refresh()
			if err != nil {
				// Silently ignore session refresh errors. Will create a new
				// one automatically if appropriate.
				session = nil
				if !removed {
					b.sessions.Remove(*sessionRef)
				}
			}
		}
		if session != nil {
			return session, nil
		}
	}

	if !register || !refresh || userEntryID == "" {
		return nil, nil
	}

	// Create new auth record with attached identity manager for userEntryID
	// and ensure that we have the required scopes to access kc. We create a
	// simple access which uses the entryID as subject.
	auth := identity.NewAuthRecord(b.identityManager, userEntryID, map[string]bool{
		oidc.ScopeOpenID: true,
		ScopeKopanoGC:    true,
	}, nil, nil)
	// Create a new access token which hopefully gets accepted by our backend.
	accessToken, err := b.oidcProvider.MakeAccessToken(ctx, "konnect", auth)
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend failed to create access token for session: %v", err)
	}

	// Logon.
	response, err := b.c.SSOLogon(ctx, kcc.KOPANO_SSO_TYPE_KCOIDC, userEntryID, []byte(accessToken), kcc.KCNoSessionID, 0)
	if err != nil {
		return nil, fmt.Errorf("kc identifier backend sso logon error: %v", err)
	}
	switch response.Er {
	case kcc.KCSuccess:
		// success.
		if response.SessionID == kcc.KCNoSessionID {
			return nil, fmt.Errorf("kc identifier backend sso logon returned no session")
		}
	default:
		return nil, fmt.Errorf("kc identifier backend sso logon failed: %v", response.Er)
	}

	// Create session instance, for internal use.
	session, err = kcc.CreateSession(b.ctx, b.c, kcc.KCSessionID(response.SessionID), response.ServerGUID, true)
	if err != nil {
		return nil, err
	}
	if register {
		// Register session instance for internal reuse.
		if ok := b.sessions.SetIfAbsent(*sessionRef, session); ok {
			b.logger.WithFields(logrus.Fields{
				"ref":     *sessionRef,
				"session": session,
				"abeid":   userEntryID,
			}).Debugln("kc identifier backend session")
		}
	}

	return session, nil
}

func (b *KCIdentifierBackend) withSessionAndRetry(ctx context.Context, session *kcc.Session, worker func(context.Context, *kcc.Session) (interface{}, error, bool)) (interface{}, error) {
	retries := 0
	for {
		if session == nil {
			// Maybe we have a global session to use?
			session = b.getGlobalSession()
		}
		if session == nil || !session.IsActive() {
			// So no session eh?
			return nil, fmt.Errorf("no server session")
		}

		var failedErr error
		for {
			result, err, shouldRetry := worker(ctx, session)
			if err != nil {
				if !shouldRetry {
					return result, err
				}

				failedErr = err
				break
			}

			// NOTE(longsleep): This is pretty crappy - is there a better way?
			kcErr := reflect.ValueOf(result).Elem().FieldByName("Er").Interface().(kcc.KCError)
			if kcErr != kcc.KCSuccess {
				if !shouldRetry {
					return result, kcErr
				}

				failedErr = kcErr
				break
			}

			return result, nil
		}

		if failedErr != nil {
			switch failedErr {
			case kcc.KCERR_END_OF_SESSION:
				session.Destroy(ctx, false)
			default:
				return nil, failedErr
			}
		}

		// If reach here, its a retry.
		select {
		case <-time.After(kcSessionRetryDelay):
			// Retry now.
		case <-ctx.Done():
			// Abort.
			return nil, ctx.Err()
		}

		retries++
		if retries > kcSessionMaxRetries {
			b.logger.WithField("retry", retries).Errorln("kc identifier backend giving up kc request")
			return nil, failedErr
		}
		b.logger.WithField("retry", retries).Debugln("kc identifier backend retry in progress")
	}
}

func (b *KCIdentifierBackend) setGlobalSession(session *kcc.Session) {
	b.globalSessionMutex.Lock()
	b.globalSession = session
	b.globalSessionMutex.Unlock()
}

func (b *KCIdentifierBackend) getGlobalSession() *kcc.Session {
	b.globalSessionMutex.RLock()
	session := b.globalSession
	b.globalSessionMutex.RUnlock()
	return session
}

func exIDToString(exID []byte) string {
	return base64.RawURLEncoding.EncodeToString(exID)
}
