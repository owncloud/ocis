package version

import (
	"time"
)

var (
	// String gets defined by the build system.
	String = "0.0.0"

	// Date indicates the build date.
	Date = "00000000"
)

// Compiled returns the compile time of this service.
func Compiled() time.Time {
	t, _ := time.Parse("20060102", Date)
	return t
}
