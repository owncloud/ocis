package icapclient

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"
	"syscall"
	"time"
)

// ICAPConnConfig is the configuration for the icap connection
type ICAPConnConfig struct {
	// Timeout is the maximum amount of time a connection will be kept open
	Timeout time.Duration
}

// ICAPConn is the one responsible for driving the transport layer operations. We have to explicitly deal with the connection because the ICAP protocol is aware of keep alive and reconnects.
type ICAPConn struct {
	tcp     net.Conn
	mu      sync.Mutex
	timeout time.Duration
}

// NewICAPConn creates a new connection to the icap server
func NewICAPConn(conf ICAPConnConfig) (*ICAPConn, error) {
	return &ICAPConn{
		timeout: conf.Timeout,
	}, nil
}

// Connect connects to the icap server
func (c *ICAPConn) Connect(ctx context.Context, address string) error {
	dialer := net.Dialer{Timeout: c.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}

	c.tcp = conn

	if dialer.Timeout == 0 {
		return nil
	}

	deadline := time.Now().UTC().Add(dialer.Timeout)

	if err := c.tcp.SetReadDeadline(deadline); err != nil {
		return err
	}

	if err := c.tcp.SetWriteDeadline(deadline); err != nil {
		return err
	}

	return nil
}

// Send sends a request to the icap server
func (c *ICAPConn) Send(in []byte) ([]byte, error) {
	if !c.ok() {
		return nil, syscall.EINVAL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	errChan := make(chan error)
	resChan := make(chan []byte)

	go func() {
		// send the message to the server
		_, err := c.tcp.Write(in)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		data := make([]byte, 0)

		for {
			tmp := make([]byte, 1096)

			// read the response from the server
			n, err := c.tcp.Read(tmp)

			// something went wrong while reading from the server,
			// send the error and exit the routine to prevent
			// sending the response to resChan
			if err != nil && err != io.EOF {
				errChan <- err
				return
			}

			// EOF detected, an entire message is received
			if err == io.EOF || n == 0 {
				break
			}

			data = append(data, tmp[:n]...)

			// explicitly breaking because the Read blocks for 100 continue message
			if bytes.Equal(data, []byte(icap100ContinueMsg)) {
				break
			}

			// EOF detected, double crlf indicates the end of the message
			if bytes.HasSuffix(data, []byte(doubleCRLF)) {
				break
			}

			// EOF detected, 204 no modifications and Double crlf indicate the end of the message
			if bytes.Contains(data, []byte(icap204NoModsMsg)) {
				break
			}
		}

		resChan <- data
	}()

	select {
	case err := <-errChan:
		return nil, err
	case res := <-resChan:
		return res, nil
	}
}

// Close closes the tcp connection
func (c *ICAPConn) Close() error {
	if !c.ok() {
		return syscall.EINVAL
	}

	return c.tcp.Close()
}

func (c *ICAPConn) ok() bool { return c != nil && c.tcp != nil }
