package scanners

import (
	"io"
	"time"

	"github.com/dutchcoders/go-clamd"
)

// NewClamAV returns an Scanner talking to clamAV via socket
func NewClamAV(socket string) *ClamAV {
	return &ClamAV{
		clamd: clamd.NewClamd(socket),
	}
}

// ClamAV is a Scanner based on clamav
type ClamAV struct {
	clamd *clamd.Clamd
}

// Scan to fulfill Scanner interface
func (s ClamAV) Scan(file io.Reader) (ScanResult, error) {
	ch, err := s.clamd.ScanStream(file, make(chan bool))
	if err != nil {
		return ScanResult{}, err
	}

	r := <-ch
	return ScanResult{
		Infected:    r.Status == clamd.RES_FOUND,
		Description: r.Description,
		Scantime:    time.Now(),
	}, nil
}
