package scanners_test

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/scanners"
)

// fakeICAPServer is a minimal ICAP server which answers OPTIONS with a 200 and
// every REQMOD with a 204. If previewBytes is set, OPTIONS advertises a preview
// and a REQMOD whose body did not fit into the preview is first answered with a
// 100 Continue.
type fakeICAPServer struct {
	addr      string
	continues atomic.Int64
}

func startFakeICAPServer(t *testing.T, previewBytes int) *fakeICAPServer {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = l.Close() })

	s := &fakeICAPServer{addr: l.Addr().String()}

	options := "ICAP/1.0 200 OK\r\nMethods: REQMOD\r\nISTag: test\r\n"
	if previewBytes > 0 {
		options += fmt.Sprintf("Preview: %d\r\n", previewBytes)
	}
	options += "\r\n"

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
				continued := false
				tmp := make([]byte, 4096)
				for {
					n, err := conn.Read(tmp)
					if n > 0 {
						buf = append(buf, tmp[:n]...)
					}
					req := string(buf)
					switch {
					case strings.HasPrefix(req, "OPTIONS") && strings.HasSuffix(req, "\r\n\r\n"):
						_, _ = conn.Write([]byte(options))
						return
					case continued && strings.HasSuffix(req, "0\r\n\r\n"):
						// the remaining body after a 100 Continue
						_, _ = conn.Write([]byte("ICAP/1.0 204 No Content\r\nISTag: test\r\n\r\n"))
						return
					case strings.HasPrefix(req, "REQMOD") && strings.HasSuffix(req, "0; ieof\r\n\r\n"):
						// the whole body fitted into the preview
						_, _ = conn.Write([]byte("ICAP/1.0 204 No Content\r\nISTag: test\r\n\r\n"))
						return
					case strings.HasPrefix(req, "REQMOD") && strings.HasSuffix(req, "0\r\n\r\n"):
						if previewBytes == 0 {
							_, _ = conn.Write([]byte("ICAP/1.0 204 No Content\r\nISTag: test\r\n\r\n"))
							return
						}
						// the preview did not contain the whole body, ask for the rest
						continued = true
						s.continues.Add(1)
						buf = nil
						_, _ = conn.Write([]byte("ICAP/1.0 100 Continue\r\n\r\n"))
					}
					if err != nil {
						return
					}
				}
			}(conn)
		}
	}()

	return s
}

// The antivirus service shares one scanner between all of its workers, so Scan
// must not share connection state between concurrent calls.
func TestICAP_ScanIsSafeForConcurrentUse(t *testing.T) {
	const (
		workers         = 10
		scansPerWorker  = 5
		body            = "hello"
		previewTooSmall = 2
	)

	tests := map[string]struct {
		previewBytes  int
		wantContinues int64
	}{
		"without preview":                 {previewBytes: 0, wantContinues: 0},
		"with preview containing body":    {previewBytes: 1024, wantContinues: 0},
		"with preview requiring continue": {previewBytes: previewTooSmall, wantContinues: workers * scansPerWorker},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := startFakeICAPServer(t, tt.previewBytes)

			scanner, err := scanners.NewICAP("icap://"+server.addr, "avscan", 5*time.Second)
			require.NoError(t, err)

			var wg sync.WaitGroup
			errs := make(chan error, workers*scansPerWorker)
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < scansPerWorker; j++ {
						res, err := scanner.Scan(scanners.Input{
							Url:  "http://example.test/file",
							Name: "file.txt",
							Size: int64(len(body)),
							Body: bytes.NewReader([]byte(body)),
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
			assert.Equal(t, tt.wantContinues, server.continues.Load())
		})
	}
}
