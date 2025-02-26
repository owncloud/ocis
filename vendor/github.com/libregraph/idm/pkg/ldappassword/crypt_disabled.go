//go:build disable_crypt || darwin || windows

/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldappassword

import (
	"errors"
)

func crypt(pass, salt string) (string, error) {
	return "", errors.New("CRYPT unsupported on this platform")
}
