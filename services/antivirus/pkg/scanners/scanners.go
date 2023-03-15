package scanners

import (
	"fmt"
	"io"
	"time"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config"
)

// ScanResult is the common scan result to all scanners
type ScanResult struct {
	Infected    bool
	Scantime    time.Time
	Description string
}

// Scanner is an abstraction for the actual virus scan
type Scanner interface {
	Scan(file io.Reader) (ScanResult, error)
}

// New returns a new scanner from config
func New(c config.Scanner) (Scanner, error) {
	switch c.Type {
	default:
		return nil, fmt.Errorf("unknown av scanner: '%s'", c.Type)
	case "clamav":
		return NewClamAV(c.ClamAV.Socket), nil
	case "icap":
		return NewICAP(c.ICAP.URL, c.ICAP.Service, time.Duration(c.ICAP.Timeout)*time.Second)
	}

}
