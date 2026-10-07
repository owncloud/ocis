// Copyright 2022 CERN
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

package ldap

// LDAP automatic reconnection mechanism, inspired by:
// https://gist.github.com/emsearcy/cba3295d1a06d4c432ab4f6173b65e4f#file-ldap_snippet-go

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	backoff "github.com/cenkalti/backoff/v5"
	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
)

// RetryPolicy controls retry behaviour for one class of LDAP operations.
//
// The predicates take the error rather than a bare result code because go-ldap raises every
// connection failure as ErrorNetwork(200) and distinguishes the pre-transmit case (safe to
// retry a write) from the post-transmit case only by the wrapped message string — and because a
// failed conn.Write arrives as a plain error with no result code at all (see sendFailedErrMsg).
// Both policies inspect the message: the write policy to decide retryability at all
// (isPreSendNetworkErr), the read policy only to catch the codeless failed-write case.
type RetryPolicy struct {
	MaxRetries     int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	isRetryable    func(err error) bool
	needsReconnect func(err error) bool
}

// NewReadPolicy returns a RetryPolicy for read operations (Search).
// Tier 1 (ErrorNetwork/ServerDown/ConnectError/Timeout/LocalError): reconnect + backoff.
// Tier 2 (Busy/Unavailable): backoff only, no reconnect (server signals transient overload).
//
// Deliberately excluded:
//   - LDAPResultTimeLimitExceeded (3): usually a too-broad query, not transient load;
//     retrying re-runs the expensive query and worsens a struggling server.
//   - LDAPResultReferral (10): a routing signal, not a transient error.
//
// A failed conn.Write (isSendFailedErr) is retried too. It carries no result code, so the switches
// below cannot match it on their own.
func NewReadPolicy(maxRetries int, baseDelay, maxDelay time.Duration) RetryPolicy {
	return RetryPolicy{
		MaxRetries: maxRetries,
		BaseDelay:  baseDelay,
		MaxDelay:   maxDelay,
		isRetryable: func(err error) bool {
			if isSendFailedErr(err) {
				return true
			}
			switch ldapErrCode(err) {
			case ldap.ErrorNetwork,
				ldap.LDAPResultServerDown,
				ldap.LDAPResultConnectError,
				ldap.LDAPResultTimeout,
				ldap.LDAPResultLocalError,
				ldap.LDAPResultBusy,
				ldap.LDAPResultUnavailable:
				return true
			}
			return false
		},
		needsReconnect: func(err error) bool {
			if isSendFailedErr(err) {
				return true
			}
			switch ldapErrCode(err) {
			case ldap.ErrorNetwork,
				ldap.LDAPResultServerDown,
				ldap.LDAPResultConnectError,
				ldap.LDAPResultTimeout,
				ldap.LDAPResultLocalError:
				return true
			}
			return false
		},
	}
}

// NewWritePolicy returns a RetryPolicy for write operations.
//
// A write is retried only when the connection was known-unusable before the request packet was
// sent, so the mutation provably never reached the server (isPreSendNetworkErr). Every other
// error — including a network drop while reading the response — may have applied the write, so
// retrying it could double-apply.
//
// This intentionally does not key on result codes: go-ldap raises all connection failures as
// ErrorNetwork(200) and never emits LDAPResultServerDown(81)/ConnectError(91), so a code-based
// write allowlist would match nothing and never retry.
func NewWritePolicy(maxRetries int, baseDelay, maxDelay time.Duration) RetryPolicy {
	return RetryPolicy{
		MaxRetries:     maxRetries,
		BaseDelay:      baseDelay,
		MaxDelay:       maxDelay,
		isRetryable:    isPreSendNetworkErr,
		needsReconnect: func(error) bool { return true }, // a pre-send network error always needs a fresh connection
	}
}

type ldapConnection struct {
	Conn  *ldap.Conn
	Error error
}

// ConnWithReconnect maintains an LDAP Connection that automatically reconnects after network errors
type ConnWithReconnect struct {
	conn    chan ldapConnection
	reset   chan *ldap.Conn
	read    RetryPolicy
	write   RetryPolicy
	sleepFn func(time.Duration)
	dialFn  func(Config) (*ldap.Conn, error)
	logger  *zerolog.Logger
}

