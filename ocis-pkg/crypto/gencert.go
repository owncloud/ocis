package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}, l log.Logger) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			l.Fatal().Err(err).Msg("Unable to marshal ECDSA private key")
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// GenCert generates TLS-Certificates
func GenCert(certName string, keyName string, l log.Logger) error {
	var pk *rsa.PrivateKey
	var err error

	pk, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// if either the key or certificate already exist skip this entire ordeal
	_, certErr := os.Stat(certName)
	_, keyErr := os.Stat(keyName)

	if certErr == nil || keyErr == nil {
		l.Debug().Msg(fmt.Sprintf("%v certificate or key already present, using these", filepath.Base(certName)))
		return nil
	}

	persistCertificate(certName, l, pk)
	persistKey(keyName, l, pk)
	return nil
}

func persistCertificate(certName string, l log.Logger, pk *rsa.PrivateKey) {
	if err := ensureExistsDir(certName); err != nil {
		l.Fatal().Err(err).Msg("creating certificate destination: " + certName)
	}

	certificate, err := generateCertificate(pk)
	if err != nil {
		l.Fatal().Err(err).Msg("creating certificate: " + filepath.Dir(certName))
	}

	certOut, err := os.Create(certName)
	if err != nil {
		l.Fatal().Err(err).Msgf("failed to open `%v` for writing", certName)
	}

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	if err != nil {
		l.Fatal().Err(err).Msg("failed to encode certificate")
	}

	err = certOut.Close()
	if err != nil {
		l.Fatal().Err(err).Msg("failed to write cert")
	}
	l.Info().Msg(fmt.Sprintf("written certificate to %v", certName))
}

func persistKey(keyName string, l log.Logger, pk *rsa.PrivateKey) {
	if err := ensureExistsDir(keyName); err != nil {
		l.Fatal().Err(err).Msg("creating certificate destination: " + keyName)
	}

	keyOut, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		l.Fatal().Err(err).Msgf("failed to open %v for writing", keyName)
	}
	err = pem.Encode(keyOut, pemBlockForKey(pk, l))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to encode key")
	}

	err = keyOut.Close()
	if err != nil {
		l.Fatal().Err(err).Msg("failed to write key")
	}
	l.Info().Msg(fmt.Sprintf("written key to %v", keyName))
}

// genCert generates a self signed certificate using a random rsa key.
func generateCertificate(pk *rsa.PrivateKey) ([]byte, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Corp"},
			CommonName:   "OCIS",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := []string{"127.0.0.1", "localhost"}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	return x509.CreateCertificate(rand.Reader, &template, &template, publicKey(pk), pk)
}

func ensureExistsDir(uri string) error {
	certPath := filepath.Dir(uri)
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		err = os.MkdirAll(certPath, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}
