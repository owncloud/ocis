package checks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	revastore "github.com/owncloud/reva/v2/pkg/store"
)

// NewNatsCheck checks the reachability of a nats server.
func NewNatsCheck(natsCluster string, options ...nats.Option) func(context.Context) error {
	return func(_ context.Context) error {
		n, err := nats.Connect(natsCluster, options...)
		if err != nil {
			return fmt.Errorf("could not connect to nats server: %v", err)
		}
		defer n.Close()
		if n.Status() != nats.CONNECTED {
			return fmt.Errorf("nats server not connected")
		}
		return nil
	}
}

// natsCheckTimeout limits a probe when the probe context has no deadline.
const natsCheckTimeout = 2 * time.Second

// NewNatsCheckWithConfig checks the reachability of the nats server at address.
// It authenticates and uses TLS like the service connections configured by cfg.
// The check uses address and ignores the endpoint of cfg. An empty address means
// that the service does not use NATS. The check then always succeeds.
//
// The check keeps one connection and uses it again for each probe. A new connection
// for each probe reads the .creds file each time. A probe can then read a rotated file
// before the write is complete. Each connection also counts against the connection
// limit of the NATS account. The check connects again only after the connection closes.
// The connection reads the .creds file again on each reconnect.
//
// The connection status stays CONNECTED until the client notices a lost server.
// Thus each probe sends a PING to the server and waits for the PONG.
func NewNatsCheckWithConfig(address string, cfg stream.NatsConfig) func(context.Context) error {
	if address == "" {
		return func(context.Context) error { return nil }
	}
	c := &natsCheck{address: address, cfg: cfg}
	return c.check
}

// NewNatsStoreCheck checks the reachability of the first node of a store.
// The check always succeeds when the store type is not a NATS type or when the store has no nodes.
func NewNatsStoreCheck(storeType string, nodes []string, cfg stream.NatsConfig) func(context.Context) error {
	if !revastore.UsesNats(storeType) || len(nodes) == 0 {
		return func(context.Context) error { return nil }
	}
	return NewNatsCheckWithConfig(nodes[0], cfg)
}

type natsCheck struct {
	address string
	cfg     stream.NatsConfig

	mu   sync.Mutex
	conn *nats.Conn
	// dial is not nil while a connect is in progress. The connect closes it when it ends.
	dial chan struct{}
	// dialErr is the error of the last connect.
	dialErr error
}

func (c *natsCheck) check(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, natsCheckTimeout)
		defer cancel()
	}

	conn, err := c.connection(ctx)
	if err != nil {
		return err
	}
	if err := conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("nats server not reachable: %v", err)
	}
	return nil
}

// connection returns the open connection. It starts a connect when no connection is open.
// The connect runs without the lock, so a slow server does not block the other probes.
// Each probe waits for the connect only until its context ends.
func (c *natsCheck) connection(ctx context.Context) (*nats.Conn, error) {
	c.mu.Lock()
	if c.conn != nil && !c.conn.IsClosed() {
		conn := c.conn
		c.mu.Unlock()
		return conn, nil
	}
	if c.dial == nil {
		c.dial = make(chan struct{})
		go c.connect(c.dial)
	}
	dial := c.dial
	c.mu.Unlock()

	select {
	case <-dial:
	case <-ctx.Done():
		return nil, fmt.Errorf("could not connect to nats server: %v", ctx.Err())
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil || c.conn.IsClosed() {
		return nil, c.dialErr
	}
	return c.conn, nil
}

func (c *natsCheck) connect(done chan struct{}) {
	conn, err := c.dialNats()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn, c.dialErr, c.dial = conn, err, nil
	close(done)
}

func (c *natsCheck) dialNats() (*nats.Conn, error) {
	options, err := stream.NatsOptions(c.cfg)
	if err != nil {
		return nil, fmt.Errorf("could not configure the nats connection: %v", err)
	}
	conn, err := nats.Connect(c.address, options...)
	if err != nil {
		return nil, fmt.Errorf("could not connect to nats server: %v", err)
	}
	return conn, nil
}
