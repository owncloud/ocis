// +build darwin windows

// Package crypt provides wrappers around functions available in crypt.h
//
// It wraps around the GNU specific extension (crypt) when the reentrant version
// (crypt_r) is unavailable. The non-reentrant version is guarded by a global lock
// so as to be safely callable from concurrent goroutines.
package crypt

import (
	"errors"
)

// Crypt does currently not provide an implementation for windows and darwin
func Crypt(pass, salt string) (string, error) {
	return "", errors.New("unsupported platform")
}
