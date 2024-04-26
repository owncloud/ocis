package envsubst

import (
	"io/ioutil"
	"os"

	"github.com/a8m/envsubst/parse"
)

// String returns the parsed template string after processing it.
// If the parser encounters invalid input, it returns an error describing the failure.
func String(s string) (string, error) {
	return StringRestricted(s, false, false)
}

// StringRestricted returns the parsed template string after processing it.
// If the parser encounters invalid input, or a restriction is violated, it returns
// an error describing the failure.
// Errors on first failure or returns a collection of failures if failOnFirst is false
func StringRestricted(s string, noUnset, noEmpty bool) (string, error) {
	return StringRestrictedNoDigit(s, noUnset, noEmpty , false)
}

// Like StringRestricted but additionally allows to ignore env variables which start with a digit.
func StringRestrictedNoDigit(s string, noUnset, noEmpty bool, noDigit bool) (string, error) {
	return parse.New("string", os.Environ(),
		&parse.Restrictions{noUnset, noEmpty, noDigit}).Parse(s)
}

// Bytes returns the bytes represented by the parsed template after processing it.
// If the parser encounters invalid input, it returns an error describing the failure.
func Bytes(b []byte) ([]byte, error) {
	return BytesRestricted(b, false, false)
}

// BytesRestricted returns the bytes represented by the parsed template after processing it.
// If the parser encounters invalid input, or a restriction is violated, it returns
// an error describing the failure.
func BytesRestricted(b []byte, noUnset, noEmpty bool) ([]byte, error) {
	return BytesRestrictedNoDigit(b, noUnset, noEmpty, false)
}

// Like BytesRestricted but additionally allows to ignore env variables which start with a digit.
func BytesRestrictedNoDigit(b []byte, noUnset, noEmpty bool, noDigit bool) ([]byte, error) {
	s, err := parse.New("bytes", os.Environ(),
		&parse.Restrictions{noUnset, noEmpty, noDigit}).Parse(string(b))
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// ReadFile call io.ReadFile with the given file name.
// If the call to io.ReadFile failed it returns the error; otherwise it will
// call envsubst.Bytes with the returned content.
func ReadFile(filename string) ([]byte, error) {
	return ReadFileRestricted(filename, false, false)
}

// ReadFileRestricted calls io.ReadFile with the given file name.
// If the call to io.ReadFile failed it returns the error; otherwise it will
// call envsubst.Bytes with the returned content.
func ReadFileRestricted(filename string, noUnset, noEmpty bool) ([]byte, error) {
	return ReadFileRestrictedNoDigit(filename, noUnset, noEmpty, false)
}

// Like ReadFileRestricted but additionally allows to ignore env variables which start with a digit.
func ReadFileRestrictedNoDigit(filename string, noUnset, noEmpty bool, noDigit bool) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return BytesRestrictedNoDigit(b, noUnset, noEmpty, noDigit)
}
