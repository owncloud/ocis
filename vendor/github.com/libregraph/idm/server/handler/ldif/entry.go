/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/libregraph/idm/pkg/ldappassword"
)

type ldifEntry struct {
	*ldap.Entry

	UserPassword *ldap.EntryAttribute
}

func (entry *ldifEntry) validatePassword(bindSimplePw string) error {
	match, err := ldappassword.Validate(bindSimplePw, entry.UserPassword.Values[0])
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("password mismatch")
	}
	return nil
}
