package cutil

import "C"

import (
	"bytes"
)

// SplitBuffer splits a byte-slice buffer, typically returned from C code,
// into a slice of strings.
// The contents of the buffer are assumed to be null-byte separated.
// If the buffer contains a sequence of null-bytes it will assume that the
// "space" between the bytes are meant to be empty strings.
func SplitBuffer(b []byte) []string {
	return splitBufStrings(b, true)
}

// SplitSparseBuffer splits a byte-slice buffer, typically returned from C code,
// into a slice of strings.
// The contents of the buffer are assumed to be null-byte separated.
// This function assumes that buffer to be "sparse" such that only non-null-byte
// strings will be returned, and no "empty" strings exist if null-bytes
// are found adjacent to each other.
func SplitSparseBuffer(b []byte) []string {
	return splitBufStrings(b, false)
}

// If keepEmpty is true, empty substrings will be returned, by default they are
// excluded from the results.
// This is almost certainly a suboptimal implementation, especially for
// keepEmpty=true case. Optimizing the functions is a job for another day.
func splitBufStrings(b []byte, keepEmpty bool) []string {
	values := make([]string, 0)
	// the final null byte should be the terminating null in C
	// we never want to preserve the empty string after it
	if len(b) > 0 && b[len(b)-1] == 0 {
		b = b[:len(b)-1]
	}
	if len(b) == 0 {
		return values
	}
	for _, s := range bytes.Split(b, []byte{0}) {
		if !keepEmpty && len(s) == 0 {
			continue
		}
		values = append(values, string(s))
	}
	return values
}
