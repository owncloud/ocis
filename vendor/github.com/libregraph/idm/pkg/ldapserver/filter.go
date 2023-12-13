// Copyright 2011 The Go Authors. All rights reserved.
// Copyright 2021 The LibreGraph Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldapserver

import (
	"strings"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/cases"
)

const (
	FilterAnd             = ldap.FilterAnd
	FilterOr              = ldap.FilterOr
	FilterNot             = ldap.FilterNot
	FilterEqualityMatch   = ldap.FilterEqualityMatch
	FilterSubstrings      = ldap.FilterSubstrings
	FilterGreaterOrEqual  = ldap.FilterGreaterOrEqual
	FilterLessOrEqual     = ldap.FilterLessOrEqual
	FilterPresent         = ldap.FilterPresent
	FilterApproxMatch     = ldap.FilterApproxMatch
	FilterExtensibleMatch = ldap.FilterExtensibleMatch
)

var (
	FilterMap = ldap.FilterMap
	casefold  = cases.Fold()
)

const (
	FilterSubstringsInitial = ldap.FilterSubstringsInitial
	FilterSubstringsAny     = ldap.FilterSubstringsAny
	FilterSubstringsFinal   = ldap.FilterSubstringsFinal
)

func CompileFilter(filter string) (*ber.Packet, error) {
	return ldap.CompileFilter(filter)
}

func DecompileFilter(packet *ber.Packet) (ret string, err error) {
	return ldap.DecompileFilter(packet)
}

func ServerApplyFilter(f *ber.Packet, entry *ldap.Entry) (bool, LDAPResultCode) {
	switch FilterMap[uint64(f.Tag)] {
	default:
		//log.Fatalf("Unknown LDAP filter code: %d", f.Tag)
		return false, ldap.LDAPResultOperationsError
	case "Equality Match":
		if len(f.Children) != 2 {
			return false, ldap.LDAPResultOperationsError
		}
		attribute := f.Children[0].Value.(string)
		value := f.Children[1].Value.(string)
		for _, a := range entry.Attributes {
			if strings.EqualFold(a.Name, attribute) {
				for _, v := range a.Values {
					if strings.EqualFold(v, value) {
						return true, ldap.LDAPResultSuccess
					}
				}
			}
		}
	case "Present":
		for _, a := range entry.Attributes {
			if strings.EqualFold(a.Name, f.Data.String()) {
				return true, ldap.LDAPResultSuccess
			}
		}
	case "And":
		for _, child := range f.Children {
			ok, exitCode := ServerApplyFilter(child, entry)
			if exitCode != ldap.LDAPResultSuccess {
				return false, exitCode
			}
			if !ok {
				return false, ldap.LDAPResultSuccess
			}
		}
		return true, ldap.LDAPResultSuccess
	case "Or":
		anyOk := false
		for _, child := range f.Children {
			ok, exitCode := ServerApplyFilter(child, entry)
			if exitCode != ldap.LDAPResultSuccess {
				return false, exitCode
			} else if ok {
				anyOk = true
			}
		}
		if anyOk {
			return true, ldap.LDAPResultSuccess
		}
	case "Not":
		if len(f.Children) != 1 {
			return false, ldap.LDAPResultOperationsError
		}
		ok, exitCode := ServerApplyFilter(f.Children[0], entry)
		if exitCode != ldap.LDAPResultSuccess {
			return false, exitCode
		} else if !ok {
			return true, ldap.LDAPResultSuccess
		}
	case "Substrings":
		if len(f.Children) != 2 {
			return false, ldap.LDAPResultOperationsError
		}
		attribute := f.Children[0].Value.(string)
		bytes := f.Children[1].Children[0].Data.Bytes()
		value := casefold.String(string(bytes))
		for _, a := range entry.Attributes {
			if strings.EqualFold(a.Name, attribute) {
				for _, v := range a.Values {
					v = casefold.String(v)
					switch f.Children[1].Children[0].Tag {
					case FilterSubstringsInitial:
						if strings.HasPrefix(v, value) {
							return true, ldap.LDAPResultSuccess
						}
					case FilterSubstringsAny:
						if strings.Contains(v, value) {
							return true, ldap.LDAPResultSuccess
						}
					case FilterSubstringsFinal:
						if strings.HasSuffix(v, value) {
							return true, ldap.LDAPResultSuccess
						}
					}
				}
			}
		}
	case "FilterGreaterOrEqual": // TODO
		return false, ldap.LDAPResultOperationsError
	case "FilterLessOrEqual": // TODO
		return false, ldap.LDAPResultOperationsError
	case "FilterApproxMatch": // TODO
		return false, ldap.LDAPResultOperationsError
	case "FilterExtensibleMatch": // TODO
		return false, ldap.LDAPResultOperationsError
	}

	return false, ldap.LDAPResultSuccess
}

func ServerFilterScope(baseDN string, scope int, entry *ldap.Entry) (bool, LDAPResultCode) {
	// constrained search scope
	parsedBaseDn, err := ldap.ParseDN(baseDN)
	if err != nil {
		return false, ldap.LDAPResultOperationsError
	}
	parsedDn, err := ldap.ParseDN(entry.DN)
	if err != nil {
		return false, ldap.LDAPResultOperationsError
	}
	switch scope {
	case ldap.ScopeWholeSubtree: // The scope is constrained to the entry named by baseObject and to all its subordinates.
	case ldap.ScopeBaseObject: // The scope is constrained to the entry named by baseObject.
		if !parsedDn.EqualFold(parsedBaseDn) {
			return false, ldap.LDAPResultSuccess
		}
	case ldap.ScopeSingleLevel: // The scope is constrained to the immediate subordinates of the entry named by baseObject.
		parts := strings.Split(entry.DN, ",")
		if len(parts) < 2 && !parsedDn.EqualFold(parsedBaseDn) {
			return false, ldap.LDAPResultSuccess
		}
		subDn := strings.Join(parts[1:], ",")
		parsedSubDn, err := ldap.ParseDN(subDn)
		if err != nil {
			return false, ldap.LDAPResultOperationsError
		}
		if !parsedSubDn.EqualFold(parsedBaseDn) {
			return false, ldap.LDAPResultSuccess
		}
	}

	return true, ldap.LDAPResultSuccess
}

func ServerFilterAttributes(attributes []string, entry *ldap.Entry) (LDAPResultCode, error) {
	// attributes
	if len(attributes) > 1 || (len(attributes) == 1 && len(attributes[0]) > 0) {
		_, err := filterAttributes(entry, attributes)
		if err != nil {
			return ldap.LDAPResultOperationsError, err
		}
	}

	return ldap.LDAPResultSuccess, nil
}
