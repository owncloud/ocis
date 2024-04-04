package scanners

import (
	"time"

	"github.com/dutchcoders/go-clamd"
)

// NewClamAV returns a Scanner talking to clamAV via socket
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
func (s ClamAV) Scan(in Input) (Result, error) {
	ch, err := s.clamd.ScanStream(in.Body, make(chan bool))
	if err != nil {
		return Result{}, err
	}

	r := <-ch
	return Result{
		Infected:    r.Status == clamd.RES_FOUND,
		Description: r.Description,
		ScanTime:    time.Now(),
	}, nil
}
