/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/go-ldap/ldap/v3"
	"github.com/go-ldap/ldif"
	"github.com/spacewander/go-suffix-tree"
)

// parseLDIFFile opens the named file for reading and parses it as LDIF.
func parseLDIFFile(fn string, options *Options) (*ldif.LDIF, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader

	if options.TemplateEngineDisabled {
		r = f
	} else {
		r, err = parseLDIFTemplate(f, options, nil)
		if err != nil {
			return nil, err
		}
	}

	return parseLDIF(r, options)
}

// parseLDIFDirectory opens all ldif files in the given path in sorted order,
// cats them all together and parses the result as LDIF.
func parseLDIFDirectory(pn string, options *Options) (*ldif.LDIF, []error, error) {
	matches, err := filepath.Glob(filepath.Join(pn, "*.ldif"))
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i] < matches[j]
	})

	var buf bytes.Buffer
	var matchErrors []error
	for _, match := range matches {
		err = func() error {
			f, openErr := os.Open(match)
			if openErr != nil {
				matchErrors = append(matchErrors, fmt.Errorf("file read error: %w", openErr))
				return nil
			}
			defer f.Close()
			if options.TemplateEngineDisabled {
				_, copyErr := io.Copy(&buf, f)
				if copyErr != nil {
					return fmt.Errorf("file read error: %w", copyErr)
				}
			} else {
				p, parseErr := parseLDIFTemplate(f, options, nil)
				if parseErr != nil {
					matchErrors = append(matchErrors, fmt.Errorf("parse error in %s: %w", match, parseErr))
					return nil
				}
				_, copyErr := io.Copy(&buf, p)
				if copyErr != nil {
					return fmt.Errorf("template read error: %w", copyErr)
				}
			}
			buf.WriteString("\n\n")
			return nil
		}()
		if err != nil {
			return nil, matchErrors, err
		}
	}

	l, err := parseLDIF(&buf, options)
	return l, matchErrors, err
}

// parseLDIFTemplate exectues the provided text template and then parses the
// result as LDIF.
func parseLDIFTemplate(r io.Reader, options *Options, m map[string]interface{}) (io.Reader, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" && t[0] == '#' {
			// Ignore commented lines.
			continue
		}
		text = append(text, scanner.Text())
	}

	if m == nil {
		m = make(map[string]interface{})
	}
	tpl, err := template.New("tpl").Funcs(TemplateFuncs(m, options)).Parse(strings.Join(text, "\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse LDIF template: %w", err)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, m)
	if err != nil {
		return nil, fmt.Errorf("failed to process LDIF template: %w", err)
	}

	if options.TemplateDebug {
		fmt.Println("---\n", buf.String(), "\n----")
	}

	return &buf, nil
}

func parseLDIF(r io.Reader, options *Options) (*ldif.LDIF, error) {
	l := &ldif.LDIF{}
	err := ldif.Unmarshal(r, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// treeFromLDIF makes a tree out of the provided LDIF and if index is not nil,
// also indexes each entry in the provided index.
func treeFromLDIF(l *ldif.LDIF, index Index, options *Options) (*suffix.Tree, error) {
	t := suffix.NewTree()

	// NOTE(longsleep): Create in memory tree records from LDIF data.
	var entry *ldap.Entry
	for _, entryRecord := range l.Entries {
		if entryRecord == nil || entryRecord.Entry == nil {
			// NOTE(longsleep): We don't use l.AllEntries as "nil" records can happen.
			continue
		}
		entry = entryRecord.Entry
		e := &ldifEntry{
			Entry: &ldap.Entry{
				DN: strings.ToLower(entry.DN),
			},
		}
		for _, a := range entry.Attributes {
			switch strings.ToLower(a.Name) {
			case "userpassword":
				// Don't include the password in the normal attributes.
				e.UserPassword = &ldap.EntryAttribute{
					Name:   a.Name,
					Values: a.Values,
				}
			default:
				// Append it.
				e.Entry.Attributes = append(e.Entry.Attributes, &ldap.EntryAttribute{
					Name:   a.Name,
					Values: a.Values,
				})
			}
			if index != nil {
				// Index equalityMatch.
				index.Add(a.Name, "eq", a.Values, e)
				// Index present.
				index.Add(a.Name, "pres", []string{""}, e)
				// Index substrings.
				index.Add(a.Name, "sub", a.Values, e)
			}
		}
		v, ok := t.Insert([]byte(e.DN), e)
		if !ok || v != nil {
			return nil, fmt.Errorf("duplicate dn value: %s", e.DN)
		}
	}

	return t, nil
}
