//go:build !disable_crypt && (linux || freebsd || netbsd)

/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldappassword

import (
	gocrypt "github.com/amoghe/go-crypt"
)

func crypt(pass, salt string) (string, error) {
	return gocrypt.Crypt(pass, salt)
}
