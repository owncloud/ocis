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
func New(c *config.Config) (Scanner, error) {
	switch c.Scanner.Type {
	default:
		return nil, fmt.Errorf("unknown av scanner: '%s'", c.Scanner.Type)
	case "clamav":
		return NewClamAV(c.Scanner.ClamAV.Socket), nil
	case "icap":
		return NewICAP(c.Scanner.ICAP.URL, c.Scanner.ICAP.Service, time.Duration(c.Scanner.ICAP.Timeout)*time.Second)
	}

}
