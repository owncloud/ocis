// Package structs provides some utility functions for dealing with structs.
package structs

// CopyOrZeroValue returns a copy of s if s is not nil otherwise the zero value of T will be returned.
func CopyOrZeroValue[T any](s *T) *T {
	cp := new(T)
	if s != nil {
		*cp = *s
	}
	return cp
}
