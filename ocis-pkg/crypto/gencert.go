package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

var (
	DEFAULT_HOSTS = []string{"127.0.0.1", "localhost"}
)

// GenCert generates TLS-Certificates. This function has side effects: it creates the respective certificate / key pair at
// the destination locations unless the tuple already exists, if that is the case, this is a noop.
func GenCert(certName string, keyName string, l log.Logger) error {
	var pk *rsa.PrivateKey
	var err error

	pk, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	_, certErr := os.Stat(certName)
	_, keyErr := os.Stat(keyName)

	if certErr == nil || keyErr == nil {
		l.Debug().Msg(
			fmt.Sprintf("%v certificate / key pair already present. skipping acme certificate generation",
				filepath.Base(certName)))
		return nil
	}

	persistCertificate(certName, l, pk)
	persistKey(keyName, l, pk)

	return nil
}

// persistCertificate generates a certificate using pk as private key and proceeds to store it into a file named certName.
func persistCertificate(certName string, l log.Logger, pk interface{}) {
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

// genCert generates a self signed certificate using a random rsa key.
func generateCertificate(pk interface{}) ([]byte, error) {
	for _, h := range DEFAULT_HOSTS {
		if ip := net.ParseIP(h); ip != nil {
			acmeTemplate.IPAddresses = append(acmeTemplate.IPAddresses, ip)
		} else {
			acmeTemplate.DNSNames = append(acmeTemplate.DNSNames, h)
		}
	}

	return x509.CreateCertificate(rand.Reader, &acmeTemplate, &acmeTemplate, publicKey(pk), pk)
}

// persistKey persists the private key used to generate the certificate at the configured location.
func persistKey(keyName string, l log.Logger, pk interface{}) {
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
