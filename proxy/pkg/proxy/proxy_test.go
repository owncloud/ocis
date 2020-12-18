package proxy

import (
	"net/url"
	"testing"

	"github.com/owncloud/ocis/proxy/pkg/config"
)

type matchertest struct {
	endpoint, target string
	matches          bool
}

func TestPrefixRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	table := []matchertest{
		{endpoint: "/foobar", target: "/foobar/baz/some/url", matches: true},
		{endpoint: "/fobar", target: "/foobar/baz/some/url", matches: false},
	}

	for _, test := range table {
		u, _ := url.Parse(test.target)
		matched := p.prefixRouteMatcher(test.endpoint, *u)
		if matched != test.matches {
			t.Errorf("PrefixRouteMatcher returned %t expected %t for endpoint: %s and target %s",
				matched, test.matches, test.endpoint, u.String())
		}
	}
}

func TestQueryRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	table := []matchertest{
		{endpoint: "/foobar?parameter=true", target: "/foobar/baz/some/url?parameter=true", matches: true},
		{endpoint: "/foobar", target: "/foobar/baz/some/url?parameter=true", matches: false},
		{endpoint: "/foobar?parameter=false", target: "/foobar/baz/some/url?parameter=true", matches: false},
		{endpoint: "/foobar?parameter=false&other=true", target: "/foobar/baz/some/url?parameter=true", matches: false},
		{
			endpoint: "/foobar?parameter=false&other=true",
			target:   "/foobar/baz/some/url?parameter=false&other=true",
			matches:  true,
		},
		{endpoint: "/fobar", target: "/foobar", matches: false},
	}

	for _, test := range table {
		u, _ := url.Parse(test.target)
		matched := p.queryRouteMatcher(test.endpoint, *u)
		if matched != test.matches {
			t.Errorf("QueryRouteMatcher returned %t expected %t for endpoint: %s and target %s",
				matched, test.matches, test.endpoint, u.String())
		}
	}
}

func TestRegexRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	table := []matchertest{
		{endpoint: ".*some\\/url.*parameter=true", target: "/foobar/baz/some/url?parameter=true", matches: true},
		{endpoint: "([\\])\\w+", target: "/foobar/baz/some/url?parameter=true", matches: false},
	}

	for _, test := range table {
		u, _ := url.Parse(test.target)
		matched := p.regexRouteMatcher(test.endpoint, *u)
		if matched != test.matches {
			t.Errorf("RegexRouteMatcher returned %t expected %t for endpoint: %s and target %s",
				matched, test.matches, test.endpoint, u.String())
		}
	}
}

func TestSingleJoiningSlash(t *testing.T) {
	type test struct {
		a, b, result string
	}

	table := []test{
		{a: "a", b: "b", result: "a/b"},
		{a: "a/", b: "b", result: "a/b"},
		{a: "a", b: "/b", result: "a/b"},
		{a: "a/", b: "/b", result: "a/b"},
	}

	for _, test := range table {
		p := singleJoiningSlash(test.a, test.b)
		if p != test.result {
			t.Errorf("SingleJoiningSlash got %s expected %s", p, test.result)
		}
	}
}
