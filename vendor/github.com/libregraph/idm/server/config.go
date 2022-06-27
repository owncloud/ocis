/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Config bundles server configuration settings.
type Config struct {
	Logger logrus.FieldLogger

	LDAPHandler string

	LDAPListenAddr  string
	LDAPSListenAddr string

	TLSCertFile string
	TLSKeyFile  string

	LDAPBaseDN  string
	LDAPAdminDN string

	LDAPAllowLocalAnonymousBind bool

	BoltDBFile string

	LDIFMain   string
	LDIFConfig string

	LDIFDefaultCompany    string
	LDIFDefaultMailDomain string
	LDIFTemplateExtraVars map[string]interface{}

	Metrics prometheus.Registerer

	OnReady func(*Server)
}
