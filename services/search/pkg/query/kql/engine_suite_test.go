package kql

import (
	"time"
)

// PatchTimeNow is here to path the package time now func,
// this only exists for the tests context
func PatchTimeNow(t func() time.Time) {
	timeNow = t
}
