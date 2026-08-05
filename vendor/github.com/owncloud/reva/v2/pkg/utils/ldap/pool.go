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

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

const (
	// defaultPoolSize is the pool size used when Config.PoolSize is <= 0. It is a plain size, not a
	// sentinel: there is no "unlimited" pool, a non-positive value always falls back to this default.
	defaultPoolSize = 5
	// defaultPoolCheckoutTimeout is the checkout timeout used when Config.PoolCheckoutTimeout is <=
	// 0. Like PoolSize, there is no "wait forever" mode: a non-positive value falls back to this
	// default rather than disabling the timeout.
	defaultPoolCheckoutTimeout = 30 * time.Second
)

var (
	// ErrPoolExhausted is returned by a checkout that couldn't obtain a connection within the
	// configured checkout timeout.
	ErrPoolExhausted = errors.New("ldap: connection pool exhausted")
	errPoolClosed    = errors.New("ldap: connection pool is closed")
)

// clientFactory dials and authenticates a new LDAP connection for the given config. ConnPool holds
// it as a field, rather than calling dialLDAP directly, purely so tests can substitute a
// network-free fake; every production pool uses dialLDAP.
type clientFactory func(Config) (ldap.Client, error)

// ConnPool is a bounded pool of authenticated LDAP connections. Connections are dialed and bound
// lazily on checkout and reused across requests; connections that fail with a network error are
// discarded and lazily re-dialed on a later checkout, instead of the pool eagerly reconnecting them.
type ConnPool struct {
	config  Config
	timeout time.Duration
	dial    clientFactory
	logger  *zerolog.Logger

	// sem bounds the number of connections checked out at once: checkout acquires it (blocking with
	// p.timeout until a slot is free), release releases it to free the slot. Its weight is the pool
	// size.
	sem *semaphore.Weighted
	// idle holds connections that are checked in and ready to be reused; it is buffered to the pool
	// size so release never blocks. checkout drains it first before dialing a new connection.
	idle chan ldap.Client

	closed atomic.Bool
}

// NewLDAPPool returns a new ConnPool initialized from config, logging through logger (or silently
// if logger is nil). No connection is dialed until the first checkout.
func NewLDAPPool(config Config, logger *zerolog.Logger) *ConnPool {
	size := config.PoolSize
	if size <= 0 {
		size = defaultPoolSize
	}
	timeout := config.PoolCheckoutTimeout
	if timeout <= 0 {
		timeout = defaultPoolCheckoutTimeout
	}
	if logger == nil {
		nop := zerolog.Nop()
		logger = &nop
	}
	return &ConnPool{
		config:  config,
		timeout: timeout,
		dial:    dialLDAP,
		logger:  logger,
		sem:     semaphore.NewWeighted(int64(size)),
		idle:    make(chan ldap.Client, size),
	}
}

// dialLDAP dials config.URI and, if config.BindDN is set, binds as that DN: it establishes a fully
// authenticated connection, not just a network connection.
func dialLDAP(config Config) (ldap.Client, error) {
	l, err := ldap.DialURL(config.URI, ldap.DialWithTLSConfig(config.TLSConfig))
	if err != nil {
		return nil, err
	}
	if config.BindDN != "" {
		if err := l.Bind(config.BindDN, config.BindPassword); err != nil {
			l.Close()
			return nil, err
		}
	}
	return l, nil
}

// checkout reserves a pool slot (blocking up to p.timeout) and returns a connection: an idle one if
// available, otherwise a freshly dialed one.
func (p *ConnPool) checkout() (ldap.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	if err := p.sem.Acquire(ctx, 1); err != nil {
		return nil, ErrPoolExhausted
	}

	if p.IsClosing() {
		p.sem.Release(1)
		return nil, errPoolClosed
	}

	select {
	case conn := <-p.idle:
		return conn, nil
	default:
	}

	p.logger.Debug().Msg("dialing new pooled LDAP connection")
	conn, err := p.dial(p.config)
	if err != nil {
		p.sem.Release(1)
		return nil, err
	}
	return conn, nil
}

// release returns conn to the pool, or closes and discards it if opErr is a network error (or the
// pool has since been closed), and frees the slot reserved by checkout.
func (p *ConnPool) release(conn ldap.Client, opErr error) {
	if p.IsClosing() || (opErr != nil && ldap.IsErrorWithCode(opErr, ldap.ErrorNetwork)) {
		if err := conn.Close(); err != nil {
			p.logger.Error().Err(err).Msg("error closing pooled LDAP connection")
		}
	} else {
		p.idle <- conn
	}
	p.sem.Release(1)
}

// do checks out a connection, runs fn, and releases the connection. Checked-out idle connections
// are not health-checked before use, so a connection that went stale while idle (server-side idle
// timeout, LB/firewall reaping) is expected to fail the first op after being idle; on a network
// error, do retries once more with a freshly checked-out connection (the failed one was evicted by
// release) instead of surfacing the failure to the caller, mirroring ConnWithReconnect's retry.
func (p *ConnPool) do(fn func(conn ldap.Client) error) error {
	var opErr error
	for try := 0; try <= defaultRetries; try++ {
		conn, err := p.checkout()
		if err != nil {
			return err
		}
		opErr = fn(conn)
		p.release(conn, opErr)
		if opErr == nil || !ldap.IsErrorWithCode(opErr, ldap.ErrorNetwork) {
			return opErr
		}
	}
	return ldap.NewError(ldap.ErrorNetwork, errMaxRetries)
}

