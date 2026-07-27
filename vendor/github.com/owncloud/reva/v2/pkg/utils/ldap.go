// Copyright 2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package utils

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/owncloud/reva/v2/pkg/logger"
	ldapReconnect "github.com/owncloud/reva/v2/pkg/utils/ldap"
	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
)

// LDAPConn holds the basic parameter for setting up an
// LDAP connection.
type LDAPConn struct {
	URI          string `mapstructure:"uri"`
	Insecure     bool   `mapstructure:"insecure"`
	CACert       string `mapstructure:"cacert"`
	BindDN       string `mapstructure:"bind_username"`
	BindPassword string `mapstructure:"bind_password"`

	// PoolEnabled switches GetLDAPClientWithPool callers to a bounded connection pool instead of
	// the single long-lived reconnecting connection. Off by default.
	PoolEnabled bool `mapstructure:"pool_enabled"`
	// PoolSize caps the number of concurrently open pooled connections. Defaults to 5 when unset.
	PoolSize int `mapstructure:"pool_size"`
	// PoolCheckoutTimeout bounds how long a checkout waits for a connection to become available
	// once the pool is at PoolSize. Defaults to 30s when unset.
	PoolCheckoutTimeout time.Duration `mapstructure:"pool_checkout_timeout"`
}

// tlsConfigFromLDAPConn builds the *tls.Config shared by all GetLDAPClient* constructors below.
func tlsConfigFromLDAPConn(c *LDAPConn) (*tls.Config, error) {
	if c.Insecure {
		logger.New().Warn().Msg("SSL Certificate verification is disabled. This is strongly discouraged for production environments.")
		return &tls.Config{
			//nolint:gosec // We need the ability to run with "insecure" (dev/testing)
			InsecureSkipVerify: true,
		}, nil
	}
	if c.CACert != "" {
		pemBytes, err := os.ReadFile(c.CACert)
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading LDAP CA Cert '%s.'", c.CACert)
		}
		rpool, _ := x509.SystemCertPool()
		rpool.AppendCertsFromPEM(pemBytes)
		return &tls.Config{
			RootCAs: rpool,
		}, nil
	}
	return nil, nil
}

// GetLDAPClientWithReconnect initializes a long-lived LDAP connection that
// automatically reconnects on connection errors. It allows to set TLS options
// e.g. to add trusted Certificates or disable Certificate verification
func GetLDAPClientWithReconnect(c *LDAPConn) (ldap.Client, error) {
	tlsConf, err := tlsConfigFromLDAPConn(c)
	if err != nil {
		return nil, err
	}

	conn := ldapReconnect.NewLDAPWithReconnect(
		ldapReconnect.Config{
			URI:          c.URI,
			BindDN:       c.BindDN,
			BindPassword: c.BindPassword,
			TLSConfig:    tlsConf,
		},
	)
	return conn, nil
}

// GetLDAPClientWithPool initializes a bounded pool of authenticated LDAP connections, dialed and
// bound lazily on first use. It is a drop-in alternative to GetLDAPClientWithReconnect intended for
// backends that need to serve concurrent requests without serializing on a single connection.
func GetLDAPClientWithPool(c *LDAPConn) (ldap.Client, error) {
	tlsConf, err := tlsConfigFromLDAPConn(c)
	if err != nil {
		return nil, err
	}

	pool := ldapReconnect.NewLDAPPool(
		ldapReconnect.Config{
			URI:                 c.URI,
			BindDN:              c.BindDN,
			BindPassword:        c.BindPassword,
			TLSConfig:           tlsConf,
			PoolSize:            c.PoolSize,
			PoolCheckoutTimeout: c.PoolCheckoutTimeout,
		},
		logger.New(),
	)
	return pool, nil
}

// GetLDAPClientFromConfig returns a connected ldap.Client for c: a bounded pool when
// c.PoolEnabled, otherwise the single long-lived reconnecting connection.
func GetLDAPClientFromConfig(c *LDAPConn) (ldap.Client, error) {
	if c.PoolEnabled {
		return GetLDAPClientWithPool(c)
	}
	return GetLDAPClientWithReconnect(c)
}

// GetLDAPClientForAuth initializes an LDAP connection. The connection is not authenticated
// when returned. The main purpose for GetLDAPClientForAuth is to get and LDAP connection that
// can be used to issue a single bind request to authenticate a user.
func GetLDAPClientForAuth(c *LDAPConn) (ldap.Client, error) {
	tlsConf, err := tlsConfigFromLDAPConn(c)
	if err != nil {
		return nil, err
	}
	l, err := ldap.DialURL(c.URI, ldap.DialWithTLSConfig(tlsConf))
	if err != nil {
		return nil, err
	}

	return l, nil
}
