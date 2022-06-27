/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

type Options struct {
	BaseDN                  string
	AdminDN                 string
	AllowLocalAnonymousBind bool

	DefaultCompany    string
	DefaultMailDomain string

	TemplateExtraVars      map[string]interface{}
	TemplateEngineDisabled bool
	TemplateDebug          bool

	templateBasePath string
}
