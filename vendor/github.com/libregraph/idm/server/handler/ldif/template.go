/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/libregraph/idm"
)

const formatAsFileSizeLimit int64 = 1024 * 1024

func TemplateFuncs(m map[string]interface{}, options *Options) template.FuncMap {
	defaults := map[string]interface{}{
		"Company":    "Default",
		"BaseDN":     idm.DefaultLDAPBaseDN,
		"MailDomain": idm.DefaultMailDomain,
	}
	for k, v := range defaults {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}
	if options != nil {
		if options.BaseDN != "" {
			m["BaseDN"] = options.BaseDN
		}
		if options.DefaultCompany != "" {
			m["Company"] = options.DefaultCompany
		}
		if options.DefaultMailDomain != "" {
			m["MailDomain"] = options.DefaultMailDomain
		}
		for k, v := range options.TemplateExtraVars {
			m[k] = v
		}
	}

	autoIncrement := uint64(1000)
	if v, ok := m["AutoIncrementMin"]; ok {
		autoIncrement = v.(uint64)
	}

	basePath := options.templateBasePath

	return template.FuncMap{
		"WithCompany": func(value string) string {
			m["Company"] = value
			return ""
		},
		"WithBaseDN": func(value string) string {
			m["BaseDN"] = value
			return ""
		},
		"WithMailDomain": func(value string) string {
			m["MailDomain"] = value
			return ""
		},
		"AutoIncrement": func(values ...uint64) uint64 {
			if len(values) > 0 {
				autoIncrement = values[0]
			} else {
				autoIncrement++
			}
			return autoIncrement
		},
		"formatAsBase64": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},
		"formatAsFileBase64": func(fn string) (string, error) {
			if basePath == "" {
				return "", fmt.Errorf("LDIF template fromFile failed, no base path")
			}
			fn = filepath.Clean(fn)
			if !filepath.IsAbs(fn) {
				fn = filepath.Join(basePath, fn)
			}
			fn, err := filepath.Abs(fn)
			if err != nil {
				return "", err
			}
			// NOTE(longsleep): Poor man base path check, should work well enough on Linux.
			// See https://github.com/golang/go/issues/18358 for details.
			if !strings.HasPrefix(fn, strings.TrimRight(basePath, "/")+"/") {
				return "", fmt.Errorf("LDIF template formatAsFile %s outside of %s is not allowed", fn, basePath)
			}

			f, err := os.Open(fn)
			if err != nil {
				return "", fmt.Errorf("LDIF template formatAsFile open failed with error: %w", err)
			}
			defer f.Close()

			reader := io.LimitReader(f, formatAsFileSizeLimit+1)

			var buf bytes.Buffer
			encoder := base64.NewEncoder(base64.StdEncoding, &buf)
			n, err := io.Copy(encoder, reader)
			if err != nil {
				return "", fmt.Errorf("LDIF template formatAsFile error: %w", err)
			}
			if n > formatAsFileSizeLimit {
				return "", fmt.Errorf("LDIF template formatAsFile size limit exceeded: %s", fn)
			}

			return buf.String(), nil
		},
	}
}
