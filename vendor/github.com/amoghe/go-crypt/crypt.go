// +build freebsd netbsd

// Package crypt provides wrappers around functions available in crypt.h
//
// It wraps around the GNU specific extension (crypt) when the reentrant version
// (crypt_r) is unavailable. The non-reentrant version is guarded by a global lock
// so as to be safely callable from concurrent goroutines.
package crypt

import (
	"sync"
	"unsafe"
)

/*
#cgo LDFLAGS: -lcrypt
#define _GNU_SOURCE
#include <stdlib.h>
#include <unistd.h>
*/
import "C"

var (
	mu sync.Mutex
)

// Crypt provides a wrapper around the glibc crypt() function.
// For the meaning of the arguments, refer to the README.
func Crypt(pass, salt string) (string, error) {
	c_pass := C.CString(pass)
	defer C.free(unsafe.Pointer(c_pass))

	c_salt := C.CString(salt)
	defer C.free(unsafe.Pointer(c_salt))

	mu.Lock()
	c_enc, err := C.crypt(c_pass, c_salt)
	mu.Unlock()

	if c_enc == nil {
		return "", err
	}
	defer C.free(unsafe.Pointer(c_enc))

	// From the crypt(3) man-page. Upon error, crypt_r writes an invalid
	// hashed passphrase to the output field of their data argument, and
	// crypt writes an invalid hash to its static storage area. This
	// string will be shorter than 13 characters, will begin with a ‘*’,
	// and will not compare equal to setting.
	hash := C.GoString(c_enc)
	if len(hash) > 0 && hash[0] == '*' {
		return "", err
	}

	// Return nil error if the string is non-nil.
	// As per the errno.h manpage, functions are allowed to set errno
	// on success. Caller should ignore errno on success.
	return C.GoString(c_enc), err
}
