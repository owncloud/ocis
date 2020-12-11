package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
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
	var priv interface{}
	var err error

	priv, err = rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		l.Fatal().Err(err).Msg("Failed to generate private key")
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to generate serial number")
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

	//template.IsCA = true
	//template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create certificate")
	}

	certPath := filepath.Dir(certName)
	l.Error().Msg("certPath: " + certPath)
	l.Error().Msg("certName: " + certName)

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		err = os.MkdirAll(certPath, 0700)
		if err != nil {
			l.Fatal().Err(err).Msg("Failed to create path " + certPath)
		}
	}

	certOut, err := os.Create(certName)
	if err != nil {
		l.Fatal().Err(err).Msgf("Failed to open %v for writing", certName)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to encode certificate")
	}
	err = certOut.Close()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to write cert")
	}
	l.Info().Msg("Written server.crt")

	keyPath := filepath.Dir(keyName)
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		err = os.MkdirAll(keyPath, 0700)
		if err != nil {
			l.Fatal().Err(err).Msg("Failed to create path " + keyPath)
		}
	}

	keyOut, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		l.Fatal().Err(err).Msgf("Failed to open %v for writing", keyName)
	}
	err = pem.Encode(keyOut, pemBlockForKey(priv, l))
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to encode key")
	}
	err = keyOut.Close()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to write key")
	}
	l.Info().Msgf("Written %v", keyName)
	return nil
}
