package proxy

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBackendTransportSettings(t *testing.T) {
	tr := newBackendTransport(&tls.Config{MinVersion: tls.VersionTLS12})

	if tr.MaxIdleConnsPerHost != 100 {
		t.Fatalf("MaxIdleConnsPerHost = %d, want 100", tr.MaxIdleConnsPerHost)
	}
	if tr.MaxIdleConns != 0 {
		t.Fatalf("MaxIdleConns = %d, want 0 (no global limit)", tr.MaxIdleConns)
	}
	if tr.IdleConnTimeout != IdleConnTimeout {
		t.Fatalf("IdleConnTimeout = %v, want %v", tr.IdleConnTimeout, IdleConnTimeout)
	}
	if IdleConnTimeout >= 60*time.Second {
		t.Fatalf("IdleConnTimeout = %v, want a value below the usual 60 s idle timeout of a load balancer", IdleConnTimeout)
	}
}

// TestBackendTransportReusesConnections sends two bursts of parallel requests.
// The second burst must reuse the idle connections of the first burst.
// The default of Go keeps 2 idle connections for each host, so the second burst
// would open new connections.
func TestBackendTransportReusesConnections(t *testing.T) {
	const parallel = 20

	var opened atomic.Int64
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond) // keep the requests parallel
		_, _ = w.Write([]byte("ok"))
	}))
	srv.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			opened.Add(1)
		}
	}
	srv.Start()
	defer srv.Close()

	tr := newBackendTransport(&tls.Config{MinVersion: tls.VersionTLS12})
	tr.Proxy = nil // do not depend on the proxy variables of the test host
	client := &http.Client{Transport: tr}

	burst := func() {
		var wg sync.WaitGroup
		for i := 0; i < parallel; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := client.Get(srv.URL)
				if err != nil {
					t.Errorf("request failed: %v", err)
					return
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}()
		}
		wg.Wait()
	}

	burst()
	afterFirst := opened.Load()
	burst()
	newInSecond := opened.Load() - afterFirst

	if newInSecond > 2 {
		t.Fatalf("the second burst opened %d new connections, want at most 2 (first burst: %d)", newInSecond, afterFirst)
	}
}
