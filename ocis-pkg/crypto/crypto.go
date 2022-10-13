// Package crypto implements utility functions for handling crypto related files.
package crypto

import (
	"bytes"
	"crypto/x509"
	"errors"
	"io"
)

// NewCertPoolFromPEM reads certificates from io.Reader and returns a x509.CertPool
// containing those certificates.
func NewCertPoolFromPEM(crts ...io.Reader) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	var buf bytes.Buffer
	for _, c := range crts {
		if _, err := io.Copy(&buf, c); err != nil {
			return nil, err
		}
		if !certPool.AppendCertsFromPEM(buf.Bytes()) {
			return nil, errors.New("failed to append cert from PEM")
		}
		buf.Reset()
	}

	return certPool, nil
}
