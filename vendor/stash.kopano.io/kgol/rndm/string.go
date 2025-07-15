/*
 * Copyright 2017 Kopano
 *
 * Use of this source code is governed by a MIT license
 * that can be found in the LICENSE.txt file.
 *
 */

package rndm

import (
	"encoding/base64"
)

// GenerateRandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system fails to provide secure random
// data.
func GenerateRandomString(s int) string {
	b := GenerateRandomBytes(s)

	return base64.RawURLEncoding.EncodeToString(b)
}
