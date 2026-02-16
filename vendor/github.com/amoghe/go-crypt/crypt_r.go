// +build linux

// Package crypt provides wrappers around functions available in crypt.h
//
// It wraps around the GNU specific extension (crypt_r) when it is available
// (i.e. where GOOS=linux). This makes the go function reentrant (and thus
// callable from concurrent goroutines).
package crypt

import (
	"syscall"
	"unsafe"
)

/*
#cgo LDFLAGS: -lcrypt

#define _GNU_SOURCE

#include <stdlib.h>
#include <string.h>
#include <crypt.h>

char *gnu_ext_crypt(char *pass, char *salt) {
  char *enc = NULL;
  char *ret = NULL;
  struct crypt_data data;
  data.initialized = 0;

  enc = crypt_r(pass, salt, &data);
  if(enc == NULL) {
    return NULL;
  }

  ret = (char *)malloc(strlen(enc)+1); // for trailing null
  strncpy(ret, enc, strlen(enc));
  ret[strlen(enc)]= '\0'; // paranoid

  return ret;
}
*/
import "C"

// Crypt provides a wrapper around the glibc crypt_r() function.
// For the meaning of the arguments, refer to the package README.
func Crypt(pass, salt string) (string, error) {
	c_pass := C.CString(pass)
	defer C.free(unsafe.Pointer(c_pass))

	c_salt := C.CString(salt)
	defer C.free(unsafe.Pointer(c_salt))

	c_enc, err := C.gnu_ext_crypt(c_pass, c_salt)
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
		// Make sure we acutally return an error, musl e.g. does not
		// set errno in all cases here.
		if err == nil {
			err = syscall.EINVAL
		}
		return "", err
	}

	// Return nil error if the string is non-nil.
	// As per the errno.h manpage, functions are allowed to set errno
	// on success. Caller should ignore errno on success.
	return hash, nil
}
