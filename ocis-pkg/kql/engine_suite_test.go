package kql

import (
	"time"
)

// PatchTimeNow is here to patch the package time now func,
// which is used in the test suite
func PatchTimeNow(t func() time.Time) {
	timeNow = t
}
