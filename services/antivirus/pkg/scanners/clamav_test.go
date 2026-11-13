package scanners_test

import (
	"io"
	"net"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/scanners"
)

// hungClamd accepts a single connection and drains it but never responds,
// reproducing the OCISDEV-1117 backend that used to block ClamAV.Scan forever
// on the result receive.
func hungClamd(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = io.Copy(io.Discard, conn)
	}()

	return ln
}

// wedgedClamd accepts a connection and never reads from it at all, so once the
// OS socket buffer fills, ScanStream's own writes (not just the result
// receive) block - a hang the plain "withholds a reply" case above doesn't
// exercise.
func wedgedClamd(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		t.Cleanup(func() { conn.Close() })
		// deliberately never read or write
	}()

	return ln
}

// drainUntilQuiet reads and discards until the peer has sent nothing for
// quietFor - ScanStream never half-closes after its own EOF marker, so a
// plain io.Copy(discard, conn) would never return on its own.
func drainUntilQuiet(conn net.Conn, quietFor time.Duration) {
	buf := make([]byte, 4096)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(quietFor))
		if _, err := conn.Read(buf); err != nil {
			return
		}
	}
}

// closingClamd accepts a connection, drains it, then closes without ever
// sending a result line - go-clamd's reader goroutine closes its channel
// without a value in this case, which used to be dereferenced as *ScanResult.
func closingClamd(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		drainUntilQuiet(conn, 200*time.Millisecond)
		conn.Close()
	}()

	return ln
}

func TestClamAV_Scan_TimesOutOnHungBackend(t *testing.T) {
	ln := hungClamd(t)
	scanner := scanners.NewClamAV("tcp://"+ln.Addr().String(), 100*time.Millisecond)

	start := time.Now()
	_, err := scanner.Scan(scanners.Input{Body: strings.NewReader("test content")})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 2*time.Second, "Scan must return once its timeout elapses, not hang forever")
}

func TestClamAV_Scan_TimesOutOnWedgedBackend(t *testing.T) {
	ln := wedgedClamd(t)
	scanner := scanners.NewClamAV("tcp://"+ln.Addr().String(), 100*time.Millisecond)

	// Large enough to exceed typical OS socket buffers so ScanStream's writes
	// actually block on a peer that never reads, not just its result receive.
	body := strings.NewReader(strings.Repeat("x", 8*1024*1024))

	start := time.Now()
	_, err := scanner.Scan(scanners.Input{Body: body})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 5*time.Second, "Scan must return once its timeout elapses even if the backend never reads the stream")
}

func TestClamAV_Scan_ConnectionClosedWithoutResult(t *testing.T) {
	ln := closingClamd(t)
	scanner := scanners.NewClamAV("tcp://"+ln.Addr().String(), 5*time.Second)

	_, err := scanner.Scan(scanners.Input{Body: strings.NewReader("test content")})

	assert.Error(t, err, "a connection closed before any result line arrives must be a plain error, not a nil-pointer panic")
}

func TestClamAV_Scan_DoesNotLeakAbortWatcherOnSuccess(t *testing.T) {
	// A fake clamd that behaves normally: drains each stream and replies once,
	// one connection per Scan call below.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				drainUntilQuiet(c, 200*time.Millisecond)
				_, _ = c.Write([]byte("stream: OK\n"))
			}(conn)
		}
	}()

	scanner := scanners.NewClamAV("tcp://"+ln.Addr().String(), 5*time.Second)

	before := runtime.NumGoroutine()
	for i := 0; i < 20; i++ {
		_, _ = scanner.Scan(scanners.Input{Body: strings.NewReader("test content")})
	}
	// let any still-finishing background goroutines settle
	for i := 0; i < 50; i++ {
		if runtime.NumGoroutine() <= before+2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	after := runtime.NumGoroutine()

	assert.LessOrEqual(t, after, before+2, "successful scans must not leak the go-clamd abort-watcher goroutine")
}
