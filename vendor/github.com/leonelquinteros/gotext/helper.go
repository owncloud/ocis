/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`%\(([a-zA-Z0-9_]+)\)[.0-9]*[svTtbcdoqXxUeEfFgGp]`)

// SimplifiedLocale simplified locale like " en_US"/"de_DE "/en_US.UTF-8/zh_CN/zh_TW/el_GR@euro/... to en_US, de_DE, zh_CN, el_GR...
func SimplifiedLocale(lang string) string {
	// en_US/en_US.UTF-8/zh_CN/zh_TW/el_GR@euro/...
	if idx := strings.Index(lang, ":"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "@"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}
	return strings.TrimSpace(lang)
}

// FormatString applies text formatting only when needed to parse variables.
func FormatString(str string, vars ...any) string {
	if len(vars) > 0 {
		return FormatStringWithArgs(str, vars)
	}

	return str
}

// FormatStringWithArgs applies text formatting as a proxy for fmt.Sprintf.
func FormatStringWithArgs(str string, vars []any) string {
	return fmt.Sprintf(str, vars...)
}

// Appendf applies text formatting only when needed to parse variables.
func Appendf(b []byte, str string, vars ...any) []byte {
	if len(vars) > 0 {
		return fmt.Appendf(b, str, vars...)
	}

	return append(b, str...)
}

// NPrintf support named format
// NPrintf("%(name)s is Type %(type)s", map[string]any{"name": "Gotext", "type": "struct"})
func NPrintf(format string, params map[string]any) {
	f, p := parseSprintf(format, params)
	fmt.Printf(f, p...)
}

// Sprintf support named format
//
//	Sprintf("%(name)s is Type %(type)s", map[string]any{"name": "Gotext", "type": "struct"})
func Sprintf(format string, params map[string]any) string {
	f, p := parseSprintf(format, params)
	return fmt.Sprintf(f, p...)
}

func parseSprintf(format string, params map[string]any) (string, []any) {
	f, n := reformatSprintf(format)
	var p []any
	if len(n) > 0 {
		p = make([]any, len(n))
		for i, v := range n {
			p[i] = params[v]
		}
	}
	return f, p
}

func reformatSprintf(f string) (string, []string) {
	matches := re.FindAllStringSubmatchIndex(f, -1)
	ord := make([]string, 0, len(matches))
	pair := make([]int, 1, 1+2*len(matches))
	pair[0] = 0

	for _, match := range matches {
		ord = append(ord, f[match[2]:match[3]])
		pair = append(pair, match[2]-1, match[3]+1)
	}
	pair = append(pair, len(f))

	var out strings.Builder
	out.Grow(len(f))
	for n := 0; n < len(pair); n += 2 {
		out.WriteString(f[pair[n]:pair[n+1]])
	}

	return out.String(), ord
}