// NewLDAPWithReconnect Returns a new ConnWithReconnect initialized from config
func NewLDAPWithReconnect(config Config) *ConnWithReconnect {
	conn := ConnWithReconnect{
		conn:    make(chan ldapConnection),
		reset:   make(chan *ldap.Conn),
		read:    NewReadPolicy(config.RetryMaxCount, config.RetryBaseDelay, config.RetryMaxDelay),
		write:   NewWritePolicy(config.RetryMaxCount, config.RetryBaseDelay, config.RetryMaxDelay),
		sleepFn: time.Sleep,
	}
	logger := zerolog.Nop()
	conn.logger = &logger
	conn.dialFn = conn.ldapConnect
	go conn.ldapAutoConnect(config)
	return &conn
}

// SetLogger sets the logger for the current instance
func (c *ConnWithReconnect) SetLogger(logger *zerolog.Logger) {
	c.logger = logger
}


// ldapErrCode extracts the LDAP result code from an error. Errors that are not an *ldap.Error carry
// no result code and map to LDAPResultOther, which neither policy's switch treats as retryable —
// so a codeless error that IS safe to retry must be matched on its message instead
// (see sendFailedErrMsg).
func ldapErrCode(err error) uint16 {
	if err == nil {
		return 0
	}
	var lerr *ldap.Error
	if errors.As(err, &lerr) {
		return lerr.ResultCode
	}
	// Unknown error type: map to a non-retryable code.
	return ldap.LDAPResultOther
}

// preSendNetworkErrMsgs are the substrings go-ldap uses for ErrorNetwork(200) failures raised in
// Conn.sendMessageWithFlags, before the request packet is handed to the write loop: the IsClosing
// check, and the "could not send" fallback. Matching one proves the request never left the client,
// so retrying a write cannot double-apply.
//
// go-ldap has no distinct result code for pre- vs post-transmit network errors — both are 200 — so
// this whitelist keys on its internal message strings. It is deliberately fail-closed: if go-ldap
// renames a message, isPreSendNetworkErr returns false and the write is surfaced to the caller
// rather than retried. Revisit this list on every go-ldap upgrade (see TestWritePolicyMatchesEmittableCode).
//
// Note these cover only the case where the reader goroutine has already noticed the drop and set
// closing; when an operation beats the reader to it the failure surfaces from the write loop
// instead — see sendFailedErrMsg.
var preSendNetworkErrMsgs = []string{
	"ldap: connection closed",
	"ldap: could not send message for unknown reason",
}

// sendFailedErrMsg is the message go-ldap's processMessages loop wraps a failed conn.Write in
// (fmt.Errorf("unable to send request: %s", err)). Unlike every other failure in the send path this
// is a plain error, not an *ldap.Error, so it carries no result code and ldapErrCode maps it to
// LDAPResultOther — meaning a code-keyed policy cannot match it.
//
// Matching it proves the request never reached the server: go-ldap adds the message to
// messageContexts only after a successful write, so an operation that failed here can never be
// matched to a response and cannot have applied a mutation. Retrying it is therefore safe even for
// a write.
//
// This is the race window a reaped idle connection hits when the operation beats the reader
// goroutine to noticing the drop: IsClosing() is still false, so the preSendNetworkErrMsgs checks
// pass, the packet reaches the write loop, and conn.Write fails with EPIPE/ECONNRESET.
const sendFailedErrMsg = "unable to send request:"

// isSendFailedErr reports whether err is go-ldap's failed-conn.Write error. See sendFailedErrMsg.
func isSendFailedErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), sendFailedErrMsg)
}

// isPreSendNetworkErr reports whether err was raised before the request was transmitted, and is
// therefore safe to retry even for a write. See preSendNetworkErrMsgs and sendFailedErrMsg.
func isPreSendNetworkErr(err error) bool {
	if err == nil {
		return false
	}
	// A failed conn.Write carries no result code, so check it before the code-keyed cases.
	if isSendFailedErr(err) {
		return true
	}
	var lerr *ldap.Error
	if !errors.As(err, &lerr) || lerr.ResultCode != ldap.ErrorNetwork {
		return false
	}
	msg := err.Error()
	for _, m := range preSendNetworkErrMsgs {
		if strings.Contains(msg, m) {
			return true
		}
	}
	return false
}

