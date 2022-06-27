/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"

	"github.com/libregraph/idm/pkg/ldapserver"
	"github.com/libregraph/idm/server/handler"
)

type ldifMiddleware struct {
	logger  logrus.FieldLogger
	fn      string
	options *Options

	baseDN string

	current atomic.Value

	next handler.Handler
}

var _ handler.Handler = (*ldifMiddleware)(nil) // Verify that *configHandler implements handler.Handler.

func NewLDIFMiddleware(logger logrus.FieldLogger, fn string, options *Options) (handler.Middleware, error) {
	if fn == "" {
		return nil, fmt.Errorf("file name is empty")
	}
	if options.BaseDN == "" {
		return nil, fmt.Errorf("base dn is empty")
	}

	fn, err := filepath.Abs(fn)
	if err != nil {
		return nil, err
	}
	logger = logger.WithField("fn", fn)

	h := &ldifMiddleware{
		logger:  logger,
		fn:      fn,
		options: options,

		baseDN: strings.ToLower(options.BaseDN),
	}

	err = h.open()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *ldifMiddleware) open() error {
	if !strings.EqualFold(h.options.BaseDN, h.baseDN) {
		return fmt.Errorf("mismatched BaseDN")
	}

	h.logger.Debugln("loading LDIF")
	l, err := parseLDIFFile(h.fn, h.options)
	if err != nil {
		return err
	}

	t, err := treeFromLDIF(l, nil, h.options)
	if err != nil {
		return err
	}

	// Store parsed data as memory value.
	value := &ldifMemoryValue{
		t: t,
	}
	h.current.Store(value)

	h.logger.WithFields(logrus.Fields{
		"version":       l.Version,
		"entries_count": len(l.Entries),
		"tree_length":   t.Len(),
		"base_dn":       h.options.BaseDN,
	}).Debugln("loaded LDIF")

	return nil
}

func (h *ldifMiddleware) load() *ldifMemoryValue {
	value := h.current.Load()
	return value.(*ldifMemoryValue)
}

func (h *ldifMiddleware) WithHandler(next handler.Handler) handler.Handler {
	h.next = next

	return h
}

func (h *ldifMiddleware) WithContext(ctx context.Context) handler.Handler {
	if ctx == nil {
		panic("nil context")
	}

	h2 := new(ldifMiddleware)
	*h2 = *h
	h2.next = h.next.WithContext(ctx)
	return h2
}

func (h *ldifMiddleware) Reload(ctx context.Context) error {
	err := h.open()
	if err != nil {
		return err
	}

	return h.next.Reload(ctx)
}

func (h *ldifMiddleware) Add(_ string, _ *ldap.AddRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifMiddleware) Bind(bindDN, bindSimplePw string, conn net.Conn) (resultCode ldapserver.LDAPResultCode, err error) {
	bindDN = strings.ToLower(bindDN)

	if bindSimplePw == "" { // Empty password means anonymous bind.
		return h.next.Bind(bindDN, bindSimplePw, conn)
	}

	current := h.load()

	entryRecord, found := current.t.Get([]byte(bindDN))
	if found {
		logger := h.logger.WithFields(logrus.Fields{
			"bind_dn":     bindDN,
			"remote_addr": conn.RemoteAddr().String(),
		})

		if !strings.HasSuffix(bindDN, h.baseDN) {
			err := fmt.Errorf("the BindDN is not in our BaseDN %s", h.baseDN)
			logger.WithError(err).Infoln("ldap bind error")
			return ldap.LDAPResultInvalidCredentials, nil
		}

		if err := entryRecord.(*ldifEntry).validatePassword(bindSimplePw); err != nil {
			logger.WithError(err).Infoln("bind error")
			return ldap.LDAPResultInvalidCredentials, nil
		}

		return ldap.LDAPResultSuccess, nil
	}

	return h.next.Bind(bindDN, bindSimplePw, conn)
}

func (h *ldifMiddleware) Delete(_ string, _ *ldap.DelRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifMiddleware) Modify(_ string, _ *ldap.ModifyRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifMiddleware) Search(bindDN string, searchReq *ldap.SearchRequest, conn net.Conn) (result ldapserver.ServerSearchResult, err error) {
	return h.next.Search(bindDN, searchReq, conn)
}

func (h *ldifMiddleware) Close(bindDN string, conn net.Conn) error {
	return h.next.Close(bindDN, conn)
}
