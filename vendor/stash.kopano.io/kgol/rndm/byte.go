/*
 * Copyright 2017 Kopano
 *
 * Use of this source code is governed by a MIT license
 * that can be found in the LICENSE.txt file.
 *
 */

package rndm

import (
	"crypto/rand"
)

// GenerateRandomBytes returns securely generated random bytes. It will panic
// when the system fails to provide enough secure random data.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := ReadRandomBytes(b)
	if err != nil {
		panic("unable to read enough random bytes")
	}

	return b
}

// ReadRandomBytes is a helper function that reads random data into the provided
// []byte. Tt returns the number of random bytes read and an error if fewer
// bytes were read.
func ReadRandomBytes(b []byte) (n int, err error) {
	return rand.Read(b)
}
