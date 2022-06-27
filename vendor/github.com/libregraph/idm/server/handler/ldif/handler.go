/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/go-ldap/ldif"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"stash.kopano.io/kgol/rndm"

	"github.com/libregraph/idm/pkg/ldapdn"
	"github.com/libregraph/idm/pkg/ldapserver"
	"github.com/libregraph/idm/server/handler"
)

type ldifHandler struct {
	logger  logrus.FieldLogger
	fn      string
	options *Options

	baseDN                  string
	adminDN                 string
	allowLocalAnonymousBind bool

	ctx context.Context

	current atomic.Value

	activeSearchPagings cmap.ConcurrentMap
}

var _ handler.Handler = (*ldifHandler)(nil) // Verify that *ldifHandler implements handler.Handler.

func NewLDIFHandler(logger logrus.FieldLogger, fn string, options *Options) (handler.Handler, error) {
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

	h := &ldifHandler{
		logger:  logger,
		fn:      fn,
		options: options,

		allowLocalAnonymousBind: options.AllowLocalAnonymousBind,

		ctx: context.Background(),

		activeSearchPagings: cmap.New(),
	}
	if h.baseDN, err = ldapdn.ParseNormalize(options.BaseDN); err != nil {
		return nil, err
	}
	if h.adminDN, err = ldapdn.ParseNormalize(options.AdminDN); err != nil {
		return nil, err
	}

	err = h.open()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *ldifHandler) open() error {
	if !strings.EqualFold(h.options.BaseDN, h.baseDN) {
		return fmt.Errorf("mismatched BaseDN")
	}

	info, err := os.Stat(h.fn)
	if err != nil {
		return fmt.Errorf("failed to open LDIF: %w", err)
	}

	var l *ldif.LDIF
	index := newIndexMapRegister()

	if info.IsDir() {
		h.logger.Debugln("loading LDIF files from folder")
		h.options.templateBasePath = h.fn
		var parseErrors []error
		l, parseErrors, err = parseLDIFDirectory(h.fn, h.options)
		if err != nil {
			return err
		}
		if len(parseErrors) > 0 {
			for _, parseErr := range parseErrors {
				h.logger.WithError(parseErr).Errorln("LDIF error")
			}
			return fmt.Errorf("error in LDIF files")
		}
	} else {
		h.logger.Debugln("loading LDIF")
		h.options.templateBasePath = filepath.Dir(h.fn)
		l, err = parseLDIFFile(h.fn, h.options)
		if err != nil {
			return err
		}
	}

	t, err := treeFromLDIF(l, index, h.options)
	if err != nil {
		return err
	}

	// Store parsed data as memory value.
	value := &ldifMemoryValue{
		l: l,
		t: t,

		index: index,
	}
	h.current.Store(value)

	h.logger.WithFields(logrus.Fields{
		"version":       l.Version,
		"entries_count": len(l.Entries),
		"tree_length":   t.Len(),
		"base_dn":       h.options.BaseDN,
		"indexes":       len(index),
	}).Debugln("loaded LDIF")

	return nil
}

func (h *ldifHandler) load() *ldifMemoryValue {
	value := h.current.Load()
	return value.(*ldifMemoryValue)
}

func (h *ldifHandler) WithContext(ctx context.Context) handler.Handler {
	if ctx == nil {
		panic("nil context")
	}

	h2 := new(ldifHandler)
	*h2 = *h
	h2.ctx = ctx
	return h2
}

func (h *ldifHandler) Reload(ctx context.Context) error {
	return h.open()
}