// Close closes the pool: currently idle connections are closed immediately, further checkouts are
// rejected, and connections already checked out are closed as they're returned.
func (p *ConnPool) Close() error {
	if !p.closed.CompareAndSwap(false, true) {
		return nil
	}

	for {
		select {
		case conn := <-p.idle:
			if err := conn.Close(); err != nil {
				p.logger.Error().Err(err).Msg("error closing pooled LDAP connection")
			}
		default:
			return nil
		}
	}
}

// Search implements the ldap.Client interface
func (p *ConnPool) Search(sr *ldap.SearchRequest) (*ldap.SearchResult, error) {
	var res *ldap.SearchResult
	err := p.do(func(conn ldap.Client) error {
		var err error
		res, err = conn.Search(sr)
		return err
	})
	return res, err
}

// Add implements the ldap.Client interface
func (p *ConnPool) Add(a *ldap.AddRequest) error {
	return p.do(func(conn ldap.Client) error {
		return conn.Add(a)
	})
}

// Del implements the ldap.Client interface
func (p *ConnPool) Del(d *ldap.DelRequest) error {
	return p.do(func(conn ldap.Client) error {
		return conn.Del(d)
	})
}

// Modify implements the ldap.Client interface
func (p *ConnPool) Modify(m *ldap.ModifyRequest) error {
	return p.do(func(conn ldap.Client) error {
		return conn.Modify(m)
	})
}

// ModifyDN implements the ldap.Client interface
func (p *ConnPool) ModifyDN(m *ldap.ModifyDNRequest) error {
	return p.do(func(conn ldap.Client) error {
		return conn.ModifyDN(m)
	})
}

// Extended implements the ldap.Client interface
func (p *ConnPool) Extended(request *ldap.ExtendedRequest) (*ldap.ExtendedResponse, error) {
	var res *ldap.ExtendedResponse
	err := p.do(func(conn ldap.Client) error {
		var err error
		res, err = conn.Extended(request)
		return err
	})
	return res, err
}

// GetLastError implements the ldap.Client interface. A pool has no single "last" connection to
// report on, so this always returns nil; it is unused by any caller of ConnPool today.
func (p *ConnPool) GetLastError() error {
	return nil
}

// IsClosing implements the ldap.Client interface
func (p *ConnPool) IsClosing() bool {
	return p.closed.Load()
}

// Remaining methods to fulfill ldap.Client interface

// Start implements the ldap.Client interface
func (p *ConnPool) Start() {}

// SetTimeout implements the ldap.Client interface
func (p *ConnPool) SetTimeout(time.Duration) {}

// StartTLS implements the ldap.Client interface
func (p *ConnPool) StartTLS(*tls.Config) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// Bind implements the ldap.Client interface
func (p *ConnPool) Bind(username, password string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// UnauthenticatedBind implements the ldap.Client interface
func (p *ConnPool) UnauthenticatedBind(username string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// SimpleBind implements the ldap.Client interface
func (p *ConnPool) SimpleBind(*ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// ExternalBind implements the ldap.Client interface
func (p *ConnPool) ExternalBind() error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// NTLMUnauthenticatedBind implements the ldap.Client interface
func (p *ConnPool) NTLMUnauthenticatedBind(domain, username string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// Unbind implements the ldap.Client interface
func (p *ConnPool) Unbind() error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// ModifyWithResult implements the ldap.Client interface
func (p *ConnPool) ModifyWithResult(m *ldap.ModifyRequest) (*ldap.ModifyResult, error) {
	var res *ldap.ModifyResult
	err := p.do(func(conn ldap.Client) error {
		var err error
		res, err = conn.ModifyWithResult(m)
		return err
	})
	return res, err
}

// Compare implements the ldap.Client interface
func (p *ConnPool) Compare(dn, attribute, value string) (bool, error) {
	return false, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// PasswordModify implements the ldap.Client interface
func (p *ConnPool) PasswordModify(m *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	var res *ldap.PasswordModifyResult
	err := p.do(func(conn ldap.Client) error {
		var err error
		res, err = conn.PasswordModify(m)
		return err
	})
	return res, err
}

// SearchWithPaging implements the ldap.Client interface
func (p *ConnPool) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// SearchAsync implements the ldap.Client interface
func (p *ConnPool) SearchAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int) ldap.Response {
	// unimplemented
	return nil
}

// TLSConnectionState implements the ldap.Client interface
func (p *ConnPool) TLSConnectionState() (tls.ConnectionState, bool) {
	return tls.ConnectionState{}, false
}

// DirSync implements the ldap.Client interface
func (p *ConnPool) DirSync(searchRequest *ldap.SearchRequest, flags, maxAttrCount int64, cookie []byte) (*ldap.SearchResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

// DirSyncAsync implements the ldap.Client interface
func (p *ConnPool) DirSyncAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, flags, maxAttrCount int64, cookie []byte) ldap.Response {
	// unimplemented
	return nil
}

// Syncrepl implements the ldap.Client interface
func (p *ConnPool) Syncrepl(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, mode ldap.ControlSyncRequestMode, cookie []byte, reloadHint bool) ldap.Response {
	// unimplemented
	return nil
}
