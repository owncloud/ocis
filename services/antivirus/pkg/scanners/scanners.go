package scanners

import (
	"io"
	"time"
)

// The Result is the common scan result to all scanners
type Result struct {
	Infected    bool
	ScanTime    time.Time
	Description string
}

// The Input is the common input to all scanners
type Input struct {
	Body io.Reader
	Size int64
	Url  string
	Name string
}
