package crypto

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

var serialNumber, _ = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

var acmeTemplate = x509.Certificate{
	SerialNumber: serialNumber,
	Subject: pkix.Name{
		Organization: []string{"Acme Corp"},
		CommonName:   "OCIS",
	},
	NotBefore: time.Now(),
	NotAfter:  time.Now().Add(24 * time.Hour * 365),

	KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	BasicConstraintsValid: true,
}
