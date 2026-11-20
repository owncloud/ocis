package command

import (
	"net/http"
	"testing"

	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy"
)

func TestNewOIDCHTTPClientKeepsConnectionsAlive(t *testing.T) {
	tr, ok := newOIDCHTTPClient(false).Transport.(*http.Transport)
	if !ok {
		t.Fatal("the transport must be an *http.Transport")
	}
	if tr.DisableKeepAlives {
		t.Fatal("DisableKeepAlives must be false")
	}
	if tr.MaxIdleConnsPerHost != 100 {
		t.Fatalf("MaxIdleConnsPerHost = %d, want 100", tr.MaxIdleConnsPerHost)
	}
	if tr.IdleConnTimeout != proxy.IdleConnTimeout {
		t.Fatalf("IdleConnTimeout = %v, want %v", tr.IdleConnTimeout, proxy.IdleConnTimeout)
	}
	if tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify must be false when insecure is false")
	}
	if !newOIDCHTTPClient(true).Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify must be true when insecure is true")
	}
}
