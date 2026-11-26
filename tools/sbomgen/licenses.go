package main

// License and copyright extraction from the text files a package ships with.

import (
	"regexp"
	"sort"
	"strings"

	"github.com/google/licensecheck"
)

// minLicenseCoverage is the share of a license file's text that must match a known license
const minLicenseCoverage = 75.0

// maxCopyrightLines bounds the copyright statement kept for one package
const maxCopyrightLines = 5

var licenseFileName = regexp.MustCompile(`(?i)^(licen[cs]e|copying|unlicen[cs]e)([-_.].*)?$`)

// noticeFileName matches NOTICE files: they carry the copyright statements of Apache-licensed projects
var noticeFileName = regexp.MustCompile(`(?i)^notice([-_.].*)?$`)

// classify returns the SPDX ids of the licenses a license text matches.
func classify(text []byte) []string {
	cov := licensecheck.Scan(text)
	if cov.Percent < minLicenseCoverage {
		return nil
	}
	var ids []string
	for _, m := range cov.Match {
		if !m.IsURL {
			ids = append(ids, m.ID)
		}
	}
	return ids
}

// joinAnd combines license ids found in the license files of one package.
func joinAnd(ids map[string]bool) string {
	if len(ids) == 0 {
		return noAssertion
	}
	list := make([]string, 0, len(ids))
	for id := range ids {
		list = append(list, id)
	}
	sort.Strings(list)
	return strings.Join(list, " AND ")
}

var (
	// a line that starts with a copyright mark; a bare "(c)" only counts with a year, because
	// "(c) You must retain ..." is an item of the Apache license text
	copyrightLine = regexp.MustCompile(`(?i)^[\s#*/>;-]*((?:(?:original code|portions|parts)\s+)?(?:copyright\b|©).*|\(c\)\s*(?:19|20)\d\d.*)$`)
	// wording of license templates and of the standard license texts themselves, not of a project
	copyrightBoilerplate = regexp.MustCompile(`(?i)\b(notices?|holders?|owners?|licen[cs]e[sd]?|statement|laws?|and related rights|yyyy|year|name of)\b|[<\[{]\s*(year|yyyy|name|owner|fullname)|fsf\.org|franklin street|temple place`)
	spaces               = regexp.MustCompile(`\s+`)
)

// copyrightLines extracts the copyright statements of a license or notice text.
func copyrightLines(text string) []string {
	var out []string
	seen := map[string]bool{}
	for _, line := range strings.Split(text, "\n") {
		m := copyrightLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		stmt := strings.TrimSpace(spaces.ReplaceAllString(m[1], " "))
		if len(stmt) < len("Copyright 1") || len(stmt) > 200 || copyrightBoilerplate.MatchString(stmt) || seen[stmt] {
			continue
		}
		seen[stmt] = true
		out = append(out, stmt)
	}
	return out
}

// mergeCopyright joins the statements found in several files, without duplicates and bounded in size.
func mergeCopyright(lists ...[]string) string {
	var out []string
	seen := map[string]bool{}
	for _, l := range lists {
		for _, s := range l {
			if !seen[s] && len(out) < maxCopyrightLines {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	if len(out) == 0 {
		return noAssertion
	}
	return strings.Join(out, "\n")
}