func (h *ldifHandler) Add(_ string, _ *ldap.AddRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldapserver.LDAPResultCode, error) {
	bindDN = strings.ToLower(bindDN)

	logger := h.logger.WithFields(logrus.Fields{
		"bind_dn":     bindDN,
		"remote_addr": conn.RemoteAddr().String(),
	})

	if err := h.validateBindDN(bindDN, conn); err != nil {
		logger.WithError(err).Debugln("ldap bind request BindDN validation failed")
		return ldap.LDAPResultInsufficientAccessRights, nil
	}

	if bindSimplePw == "" {
		logger.Debugf("ldap anonymous bind request")
		if bindDN == "" {
			return ldap.LDAPResultSuccess, nil
		} else {
			return ldap.LDAPResultUnwillingToPerform, nil
		}
	} else {
		logger.Debugf("ldap bind request")
	}

	current := h.load()

	entryRecord, found := current.t.Get([]byte(bindDN))
	if !found {
		err := fmt.Errorf("user not found")
		logger.WithError(err).Debugf("ldap bind error")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	entry := entryRecord.(*ldifEntry)

	if err := entry.validatePassword(bindSimplePw); err != nil {
		logger.WithError(err).Debugf("ldap bind credentials error")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	return ldap.LDAPResultSuccess, nil
}

func (h *ldifHandler) Delete(_ string, _ *ldap.DelRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifHandler) Modify(_ string, _ *ldap.ModifyRequest, _ net.Conn) (ldapserver.LDAPResultCode, error) {
	return ldap.LDAPResultUnwillingToPerform, errors.New("unsupported operation")
}

func (h *ldifHandler) Search(bindDN string, searchReq *ldap.SearchRequest, conn net.Conn) (ldapserver.ServerSearchResult, error) {
	bindDN = strings.ToLower(bindDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)
	logger := h.logger.WithFields(logrus.Fields{
		"bind_dn":        bindDN,
		"search_base_dn": searchBaseDN,
		"remote_addr":    conn.RemoteAddr().String(),
		"controls":       searchReq.Controls,
		"size_limit":     searchReq.SizeLimit,
	})

	logger.Debugf("ldap search request for %s", searchReq.Filter)

	if err := h.validateBindDN(bindDN, conn); err != nil {
		logger.WithError(err).Debugln("ldap search request BindDN validation failed")
		return ldapserver.ServerSearchResult{
			ResultCode: ldap.LDAPResultInsufficientAccessRights,
		}, err
	}

	indexFilter, _ := parseFilterToIndexFilter(searchReq.Filter)

	if !strings.HasSuffix(searchBaseDN, h.baseDN) {
		err := fmt.Errorf("ldap search BaseDN is not in our BaseDN %s", h.baseDN)
		return ldapserver.ServerSearchResult{
			ResultCode: ldap.LDAPResultInsufficientAccessRights,
		}, err
	}

	doneControls := []ldap.Control{}
	var pagingControl *ldap.ControlPaging
	var pagingCookie []byte
	if paging := ldap.FindControl(searchReq.Controls, ldap.ControlTypePaging); paging != nil {
		pagingControl = paging.(*ldap.ControlPaging)
		if searchReq.SizeLimit > 0 && pagingControl.PagingSize >= uint32(searchReq.SizeLimit) {
			pagingControl = nil
		} else {
			pagingCookie = pagingControl.Cookie
		}
	}

	pumpCh, resultCode := func() (<-chan *ldifEntry, ldapserver.LDAPResultCode) {
		var pumpCh chan *ldifEntry
		var start = true
		if pagingControl != nil {
			if len(pagingCookie) == 0 {
				pagingCookie = []byte(base64.RawStdEncoding.EncodeToString(rndm.GenerateRandomBytes(8)))
				pagingControl.Cookie = pagingCookie
				pumpCh = make(chan *ldifEntry)
				h.activeSearchPagings.Set(string(pagingControl.Cookie), pumpCh)
				logger.WithField("paging_cookie", string(pagingControl.Cookie)).Debugln("ldap search paging pump start")
			} else {
				pumpChRecord, ok := h.activeSearchPagings.Get(string(pagingControl.Cookie))
				if !ok {
					return nil, ldap.LDAPResultUnwillingToPerform
				}
				if pagingControl.PagingSize > 0 {
					logger.WithField("paging_cookie", string(pagingControl.Cookie)).Debugln("ldap search paging pump continue")
					pumpCh = pumpChRecord.(chan *ldifEntry)
					start = false
				} else {
					// No paging size with cookie, means abandon.
					start = false
					logger.WithField("paging_cookie", string(pagingControl.Cookie)).Debugln("search paging pump abandon")
					// TODO(longsleep): Cancel paging pump context.
					h.activeSearchPagings.Remove(string(pagingControl.Cookie))
				}
			}
		} else {
			pumpCh = make(chan *ldifEntry)
		}
		if start {
			current := h.load()
			go h.searchEntriesPump(h.ctx, current, pumpCh, searchReq, pagingControl, indexFilter)
		}

		return pumpCh, ldap.LDAPResultSuccess
	}()
	if resultCode != ldap.LDAPResultSuccess {
		err := fmt.Errorf("search unable to perform: %d", resultCode)
		return ldapserver.ServerSearchResult{
			ResultCode: resultCode,
		}, err
	}

	filterPacket, err := ldapserver.CompileFilter(searchReq.Filter)
	if err != nil {
		return ldapserver.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, err
	}

	var entryRecord *ldifEntry
	var entries []*ldap.Entry
	var entry *ldap.Entry
	var count uint32
	var keep bool
results:
	for {
		select {
		case entryRecord = <-pumpCh:
			if entryRecord == nil {
				// All done, set cookie to empty.
				pagingCookie = []byte{}
				break results

			} else {
				entry = entryRecord.Entry

				// Apply filter.
				keep, resultCode = ldapserver.ServerApplyFilter(filterPacket, entry)
				if resultCode != ldap.LDAPResultSuccess {
					return ldapserver.ServerSearchResult{
						ResultCode: resultCode,
					}, errors.New("search filter apply error")
				}
				if !keep {
					continue
				}

				// Filter scope.
				keep, resultCode = ldapserver.ServerFilterScope(searchReq.BaseDN, searchReq.Scope, entry)
				if resultCode != ldap.LDAPResultSuccess {
					return ldapserver.ServerSearchResult{
						ResultCode: resultCode,
					}, errors.New("search scope apply error")
				}
				if !keep {
					continue
				}

				// Make a copy, before filtering attributes.
				e := &ldap.Entry{
					DN:         entry.DN,
					Attributes: make([]*ldap.EntryAttribute, len(entry.Attributes)),
				}
				copy(e.Attributes, entry.Attributes)

				// Filter attributes from entry.
				resultCode, err = ldapserver.ServerFilterAttributes(searchReq.Attributes, e)
				if err != nil {
					return ldapserver.ServerSearchResult{
						ResultCode: resultCode,
					}, err
				}

				// Append entry as result.
				entries = append(entries, e)

				// Count and more.
				count++
				if pagingControl != nil {
					if count >= pagingControl.PagingSize {
						break results
					}
				}
				if searchReq.SizeLimit > 0 && count >= uint32(searchReq.SizeLimit) {
					// TODO(longsleep): handle total sizelimit for paging.
					break results
				}
			}
		}
	}

	if pagingControl != nil {
		doneControls = append(doneControls, &ldap.ControlPaging{
			PagingSize: 0,
			Cookie:     pagingCookie,
		})
	}

	return ldapserver.ServerSearchResult{
		Entries:    entries,
		Referrals:  []string{},
		Controls:   doneControls,
		ResultCode: ldap.LDAPResultSuccess,
	}, nil
}

func (h *ldifHandler) searchEntriesPump(ctx context.Context, current *ldifMemoryValue, pumpCh chan<- *ldifEntry, searchReq *ldap.SearchRequest, pagingControl *ldap.ControlPaging, indexFilter [][]string) {
	defer func() {
		if pagingControl != nil {
			h.activeSearchPagings.Remove(string(pagingControl.Cookie))
			close(pumpCh)
			h.logger.WithField("paging_cookie", string(pagingControl.Cookie)).Debugln("ldap search paging pump ended")
		} else {
			close(pumpCh)
		}
	}()

	pump := func(entryRecord *ldifEntry) bool {
		select {
		case pumpCh <- entryRecord:
		case <-ctx.Done():
			if pagingControl != nil {
				h.logger.WithField("paging_cookie", string(pagingControl.Cookie)).Warnln("ldap search paging pump context done")
			} else {
				h.logger.Warnln("ldap search pump context done")
			}
			return false
		case <-time.After(1 * time.Minute):
			if pagingControl != nil {
				h.logger.WithField("paging_cookie", string(pagingControl.Cookie)).Warnln("ldap search paging pump timeout")
			} else {
				h.logger.Warnln("ldap search pump timeout")
			}
			return false
		}
		return true
	}

	searchBaseDN := strings.ToLower(searchReq.BaseDN)

	load := true
	if len(indexFilter) > 0 {
		// Get entries with help of index.
		load = false
		var results []*[]*ldifEntry
		for _, f := range indexFilter {
			indexed, found := current.index.Load(f[0], f[1], f[2:]...)
			if !found {
				load = true
				break
			}
			results = append(results, &indexed)
		}
		if !load {
			cache := make(map[*ldifEntry]struct{})
			for _, indexed := range results {
				for _, entryRecord := range *indexed {
					if _, cached := cache[entryRecord]; cached {
						// Prevent duplicates.
						continue
					}
					if strings.HasSuffix(entryRecord.DN, searchBaseDN) {
						if ok := pump(entryRecord); !ok {
							return
						}
					}
					cache[entryRecord] = struct{}{}
				}
			}
		}
	}
	if load {
		// Walk through all entries (this is slow).
		h.logger.WithField("filter", searchReq.Filter).Warnln("ldap search filter does not match any index, using slow walk")
		current.t.WalkSuffix([]byte(searchBaseDN), func(key []byte, entryRecord interface{}) bool {
			if ok := pump(entryRecord.(*ldifEntry)); !ok {
				return true
			}
			return false
		})
	}
}

func (h *ldifHandler) validateBindDN(bindDN string, conn net.Conn) error {
	if bindDN == "" {
		if h.allowLocalAnonymousBind {
			host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
			if net.ParseIP(host).IsLoopback() {
				return nil
			}
			return fmt.Errorf("anonymous BindDN rejected")
		}
		return fmt.Errorf("anonymous BindDN not allowed")
	}

	if strings.HasSuffix(bindDN, h.baseDN) {
		return nil
	}
	return fmt.Errorf("the BindDN is not in our BaseDN: %s", h.baseDN)
}

func (h *ldifHandler) Close(bindDN string, conn net.Conn) error {
	h.logger.WithFields(logrus.Fields{
		"bind_dn":     bindDN,
		"remote_addr": conn.RemoteAddr().String(),
	}).Debugln("ldap close")

	return nil
}
