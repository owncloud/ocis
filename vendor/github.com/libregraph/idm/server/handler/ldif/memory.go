/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"github.com/go-ldap/ldif"
	"github.com/spacewander/go-suffix-tree"
)

type ldifMemoryValue struct {
	l *ldif.LDIF
	t *suffix.Tree

	index Index
}
