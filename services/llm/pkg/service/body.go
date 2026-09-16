package service

import (
	"errors"
	"io"
)

var errBodyTooLarge = errors.New("request body too large")

// readLimited reads r fully, returning errBodyTooLarge if the body exceeds
// maxBytes.
func readLimited(r io.Reader, maxBytes int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, errBodyTooLarge
	}
	return data, nil
}
