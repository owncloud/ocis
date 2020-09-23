package proxy

import (
	"net/url"
	"testing"

	"github.com/owncloud/ocis/proxy/pkg/config"
)

func TestPrefixRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar"
	u, _ := url.Parse("/foobar/baz/some/url")

	matched := p.prefixRouteMatcher(endpoint, *u)
	if !matched {
		t.Errorf("Endpoint %s and URL %s should match", endpoint, u.String())
	}
}

func TestQueryRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar?parameter=true"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.queryRouteMatcher(endpoint, *u)
	if !matched {
		t.Errorf("Endpoint %s and URL %s should match", endpoint, u.String())
	}
}

func TestQueryRouteMatcherWithoutParameters(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.queryRouteMatcher(endpoint, *u)
	if matched {
		t.Errorf("Endpoint %s and URL %s should not match", endpoint, u.String())
	}
}

func TestQueryRouteMatcherWithDifferingParameters(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar?parameter=false"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.queryRouteMatcher(endpoint, *u)
	if matched {
		t.Errorf("Endpoint %s and URL %s should not match", endpoint, u.String())
	}
}

func TestQueryRouteMatcherWithMultipleDifferingParameters(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar?parameter=false&other=true"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.queryRouteMatcher(endpoint, *u)
	if matched {
		t.Errorf("Endpoint %s and URL %s should not match", endpoint, u.String())
	}
}

func TestQueryRouteMatcherWithMultipleParameters(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "/foobar?parameter=false&other=true"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=false&other=true")

	matched := p.queryRouteMatcher(endpoint, *u)
	if !matched {
		t.Errorf("Endpoint %s and URL %s should match", endpoint, u.String())
	}
}

func TestRegexRouteMatcher(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := ".*some\\/url.*parameter=true"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.regexRouteMatcher(endpoint, *u)
	if !matched {
		t.Errorf("Endpoint %s and URL %s should match", endpoint, u.String())
	}
}

func TestRegexRouteMatcherWithInvalidPattern(t *testing.T) {
	cfg := config.New()
	p := NewMultiHostReverseProxy(Config(cfg))

	endpoint := "([\\])\\w+"
	u, _ := url.Parse("/foobar/baz/some/url?parameter=true")

	matched := p.regexRouteMatcher(endpoint, *u)
	if matched {
		t.Errorf("Endpoint %s and URL %s should not match", endpoint, u.String())
	}
}
