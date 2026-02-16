// Copyright 2011 The Go Authors. All rights reserved.
// Copyright 2021 The LibreGraph Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldapserver

type LDAPResultCode uint8

const (
	LDAPBindAuthSimple = 0
	LDAPBindAuthSASL   = 3
)
