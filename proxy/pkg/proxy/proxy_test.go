package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/config/defaults"
)

type matchertest struct {
	method, endpoint, target string
	matches                  bool
}

func TestPrefixRouteMatcher(t *testing.T) {
	cfg := defaults.DefaultConfig()
	cfg.Policies = defaults.DefaultPolicies()
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
	cfg := defaults.DefaultConfig()
	cfg.Policies = defaults.DefaultPolicies()
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
	cfg := defaults.DefaultConfig()
	cfg.Policies = defaults.DefaultPolicies()
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

func TestDirectorSelectionDirector(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))
	defer svr.Close()

	p := NewMultiHostReverseProxy(Config(&config.Config{
		PolicySelector: &config.PolicySelector{
			Static: &config.StaticSelectorConf{
				Policy: "default",
			},
		},
	}))
	p.AddHost("default", &url.URL{Host: "ocdav"}, config.Route{Type: config.PrefixRoute, Method: "", Endpoint: "/dav", Backend: "ocdav"})
	p.AddHost("default", &url.URL{Host: "ocis-webdav"}, config.Route{Type: config.PrefixRoute, Method: "REPORT", Endpoint: "/dav", Backend: "ocis-webdav"})

	table := []matchertest{
		{method: "PROPFIND", endpoint: "/dav/files/demo/", target: "ocdav"},
		{method: "REPORT", endpoint: "/dav/files/demo/", target: "ocis-webdav"},
	}

	for _, test := range table {
		r := httptest.NewRequest(http.MethodGet, "/dav/files/demo/", nil)
		p.directorSelectionDirector(r)
		if r.Host != test.target {
			t.Errorf("TestDirectorSelectionDirector got host %s expected %s", r.Host, test.target)

		}
	}
}