// RetryOp executes fn under the given retry policy.
func (c *ConnWithReconnect) RetryOp(policy RetryPolicy, fn func(*ldap.Conn) error) error {
	if policy.MaxRetries < 1 {
		policy.MaxRetries = 1
	}
	conn, err := c.getConnection()
	if err != nil {
		return err
	}

	var bo *backoff.ExponentialBackOff
	for try := 0; ; try++ {
		err = fn(conn)
		if err == nil {
			return nil
		}
		if !policy.isRetryable(err) || try >= policy.MaxRetries {
			break
		}

		if policy.BaseDelay > 0 {
			if bo == nil {
				bo = backoff.NewExponentialBackOff()
				bo.InitialInterval = policy.BaseDelay
				if policy.MaxDelay > 0 {
					bo.MaxInterval = policy.MaxDelay
				}
				bo.Reset()
			}
			c.sleepFn(bo.NextBackOff())
		}

		if policy.needsReconnect(err) {
			conn, err = c.reconnect(conn)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// Search implements the ldap.Client interface
func (c *ConnWithReconnect) Search(sr *ldap.SearchRequest) (*ldap.SearchResult, error) {
	var res *ldap.SearchResult
	err := c.RetryOp(c.read, func(conn *ldap.Conn) error {
		var e error
		res, e = conn.Search(sr)
		return e
	})
	return res, err
}

// Add implements the ldap.Client interface
func (c *ConnWithReconnect) Add(a *ldap.AddRequest) error {
	return c.RetryOp(c.write, func(conn *ldap.Conn) error {
		return conn.Add(a)
	})
}

// Del implements the ldap.Client interface
func (c *ConnWithReconnect) Del(d *ldap.DelRequest) error {
	return c.RetryOp(c.write, func(conn *ldap.Conn) error {
		return conn.Del(d)
	})
}

// Modify implements the ldap.Client interface
func (c *ConnWithReconnect) Modify(m *ldap.ModifyRequest) error {
	return c.RetryOp(c.write, func(conn *ldap.Conn) error {
		return conn.Modify(m)
	})
}

// ModifyDN implements the ldap.Client interface
func (c *ConnWithReconnect) ModifyDN(m *ldap.ModifyDNRequest) error {
	return c.RetryOp(c.write, func(conn *ldap.Conn) error {
		return conn.ModifyDN(m)
	})
}

// Extended implements the ldap.Client interface
func (c *ConnWithReconnect) Extended(request *ldap.ExtendedRequest) (*ldap.ExtendedResponse, error) {
	var res *ldap.ExtendedResponse
	err := c.RetryOp(c.write, func(conn *ldap.Conn) error {
		var e error
		res, e = conn.Extended(request)
		return e
	})
	return res, err
}

// ModifyWithResult implements the ldap.Client interface
func (c *ConnWithReconnect) ModifyWithResult(m *ldap.ModifyRequest) (*ldap.ModifyResult, error) {
	var res *ldap.ModifyResult
	err := c.RetryOp(c.write, func(conn *ldap.Conn) error {
		var e error
		res, e = conn.ModifyWithResult(m)
		return e
	})
	return res, err
}

// PasswordModify implements the ldap.Client interface
func (c *ConnWithReconnect) PasswordModify(m *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	var res *ldap.PasswordModifyResult
	err := c.RetryOp(c.write, func(conn *ldap.Conn) error {
		var e error
		res, e = conn.PasswordModify(m)
		return e
	})
	return res, err
}

func (c *ConnWithReconnect) getConnection() (*ldap.Conn, error) {
	conn := <-c.conn
	if conn.Conn != nil && !ldap.IsErrorWithCode(conn.Error, ldap.ErrorNetwork) {
		c.logger.Debug().Msg("using existing Connection")
		return conn.Conn, conn.Error
	}
	return c.reconnect(conn.Conn)
}

func (c *ConnWithReconnect) ldapAutoConnect(config Config) {
	var (
		l   *ldap.Conn
		err error
	)

	for {
		select {
		case resConn := <-c.reset:
			// Only close the connection and reconnect if the current
			// connection, matches the one we got via the reset channel.
			// If they differ we already reconnected
			switch {
			case l == nil:
				c.logger.Debug().Msg("reconnecting to LDAP")
				l, err = c.dialFn(config)
			case l != resConn:
				c.logger.Debug().Msg("already reconnected")
				continue
			default:
				c.logger.Debug().Msg("closing and reconnecting to LDAP")
				l.Close()
				l, err = c.dialFn(config)
			}
		case c.conn <- ldapConnection{l, err}:
		}
	}
}

func (c *ConnWithReconnect) ldapConnect(config Config) (*ldap.Conn, error) {
	c.logger.Debug().Msgf("Connecting to %s", config.URI)

	l, err := ldap.DialURL(config.URI, ldap.DialWithTLSConfig(config.TLSConfig))
	if err != nil {
		c.logger.Error().Err(err).Msg("could not get ldap Connection")
		return nil, err
	}
	c.logger.Debug().Msg("LDAP Connected")
	if config.BindDN != "" {
		c.logger.Debug().Msgf("Binding as %s", config.BindDN)
		err = l.Bind(config.BindDN, config.BindPassword)
		if err != nil {
			c.logger.Debug().Err(err).Msg("Bind failed")
			l.Close()
			return nil, err
		}
	}
	return l, err
}

func (c *ConnWithReconnect) reconnect(resetConn *ldap.Conn) (*ldap.Conn, error) {
	c.logger.Debug().Msg("LDAP connection reset")
	c.reset <- resetConn
	c.logger.Debug().Msg("Waiting for new connection")
	result := <-c.conn
	return result.Conn, result.Error
}

// Remaining methods to fulfill ldap.Client interface

// Start implements the ldap.Client interface
func (c *ConnWithReconnect) Start() {}

// StartTLS implements the ldap.Client interface
func (c *ConnWithReconnect) StartTLS(*tls.Config) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// Close implements the ldap.Client interface
func (c *ConnWithReconnect) Close() (err error) {
	conn, err := c.getConnection()
	if err != nil {
		return err
	}
	return conn.Close()
}

func (c *ConnWithReconnect) GetLastError() error {
	conn, err := c.getConnection()
	if err != nil {
		return err
	}
	return conn.GetLastError()
}

// IsClosing implements the ldap.Client interface
func (c *ConnWithReconnect) IsClosing() bool {
	return false
}

// SetTimeout implements the ldap.Client interface
func (c *ConnWithReconnect) SetTimeout(time.Duration) {}

// Bind implements the ldap.Client interface
func (c *ConnWithReconnect) Bind(username, password string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// UnauthenticatedBind implements the ldap.Client interface
func (c *ConnWithReconnect) UnauthenticatedBind(username string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// SimpleBind implements the ldap.Client interface
func (c *ConnWithReconnect) SimpleBind(*ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// ExternalBind implements the ldap.Client interface
func (c *ConnWithReconnect) ExternalBind() error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// Compare implements the ldap.Client interface
func (c *ConnWithReconnect) Compare(dn, attribute, value string) (bool, error) {
	return false, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// SearchWithPaging implements the ldap.Client interface
func (c *ConnWithReconnect) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// SearchAsync implements the ldap.Client interface
func (c *ConnWithReconnect) SearchAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int) ldap.Response {
	return nil
}

// NTLMUnauthenticatedBind implements the ldap.Client interface
func (c *ConnWithReconnect) NTLMUnauthenticatedBind(domain, username string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// TLSConnectionState implements the ldap.Client interface
func (c *ConnWithReconnect) TLSConnectionState() (tls.ConnectionState, bool) {
	return tls.ConnectionState{}, false
}

// Unbind implements the ldap.Client interface
func (c *ConnWithReconnect) Unbind() error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// DirSync implements the ldap.Client interface
func (c *ConnWithReconnect) DirSync(searchRequest *ldap.SearchRequest, flags, maxAttrCount int64, cookie []byte) (*ldap.SearchResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// DirSyncAsync implements the ldap.Client interface
func (c *ConnWithReconnect) DirSyncAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, flags, maxAttrCount int64, cookie []byte) ldap.Response {
	return nil
}

// Syncrepl implements the ldap.Client interface
func (c *ConnWithReconnect) Syncrepl(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, mode ldap.ControlSyncRequestMode, cookie []byte, reloadHint bool) ldap.Response {
	return nil
}
