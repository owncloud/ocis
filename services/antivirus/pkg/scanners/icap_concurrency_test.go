package scanners_test

import (
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/scanners"
)

// startFakeICAPServer starts a minimal ICAP server which answers OPTIONS with a
// 200 and every REQMOD with a 204. It returns the address it listens on.
func startFakeICAPServer(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = l.Close() })

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

				var buf []byte
				tmp := make([]byte, 4096)
				for {
					n, err := conn.Read(tmp)
					if n > 0 {
						buf = append(buf, tmp[:n]...)
					}
					req := string(buf)
					switch {
					case strings.HasPrefix(req, "OPTIONS") && strings.HasSuffix(req, "\r\n\r\n"):
						_, _ = conn.Write([]byte("ICAP/1.0 200 OK\r\nMethods: REQMOD\r\nISTag: test\r\n\r\n"))
						return
					case strings.HasPrefix(req, "REQMOD") && strings.HasSuffix(req, "0\r\n\r\n"):
						_, _ = conn.Write([]byte("ICAP/1.0 204 No Content\r\nISTag: test\r\n\r\n"))
						return
					}
					if err != nil {
						return
					}
				}
			}(conn)
		}
	}()

	return l.Addr().String()
}

// The antivirus service shares one scanner between all of its workers, so Scan
// must not share connection state between concurrent calls.
func TestICAP_ScanIsSafeForConcurrentUse(t *testing.T) {
	addr := startFakeICAPServer(t)

	scanner, err := scanners.NewICAP("icap://"+addr, "avscan", 5*time.Second)
	require.NoError(t, err)

	const workers = 10
	var wg sync.WaitGroup
	errs := make(chan error, workers*5)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				res, err := scanner.Scan(scanners.Input{
					Url:  "http://example.test/file",
					Name: "file.txt",
					Size: 5,
					Body: bytes.NewReader([]byte("hello")),
				})
				if err != nil {
					errs <- err
					continue
				}
				if res.Infected {
					errs <- assert.AnError
				}
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		assert.NoError(t, err)
	}
}
