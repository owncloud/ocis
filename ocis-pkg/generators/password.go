package generators

import (
	"crypto/rand"
	"math/big"
)

const (
	// PasswordChars contains alphanumeric chars (0-9, A-Z, a-z), plus "-=+!@#$%^&*."
	PasswordChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-=+!@#$%^&*."
	// AlphaNumChars contains alphanumeric chars (0-9, A-Z, a-z)
	AlphaNumChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenerateRandomPassword generates a random password with the given length.
// The password will contain chars picked from the `PasswordChars` constant.
// If an error happens, the string will be empty and the error will be non-nil.
//
// This is equivalent to `GenerateRandomString(PasswordChars, length)`
func GenerateRandomPassword(length int) (string, error) {
	return generateString(PasswordChars, length)
}

// GenerateRandomString generates a random string with the given length
// based on the chars provided. You can use `PasswordChars` or `AlphaNumChars`
// constants, or even any other string.
//
// Chars from the provided string will be picked uniformly. The provided
// constants have unique chars, which means that all the chars will have the
// same probability of being picked.
// You can use your own strings to change that probability. For example, using
// "AAAB" you'll have a 75% of probability of getting "A" and 25% of "B"
func GenerateRandomString(chars string, length int) (string, error) {
	return generateString(chars, length)
}

func generateString(chars string, length int) (string, error) {
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		ret[i] = chars[num.Int64()]
	}

	return string(ret), nil
}
