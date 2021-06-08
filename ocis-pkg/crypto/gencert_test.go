package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"path/filepath"
	"testing"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

func Test_ensureExistsDir(t *testing.T) {
	var tmpDir = t.TempDir()

	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "creates a dir if it does not exist",
			args: args{
				uri: filepath.Join(tmpDir, "example"),
			},
			wantErr: false,
		},
		{
			name: "noop if the target directory exists",
			args: args{
				uri: filepath.Join(tmpDir, "example"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureExistsDir(tt.args.uri); (err != nil) != tt.wantErr {
				t.Errorf("ensureExistsDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_persistKey(t *testing.T) {
	p256 := elliptic.P256()
	var (
		tmpDir     = t.TempDir()
		keyPath    = filepath.Join(tmpDir, "ocis", "testKey")
		rsaPk, _   = rsa.GenerateKey(rand.Reader, 2048)
		ecdsaPk, _ = ecdsa.GenerateKey(p256, rand.Reader)
	)

	type args struct {
		keyName string
		l       log.Logger
		pk      interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "writes a private key (rsa) to the specified location",
			args: args{
				keyName: keyPath,
				l:       log.NewLogger(),
				pk:      rsaPk,
			},
		},
		{
			name: "writes a private key (ecdsa) to the specified location",
			args: args{
				keyName: keyPath,
				l:       log.NewLogger(),
				pk:      ecdsaPk,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := persistKey(tt.args.keyName, tt.args.l, tt.args.pk); err != nil {
				t.Error(err)
			}
		})

		// side effect: tt.args.keyName is created
		if _, err := os.Stat(tt.args.keyName); err != nil {
			t.Errorf("persistKey() error = %v", err)
		}
	}
}

func Test_persistCertificate(t *testing.T) {
	p256 := elliptic.P256()
	var (
		tmpDir     = t.TempDir()
		certPath   = filepath.Join(tmpDir, "ocis", "testCert")
		rsaPk, _   = rsa.GenerateKey(rand.Reader, 2048)
		ecdsaPk, _ = ecdsa.GenerateKey(p256, rand.Reader)
	)

	type args struct {
		certName string
		l        log.Logger
		pk       interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "store a certificate with an rsa private key",
			args: args{
				certName: certPath,
				l:        log.NewLogger(),
				pk:       rsaPk,
			},
			wantErr: false,
		},
		{
			name: "store a certificate with an ecdsa private key",
			args: args{
				certName: certPath,
				l:        log.NewLogger(),
				pk:       ecdsaPk,
			},
			wantErr: false,
		},
		{
			name: "should fail",
			args: args{
				certName: certPath,
				l:        log.NewLogger(),
				pk:       42,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				if err := persistCertificate(tt.args.certName, tt.args.l, tt.args.pk); err != nil {
					if !tt.wantErr {
						t.Error(err)
					}
				}
			})

			// side effect: tt.args.keyName is created
			if _, err := os.Stat(tt.args.certName); err != nil {
				t.Errorf("persistCertificate() error = %v", err)
			}
		})
	}
}
